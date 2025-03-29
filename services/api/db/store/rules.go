package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type RuleStore interface {
	CreateRule(ctx context.Context, params models.CreateRuleParams) error
	GetRules(ctx context.Context, params models.FilterRuleParams) (map[string]map[string][]models.Rule, error)
	GetRuleById(ctx context.Context, ruleId uuid.UUID) (models.Rule, error)
	GetRuleByIds(ctx context.Context, ruleIds []uuid.UUID) ([]models.Rule, error)
	UpdateRule(ctx context.Context, ruleId uuid.UUID, params models.UpdateRuleParams) error
	UpdateRulePriority(ctx context.Context, params models.UpdateRulePriorityParams) error
	DeleteRule(ctx context.Context, params models.DeleteRuleParams) error
}

func (s *appStore) CreateRule(ctx context.Context, params models.CreateRuleParams) error {
	filterConfig, err := json.Marshal(params.FilterConfig)
	if err != nil {
		return err
	}

	rule := models.Rule{
		ID:             params.Id,
		OrganizationId: params.OrganizationId,
		DatasetId:      params.DatasetId,
		Column:         params.Column,
		Value:          params.Value,
		FilterConfig:   filterConfig,
		Title:          params.Title,
		Description:    params.Description,
		Priority:       1,
		CreatedAt:      time.Now(),
		CreatedBy:      params.CreatedBy,
		UpdatedAt:      time.Now(),
		UpdatedBy:      params.CreatedBy,
	}

	return s.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&rule).
			Where("organization_id = ?", params.OrganizationId).
			Where("dataset_id = ?", params.DatasetId).
			Where("\"column\" = ?", params.Column).
			Where("deleted_at IS NULL").
			Updates(map[string]interface{}{
				"priority":   gorm.Expr("priority + 1"),
				"updated_at": time.Now(),
				"updated_by": params.CreatedBy,
			}).Error; err != nil {
			return err
		}

		return tx.Create(&rule).Error
	})

}

func (s *appStore) GetRules(ctx context.Context, params models.FilterRuleParams) (map[string]map[string][]models.Rule, error) {
	db := s.client.WithContext(ctx)

	var rules []models.Rule

	db = db.Where("organization_id = ?", params.OrganizationId)

	var datasetIds []uuid.UUID
	var columns []string
	for _, datasetColumn := range params.DatasetColumns {
		datasetIds = append(datasetIds, datasetColumn.DatasetId)
		columns = append(columns, datasetColumn.Columns...)
	}

	if len(datasetIds) > 0 {
		db = db.Where("dataset_id IN (?)", datasetIds)
	}

	if len(columns) > 0 {
		db = db.Where("\"column\" IN (?)", columns)
	}

	db = db.Where("deleted_at IS NULL").Order("priority asc")

	err := db.Find(&rules).Error
	if err != nil {
		return nil, err
	}

	rulesMap := make(map[string]map[string][]models.Rule)

	for _, rule := range rules {
		if _, ok := rulesMap[rule.DatasetId.String()]; !ok {
			rulesMap[rule.DatasetId.String()] = make(map[string][]models.Rule)
		}
		rulesMap[rule.DatasetId.String()][rule.Column] = append(rulesMap[rule.DatasetId.String()][rule.Column], rule)
	}

	return rulesMap, nil
}

func (s *appStore) GetRuleByIds(ctx context.Context, ruleIds []uuid.UUID) ([]models.Rule, error) {
	db := s.client.WithContext(ctx)

	var rules []models.Rule
	db = db.Where("rule_id IN (?)", ruleIds).Find(&rules)

	if db.Error != nil {
		return nil, db.Error
	}

	return rules, nil
}

func (s *appStore) GetRuleById(ctx context.Context, ruleId uuid.UUID) (models.Rule, error) {
	db := s.client.WithContext(ctx)

	rule := models.Rule{}
	db = db.Where("rule_id = ?", ruleId).First(&rule)

	if db.Error != nil {
		return models.Rule{}, db.Error
	}

	return rule, nil
}

func (s *appStore) UpdateRule(ctx context.Context, ruleId uuid.UUID, params models.UpdateRuleParams) error {
	db := s.client.WithContext(ctx)

	rule, err := s.GetRuleById(ctx, ruleId)
	if err != nil {
		return err
	}

	filterConfig, err := json.Marshal(params.FilterConfig)
	if err != nil {
		return err
	}

	db = db.Model(rule).Where("rule_id = ?", ruleId).Updates(map[string]interface{}{
		"title":         params.Title,
		"description":   params.Description,
		"value":         params.Value,
		"filter_config": filterConfig,
		"updated_by":    params.UpdatedBy,
		"updated_at":    time.Now(),
	})

	return db.Error
}

func (s *appStore) UpdateRulePriority(ctx context.Context, params models.UpdateRulePriorityParams) error {
	// Build CASE expression for priority update
	var cases string
	var ruleIds []uuid.UUID
	for i, rulePriority := range params.RulePriority {
		if i > 0 {
			cases += " "
		}
		cases += "WHEN rule_id = '" + rulePriority.RuleId.String() + "' THEN " + fmt.Sprint(rulePriority.Priority)
		ruleIds = append(ruleIds, rulePriority.RuleId)
	}

	ruleModel := &models.Rule{DatasetId: params.DatasetId}

	// Single UPDATE query with CASE statement
	return s.client.WithContext(ctx).
		Model(&ruleModel).
		Where("rule_id IN (?)", ruleIds).
		Updates(map[string]interface{}{
			"priority":   gorm.Expr("CASE " + cases + " END"),
			"updated_at": time.Now(),
			"updated_by": params.UpdatedBy,
		}).
		Error
}

func (s *appStore) DeleteRule(ctx context.Context, params models.DeleteRuleParams) error {
	return s.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rule models.Rule
		if err := tx.Model(&models.Rule{}).Where("rule_id = ?", params.RuleId).First(&rule).Error; err != nil {
			return err
		}

		if rule.DeletedAt != nil {
			return errors.New("rule already deleted")
		}

		if err := tx.Model(&rule).
			Where("organization_id = ?", rule.OrganizationId).
			Where("dataset_id = ?", rule.DatasetId).
			Where("\"column\" = ?", rule.Column).
			Where("priority > ?", rule.Priority).
			Where("deleted_at IS NULL").
			Updates(map[string]interface{}{
				"priority": gorm.Expr("priority - 1"),
			}).Error; err != nil {
			return err
		}

		tx.Model(&rule).Where("rule_id = ?", params.RuleId).Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"deleted_by": params.DeletedBy,
		})

		return tx.Error
	})
}
