package main

import (
	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	dataplatform "github.com/Zampfi/application-platform/services/api/core/dataplatform"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	fileimportsservice "github.com/Zampfi/application-platform/services/api/core/fileimports"
	ruleservice "github.com/Zampfi/application-platform/services/api/core/rules/service"
	cloudservice "github.com/Zampfi/application-platform/services/api/pkg/cloudservices/service"
	"github.com/Zampfi/application-platform/services/api/pkg/logging"
	querybuilderservice "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/service"
	"github.com/Zampfi/application-platform/services/api/server"
)

func main() {

	// initialize logger
	logger, err := logging.GetLogger()
	if err != nil {
		panic(err)
	}
	if logger == nil {
		panic("logger is nil")
	}

	// initialize server config
	serverConfig, cleanup, err := serverconfig.Createserverconfig(logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if serverConfig == nil {
		panic("serverConfig is nil")
	}

	queryBuilderService := querybuilderservice.NewQueryBuilder()

	dataPlatformService, err := dataplatform.InitDataPlatformService(serverConfig.DataPlatformConfig)
	if err != nil {
		panic(err)
	}
	if dataPlatformService == nil {
		panic("DataPlatformService is nil")
	}

	ruleService := ruleservice.NewRuleService(serverConfig.Store)
	if ruleService == nil {
		panic("RuleService is nil")
	}

	fileImportService := fileimportsservice.NewFileImportService(serverConfig.DefaultS3Client, serverConfig.Store, serverConfig.Env.AWSDefaultBucketName)
	if fileImportService == nil {
		panic("FileImportService is nil")
	}

	cloudService, err := cloudservice.NewCloudService("GCP", *serverConfig.Env)
	if err != nil {
		panic(err)
	}
	if cloudService == nil {
		panic("CloudService is nil")
	}

	datasetService := datasetservice.NewDatasetService(serverConfig.Store, queryBuilderService, dataPlatformService, ruleService, fileImportService, serverConfig.TemporalSdk, cloudService, serverConfig.DefaultS3Client, *serverConfig.DatasetConfig, serverConfig.CacheClient)
	if datasetService == nil {
		panic("DatasetService is nil")
	}
	// TODO: Graceful shutdown

	// run application server
	server.RunServer(serverConfig, logger, queryBuilderService, dataPlatformService, ruleService, datasetService)
}
