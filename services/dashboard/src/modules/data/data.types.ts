import { RowClickedEvent } from 'ag-grid-community';

export enum DATASET_ACCESS_PRIVILEGES {
  ADMIN = 'admin',
  VIEWER = 'viewer',
}

export type UserAccessToDataSetType = {
  name: string;
  privilege: string;
  resource_type: string;
}[];

export type ShareDatasetPopupPropsType = {
  datasetId: string;
};

export type DatasetAccessPrivilegesType = {
  label: string;
  value: DATASET_ACCESS_PRIVILEGES;
};

export enum DATASET_ACTION_STATUS {
  INITIATED = 'INITIATED',
  SUCCESSFUL = 'SUCCESSFUL',
  FAILED = 'FAILED',
}

export type DatasetAccessToAudiencesPropsType = {
  name?: string;
  resource_type: string;
  privilege?: string;
  datasetId: string;
  resource_audience_id: string;
  resource_audience_type: string;
  user?: {
    email: string;
    name?: string;
  };
  userPrivilege: string;
  orgName?: string;
  customerName?: string;
  teamInfo?: {
    name?: string;
    color?: string;
  };
};

export type DatasetColumnRequest = {
  dataset_id: string;
  columns: string[];
};

export enum LOADER_STATUS {
  ALIGNMENT_PENDING = 'allignment_pending',
  ALIGNMENT_COMPLETED = 'allignment_completed',
  INITIATED = 'initiated',
  LOADING = 'loading',
  SUCCESS = 'success',
  ERROR = 'error',
}

export type ListingPropsType = {
  onRowClicked: (event: RowClickedEvent) => void;
};

export type ColumnOrderingVisibilityType = {
  colId: string;
  isVisible: boolean;
  width: number;
};
