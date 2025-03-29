import React from 'react';
import { ICellRendererParams } from 'ag-grid-community';

const CustomAmountRenderer = (props: ICellRendererParams) => {
  const { colDef, data, value, valueFormatted } = props;
  const valueToDisplay = valueFormatted ?? value;
  const prefix = data[colDef?.cellRendererParams?.config?.currency_column]?.toUpperCase();

  const formattedValue = prefix && valueToDisplay ? `${prefix} ${valueToDisplay}` : valueToDisplay;

  return <div>{formattedValue}</div>;
};

export default CustomAmountRenderer;
