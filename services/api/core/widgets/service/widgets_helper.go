package widgets

import (
	"fmt"
	"strings"

	dataplatformdataConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	datasetmodels "github.com/Zampfi/application-platform/services/api/core/datasets/models"
	"github.com/Zampfi/application-platform/services/api/core/widgets/constants"
	dataplatformmodels "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
)

func (s *widgetsService) flattenTags(data *datasetmodels.DatasetData, tagColumn string) datasetmodels.DatasetData {

	maxTagDepth := 0
	for _, row := range data.Rows {
		if tagValue, ok := row[tagColumn]; ok {
			if tagValue == nil {
				row[tagColumn] = "__UNTAGGED__"
			} else {
				tag := tagValue.(string)
				tagLevels := strings.Split(tag, ".")
				if len(tagLevels) > maxTagDepth {
					maxTagDepth = len(tagLevels)
				}
			}
		}
	}

	for i := 0; i < maxTagDepth; i++ {
		data.Columns = append(data.Columns, dataplatformmodels.ColumnMetadata{
			Name:         fmt.Sprintf("__%s_%s_%d", tagColumn, constants.HEIRARCHY_SUFFIX, i+1),
			DatabaseType: string(dataplatformdataConstants.StringDataType),
		})
	}

	for _, row := range data.Rows {
		if tagValue, ok := row[tagColumn]; ok {
			if tagValue == nil {
				row[tagColumn] = "__UNTAGGED__"
			} else {
				tag := tagValue.(string)
				tagLevels := strings.Split(tag, ".")
				for i, level := range tagLevels {
					row[fmt.Sprintf("__%s_%s_%d", tagColumn, constants.HEIRARCHY_SUFFIX, i+1)] = level
				}
			}

			// See if it's required
			delete(row, tagColumn)
		}
	}

	return *data
}

func (s *widgetsService) addRefToDataResults(data *datasetmodels.DatasetData, ref *string, populateEmptyRowsWithRef bool) datasetmodels.DatasetData {
	if ref == nil {
		return *data
	}

	data.Columns = append(data.Columns, dataplatformmodels.ColumnMetadata{
		Name:         constants.REF_PREFIX,
		DatabaseType: string(dataplatformdataConstants.StringDataType),
	})

	for _, row := range data.Rows {
		row[constants.REF_PREFIX] = *ref
	}

	if populateEmptyRowsWithRef {
		if len(data.Rows) == 0 {
			data.Rows = []map[string]interface{}{
				{
					constants.REF_PREFIX: *ref,
				},
			}
		}
	}

	return *data
}
