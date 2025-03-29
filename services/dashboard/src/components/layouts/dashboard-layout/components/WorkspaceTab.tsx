import React from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { defaultFnType } from 'types/commonTypes';
import { cn } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface WorkspaceTabProps {
  label: string;
  isSelected?: boolean;
  onClick?: defaultFnType;
  className?: string;
  color?: string;
}

const WorkspaceTab = ({ label, isSelected, onClick, className, color }: WorkspaceTabProps) => {
  return (
    <div
      onClick={(e) => {
        e.preventDefault();
        onClick?.();
      }}
      className={cn(
        'flex items-center gap-1 px-2 py-2.5 f-13-500 select-none cursor-pointer rounded-md',
        onClick ? 'hover:bg-GRAY_20' : '',
        isSelected ? 'bg-GRAY_100' : '',
        className,
      )}
    >
      <div
        className={cn('w-3.5 h-3.5 flex items-center justify-center rounded-sm f-9-600 text-white mr-1.5')}
        style={{ backgroundColor: color }}
      >
        {label.charAt(0).toUpperCase()}
      </div>
      <div className='flex-1'>{label}</div>
      {isSelected && (
        <SvgSpriteLoader
          id='check'
          iconCategory={ICON_SPRITE_TYPES.GENERAL}
          width={14}
          height={14}
          className='min-w-4 float-right'
        />
      )}
    </div>
  );
};

export default WorkspaceTab;
