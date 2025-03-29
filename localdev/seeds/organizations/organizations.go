package organizations

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func RunOrganizationSeeds(pgClient *sql.DB) {
	setMemberOrgIdAsAdminOrgId(pgClient)
}

func setMemberOrgIdAsAdminOrgId(pgClient *sql.DB) {
	// get user ids from users_with_traits where email = 'member@zamp.ai' and 'admin@zamp.ai'

	type user struct {
		UserId uuid.UUID `json:"user_id"`
		Email  string    `json:"email"`
	}

	query := `SELECT user_id, email FROM app.users_with_traits`
	rows, err := pgClient.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user
		err := rows.Scan(&u.UserId, &u.Email)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}

	var adminUserId uuid.UUID
	var memberUserId uuid.UUID
	for _, u := range users {
		if u.Email == "admin@zamp.ai" {
			adminUserId = u.UserId
		}
		if u.Email == "member@zamp.ai" {
			memberUserId = u.UserId
		}
	}

	// get the organizations that they're a part of
	raQuery := `select resource_audience_policy_id, resource_audience_type, resource_audience_id, privilege, resource_type, resource_id from app.resource_audience_policies where resource_type = 'organization' AND resource_audience_id in ($1, $2) and resource_audience_type = 'user'`

	rows, err = pgClient.Query(raQuery, adminUserId, memberUserId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	type resourceAudiencePolicy struct {
		ID                   uuid.UUID `json:"resource_audience_policy_id"`
		ResourceAudienceType string    `json:"resource_audience_type"`
		ResourceAudienceID   uuid.UUID `json:"resource_audience_id"`
		Privilege            string    `json:"privilege"`
		ResourceType         string    `json:"resource_type"`
		ResourceID           uuid.UUID `json:"resource_id"`
	}

	var raPolicies []resourceAudiencePolicy
	for rows.Next() {
		var ra resourceAudiencePolicy
		err := rows.Scan(&ra.ID, &ra.ResourceAudienceType, &ra.ResourceAudienceID, &ra.Privilege, &ra.ResourceType, &ra.ResourceID)
		if err != nil {
			panic(err)
		}
		raPolicies = append(raPolicies, ra)
	}

	// get the organization ID for which admin is a part of
	var orgId uuid.UUID
	for _, ra := range raPolicies {
		if ra.ResourceAudienceID == adminUserId {
			orgId = ra.ResourceID
		}
	}

	// update the organization ID for member to be the same as admin, but privilege to be "member"
	updateQuery := `update app.resource_audience_policies set resource_id = $1, privilege = 'member' where resource_audience_id = $2 and resource_type = 'organization'`

	_, err = pgClient.Exec(updateQuery, orgId, memberUserId)

	if err != nil {
		fmt.Println("Failed updating member organization id", err)
	} else {
		fmt.Println("Successfully updated member organization id")
	}
}
