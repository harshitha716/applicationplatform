import { ElementType } from 'react';
import { MenuItem, TAB_TYPES } from 'types/common/components';

export type TabsPropsType = {
  list: Array<MenuItem>;
  id: string;
  type: TAB_TYPES;
  customSelectedIndex?: number;
  onSelect?: (item?: MenuItem) => void;
  TabItemComponent?: ElementType | null;
  tabItemWrapperClassName?: string;
  scrollWrapperClassName?: string;
  wrapperClassName?: string;
  indicatorClassName?: string;
  selectedTabIndicatorClassName?: string;
  scrollWrapperStyle?: string;
  wrapperStyle?: string;
  tabItemWrapperStyle?: string;
  tabItemStyle?: string;
  tabItemSelectedStyle?: string;
  tabItemDefaultStyle?: string;
  indicatorStyle?: string;
  selectedTabIndicatorStyle?: string;
  disabled?: boolean;
  showSingleAsWell?: boolean;
};
