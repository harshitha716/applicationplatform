package rosetta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"
	rosettaConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/rosetta/constants"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/rosetta/models"
	"github.com/Zampfi/application-platform/services/api/helper"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	"go.uber.org/zap"
)

type RosettaService interface {
	TranslateQuery(ctx context.Context, query string, providerType constants.ProviderType) (string, error)
}

type rosettaService struct {
	baseUrl    string
	httpClient helper.HTTPClient
}

func InitRosettaService(baseUrl string) RosettaService {
	httpClient := http.Client{}
	return &rosettaService{
		baseUrl:    baseUrl,
		httpClient: &httpClient,
	}
}

// TranslateQuery Translates the query to the provider's SQL dialect
// Falls back to the original query if translation fails
func (s *rosettaService) TranslateQuery(ctx context.Context, query string, providerType constants.ProviderType) (string, error) {
	startTime := time.Now()
	logger := apicontext.GetLoggerFromCtx(ctx)
	logger.Info("ROSETTA_INPUT_QUERY", zap.String("QUERY", query))

	requestBody := models.TranslateQueryRequest{
		Query:        query,
		OutputFormat: string(providerType),
	}

	url := s.baseUrl + string(rosettaConstants.RosettaTranslateSqlApiPath)
	body, resp, err := helper.HttpPost(s.httpClient, url, requestBody, map[string]string{})
	if err != nil {
		logger.Error(errors.RosettaHttpCallFailedErrorErrMessage, zap.Error(err))
		return query, errors.ErrRosettaHttpCallFailedError
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("%s: %d, message: %s", errors.RosettaHttpErrorErrMessage, resp.StatusCode, body))
		return query, errors.ErrRosettaHttpError
	}

	var queryResponse models.TranslateQueryResponse
	err = json.Unmarshal(body, &queryResponse)
	if err != nil {
		logger.Error(errors.JSONUnmarshallingFailedErrMessage, zap.Error(err))
		return query, errors.ErrJSONUnmarshallingFailed
	}

	logger.Info("ROSETTA_OUTPUT_QUERY", zap.String("QUERY", queryResponse.Query), zap.Any("ROSETTA_TRANSLATE_SQL_TIME_MS", time.Since(startTime).Milliseconds()))
	return queryResponse.Query, nil
}
