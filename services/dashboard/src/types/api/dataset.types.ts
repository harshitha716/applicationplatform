import { DATASET_ACTION_STATUS } from 'modules/data/data.types';
import { ResourceAudienceType } from 'types/api/auth.types';
import { MapAny } from 'types/commonTypes';
import { FilterModelType } from 'types/components/table.type';
import { CUSTOM_COLUMNS_TYPE, VALUE_FORMAT_TYPE } from 'components/common/table/table.types';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export type ValueFormatType = {
  type: VALUE_FORMAT_TYPE;
  value: number | string;
};

export type DatasetFilterConfigMetadataType = {
  is_hidden?: boolean;
  custom_type?: CUSTOM_COLUMNS_TYPE;
  config?: {
    currency_column?: string;
    amount_column?: string;
    currency_value?: string;
    format?: string;
    value_format?: ValueFormatType[];
  };
  is_editable?: boolean;
};

export type DatasetFilterConfigResponseType = {
  column: string;
  type: FILTER_TYPES;
  options: string[];
  datatype: string;
  alias: string;
  metadata?: DatasetFilterConfigMetadataType;
};

export type DatasetDataResponseType = {
  data: {
    rows: MapAny[];
    columns: MapAny[];
    config: {
      is_drilldown_enabled: boolean;
    };
    total_count: number;
  };
  description: string;
  title: string;
};

export type DatasetExportResponseType = {
  workflow_id: string;
};

export type DatasetDataRequestType = {
  datasetId: string;
  query_config?: string;
};

export type DatasetExportsSignedUrlRequestType = {
  datasetId: string;
  workflowId: string;
};

export type DatasetExportsSignedUrlResponseType = {
  signed_url: string;
};

export type DatasetDrilldownRequestType = {
  datasetId: string;
  rowId: string;
};

export type DatasetDrilldownResponseType = {
  tabs: {
    dataset_id: string;
    dataset_title: string;
    filters: FilterModelType;
  }[];
};

export type DatasetType = {
  id: string;
  title: string;
  description: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  organization_id: string;
  metadata: MapAny;
};

export type DatasetListingResponseType = {
  datasets: DatasetType[];
  total_count: number;
};

export type DatasetListingRequestType = {
  page: number;
  pageSize: number;
  sort?: string;
};

export type DatasetUpdateRequestType = {
  datasetId: string;
  data: {
    filters: FilterModelType | null;
    update: {
      column: string;
      value: string;
    };
    save_as_rule?: boolean;
    rule_title?: string;
    rule_description?: string;
  };
};

export type DatasetUpdateResponseType = {
  action_id: string;
  action_type: string;
  dataset_id: string;
  status: string;
  config: MapAny;
  action_by: string;
  is_completed: boolean;
};

export type DatasetActionStatusRequestType = {
  datasetId: string;
  params: {
    action_ids?: string[];
    status?: DATASET_ACTION_STATUS;
  };
};

export type DatasetActionStatusResponseType = {
  action_id: string;
  action_type: string;
  dataset_id: string;
  status: DATASET_ACTION_STATUS;
  is_completed: boolean;
  config: MapAny;
  action_by: string;
};

export type AudiencesByDatasetIdRequestType = {
  datasetId: string;
};

export type AudiencesByDatasetIdResponseType = {
  user: {
    email: string;
    name?: string;
  };
  privilege: string;
  resource_audience_type: ResourceAudienceType;
  resource_audience_id: string;
  resource_type: string;
};

export type AudiencesDatasetShareData = {
  audiences: {
    audience_type: string;
    audience_id: string;
    role: string;
  }[];
};

export type PostShareDatasetToAudiencesByDatasetIdType = { datasetId: string; body: AudiencesDatasetShareData };

export type PatchChangeAudienceRoleInDatasetType = { datasetId: string; body: { audience_id: string; role: string } };

export type DeleteAudienceFromDatasetAccessType = { datasetId: string; body: { audience_id: string } };

export type GetRulesByDatasetColumnsRequestType = {
  dataset_columns: string;
};

export type ConditionType = {
  logical_operator: string;
  column: {
    column: string;
    datatype: string;
    custom_data_config?: MapAny;
    alias?: string;
  };
  operator: CONDITION_OPERATOR_TYPE;
  value: string[] | string;
  conditions?: ConditionType[];
};

export type RuleFilters = {
  logical_operator: string;
  conditions: ConditionType[];
};

export type RuleType = {
  rule_id: string;
  organization_id: string;
  dataset_id: string;
  column: string;
  value: string;
  filter_config: {
    query_config: {
      table_config: {
        dataset_id: string;
        columns: string[];
        update: {
          column: {
            column: string;
            datatype: string;
            custom_data_config: null;
            alias: null;
          };
          value: string;
        }[];
      };
      filters: RuleFilters;
    };
  };
  title: string;
  description: string;
  priority: number;
  created_at: string;
  created_by: string;
  updated_at: string;
  updated_by: string;
  deleted_at: string;
  deleted_by: string;
};

export type RulesByDatasetColumnType = {
  [column: string]: RuleType[];
};

export type GetRulesByDatasetColumnsResponseType = {
  [datasetId: string]: RulesByDatasetColumnType;
};

export type UploadFileResponseType = {
  identifier: string;
  url: string;
  fileName: string;
  downloadableUrl: string;
  rawFile: File | null;
};

export type TableData = {
  columns: string[];
  rows: { [key: string]: string | number | boolean | null }[];
};

export type TransformationPreviewMetadata = {
  data_preview: TableData;
};

export type RawMetadata = {
  columns: string[];
  rows: { [key: string]: string | number | boolean | null }[];
};

export type GetRulesByRuleIdsRequestType = {
  rule_ids: string[];
};

export type RulePriorityType = {
  rule_id: string;
  priority: number;
};

export type RulePrioritiesType = {
  updated_by: string;
  rule_priority: RulePriorityType[];
};

export type UpdateRulePriorityRequestType = {
  dataset_id: string;
  column: string;
  rule_priorities: RulePrioritiesType;
};

export type SignedUrlBodyType = {
  file_name: string;
  file_type: string;
};
export type PreviewTransformationRequest = {
  file_upload_id: string;
  dataset_id?: string;
};

export type PreviewTransformationResponse = {
  dataset_action_id: string;
};

export type PostAiTransformationConfirmResponseType = {
  message: string;
};

export type GetAiTransformationResponseType = {
  data_preview: TableData;
};

export type GetAiTransformationRequestType = {
  file_upload_id: string;
};

export type PostAiTransformationConfirmRequestType = {
  file_upload_id: string;
  dataset_id?: string;
};

export type GetFileImportHistoryRequestType = {
  datasetId: string;
};

export type userDetailsType = {
  email: string;
  name: string;
};

export type GetFileImportHistoryResponseType = {
  file_uploads: {
    id: string;
    dataset_id: string;
    file_id: string;
    file_name: string;
    file_upload_status: string;
    file_upload_created_at: string;
    uploaded_by_user: userDetailsType;
  }[];
};
