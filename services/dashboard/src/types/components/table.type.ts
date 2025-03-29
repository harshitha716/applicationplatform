import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export enum OrderType {
  ASC = 'asc',
  DESC = 'desc',
}

export enum LogicalOperatorType {
  OperatorLogicalAnd = 'AND',
  OperatorLogicalOr = 'OR',
}

export enum AggregationFunctionType {
  AggregationFunctionSum = 'SUM',
  AggregationFunctionAvg = 'AVG',
  AggregationFunctionMin = 'MIN',
  AggregationFunctionMax = 'MAX',
  AggregationFunctionCount = 'COUNT',
}

export type FilterType = {
  logical_operator?: LogicalOperatorType;
  column?: string;
  operator?: CONDITION_OPERATOR_TYPE;
  value?: any;
  conditions?: FilterType[];
};

export type AggregationType = {
  column: string;
  alias: string;
  function: AggregationFunctionType;
};

export type GroupByType = {
  column: string;
  alias: string;
};

export type OrderByType = {
  column: string;
  order: OrderType;
};

export type PaginationType = {
  page: number;
  page_size: number;
};

export type FilterModelType = {
  logical_operator?: LogicalOperatorType;
  conditions?: FilterType[];
};

export type RequestType = {
  filters: FilterModelType | null;
  aggregations: AggregationType[];
  group_by: GroupByType[];
  order_by: OrderByType[];
  pagination: PaginationType;
  get_total_records: boolean;
  fx_currency?: string | undefined;
};
