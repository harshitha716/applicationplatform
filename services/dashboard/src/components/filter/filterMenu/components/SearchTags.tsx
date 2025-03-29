import { FC } from 'react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { MapAny } from 'types/commonTypes';
import { stopPropagationAction } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export enum DESCRIPTION_TAGS {
  DESCRIPTION_PROPERTY = 'description_property',
  DESCRIPTION_VALUE = 'description_value',
}

interface SearchTagsProps {
  tags: MapAny[];
  onDeleteTag: (index: number) => void;
}

const SearchTags: FC<SearchTagsProps> = ({ tags, onDeleteTag }) => {
  return tags?.map((tag, index) => (
    <div
      key={index}
      onClick={stopPropagationAction}
      className={`h-[20px] whitespace-nowrap w-auto py-1 px-2 f-12-400 flex items-center justify-between gap-2 ${
        tag?.type === DESCRIPTION_TAGS.DESCRIPTION_VALUE ? 'bg-GRAY_100' : 'bg-BG_TAB_SELECTION rounded-[26px]'
      }`}
    >
      {tag?.type === DESCRIPTION_TAGS.DESCRIPTION_PROPERTY && (
        <SvgSpriteLoader
          iconCategory={ICON_SPRITE_TYPES.WEATHER}
          id='stars-02'
          width={12}
          height={12}
          color={COLORS.BLUE_150}
        />
      )}
      {tag?.label}
      <SvgSpriteLoader
        id='x-close'
        onClick={() => onDeleteTag(index)}
        className='cursor-pointer'
        iconCategory={ICON_SPRITE_TYPES.GENERAL}
        height={9}
        width={9}
      />
    </div>
  ));
};

export default SearchTags;
