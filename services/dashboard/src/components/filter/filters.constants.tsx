import { FILTER_TYPES } from 'components/filter/filter.types';
import AmountRangeFilterMenuItem from 'components/filter/filterMenu/AmountRangeFilterMenuItem';
import DateRangeFilterMenuItem from 'components/filter/filterMenu/DateRangeFilterMenuItem';
import MultiSearchFilterMenuItem from 'components/filter/filterMenu/MultiSearchFilterMenuItem';
import MultiSelectFilterMenuItem from 'components/filter/filterMenu/MultiSelectFilterMenuItem';
import SearchFilterMenuItem from 'components/filter/filterMenu/SearchFilterMenuItem';
import SingleSelectFilterMenuItem from 'components/filter/filterMenu/SingleSelectFilterMenuItem';
import TagsSelectFilterMenuItem from 'components/filter/filterMenu/TagsSelectFilterMenuItem';

export enum CONDITION_OPERATOR_TYPE {
  IN = 'in',
  NOT_IN = 'nin',
  CONTAINS = 'contains',
  ARRAY_CONTAINS = 'array_contains',
  ARRAY_IN = 'array_in',
  IS_NULL = 'is_null',
  NOT_CONTAINS = 'ncontains',
  EQUAL = 'eq',
  NOT_EQUAL = 'neq',
  GREATER_THAN = 'gt',
  GREATER_THAN_EQUAL = 'gte',
  LESS_THAN = 'lt',
  LESS_THAN_EQUAL = 'lte',
  ONE_OF = 'iof',
  DEBIT = 'debit',
  CREDIT = 'credit',
  ANY = 'any',
  STARTS_WITH = 'startswith',
  ENDS_WITH = 'endswith',
  IN_BETWEEN = 'inbetween',
}

export enum FILTER_PERIODICITIES {
  YEARLY = 'yearly',
  QUARTERLY = 'quarterly',
  MONTHLY = 'monthly',
  WEEKLY = 'weekly',
  DAILY = 'daily',
}

export enum FILTER_KEYS {
  DATE_RANGE = 'date_range',
}

export const AMOUNT_RANGE_TYPE_SYMBOL_MAP = {
  [CONDITION_OPERATOR_TYPE.EQUAL]: '=',
  [CONDITION_OPERATOR_TYPE.NOT_EQUAL]: '!=',
  [CONDITION_OPERATOR_TYPE.GREATER_THAN]: '>',
  [CONDITION_OPERATOR_TYPE.LESS_THAN]: '<',
  [CONDITION_OPERATOR_TYPE.GREATER_THAN_EQUAL]: '>=',
  [CONDITION_OPERATOR_TYPE.LESS_THAN_EQUAL]: '<=',
  [CONDITION_OPERATOR_TYPE.IN_BETWEEN]: 'in between',
};

export const AG_GRID_FILTER_TYPES = {
  [FILTER_TYPES.SEARCH]: SearchFilterMenuItem,
  [FILTER_TYPES.ARRAY_SEARCH]: MultiSearchFilterMenuItem,
  [FILTER_TYPES.DATE_RANGE]: DateRangeFilterMenuItem,
  [FILTER_TYPES.AMOUNT_RANGE]: AmountRangeFilterMenuItem,
  [FILTER_TYPES.SINGLE_SELECT]: SingleSelectFilterMenuItem,
  [FILTER_TYPES.MULTI_SELECT]: MultiSelectFilterMenuItem,
  [FILTER_TYPES.TAGS]: TagsSelectFilterMenuItem,
};

export const AG_GRID_FILTER_OPERATORS = {
  [FILTER_TYPES.SEARCH]: 'agTextColumnFilter',
  [FILTER_TYPES.DATE_RANGE]: 'agDateColumnFilter',
  [FILTER_TYPES.AMOUNT_RANGE]: 'agNumberColumnFilter',
  [FILTER_TYPES.MULTI_SELECT]: 'agMultiSelectColumnFilter',
};

export const AG_GRID_FILTER_OPTIONS = {
  [FILTER_TYPES.AMOUNT_RANGE]: [
    'equals',
    'notEqual',
    'lessThan',
    'lessThanOrEqual',
    'greaterThan',
    'greaterThanOrEqual',
    'inRange',
  ],
  [FILTER_TYPES.SEARCH]: ['contains', 'notContains', 'equals', 'notEqual', 'startsWith', 'endsWith'],
};

export const AMOUNT_RANGE_FILTER_OPTIONS = [
  { label: 'is equal to', value: CONDITION_OPERATOR_TYPE.EQUAL },
  { label: 'does not equal to', value: CONDITION_OPERATOR_TYPE.NOT_EQUAL },
  { label: 'is greater than', value: CONDITION_OPERATOR_TYPE.GREATER_THAN },
  { label: 'is greater than or equal to', value: CONDITION_OPERATOR_TYPE.GREATER_THAN_EQUAL },
  { label: 'is less than', value: CONDITION_OPERATOR_TYPE.LESS_THAN },
  { label: 'is less than or equal to', value: CONDITION_OPERATOR_TYPE.LESS_THAN_EQUAL },
  { label: 'is between', value: CONDITION_OPERATOR_TYPE.IN_BETWEEN },
  { label: 'is blank', value: CONDITION_OPERATOR_TYPE.IS_NULL },
];

