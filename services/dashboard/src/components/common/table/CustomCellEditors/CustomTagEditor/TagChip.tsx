import { FC, useMemo } from 'react';
import { cn, getTagColor } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

type TagChipProps = { item: string; existingList?: string[]; showIcon?: boolean; externalColor?: string };

const TagChip: FC<TagChipProps> = ({ item, existingList, showIcon = false, externalColor }) => {
  const isExisting = useMemo(() => existingList?.includes(item), [existingList, item]);

  const backgroundColor = useMemo(getTagColor, [item]);

  return (
    <span
      className={cn(isExisting ? '' : 'f-11-400 py-1 px-1.5 rounded-md text-GRAY_1000 flex items-center w-fit')}
      style={isExisting ? {} : { backgroundColor: externalColor ?? backgroundColor }}
    >
      {showIcon && <SvgSpriteLoader id='lightning-01' className='mr-1' height={12} width={12} />}
      <span>{item}</span>
    </span>
  );
};

export default TagChip;
