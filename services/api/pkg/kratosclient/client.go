package kratosclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	kratos "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int
	Message string
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	publicUrl url.URL
	kratos    *kratos.APIClient
}

type KratosClient interface {
	GetSessionInfo(ctx context.Context, logger *zap.Logger, cookie string) (*kratos.Session, *http.Response, *Error)
	GetAuthProxy() (*httputil.ReverseProxy, url.URL)
	CreateLoginFlow(ctx context.Context, logger *zap.Logger, email string) (*kratos.LoginFlow, *http.Response, *Error)
	SignupUserEmailPassword(ctx context.Context, logger *zap.Logger, email string, password string) (*kratos.Identity, *http.Response, *Error)
}

func NewClient(publicUrl string) (*Client, error) {

	publicUrlParsed, err := url.ParseRequestURI(publicUrl)
	if err != nil || publicUrlParsed == nil {
		return nil, fmt.Errorf("invalid public url: %w", err)
	}

	kratosCfg := kratos.NewConfiguration()
	kratosCfg.Servers = kratos.ServerConfigurations{
		{
			URL:         publicUrl,
			Description: "public",
		},
	}

	kratosClient := kratos.NewAPIClient(kratosCfg)

	return &Client{
		publicUrl: *publicUrlParsed,
		kratos:    kratosClient,
	}, nil
}

func (c *Client) GetSessionInfo(ctx context.Context, logger *zap.Logger, cookie string) (*kratos.Session, *http.Response, *Error) {

	session, httpResp, err := c.kratos.FrontendAPI.ToSession(ctx).Cookie(cookie).Execute()
	if err != nil {
		logger.Error("failed to get session info", zap.Error(err))
		if httpResp == nil {
			return nil, nil, &Error{
				Code:    http.StatusInternalServerError,
				Message: "internal error",
			}
		}
		return nil, httpResp, &Error{
			Code:    httpResp.StatusCode,
			Message: err.Error(),
		}
	}

	return session, httpResp, nil
}

func (c *Client) GetAuthProxy() (*httputil.ReverseProxy, url.URL) {
	return httputil.NewSingleHostReverseProxy(&c.publicUrl), c.publicUrl
}

func (c *Client) CreateLoginFlow(ctx context.Context, logger *zap.Logger, email string) (*kratos.LoginFlow, *http.Response, *Error) {
	loginFlow, httpResp, err := c.kratos.FrontendAPI.CreateBrowserLoginFlow(ctx).Execute()

	if err != nil {
		logger.Error("failed to create login flow", zap.Error(err))
		statusCode := http.StatusInternalServerError
		if httpResp != nil {
			statusCode = httpResp.StatusCode
		}
		return nil, httpResp, &Error{
			Code:    statusCode,
			Message: err.Error(),
		}
	}
	return loginFlow, httpResp, nil
}

func (c *Client) SignupUserEmailPassword(ctx context.Context, logger *zap.Logger, email string, password string) (*kratos.Identity, *http.Response, *Error) {
	flow, _, err := c.kratos.FrontendAPI.CreateNativeRegistrationFlow(ctx).Execute()
	if err != nil {
		logger.Error("failed to create registration flow", zap.Error(err))
		return nil, nil, &Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	result, httpResp, err := c.kratos.FrontendAPI.UpdateRegistrationFlow(ctx).Flow(flow.Id).UpdateRegistrationFlowBody(
		kratos.UpdateRegistrationFlowWithPasswordMethodAsUpdateRegistrationFlowBody(&kratos.UpdateRegistrationFlowWithPasswordMethod{
			Method:   "password",
			Password: password,
			Traits:   map[string]interface{}{"email": email},
		}),
	).Execute()

	if err != nil {
		logger.Error("failed to update registration flow", zap.Error(err))
		return nil, nil, &Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return &result.Identity, httpResp, nil
}
