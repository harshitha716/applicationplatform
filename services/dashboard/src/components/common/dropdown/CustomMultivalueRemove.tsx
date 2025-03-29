import React from 'react';
import { components, MultiValueRemoveProps } from 'react-select';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { OptionsType } from 'types/commonTypes';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export const CustomMultivalueRemove = (props: MultiValueRemoveProps<OptionsType>) => {
  return (
    <components.MultiValueRemove {...props}>
      <SvgSpriteLoader
        id='x-close'
        iconCategory={ICON_SPRITE_TYPES.GENERAL}
        color={COLORS.TEXT_TERTIARY}
        width={16}
        height={16}
      />
    </components.MultiValueRemove>
  );
};
