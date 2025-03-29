import React, { FC, memo } from 'react';
import Image from 'next/image';
import { cn, getFirstLetters, getRandomColor } from 'utils/common';

type AvatarProps = {
  className?: string;
  name: string;
  imageUrl?: string;
  backgroundColor?: string;
};

const Avatar: FC<AvatarProps> = ({
  className = 'f-12-300 font-normal h-8 min-w-[32px]',
  name,
  imageUrl,
  backgroundColor,
}) => (
  <div
    className={cn('flex justify-center items-center rounded-full relative text-white', className)}
    style={{ backgroundColor: backgroundColor || getRandomColor() }}
  >
    {imageUrl ? <Image src={imageUrl} alt={name} fill /> : <>{getFirstLetters(name)}</>}
  </div>
);

export default memo(Avatar);
