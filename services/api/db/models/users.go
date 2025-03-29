package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID    uuid.UUID `json:"user_id" gorm:"column:user_id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

func (u User) TableName() string {
	return "users_with_traits"
}

func (u *User) GetQueryFilters(db *gorm.DB, currentUserId uuid.UUID, organizationIds []uuid.UUID) *gorm.DB {

	// TODO: Terrible query, need to optimize
	return db.Where(
		`EXISTS (
			SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
			WHERE frap.resource_type = 'organization'
			AND frap.resource_id IN (?)
			AND (
				frap.user_id = users_with_traits.user_id
				OR (
					EXISTS (
						SELECT 1 FROM "app"."organization_membership_requests" omr
						WHERE omr.organization_id = frap.resource_id
						AND omr.user_id = users_with_traits.user_id
					)
				)
			)
			AND frap.deleted_at IS NULL
		) OR users_with_traits.user_id = ?
		`,
		organizationIds,
		currentUserId,
	)
}

func (u *User) BeforeCreate(db *gorm.DB) error {

	return fmt.Errorf("insert forbidden")

}
