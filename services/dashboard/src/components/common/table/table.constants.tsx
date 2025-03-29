import { CellSelectionOptions, themeQuartz } from 'ag-grid-community';
import { COLORS } from 'constants/colors';
import { AggregationFunctionType, LogicalOperatorType } from 'types/components/table.type';
import { DisplayOptionItemProps } from 'components/common/table/DisplayOptions/DisplayOptionItem';
import { DISPLAY_OPTIONS } from 'components/common/table/table.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export const myTheme = themeQuartz.withParams({
  fontFamily: { googleFont: 'Inter' },
  headerFontSize: 12,
  headerFontWeight: 600,
  rowHeight: 32,
  rowBorder: { style: 'solid', width: 1, color: COLORS.GRAY_100 },
  columnBorder: { style: 'solid', width: 1, color: COLORS.GRAY_100 },
  headerHeight: 48,
  headerRowBorder: { style: 'solid', width: 1, color: COLORS.GRAY_400 },
  headerColumnBorder: { style: 'solid', width: 1, color: COLORS.GRAY_100 },
  headerBackgroundColor: COLORS.WHITE,
  wrapperBorderRadius: 0,
  iconSize: 12,
  rowHoverColor: COLORS.BACKGROUND_GRAY_1,
  checkboxBorderRadius: 2,
  checkboxCheckedBackgroundColor: COLORS.GRAY_1000,
  checkboxCheckedBorderColor: COLORS.GRAY_600,
  checkboxCheckedShapeColor: COLORS.WHITE,
  checkboxUncheckedBackgroundColor: COLORS.WHITE,
  checkboxUncheckedBorderColor: COLORS.GRAY_400,
  sideBarBackgroundColor: COLORS.WHITE,
  headerColumnResizeHandleColor: COLORS.WHITE,
  menuBorder: { style: 'solid', width: 1, color: COLORS.GRAY_500 },
  menuBackgroundColor: COLORS.WHITE,
  wrapperBorder: { width: 1, style: 'solid', color: COLORS.GRAY_400 },
  rowLoadingSkeletonEffectColor: COLORS.GRAY_50,
  selectCellBorder: { style: 'solid', width: 1, color: COLORS.BLUE_700 },
  rangeSelectionBorderColor: COLORS.BLUE_700,
  cellEditingBorder: { style: 'solid', width: 1, color: COLORS.BLUE_700 },
  menuShadow: '1px 2px 20px 0px #0000001A',
});

export const myIcons = {
  groupExpanded: `<svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M3 4.5L6 7.5L9 4.5" fill="#8F8F8F"/>
</svg>
`,
  groupContracted: `<svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M4.5 9L7.5 6L4.5 3" fill="#8F8F8F"/>
</svg>
`,
  sortDescending: `<svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M6 2V10M6 10L9 7M6 10L3 7" stroke="#2546F5" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
`,
  sortAscending: `<svg width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M12 20V4M12 4L6 10M12 4L18 10" stroke="#2546F5" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
`,
};

export const PAGE_SIZE = 1000;

export const sideBarConfig = {
  toolPanels: [
    {
      id: 'columns',
      labelDefault: 'Columns',
      labelKey: 'columns',
      iconKey: 'columns',
      toolPanel: 'agColumnsToolPanel',
      toolPanelParams: {
        suppressPivotMode: true, // This removes the "Pivot Mode" toggle
      },
    },
  ],
};

export const DATA_TABLE_THEME_PARAMS = {
  fontFamily: { googleFont: 'Inter' },
  wrapperBorderRadius: 0,
  wrapperBorder: { width: 0 },
  headerFontSize: 11,
  headerFontWeight: 400,
  headerTextColor: '#8F8F8F',
  headerHeight: 36,
  headerRowBorder: { style: 'solid', width: 0.5, color: '#EBEBEB' },
  headerColumnBorder: { width: 0 },
  headerColumnResizeHandleWidth: 0,
  headerBackgroundColor: COLORS.WHITE,
  rowHeight: 60,
  rowBorder: { style: 'solid', width: 0.5, color: '#EBEBEB' },
  rowHoverColor: '#FBFBFB',
  columnBorder: { width: 0 },
  cellHorizontalPadding: 24,
  rowLoadingSkeletonEffectColor: COLORS.GRAY_50,
};

export const DATA_TABLE_CONFIG = {
  filter: undefined,
  headerClass: '',
  cellClass: 'f-12-400 cursor-pointer content-center',
  flex: 1,
};

export const OperatorMap: Record<string, CONDITION_OPERATOR_TYPE> = {
  contains: CONDITION_OPERATOR_TYPE.CONTAINS,
  notContains: CONDITION_OPERATOR_TYPE.NOT_CONTAINS,
  equals: CONDITION_OPERATOR_TYPE.EQUAL,
  notEqual: CONDITION_OPERATOR_TYPE.NOT_EQUAL,
  startsWith: CONDITION_OPERATOR_TYPE.STARTS_WITH,
  endsWith: CONDITION_OPERATOR_TYPE.ENDS_WITH,
};

export const LogicalOperatorMap: Record<string, LogicalOperatorType> = {
  AND: LogicalOperatorType.OperatorLogicalAnd,
  OR: LogicalOperatorType.OperatorLogicalOr,
};

export const AggregationFunctionMap: Record<string, AggregationFunctionType> = {
  sum: AggregationFunctionType.AggregationFunctionSum,
  avg: AggregationFunctionType.AggregationFunctionAvg,
  min: AggregationFunctionType.AggregationFunctionMin,
  max: AggregationFunctionType.AggregationFunctionMax,
  count: AggregationFunctionType.AggregationFunctionCount,
};

export const ArrayFilters = [
  CONDITION_OPERATOR_TYPE.IN,
  CONDITION_OPERATOR_TYPE.NOT_IN,
  CONDITION_OPERATOR_TYPE.NOT_CONTAINS,
  CONDITION_OPERATOR_TYPE.IN_BETWEEN,
  CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS,
  CONDITION_OPERATOR_TYPE.CONTAINS,
];

export const cellSelectionConfig: CellSelectionOptions<any> = {
  handle: {
    mode: 'fill',
    direction: 'y',
  },
};

export const DisplayOptionsList: DisplayOptionItemProps[] = [
  {
    id: DISPLAY_OPTIONS.COLUMNS,
    label: 'Columns',
    iconId: 'columns-03',
  },
  {
    id: DISPLAY_OPTIONS.GROUP_BY,
    label: 'Group By',
    iconId: 'left-indent-02',
  },
  // TODO: Add currency option
  // {
  //   id: DISPLAY_OPTIONS.CURRENCY,
  //   label: 'Currency',
  //   iconId: 'coins-swap-02',
  // },
];
