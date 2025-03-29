package service

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/core/rules/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RuleServiceStore interface {
	store.RuleStore
}

type RuleService interface {
	CreateRule(ctx context.Context, params dbmodels.CreateRuleParams) error
	GetRules(ctx context.Context, params dbmodels.FilterRuleParams) (map[string]map[string][]models.Rule, error)
	GetRuleByIds(ctx context.Context, ruleIds []uuid.UUID) ([]models.Rule, error)
	UpdateRule(ctx context.Context, ruleId uuid.UUID, params dbmodels.UpdateRuleParams) error
	UpdateRulePriority(ctx context.Context, params dbmodels.UpdateRulePriorityParams) error
	DeleteRule(ctx context.Context, params dbmodels.DeleteRuleParams) error
}

type ruleService struct {
	store RuleServiceStore
}

func NewRuleService(store RuleServiceStore) RuleService {
	return &ruleService{
		store: store,
	}
}

func (s *ruleService) CreateRule(ctx context.Context, params dbmodels.CreateRuleParams) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	err := s.store.CreateRule(ctx, params)
	if err != nil {
		logger.Error("Failed to create rule", zap.String("organizationId", params.OrganizationId.String()), zap.String("datasetId", params.DatasetId.String()), zap.Error(err))
		return err
	}

	return nil
}

func (s *ruleService) GetRules(ctx context.Context, params dbmodels.FilterRuleParams) (map[string]map[string][]models.Rule, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	rules, err := s.store.GetRules(ctx, params)
	if err != nil {
		logger.Error("Failed to get rules", zap.Any("params", params), zap.Error(err))
		return nil, err
	}

	result := make(map[string]map[string][]models.Rule)

	for datasetId, datasetRulesMap := range rules {
		datasetRules := make(map[string][]models.Rule)
		for column, rules := range datasetRulesMap {
			var rulesModels []models.Rule
			for _, rule := range rules {
				var ruleModel models.Rule
				err := ruleModel.FromSchema(&rule)
				if err != nil {
					logger.Error("Failed to convert rule schema to rule", zap.Any("rule", rule), zap.Error(err))
					return nil, err
				}
				rulesModels = append(rulesModels, ruleModel)
			}
			datasetRules[column] = rulesModels
		}

		result[datasetId] = datasetRules
	}

	return result, nil
}

func (s *ruleService) GetRuleByIds(ctx context.Context, ruleIds []uuid.UUID) ([]models.Rule, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	ruleSchema, err := s.store.GetRuleByIds(ctx, ruleIds)
	if err != nil {
		logger.Error("Failed to get rules by ids", zap.Any("ruleIds", ruleIds), zap.Error(err))
		return nil, err
	}

	var rules []models.Rule
	for _, rule := range ruleSchema {
		var ruleModel models.Rule
		err := ruleModel.FromSchema(&rule)
		if err != nil {
			logger.Error("Failed to convert rule schema to rule", zap.String("ruleId", rule.ID.String()), zap.Error(err))
			return nil, err
		}
		rules = append(rules, ruleModel)
	}

	return rules, nil
}

func (s *ruleService) UpdateRule(ctx context.Context, ruleId uuid.UUID, params dbmodels.UpdateRuleParams) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	err := s.store.UpdateRule(ctx, ruleId, params)
	if err != nil {
		logger.Error("Failed to update rule", zap.String("ruleId", ruleId.String()), zap.Any("params", params), zap.Error(err))
		return err
	}

	return nil
}

func (s *ruleService) UpdateRulePriority(ctx context.Context, params dbmodels.UpdateRulePriorityParams) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	err := s.store.UpdateRulePriority(ctx, params)
	if err != nil {
		logger.Error("Failed to update rule priority", zap.Any("params", params), zap.Error(err))
		return err
	}

	return nil
}

func (s *ruleService) DeleteRule(ctx context.Context, params dbmodels.DeleteRuleParams) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	err := s.store.DeleteRule(ctx, params)
	if err != nil {
		logger.Error("Failed to delete rule", zap.Any("params", params), zap.Error(err))
		return err
	}

	return nil
}
