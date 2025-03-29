import { CustomHeaderMenuOptionTypes } from 'components/common/table/CustomHeader/customHeader.types';

export const CustomHeaderMenuOptions = [
  {
    label: 'Rules',
    value: CustomHeaderMenuOptionTypes.RULES,
    iconId: 'lightning-01',
  },
  {
    label: 'Sort Ascending',
    value: CustomHeaderMenuOptionTypes.SORT_ASC,
    iconId: 'arrow-up',
  },
  {
    label: 'Sort Descending',
    value: CustomHeaderMenuOptionTypes.SORT_DESC,
    iconId: 'arrow-down',
  },
  {
    label: 'Remove Sort',
    value: CustomHeaderMenuOptionTypes.REMOVE_SORT,
    iconId: 'x-close',
  },
  {
    label: 'Filter',
    value: CustomHeaderMenuOptionTypes.FILTER,
    iconId: 'filter-lines',
  },
];
