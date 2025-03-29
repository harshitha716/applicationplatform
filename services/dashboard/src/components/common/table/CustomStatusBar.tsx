import React, { useEffect, useState } from 'react';
import { ColumnHeaderClickedEvent } from 'ag-grid-community';
import { CustomStatusPanelProps } from 'ag-grid-react';
import { MapAny } from 'types/commonTypes';
import { getCommaSeparatedNumber, sentenceCase } from 'utils/common';

const CustomStatusBar = (props: CustomStatusPanelProps & { totalRows?: number; columnLevelStats?: MapAny }) => {
  const [statusBar, setStatusBar] = useState<MapAny>();

  const handleColumnHeaderClicked = (event: ColumnHeaderClickedEvent) => {
    const columnLevelStatsData = props?.columnLevelStats?.[event?.column?.getId()];

    setStatusBar(columnLevelStatsData);
  };

  // Track column header clicked
  useEffect(() => {
    props?.api?.addEventListener('columnHeaderClicked', handleColumnHeaderClicked);

    return () => {
      props?.api?.removeEventListener('columnHeaderClicked', handleColumnHeaderClicked);
    };
  }, [props?.api, handleColumnHeaderClicked]);

  return (
    <div className='flex gap-2 f-11-500 py-2'>
      {statusBar ? (
        <>
          {Object.entries(statusBar).map(([key, value]) => {
            return (
              <div key={key}>
                {sentenceCase(key?.toLowerCase())}: {getCommaSeparatedNumber(value, 2)}
              </div>
            );
          })}
        </>
      ) : null}
      <div>Total Rows: {getCommaSeparatedNumber(props?.totalRows)}</div>
    </div>
  );
};

export default CustomStatusBar;
