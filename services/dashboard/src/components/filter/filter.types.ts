import { ReactNode } from 'react';
import { MapAny } from 'types/commonTypes';

export type FilterValueTypes = string | null | Array<MapAny> | MapAny;

export enum FILTER_TYPES {
  SEARCH = 'search',
  MULTI_SELECT = 'multi-select',
  SINGLE_SELECT = 'single-select',
  DATE_RANGE = 'date-range',
  AMOUNT_RANGE = 'amount-range',
  TAGS = 'tags',
  ARRAY_SEARCH = 'array-search',
}

export interface FilterMenuType {
  id?: string;
  label?: string | ReactNode;
  value: string | number;
  type?: string;
  is_equal_to?: number;
  is_not_equal_to?: number;
  is_greater_than?: number;
  is_less_than?: number;
}

export enum FILTER_LABEL_TYPES {
  LABEL = 'LABEL',
  COUNT = 'COUNT',
}

export interface FilterConfigType {
  key: string;
  title?: string;
  label: string;
  values: string[];
  type: string;
  datatype: string;
  widgetsInScope: string[];
  targets: {
    dataset_id: string;
    column: string;
  }[];
}

export interface FilterEntityMenuType extends FilterMenuType {
  id: string;
  display_name?: string;
  amount_range_currencies?: FilterMenuType[];
  account_types?: FilterMenuType[];
}
