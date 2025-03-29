import { FC, useEffect, useMemo, useState } from 'react';
import { Responsive, WidthProvider } from 'react-grid-layout';
import { useGetRulesByDatasetColumnsQuery, useUpdateRulePriorityMutation } from 'apis/dataset';
import { DRAG_ICON } from 'constants/icons';
import { ZAMP_LOGO_LOADER } from 'constants/lottie/zamp-logo-loader';
import { DatasetColumnRequest } from 'modules/data/data.types';
import RuleCard, { RuleCardProps } from 'modules/data/RulesListing/RuleCard';
import { searchRules } from 'modules/data/RulesListing/ruleListing.utils';
import Image from 'next/image';
import { DatasetUpdateResponseType } from 'types/api/dataset.types';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { BUTTON_TYPES, ICON_POSITION_TYPES } from 'types/components/button.type';
import { OrderType } from 'types/components/table.type';
import { getUserId } from 'utils/accessPermission/accessPermission.utils';
import { Button } from 'components/common/button/Button';
import Input from 'components/common/input';
import Popup from 'components/common/popup/Popup';
import SideDrawer from 'components/common/SideDrawer/SideDrawer';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import DynamicLottiePlayer from 'components/DynamicLottiePlayer';
import { getTagLabel } from 'components/filter/filter.utils';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const ResponsiveGridLayout = WidthProvider(Responsive);

type RulesListingSideDrawerProps = {
  onClose: defaultFnType;
  datasetId: string;
  column: string;
  handleSuccessfulUpdate: (data: DatasetUpdateResponseType) => void;
};

