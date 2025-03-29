import { WidgetDataValueType } from 'modules/widgets/widgets.constant';
import { WidgetInstanceType } from 'types/api/widgets.types';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export type PageResponseType = {
  page_id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
  fractional_index: number;
  organization_id: string;
};

export type SheetResponseType = {
  page_id: string;
  name: string;
  description: string;
  sheets: {
    sheet_id: string;
    name: string;
    description: string;
    created_at: string;
    updated_at: string;
    fractional_index: number;
    page_id: string;
  }[];
  created_at: string;
  updated_at: string;
  fractional_index: number;
  organization_id: string;
};
export interface SheetDetailsResponseType {
  sheet_id: string;
  name: string;
  description: string;
  widget_instances: WidgetInstanceType[];
  sheet_config: SheetConfigType;
  created_at: string;
  updated_at: string;
  fractional_index: number;
  page_id: string;
}

export interface LayoutType {
  x: number;
  y: number;
  w: number;
  h: number;
}

export interface WidgetType {
  name: string;
  layout: LayoutType;
  default_widget: string;
  widget_group: string[];
}

export interface SheetConfigType {
  version: string;
  sheet_layout: WidgetType[];
  currency: {
    hide_currency_filter: true;
    default_currency: string;
  };
}

export type SheetDetailsRequestType = {
  sheetId: string;
  pageId?: string;
};

export type SheetFilterConfigResponseType = {
  native_filter_config: SheetFilterType[];
};

export type SheetFilterType = {
  name: string; // Name of the filter
  filter_type: FILTER_TYPES; // Type of filter
  data_type?: WidgetDataValueType; // Data type of the filter
  widgets_in_scope?: string[]; // List of widgets affected by this filter
  targets: TargetType[]; // List of dataset targets for this filter
  default_value?: FilterDefaultValueType; // Default value for the filter
  options?: string[]; // Options for multi-select filters (if applicable)
};

export interface TargetType {
  dataset_id: string; // Unique identifier for the dataset
  column: string; // Column name in the dataset
}

export interface FilterDefaultValueType {
  operator: CONDITION_OPERATOR_TYPE; // Operator for range-based filters
  value: string[];
  from?: number | string; // Start value or single value
  to?: number | string; // End value (if applicable)
}

export type AudiencesByPageIdRequest = {
  pageId: string;
};

export type AudiencesByPageIdResponse = {
  resource_audience_type: string;
  resource_audience_id: string;
  privilege: string;
  resource_type: string;
  resource_id: string;
  user?: {
    role?: string;
    email?: string;
  };
};

export type AudiencesPageShareData = {
  audiences: {
    audience_type: string;
    audience_id: string;
    role: string;
  }[];
};

export type PostPagesToAudiencesByPageIdType = { pageId: string; body: AudiencesPageShareData };
export type PatchChangeAudienceRoleInPageType = { pageId: string; body: { audience_id: string; role: string } };
export type DeleteAudienceFromPageAccessType = { pageId: string; body: { audience_id: string } };
