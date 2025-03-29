import React, { FC } from 'react';
import { cn } from 'utils/common';
import { HtmlTablePropsType } from 'components/common/htmlTable/HtmlTable.types';
import CustomNoRowsOverlay from 'components/common/table/CustomNoRowsOverlay';
import CommonWrapper from 'components/commonWrapper';

const HtmlTable: FC<HtmlTablePropsType> = ({ rows, columns, wrapperClassName, colCellClassName, rowCellClassName }) => {
  return (
    <CommonWrapper
      isNoData={!rows?.length || !columns?.length}
      noDataBanner={
        <div className='h-full w-full flex items-center justify-center'>
          <CustomNoRowsOverlay />
        </div>
      }
      className='h-full w-full'
    >
      <div className={cn('h-full w-full overflow-auto', wrapperClassName)}>
        <table className='border-collapse'>
          <thead>
            <tr>
              {columns?.map((col, index) => (
                <th
                  key={index}
                  className={cn(
                    'border text-start overflow-hidden whitespace-nowrap py-4 px-3.5 text-GRAY_950 border-GRAY_100 f-12-500',
                    colCellClassName,
                  )}
                >
                  {col}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows?.map((row, rowIndex) => (
              <tr key={rowIndex} className='hover:bg-GRAY_20'>
                {columns?.map((col, colIndex) => (
                  <td
                    key={colIndex}
                    className={cn(
                      'border overflow-hidden whitespace-nowrap text-start px-3 py-2 text-GRAY_950 border-GRAY_100 f-11-400',
                      rowCellClassName,
                    )}
                  >
                    {row[col] ?? ''}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </CommonWrapper>
  );
};

export default HtmlTable;
