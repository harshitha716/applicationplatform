import { ColumnContext } from 'modules/widgets/Pivot/pivot.utils';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export type PivotColumnMetadata =
  | {
      kind: 'group';
      name: string;
      dataType: 'string' | 'number' | 'date';
      sourceName: string;
      alias: string;
      mappingName: string;
      heirarchy: number;
      hasChildren: boolean;
      maxHeirarchy: number;
    }
  | {
      kind: 'pivot';
      name: string;
      dataType: 'string' | 'number' | 'date';
      sourceName: string;
      alias: string;
      mappingName: string;
    }
  | {
      kind: 'aggregate';
      name: string;
      dataType: 'string' | 'number' | 'date';
      aggregation: string;
      sourceName: string;
      alias: string;
      mappingName: string;
    };

export enum PIVOT_DATA_TYPES {
  STRING = 'string',
  NUMBER = 'number',
  DATE = 'date',
  STATUS = 'status',
  TIMESTAMP = 'timestamp',
  COUNTRY = 'country',
  BANK = 'bank',
  TAG = 'tag',
  BOOLEAN = 'boolean',
  AMOUNT = 'amount',
}

export type MappingDetails = {
  column?: string;
  drilldown_filter_type?: string;
  drilldown_filter_operator?: string;
};

export type ParentMappingDetail = {
  key: string;
  tag: boolean;
  mappingDetails: MappingDetails | null;
};

export enum UNTAGGED_TAGS {
  UNTAGGED = '__UNTAGGED__',
}

export enum UNTAGGED_TAGS_FRONTEND_MAPPING {
  UNTAGGED = 'Untagged',
}

export type FilterConfig = {
  filterType?: string;
  type?: string;
  values?: string[];
  dateFrom?: string;
  dateTo?: string;
  column?: string;
  targets?: string[];
};

export type ParentFilters = Record<string, FilterConfig>;

export type PivotContext = {
  filterContext: Record<string, ColumnFilterConfig[]>;
  widgetMappingDatasets: Record<string, string>;
  columnContextMapping: Record<string, Record<string, ColumnContext>>;
};

export type ColumnFilterConfig = {
  column: string;
} & (
  | {
      filterType: FILTER_TYPES.MULTI_SELECT;
      type: CONDITION_OPERATOR_TYPE.IN;
    }
  | {
      filterType: FILTER_TYPES.DATE_RANGE;
      type: CONDITION_OPERATOR_TYPE.IN_BETWEEN;
    }
  | {
      filterType: FILTER_TYPES.SEARCH;
      type: CONDITION_OPERATOR_TYPE.STARTS_WITH;
    }
);
