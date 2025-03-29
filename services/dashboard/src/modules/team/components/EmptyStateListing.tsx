import React, { FC } from 'react';
import { EmptyStateListingPropsType } from 'modules/team/people.types';

const EmptyStateListing: FC<EmptyStateListingPropsType> = ({ title = 'Nothing to show up' }) => {
  return <span className='f-16-500 flex justify-center items-center h-3/5 w-full text-GRAY_600'>{title}</span>;
};

export default EmptyStateListing;