export const SEARCH_FILTER_OPTIONS = [
  { label: 'contains', value: CONDITION_OPERATOR_TYPE.CONTAINS },
  { label: 'does not contain', value: CONDITION_OPERATOR_TYPE.NOT_CONTAINS },
  { label: 'equals', value: CONDITION_OPERATOR_TYPE.EQUAL },
  { label: 'does not equal', value: CONDITION_OPERATOR_TYPE.NOT_EQUAL },
  { label: 'begins with', value: CONDITION_OPERATOR_TYPE.STARTS_WITH },
  { label: 'ends with', value: CONDITION_OPERATOR_TYPE.ENDS_WITH },
];

export const MULTI_SELECT_FILTER_OPTIONS = [
  { label: 'contains', value: CONDITION_OPERATOR_TYPE.CONTAINS },
  { label: 'does not contain', value: CONDITION_OPERATOR_TYPE.NOT_CONTAINS },
  { label: 'is blank', value: CONDITION_OPERATOR_TYPE.IS_NULL },
];

export const TAGS_SELECT_FILTER_OPTIONS = [
  { label: 'contains', value: CONDITION_OPERATOR_TYPE.CONTAINS },
  { label: 'Untagged', value: CONDITION_OPERATOR_TYPE.IS_NULL },
];

export const CONDITION_OPERATOR_TYPE_LABEL_MAP = {
  [CONDITION_OPERATOR_TYPE.CONTAINS]: 'contains',
  [CONDITION_OPERATOR_TYPE.NOT_CONTAINS]: 'does not contain',
  [CONDITION_OPERATOR_TYPE.IS_NULL]: 'is blank',
  [CONDITION_OPERATOR_TYPE.EQUAL]: 'is equal to',
  [CONDITION_OPERATOR_TYPE.NOT_EQUAL]: 'does not equal to',
  [CONDITION_OPERATOR_TYPE.GREATER_THAN]: 'is greater than',
  [CONDITION_OPERATOR_TYPE.LESS_THAN]: 'is less than',
  [CONDITION_OPERATOR_TYPE.GREATER_THAN_EQUAL]: 'is greater than or equal to',
  [CONDITION_OPERATOR_TYPE.LESS_THAN_EQUAL]: 'is less than or equal to',
  [CONDITION_OPERATOR_TYPE.IN_BETWEEN]: 'is between',
  [CONDITION_OPERATOR_TYPE.STARTS_WITH]: 'begins with',
  [CONDITION_OPERATOR_TYPE.ENDS_WITH]: 'ends with',
  [CONDITION_OPERATOR_TYPE.ONE_OF]: 'is one of',
  [CONDITION_OPERATOR_TYPE.DEBIT]: 'is debit',
  [CONDITION_OPERATOR_TYPE.CREDIT]: 'is credit',
  [CONDITION_OPERATOR_TYPE.ANY]: 'is any',
};

export const OPERATOR = {
  InOperator: { label: 'equals', value: CONDITION_OPERATOR_TYPE.IN },
  EqualOperator: { label: 'equals', value: CONDITION_OPERATOR_TYPE.EQUAL },
  NotInOperator: { label: 'is not equal', value: CONDITION_OPERATOR_TYPE.NOT_IN },
  GreaterThanOperator: { label: 'is greater than', value: CONDITION_OPERATOR_TYPE.GREATER_THAN },
  GreaterThanOrEqualOperator: { label: 'is greater than or equal', value: CONDITION_OPERATOR_TYPE.GREATER_THAN_EQUAL },
  LessThanOperator: { label: 'is less than', value: CONDITION_OPERATOR_TYPE.LESS_THAN },
  LessThanOrEqualOperator: { label: 'is less than or equal', value: CONDITION_OPERATOR_TYPE.LESS_THAN_EQUAL },
  ContainsOperator: { label: 'contains', value: CONDITION_OPERATOR_TYPE.CONTAINS },
  NotContainsOperator: { label: 'does not contain', value: CONDITION_OPERATOR_TYPE.NOT_CONTAINS },
  OneOfOperator: { label: 'is one of', value: CONDITION_OPERATOR_TYPE.ONE_OF },
  Debit: { label: 'Debit', value: CONDITION_OPERATOR_TYPE.DEBIT },
  Credit: { label: 'Credit', value: CONDITION_OPERATOR_TYPE.CREDIT },
  Any: { label: 'Any', value: CONDITION_OPERATOR_TYPE.ANY },
  StartsWithOperator: { label: 'starts with', value: CONDITION_OPERATOR_TYPE.STARTS_WITH },
  EndsWithOperator: { label: 'ends with', value: CONDITION_OPERATOR_TYPE.ENDS_WITH },
  InBetween: { label: 'is in between', value: CONDITION_OPERATOR_TYPE.IN_BETWEEN },
  ArrayContains: { label: 'contains', value: CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS },
  ArrayIn: { label: 'is in', value: CONDITION_OPERATOR_TYPE.ARRAY_IN },
};

export const DESCRIPTION_OPERATORS = [OPERATOR.ArrayContains, OPERATOR.ArrayIn];
