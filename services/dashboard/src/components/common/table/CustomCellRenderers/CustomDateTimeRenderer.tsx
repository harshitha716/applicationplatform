import React from 'react';
import { ICellRendererParams } from 'ag-grid-community';
import { DATE_FORMATS, VALID_DATE_FORMATS } from 'constants/date.constants';
import { format, isValid } from 'date-fns';

const CustomDateTimeRenderer = (props: ICellRendererParams) => {
  const { colDef, value } = props;
  const dateFormat = colDef?.cellRendererParams?.config?.format;
  const validDateFormat = VALID_DATE_FORMATS.includes(dateFormat) ? dateFormat : DATE_FORMATS.ddMMMyyyy;
  const date = new Date(value);
  const formattedValue = isValid(date) ? format(date, validDateFormat) : value;

  return <div>{formattedValue}</div>;
};

export default CustomDateTimeRenderer;