const RulesListingSideDrawer: FC<RulesListingSideDrawerProps> = ({
  onClose,
  datasetId,
  column,
  handleSuccessfulUpdate,
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [rules, setRules] = useState<RuleCardProps[]>([]);
  const [prioritySorting, setPrioritySorting] = useState<OrderType>(OrderType.DESC);
  const [isApplyChangesPopupOpen, setIsApplyChangesPopupOpen] = useState(false);
  const [isApplyChangesEnabled, setIsApplyChangesEnabled] = useState(false);
  // State for grid layout
  const [layout, setLayout] = useState(
    rules?.map((rule, index) => ({
      i: rule?.id,
      x: 0,
      y: index,
      w: 1,
      h: 1,
    })),
  );
  const { data, isLoading, isError } = useGetRulesByDatasetColumnsQuery(
    {
      dataset_columns: JSON.stringify([{ dataset_id: datasetId, columns: [column] }] as DatasetColumnRequest[]),
    },
    { skip: !datasetId || !column },
  );
  const [updateRulePriority, { isLoading: isUpdating }] = useUpdateRulePriorityMutation();

  const listOfRules = useMemo(() => {
    return (
      data?.[datasetId]?.[column]?.map((rule) => {
        return {
          filters: rule?.filter_config?.query_config?.filters,
          value: rule?.value,
          createdOn: rule?.created_at,
          id: rule?.rule_id,
          priority: rule?.priority,
        };
      }) ?? []
    );
  }, [data, datasetId, column]);

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;

    setSearchTerm(value);

    const filteredRules = value ? searchRules(listOfRules, value) : listOfRules;

    setRules(filteredRules);
    setLayout(
      filteredRules?.map((rule, index) => ({
        i: rule?.id,
        x: 0,
        y: index,
        w: 1,
        h: 1,
      })),
    );
  };

  // Handle layout change
  const handleLayoutChange = (newLayout: any) => {
    setLayout(newLayout);
    // Optional: Update item order based on layout
    const orderedItems: RuleCardProps[] = newLayout
      .slice()
      .sort((a: any, b: any) => a.y - b.y)
      .map((l: any) => rules.find((rule) => rule?.id === l.i)!);

    const isOrderSame = orderedItems.every((rule, index) => rule?.id === listOfRules[index]?.id);

    setIsApplyChangesEnabled(!isOrderSame);

    setRules(orderedItems);
  };

  const handlePrioritySorting = () => {
    const updatedSortingOrder = prioritySorting === OrderType.ASC ? OrderType.DESC : OrderType.ASC;
    const orderedItems = rules.sort((a, b) => {
      if (updatedSortingOrder === OrderType.ASC) {
        return (a.priority ?? 0) - (b.priority ?? 0);
      } else {
        return (b.priority ?? 0) - (a.priority ?? 0);
      }
    });

    setLayout(
      orderedItems?.map((rule, index) => ({
        i: rule?.id,
        x: 0,
        y: index,
        w: 1,
        h: 1,
      })),
    );
    setRules(orderedItems);
    setPrioritySorting(updatedSortingOrder);
  };

  const handleApplyChanges = () => {
    // check if  order of orderedItems are not same as list of rules
    const isOrderSame = rules.every(
      (rule, index) =>
        rule?.id === listOfRules[prioritySorting === OrderType.DESC ? index : listOfRules.length - index]?.id,
    );

    if (!isOrderSame) {
      updateRulePriority({
        dataset_id: datasetId,
        column,
        rule_priorities: {
          updated_by: getUserId(),
          rule_priority: rules.map((rule, index) => ({
            rule_id: rule?.id ?? '',
            priority: prioritySorting === OrderType.ASC ? index + 1 : listOfRules.length - index,
          })),
        },
      }).then((res) => {
        handleApplyChangesPopupClose();
        onClose();
        handleSuccessfulUpdate(res?.data as DatasetUpdateResponseType);
      });
    }
  };

  const handleApplyChangePopupOpen = () => {
    setIsApplyChangesPopupOpen(true);
  };

  const handleApplyChangesPopupClose = () => {
    setIsApplyChangesPopupOpen(false);
  };

  const handleDiscardChanges = () => {
    handleApplyChangesPopupClose();
    const orderedItems = rules.sort((a, b) => {
      if (prioritySorting === OrderType.ASC) {
        return (a.priority ?? 0) - (b.priority ?? 0);
      } else {
        return (b.priority ?? 0) - (a.priority ?? 0);
      }
    });

    setLayout(
      orderedItems?.map((rule, index) => ({
        i: rule?.id,
        x: 0,
        y: index,
        w: 1,
        h: 1,
      })),
    );
    setRules(orderedItems);
  };

  const handleExpandRule = (id: string) => {
    setLayout((prev) => prev.map((item) => (item.i === id ? { ...item, h: 2.2 } : item)));
  };

  const handleCollapseRule = (id: string) => {
    setLayout((prev) => prev.map((item) => (item.i === id ? { ...item, h: 1 } : item)));
  };

  useEffect(() => {
    setRules(listOfRules);
  }, [listOfRules]);

  return (
    <SideDrawer
      isOpen
      id='rules-listing-side-drawer'
      onClose={onClose}
      hideCloseButton
      childrenWrapperClassName='!px-0'
    >
      <div className='h-full mt-2'>
        <div className='px-6'>
          <div className='flex justify-between items-center'>
            <div className='f-16-600'>{column}</div>
            {isApplyChangesEnabled && (
              <Button
                type={BUTTON_TYPES.SECONDARY}
                size={SIZE_TYPES.XSMALL}
                onClick={handleApplyChangePopupOpen}
                id='apply-priority-changes'
                iconProps={{
                  id: 'check',
                }}
                iconPosition={ICON_POSITION_TYPES.LEFT}
              >
                Apply Changes
              </Button>
            )}
          </div>
          <div className='flex justify-between items-center'>
            <Input
              placeholder='Search'
              size={SIZE_TYPES.XSMALL}
              noBorders
              focusClassNames='mt-6 mb-3.5 !px-0'
              onChange={handleSearch}
              value={searchTerm}
            />
            <div className='flex items-center gap-1 text-GRAY_700 cursor-pointer' onClick={handlePrioritySorting}>
              <SvgSpriteLoader
                id={prioritySorting === OrderType.DESC ? 'arrow-narrow-down' : 'arrow-narrow-up'}
                width={14}
                height={14}
              />
              <div className='f-12-400 select-none'>Priority</div>
            </div>
          </div>
        </div>
        <CommonWrapper
          isLoading={isLoading}
          isError={isError}
          className='h-[calc(100vh-120px)] overflow-auto pl-1.5'
          skeletonType={SkeletonTypes.CUSTOM}
          loader={
            <div className='flex justify-center items-center h-full'>
              <DynamicLottiePlayer
                src={ZAMP_LOGO_LOADER}
                className='lottie-player h-[140px]'
                autoplay
                loop
                keepLastFrame
              />
            </div>
          }
        >
          <ResponsiveGridLayout
            className='layout'
            layouts={{ lg: layout }}
            breakpoints={{ lg: 1200 }}
            cols={{ lg: 1 }} // Single-column layout
            isResizable={false} // Disable resizing
            onLayoutChange={handleLayoutChange} // Handle drag-and-drop reordering
            draggableHandle='.drag-handle' // Restrict drag to the handle
            rowHeight={118}
            margin={[14, 14]}
            containerPadding={[0, 0]}
          >
            {rules?.map((rule) => (
              <div
                key={rule?.id}
                data-grid={layout?.find((layout) => layout.i === rule?.id)}
                className='flex items-center gap-1'
              >
                <div className='drag-handle cursor-grab min-w-[14px]'>
                  <Image src={DRAG_ICON} width={14} height={14} alt='drag icon' className='rotate-90' priority />
                </div>
                <RuleCard
                  filters={rule?.filters}
                  value={getTagLabel(rule?.value ?? '')}
                  createdOn={rule?.createdOn}
                  key={rule?.id}
                  className='w-[428px]'
                  onExpand={handleExpandRule}
                  onCollapse={handleCollapseRule}
                  id={rule?.id}
                />
              </div>
            ))}
          </ResponsiveGridLayout>
        </CommonWrapper>
      </div>
      <Popup
        isOpen={isApplyChangesPopupOpen}
        onClose={handleApplyChangesPopupClose}
        title='Apply Changes ?'
        iconId='x-close'
        className='w-[344px] border-2 border-GRAY_400 rounded-3.5 bg-white !p-0 shadow-menuList'
        titleClassName='f-16-600 text-GRAY_950'
        showIcon
      >
        <div className='f-13-400 text-GRAY_900 px-5 py-6'>
          You&apos;ve updated the priority of rules. Do you want to apply these changes before leaving?
        </div>
        <div className='flex justify-end gap-2 px-5 py-4 border-t border-GRAY_400'>
          <Button
            type={BUTTON_TYPES.SECONDARY}
            size={SIZE_TYPES.MEDIUM}
            id='cancel-apply-changes'
            onClick={handleDiscardChanges}
          >
            Discard
          </Button>
          <Button
            type={BUTTON_TYPES.PRIMARY}
            size={SIZE_TYPES.MEDIUM}
            id='apply-changes'
            onClick={handleApplyChanges}
            isLoading={isUpdating}
          >
            Yes, Apply
          </Button>
        </div>
      </Popup>
    </SideDrawer>
  );
};

export default RulesListingSideDrawer;
