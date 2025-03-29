import { FC } from 'react';
import { CustomMenuItemProps } from 'ag-grid-react';
import { COLORS } from 'constants/colors';
import { defaultFnType } from 'types/commonTypes';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface CustomContextMenuItemProps extends CustomMenuItemProps {
  action: defaultFnType;
}

const CustomContextMenuItem: FC<CustomContextMenuItemProps> = ({ action, menuItemParams, name, closeMenu }) => {
  const handleClick = () => {
    action();
    closeMenu();
  };

  return (
    <div
      className='flex items-center gap-1.5 group hover:bg-GRAY_100 py-2 px-2.5 cursor-pointer mx-1 rounded-md'
      onClick={handleClick}
    >
      <SvgSpriteLoader id={menuItemParams.iconId} color={COLORS.GRAY_900} width={12} height={12} />
      <span className='text-GRAY_900 f-12-500 group-hover:text-GRAY_1000'>{name}</span>
    </div>
  );
};

export default CustomContextMenuItem;
