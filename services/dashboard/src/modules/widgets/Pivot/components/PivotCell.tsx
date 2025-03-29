import { FC, memo, useMemo, useRef, useState } from 'react';
import { Column, GridApi, IRowNode } from 'ag-grid-community';
import { CURRENCY_SYMBOLS } from 'modules/page/pages.constants';
import { MapAny } from 'types/commonTypes';
import { cn, getCommaSeparatedNumber } from 'utils/common';

interface PivotCellProps {
  node: IRowNode;
  value: string | number;
  maxGroupingLevel: number;
  showPercentage?: MapAny;
  column?: Column;
  api?: GridApi;
  currency?: string;
}

const PivotCell: FC<PivotCellProps> = ({ node, value, maxGroupingLevel, showPercentage, api, column, currency }) => {
  const [toggledRows, setToggledRows] = useState<Record<string, boolean>>({});
  const clickTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const formattedValue = useMemo(() => {
    const parsedValue = parseFloat(value?.toString().replace(/[^0-9.-]/g, '') || '0');
    const numericValue = typeof value === 'number' ? value : parsedValue;

    if (isNaN(numericValue)) return '-';

    const currencySymbol = CURRENCY_SYMBOLS[currency as keyof typeof CURRENCY_SYMBOLS] ?? currency;

    return currency ? `${currencySymbol} ${getCommaSeparatedNumber(numericValue, 2)}` : numericValue;
  }, [currency, value]);

  const { isLastNode, isTopNode, isRootLevel } = useMemo(
    () => ({
      isLastNode: node.level === maxGroupingLevel,
      isTopNode: node.level === 0,
      isRootLevel: node.level === -1,
    }),
    [node.level, maxGroupingLevel],
  );

  const isLastCell = useMemo(() => {
    const displayedColumns = api?.getAllDisplayedColumns();

    return displayedColumns?.[displayedColumns.length - 1] === column;
  }, [column, api]);

  const numericValue = useMemo(() => {
    return typeof value === 'number' ? value : parseFloat(value.toString().replace(/[$,]/g, ''));
  }, [value]);

  const aggData = node?.aggData || {};
  const matchingField = Object.keys(aggData).find(
    (key) => parseFloat(aggData[key]?.toFixed(2)) === parseFloat(numericValue?.toFixed(2)),
  );

  const totalValue = matchingField ? parseFloat(node?.parent?.aggData?.[matchingField]?.toFixed(2)) || 0 : 0;
  const isToggled = toggledRows[node?.id || node?.key || ''];
  const { only_parent = false } = showPercentage || {};

  const percentageValue = useMemo(() => {
    return totalValue > 0 ? `${getCommaSeparatedNumber(((numericValue || 0) / totalValue) * 100, 2)}%` : '0.00%';
  }, [numericValue, totalValue]);

  const shouldShowPercentage = showPercentage && !isRootLevel && (only_parent ? isTopNode : true);
  const displayValue = shouldShowPercentage ? (isToggled ? formattedValue : percentageValue) : formattedValue;

  const handleToggle = () => {
    if (shouldShowPercentage) {
      if (clickTimeoutRef.current) {
        clearTimeout(clickTimeoutRef.current);
        clickTimeoutRef.current = null;

        return;
      }

      clickTimeoutRef.current = setTimeout(() => {
        setToggledRows((prev) => ({
          ...prev,
          [node?.id || node?.key || '']: !prev[node?.id || node?.key || ''],
        }));
        clickTimeoutRef.current = null;
      }, 200);
    }
  };

  return (
    <div
      className={cn(
        'h-full w-full flex items-center justify-end gap-3 px-3 py-2 text-GRAY_950 border-b-0.5 border-b-GRAY_400 border-r-0.5 border-r-GRAY_400 f-13-450 cursor-pointer select-none hover:bg-GRAY_100',
        {
          'bg-BACKGROUND_GRAY_1': isLastNode || isRootLevel,
          'border-r-0': isLastCell,
          'border-b-0': isRootLevel,
        },
      )}
      onClick={handleToggle}
    >
      {displayValue}
    </div>
  );
};

export default memo(PivotCell);
