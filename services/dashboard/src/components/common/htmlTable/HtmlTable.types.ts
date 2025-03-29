import { MapAny } from 'types/commonTypes';

export type HtmlTablePropsType = {
  rows: MapAny[];
  columns: (string | number)[];
  wrapperClassName?: string;
  colCellClassName?: string;
  rowCellClassName?: string;
};
