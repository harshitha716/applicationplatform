import { FC, Fragment } from 'react';
import { ColDef } from 'ag-grid-community';
import { useGetAudiencesByOrganisationIdQuery } from 'apis/people';
import { useAppSelector } from 'hooks/toolkit';
import PropertyRow from 'modules/data/RowProperties/PropertyRow';
import { RootState } from 'store';
import { MapAny } from 'types/commonTypes';

type PropertiesProps = {
  data: MapAny;
  columns: ColDef[];
  onRuleClick: (ruleId: string) => void;
};

const Properties: FC<PropertiesProps> = ({ data, columns, onRuleClick }) => {
  const organizationId = useAppSelector((state: RootState) => state?.user?.user?.orgs?.[0]?.organization_id) ?? '';
  const { data: teamMembersData } = useGetAudiencesByOrganisationIdQuery({ organizationId }, { skip: !organizationId });

  return (
    <div className='grid grid-cols-2 gap-2.5'>
      {Object.entries(data).map(([key, value]) => {
        const column = columns.find((column) => column.field === key);

        return value && column ? (
          <Fragment key={key}>
            <PropertyRow
              value={value}
              column={column}
              data={data}
              teamMembersData={teamMembersData}
              onRuleClick={onRuleClick}
            />
          </Fragment>
        ) : null;
      })}
    </div>
  );
};

export default Properties;
