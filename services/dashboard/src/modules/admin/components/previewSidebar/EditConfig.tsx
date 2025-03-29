import React from 'react';
import { COLORS } from 'constants/colors';
import { cn } from 'utils/common';
import Input from 'components/common/input';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const EditConfig = ({ isEditing, onClick }: { isEditing: boolean; onClick: () => void }) => {
  return (
    <div onClick={onClick}>
      <span className='flex bg-GRAY_100 rounded-md p-2 cursor-pointer h-fit w-fit hover:bg-GRAY_200'>
        <SvgSpriteLoader id={isEditing ? 'check' : 'pencil-02'} height={14} width={14} color={COLORS.TEXT_PRIMARY} />
      </span>
    </div>
  );
};

const EditableConfigField = ({
  value,
  isEditing,
  onEditToggle,
  onChange,
  firstColumn,
}: {
  value: string;
  isEditing: boolean;
  onEditToggle: () => void;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  firstColumn?: boolean;
}) => {
  return (
    <div
      className={cn(firstColumn && 'border-l', 'flex items-center gap-2 border-r border-GRAY_400 p-2 overflow-hidden')}
    >
      {value && <EditConfig isEditing={isEditing} onClick={onEditToggle} />}
      {isEditing ? <Input type='text' value={value} onChange={onChange} /> : <span>{value}</span>}
    </div>
  );
};

export default EditableConfigField;
