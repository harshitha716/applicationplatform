import { PERIODICITY_TYPES } from 'constants/date.constants';
import { MapAny } from 'types/commonTypes';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export enum WIDGET_TYPES {
  BAR_CHART = 'bar_chart',
  LINE_CHART = 'line_chart',
  PIE_CHART = 'pie_chart',
  DONUT_CHART = 'donut_chart',
  KPI = 'kpi',
  TABLE = 'table',
  PIVOT_TABLE = 'pivot_table',
}

export type WidgetColumnType = {
  column_name: string;
  column_type: string;
};

export type WidgetDataRowType = {
  CurrencyCode: string;
  'SUM(IntegerAmount)': number;
};

export type WidgetDataType = {
  status: string;
  error: MapAny;
  rowcount: number;
  columns: WidgetColumnType[];
  data: MapAny[];
};

export type WidgetDataResponseType = {
  result: WidgetDataType[];
  periodicity?: PERIODICITY_TYPES;
  currency?: string;
};

export type WidgetDataRequestType = {
  widgetId: string;
  payload: {
    filters: string;
    time_columns?: string;
    periodicity?: PERIODICITY_TYPES;
    currency?: string;
  };
};

export enum AGGREGATION_TYPES {
  SUM = 'sum',
  AVG = 'avg',
  COUNT = 'count',
  MIN = 'min',
  MAX = 'max',
}

export enum FIELD_TYPES {
  DIMENSION = 'dimension',
  MEASURE = 'measure',
}

export type KPITagWidgetMapping = {
  dataset_id: string;
  fields: {
    primary_value: {
      type: string;
      column: string;
      field_type: FIELD_TYPES;
      aggregation?: AGGREGATION_TYPES;
    }[];
  };
  default_filters?: Partial<DefaultFilterType>;
};

export type PivotTableWidgetMapping = {
  dataset_id: string;
  ref: string;
  fields: {
    columns: {
      column: string;
      alias?: string;
      type: string;
      field_type: FIELD_TYPES;
      drilldown_filter_type?: FILTER_TYPES; // todo -- use discrimnated unions
      drilldown_filter_operator?: CONDITION_OPERATOR_TYPE;
    }[];
    rows?: {
      column: string;
      alias?: string;
      type: string;
      field_type: FIELD_TYPES;
      drilldown_filter_type?: FILTER_TYPES; // todo -- use discrimnated unions
      drilldown_filter_operator?: CONDITION_OPERATOR_TYPE;
    }[];
    values: {
      column: string;
      alias?: string;
      aggregation: AGGREGATION_TYPES;
      type: string;
      field_type: FIELD_TYPES;
      drilldown_filter_type?: FILTER_TYPES; // todo -- use discrimnated unions
      drilldown_filter_operator?: CONDITION_OPERATOR_TYPE;
    }[];
  };
  default_filters?: Partial<DefaultFilterType>;
};

export type DefaultFilterType = {
  logical_operator: string;
  conditions: {
    logical_operator: string;
    conditions: string;
    column: string;
    type: string;
    operator: string;
    value: string[];
  }[];
};

export type FieldsMappingType = {
  x_axis: AxisMappingType[];
  y_axis: AxisMappingType[];
  group_by?: AxisMappingType[];
};

export type BarLineChartWidgetMapping = {
  dataset_id: string;
  fields: FieldsMappingType;
  default_filters?: Partial<DefaultFilterType>;
};

export type PieDonutChartFieldsMappingType = {
  slices?: AxisMappingType[];
  values?: AxisMappingType[];
};

export interface PieDonutChartWidgetMapping {
  dataset_id: string;
  fields: PieDonutChartFieldsMappingType;
  default_filters?: Partial<DefaultFilterType>;
}

export interface AxisMappingType {
  type: string;
  column: string;
  field_type: FIELD_TYPES;
  alias?: string;
  aggregation?: AGGREGATION_TYPES;
  drilldown_filter_type?: string;
  drilldown_filter_operator?: CONDITION_OPERATOR_TYPE;
}

export type WidgetInstanceResponseType = WidgetInstanceType;

export interface WidgetInstanceBaseType {
  widget_instance_id: string;
  widget_id: string;
  sheet_id: string;
  title: string;
  dataset_id: string;
  created_at: string;
  updated_at: string;
  display_config?: MapAny;
}

export interface LineBarChartWidgetInstanceType extends WidgetInstanceBaseType {
  widget_type: WIDGET_TYPES.BAR_CHART | WIDGET_TYPES.LINE_CHART;
  data_mappings: {
    version: string;
    datasets: { id: string }[];
    mappings: BarLineChartWidgetMapping[];
  };
}

export interface PieDonutChartWidgetInstanceType extends WidgetInstanceBaseType {
  widget_type: WIDGET_TYPES.PIE_CHART | WIDGET_TYPES.DONUT_CHART;
  data_mappings: {
    version: string;
    datasets: { id: string }[];
    mappings: PieDonutChartWidgetMapping[];
  };
}

export interface PivotTableWidgetInstanceType extends WidgetInstanceBaseType {
  widget_type: WIDGET_TYPES.PIVOT_TABLE;
  data_mappings: {
    version: string;
    datasets: { id: string }[];
    mappings: PivotTableWidgetMapping[];
    display_config?: MapAny;
  };
}

export interface KPITagWidgetInstanceType extends WidgetInstanceBaseType {
  widget_type: WIDGET_TYPES.KPI;
  data_mappings: {
    version: string;
    mappings: KPITagWidgetMapping[];
  };
}

export type WidgetInstanceType =
  | LineBarChartWidgetInstanceType
  | PieDonutChartWidgetInstanceType
  | PivotTableWidgetInstanceType
  | KPITagWidgetInstanceType;
