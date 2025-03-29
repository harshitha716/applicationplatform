import { FC, useEffect, useRef, useState } from 'react';
import { TAB_TYPES } from 'types/common/components';
import { defaultFn } from 'types/commonTypes';
import { cn } from 'utils/common';
import { TAB_STYLES } from 'components/common/tabs/tabs.constants';
import { TabsPropsType } from 'components/common/tabs/tabs.types';

export const Tabs: FC<TabsPropsType> = ({
  list = [],
  customSelectedIndex = 0,
  onSelect = defaultFn,
  TabItemComponent = null,
  scrollWrapperClassName = 'pb-0.5',
  wrapperClassName = '',
  tabItemWrapperClassName = 'relative flex justify-center items-center w-fit',
  scrollWrapperStyle = '',
  wrapperStyle = '',
  tabItemWrapperStyle = '',
  tabItemStyle = '',
  tabItemSelectedStyle = '',
  tabItemDefaultStyle = '',
  id = '',
  disabled,
  type = TAB_TYPES.FILLED,
  showSingleAsWell = false,
}) => {
  const [selectedIndex, setSelectedIndex] = useState<number>(customSelectedIndex);
  const firstLoad = useRef(true);

  const handleSelect = (itemIndex: number) => {
    if (!disabled) {
      setSelectedIndex(itemIndex);
      onSelect(list[itemIndex]);
    }
  };

  useEffect(() => {
    firstLoad.current = false;
  }, [selectedIndex, id, list]);

  useEffect(() => {
    setSelectedIndex((prevIndex) => (customSelectedIndex !== prevIndex ? customSelectedIndex : prevIndex));
  }, [customSelectedIndex]);

  const TabItem = TabItemComponent ?? null;
  const hasMultipleItems = list?.length > 1;
  const { tabItemSelectedClassName, tabItemDefaultClassName, tabItemGapClassName, tabItemClassName } = TAB_STYLES[type];

  return (
    <>
      {(hasMultipleItems || showSingleAsWell) && (
        <div className={`overflow-hidden no-scrollbar select-none ${scrollWrapperClassName} ${scrollWrapperStyle}`}>
          <div className={`flex w-full ${wrapperClassName} ${wrapperStyle}`}>
            {list?.map((tabItem, index) => {
              const selected = index === selectedIndex;
              const last = index === list.length - 1;

              return (
                !tabItem.isHidden && (
                  <div
                    role='presentation'
                    key={index}
                    onClick={() => handleSelect(index)}
                    className={`cursor-pointer ${
                      !last ? tabItemGapClassName : ''
                    } ${tabItemWrapperClassName} ${tabItemWrapperStyle}`}
                    data-testid={`TAB_ITEM_${tabItem.label}`}
                  >
                    {TabItem ? (
                      <TabItem
                        data={tabItem}
                        selected={selected}
                        className={`${tabItemClassName} ${tabItemStyle}`}
                        selectedClassName={`${tabItemSelectedClassName} ${tabItemSelectedStyle}`}
                        defaultClassName={`${tabItemDefaultClassName} ${tabItemDefaultStyle}`}
                        index={index}
                      />
                    ) : (
                      <div
                        className={cn(
                          'flex gap-1 w-full justify-center items-center f-12-500',
                          tabItemClassName,
                          tabItemStyle,
                          selected
                            ? `${tabItemSelectedClassName} ${tabItemSelectedStyle}`
                            : `${tabItemDefaultClassName} ${tabItemDefaultStyle}`,
                        )}
                      >
                        {tabItem?.label}
                        {!!tabItem?.metadata?.count && (
                          <div className='f-12-300 border border-BORDER_7 text-GRAY_600 px-1 bg-white h-4 !leading-3.5'>
                            {tabItem?.metadata?.count}
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                )
              );
            })}
          </div>
        </div>
      )}
    </>
  );
};
