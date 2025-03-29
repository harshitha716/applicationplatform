package auth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	kratos "github.com/ory/kratos-client-go"
)

func RunAuthSeeds(authBase string, apiBase string, adminSecret string, pgClient *sql.DB) {
	kratosConfig := kratos.NewConfiguration()
	kratosConfig.Servers = kratos.ServerConfigurations{
		{
			URL: authBase,
		},
	}

	client := kratos.NewAPIClient(kratosConfig)

	// create an admin user and their organization

	/*
	   {"OrganizationName":"Test org 2","OrganizationDescription":"Test org 2","AdminEmail":"game@zamp.ai","AdminPassword":"Zamp@123","SSODomain":"zamp.ai","SSOProviderName":"google-id","SSOProviderID":"google-id"}
	*/

	var createOrgPayload struct {
		OrganizationName        string `json:"OrganizationName"`
		OrganizationDescription string `json:"OrganizationDescription"`
		AdminEmail              string `json:"AdminEmail"`
		AdminPassword           string `json:"AdminPassword"`
		SSODomain               string `json:"SSODomain"`
		SSOProviderName         string `json:"SSOProviderName"`
		SSOProviderID           string `json:"SSOProviderID"`
	}

	createOrgPayload.OrganizationName = "Test organization"
	createOrgPayload.OrganizationDescription = "Test organization"
	createOrgPayload.AdminEmail = "admin@zamp.ai"
	createOrgPayload.AdminPassword = "MeZamp@123"
	createOrgPayload.SSODomain = "zamp.ai"
	createOrgPayload.SSOProviderName = "google-id"
	createOrgPayload.SSOProviderID = "google-id"

	createOrgPayloadJson, err := json.Marshal(createOrgPayload)
	if err != nil {
		fmt.Println("Error in creating seed organization payload", err)
		return
	}

	httpReq, err := http.NewRequest("POST", apiBase+"/admin/create-organization", bytes.NewReader(createOrgPayloadJson))
	if err != nil {
		fmt.Println("Error in creating seed organization request", err.Error())
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Zamp-Admin-Secret", adminSecret)
	httpReq.Header.Set("X-Zamp-User-Id", uuid.Nil.String())
	httpReq.Header.Set("X-Zamp-Organization-Ids", uuid.Nil.String())

	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Println("Error in creating seed organization", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error in creating seed organization", resp.StatusCode)
		return
	}

	createOrgResp := struct {
		Organization struct {
			OrganizationId string `json:"organization_id"`
			OwnerId        string `json:"owner_id"`
		} `json:"organization"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&createOrgResp)
	if err != nil {
		fmt.Println("Error in creating seed organization", err)
		return
	}

	fmt.Println("Created organization", createOrgResp)

	// invite other users to the organization
	inviteUserPayload := struct {
		Invitations []struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}
	}{
		Invitations: []struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}{
			{
				Email: "manager@zamp.ai",
				Role:  "system_admin",
			},
			{
				Email: "member@zamp.ai",
				Role:  "member",
			},
		},
	}

	inviteUserPayloadJson, err := json.Marshal(inviteUserPayload)
	if err != nil {
		fmt.Println("Error in creating seed organization payload", err)
		return
	}

	httpReq, err = http.NewRequest("POST", apiBase+"/organizations/"+createOrgResp.Organization.OrganizationId+"/audiences/invitations", bytes.NewReader(inviteUserPayloadJson))
	if err != nil {
		fmt.Println("Error in creating seed organization payload", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Zamp-Admin-Secret", adminSecret)
	httpReq.Header.Set("X-Zamp-User-Id", createOrgResp.Organization.OwnerId)
	httpReq.Header.Set("X-Zamp-Organization-Ids", createOrgResp.Organization.OrganizationId)

	resp, err = httpClient.Do(httpReq)
	if err != nil {
		fmt.Println("Error in inviting users to organization", err)
		return
	}

	respJson := map[string]interface{}{}

	respString, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in inviting users to organization. Could not read response", err)
		return
	}

	fmt.Println("Organization invite response", string(respString))

	err = json.Unmarshal(respString, &respJson)
	if err != nil {
		fmt.Println("Error in inviting users to organization. Could not decode response", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error in inviting users to organization. Status code", resp.StatusCode)
	}

	fmt.Println("Invited users to organization", respJson)

	err = CreateIdentityWithPassword(client, "manager@zamp.ai", "MeZamp@123")
	if err != nil {
		fmt.Println("Error in creating seed user 2", err)
	}

	err = CreateIdentityWithPassword(client, "member@zamp.ai", "MeZamp@123")
	if err != nil {
		fmt.Println("Error in creating seed user 3", err)
	}

	// create a team in the organization
	createTeamPayload := struct {
		TeamName     string `json:"name"`
		ColorHexCode string `json:"color_hex_code"`
		Description  string `json:"description"`
	}{
		TeamName:     "Test team",
		ColorHexCode: "#000000",
		Description:  "Test team",
	}

	createTeamPayloadJson, err := json.Marshal(createTeamPayload)
	if err != nil {
		fmt.Println("Error in creating seed team payload", err)
		return
	}

	httpReq, err = http.NewRequest("POST", apiBase+"/organizations/"+createOrgResp.Organization.OrganizationId+"/teams", bytes.NewReader(createTeamPayloadJson))
	if err != nil {
		fmt.Println("Error in creating seed team request", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Zamp-Admin-Secret", adminSecret)
	httpReq.Header.Set("X-Zamp-User-Id", createOrgResp.Organization.OwnerId)
	httpReq.Header.Set("X-Zamp-Organization-Ids", createOrgResp.Organization.OrganizationId)

	resp, err = httpClient.Do(httpReq)
	if err != nil {
		fmt.Println("Error in creating seed team", err)
		return
	}

	respJson = map[string]interface{}{}

	respString, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in creating seed team. Could not read response", err)
		return
	}

	fmt.Println("Created team", string(respString))
}

func CreateIdentityWithPassword(c *kratos.APIClient, email, password string) error {

	ctx := context.Background()

	// Initialize a registration flow

	flow, _, err := c.FrontendAPI.CreateNativeRegistrationFlow(ctx).Execute()
	if err != nil {
		fmt.Println("Failed to create a registration flow")
		return err
	}

	// Submit the registration flow

	result, _, err := c.FrontendAPI.UpdateRegistrationFlow(ctx).Flow(flow.Id).UpdateRegistrationFlowBody(
		kratos.UpdateRegistrationFlowWithPasswordMethodAsUpdateRegistrationFlowBody(&kratos.UpdateRegistrationFlowWithPasswordMethod{
			Method:   "password",
			Password: password,
			Traits:   map[string]interface{}{"email": email},
		}),
	).Execute()

	if err != nil {
		fmt.Println("Error in creating user", err.Error(), err)
		return err
	}

	fmt.Println("Created identity: ", result.Identity.Traits)
	fmt.Println("Password: ", password)

	return nil
}
