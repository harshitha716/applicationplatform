import { FC, useRef, useState } from 'react';
import { COLORS } from 'constants/colors';
import { useOnClickOutside } from 'hooks';
import { OptionsType } from 'types/common/components/dropdown/dropdown.types';
import { MapAny } from 'types/commonTypes';
import { CONDITION_OPERATOR_TYPE, DESCRIPTION_OPERATORS } from 'components/filter/filters.constants';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface DescriptionOperatorsDropdownProps {
  operator?: MapAny;
  updateOperator: (operator: MapAny) => void;
  isLoading?: boolean;
  label?: string;
  operatorOptions?: OptionsType[];
}

const DescriptionOperatorsDropdown: FC<DescriptionOperatorsDropdownProps> = ({
  operator,
  updateOperator,
  isLoading,
  label = 'Description',
  operatorOptions = DESCRIPTION_OPERATORS,
}) => {
  const ref = useRef(null);
  const [isOpen, setIsOpen] = useState(false);

  const onSelect = (operator: MapAny) => {
    if (isLoading) return;

    setIsOpen(false);
    updateOperator(operator);
  };

  const onToggleDropdown = () => {
    if (isLoading) return;

    setIsOpen(!isOpen);
  };

  useOnClickOutside(ref, () => setIsOpen(false));

  return (
    <div className='flex items-center'>
      <div className=' text-GRAY_700 f-11-400 mr-1'>{label} </div>
      <div className=''>
        <div className='flex items-center cursor-pointer relative' onClick={onToggleDropdown}>
          <div className='text-BLUE_700 f-11-500 mr-1'>{operator?.label ?? CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS}</div>
          <SvgSpriteLoader id='chevron-down' width={12} height={12} color={COLORS.GRAY_700} />
          {isOpen && (
            <div
              ref={ref}
              className='p-1 z-10 absolute top-full left-0 min-w-[120px] bg-white text-GRAY_900 border border-GRAY_400 shadow-tableFilterMenu rounded-md'
            >
              {operatorOptions.map((option) => (
                <div
                  className='hover:bg-GRAY_100 f-12-500 py-2 px-2.5 rounded-md'
                  key={option.value}
                  onClick={() => onSelect(option)}
                >
                  {option.label}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default DescriptionOperatorsDropdown;
