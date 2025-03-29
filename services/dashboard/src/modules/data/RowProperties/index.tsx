import { FC, useMemo, useState } from 'react';
import { ColDef } from 'ag-grid-community';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { ROUTES_PATH } from 'constants/routeConfig';
import Properties from 'modules/data/RowProperties/Properties';
import { ROW_PROPERTIES_TABS } from 'modules/data/RowProperties/rowProperties.constants';
import { ROW_PROPERTIES_TABS_TYPES, TAG_SOURCE_TYPES } from 'modules/data/RowProperties/rowProperties.types';
import Rules from 'modules/data/RowProperties/Rules';
import { useRouter } from 'next/router';
import { MenuItem, SIZE_TYPES, TAB_TYPES } from 'types/common/components';
import { defaultFnType, MapAny } from 'types/commonTypes';
import { BUTTON_TYPES, ICON_POSITION_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';
import SideDrawer from 'components/common/SideDrawer/SideDrawer';
import { CUSTOM_COLUMNS_TYPE } from 'components/common/table/table.types';
import { Tabs } from 'components/common/tabs/Tabs';

type RowPropertiesSideDrawerProps = {
  onClose: defaultFnType;
  data: MapAny;
  datasetId: string;
  isDrillDownEnabled?: boolean;
  columns: ColDef[];
};

const RowPropertiesSideDrawer: FC<RowPropertiesSideDrawerProps> = ({
  onClose,
  data,
  datasetId,
  isDrillDownEnabled = false,
  columns,
}) => {
  const router = useRouter();

  const [selectedTab, setSelectedTab] = useState<ROW_PROPERTIES_TABS_TYPES>(ROW_PROPERTIES_TABS_TYPES.PROPERTIES);
  const [selectedRuleId, setSelectedRuleId] = useState<string>('');

  const handleSourceDrillDownClick = () => {
    router.push(ROUTES_PATH.DRILLDOWN.replace(':datasetId', datasetId).replace(':rowId', data?._zamp_id as string));
  };

  const handleTabChange = (item?: MenuItem) => {
    setSelectedTab(item?.value as ROW_PROPERTIES_TABS_TYPES);
  };

  const ruleIds = useMemo(() => {
    const ruleIds: string[] = [];
    const tagColumns = columns.filter(
      (column) => column.headerComponentParams?.metadata?.custom_type === CUSTOM_COLUMNS_TYPE.TAG,
    );

    tagColumns.forEach((column) => {
      const sourceColumnId = `_zamp_source_json_${column?.field}`;
      const sourceValue = JSON.parse(data[sourceColumnId] ?? '{}');

      if (sourceValue?.source_type === TAG_SOURCE_TYPES.RULE) {
        ruleIds.push(sourceValue?.source_id);
      }
    });

    return ruleIds;
  }, [columns, data]);

  const getTabContent = () => {
    switch (selectedTab) {
      case ROW_PROPERTIES_TABS_TYPES.PROPERTIES:
        return (
          <Properties
            data={data}
            columns={columns}
            onRuleClick={(ruleId: string) => {
              setSelectedTab(ROW_PROPERTIES_TABS_TYPES.RULES);
              setSelectedRuleId(ruleId);
            }}
          />
        );
      case ROW_PROPERTIES_TABS_TYPES.RULES:
        return <Rules ruleIds={ruleIds} selectedRuleId={selectedRuleId} />;
    }
  };

  return (
    <SideDrawer
      isOpen
      id='row-properties-side-drawer'
      onClose={onClose}
      hideCloseButton
      headerClassName='!p-6'
      topBar={
        <div className='flex items-center justify-between flex-1'>
          <Tabs
            id='row-properties-tabs'
            list={ROW_PROPERTIES_TABS}
            type={TAB_TYPES.FILLED}
            onSelect={handleTabChange}
            customSelectedIndex={selectedTab === ROW_PROPERTIES_TABS_TYPES.RULES ? 1 : 0}
            showSingleAsWell
          />
          {isDrillDownEnabled && (
            <Button
              type={BUTTON_TYPES.SECONDARY}
              id='row-properties-source-drill-down-button'
              className='border-none !text-GRAY_900'
              iconProps={{
                id: 'arrow-up-left',
                iconCategory: ICON_SPRITE_TYPES.ARROWS,
                width: 12,
                height: 12,
              }}
              iconPosition={ICON_POSITION_TYPES.LEFT}
              size={SIZE_TYPES.SMALL}
              onClick={handleSourceDrillDownClick}
            >
              Source drill down
            </Button>
          )}
        </div>
      }
    >
      {getTabContent()}
    </SideDrawer>
  );
};

export default RowPropertiesSideDrawer;
