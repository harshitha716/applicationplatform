export interface NavigationItemSchema {
  label: string;
  iconId: string;
  path: string;
  children?: NavigationItemSchema[];
  isHidden?: boolean;
}
