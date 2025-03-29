import { FC, useState } from 'react';
import { ColDef } from 'ag-grid-community';
import { DATE_FORMATS } from 'constants/date.constants';
import { format, isValid } from 'date-fns';
import { TAG_SOURCE_TYPES } from 'modules/data/RowProperties/rowProperties.types';
import { AudiencesByOrganisationIdResponse } from 'types/api/people.types';
import { MapAny } from 'types/commonTypes';
import { copyToClipBoard } from 'utils/common';
import { Label } from 'components/common/Label';
import TagChip from 'components/common/table/CustomCellEditors/CustomTagEditor/TagChip';
import { CUSTOM_COLUMNS_TYPE } from 'components/common/table/table.types';
import { Tooltip } from 'components/common/tooltip';
import { getTagLabel } from 'components/filter/filter.utils';

type PropertyRowProps = {
  value: any;
  column: ColDef;
  data: MapAny;
  teamMembersData?: AudiencesByOrganisationIdResponse[];
  onRuleClick: (ruleId: string) => void;
};

const PropertyRow: FC<PropertyRowProps> = ({ value, column, data, teamMembersData, onRuleClick }) => {
  const [showCopyTooltip, setShowCopyTooltip] = useState(false);

  const handleCopy = () => {
    copyToClipBoard(value);
    setShowCopyTooltip(true);
    setTimeout(() => {
      setShowCopyTooltip(false);
    }, 2000);
  };

  const getValue = (column: ColDef, value: any) => {
    if (column.cellRenderer) {
      if (column.headerComponentParams?.metadata?.custom_type === CUSTOM_COLUMNS_TYPE.TAG) {
        const sourceColumnId = `_zamp_source_json_${column?.field}`;
        const sourceValue = JSON.parse(data[sourceColumnId] ?? '{}');

        let tooltipTitle: string = '';
        let tooltipDescription: string = '';

        if (sourceValue?.source_type === TAG_SOURCE_TYPES.USER) {
          const teamMember = teamMembersData?.find((member) => member.user?.user_id === sourceValue?.source_id);
          const date = new Date(sourceValue?.updated_at);
          const formattedDate = isValid(date) ? format(date, DATE_FORMATS.ddMMMyyyy) : value;

          tooltipTitle = `Applied Manually by ${teamMember?.user?.name ?? ''}`;
          tooltipDescription = `on ${formattedDate}`;
        } else {
          tooltipTitle = 'Applied by rule';
          tooltipDescription = 'Click to view more info';
        }

        return (
          <Tooltip
            tooltipBody={
              <Label
                title={tooltipTitle}
                description={tooltipDescription}
                titleClassName='f-10-400 text-white mb-[4px]'
                descriptionClassName='f-10-400 text-GRAY_700 text-wrap break-keep'
              />
            }
            tooltipBodyClassName='f-12-300 px-3 py-2 rounded-md whitespace-nowrap z-999 bg-black text-white w-[102px]'
            className='z-1'
          >
            <div
              onClick={
                sourceValue.source_type === TAG_SOURCE_TYPES.RULE
                  ? () => onRuleClick(sourceValue?.source_id)
                  : undefined
              }
            >
              <TagChip item={getTagLabel(value)} showIcon />
            </div>
          </Tooltip>
        );
      }

      return column.cellRenderer({ colDef: column, data, value });
    }

    return (
      <Tooltip
        tooltipBody={showCopyTooltip ? 'Copied' : 'Click to copy'}
        tooltipBodyClassName='f-12-300 px-3 py-1.5 rounded-md whitespace-nowrap z-999 bg-black text-white'
        className='z-1'
      >
        <div className='hover:bg-GRAY_100 p-1 rounded-md' onClick={handleCopy}>
          {value}
        </div>
      </Tooltip>
    );
  };

  return (
    <>
      <div className='f-12-400 text-GRAY_700 h-6 flex items-center'>
        <p>{column?.headerName ?? column?.field}</p>
      </div>
      <div className='f-11-400 text-GRAY_1000 min-h-6 h-fit flex items-center break-all'>{getValue(column, value)}</div>
    </>
  );
};

export default PropertyRow;
