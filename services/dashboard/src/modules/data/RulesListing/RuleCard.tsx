import { FC, useState } from 'react';
import { DATE_FORMATS } from 'constants/date.constants';
import { ICON_SPRITE_TYPES, RULE_ICON } from 'constants/icons';
import { format } from 'date-fns';
import RuleStatement from 'modules/data/RulesListing/RuleStatement';
import Image from 'next/image';
import { RuleFilters } from 'types/api/dataset.types';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { cn } from 'utils/common';
import { Button } from 'components/common/button/Button';
import TagChip from 'components/common/table/CustomCellEditors/CustomTagEditor/TagChip';
import { getFilterStatementValues, getTagLabel } from 'components/filter/filter.utils';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export type RuleCardProps = {
  filters: RuleFilters;
  value?: string;
  createdOn?: string;
  defaultExpanded?: boolean;
  className?: string;
  id?: string;
  priority?: number;
  onExpand?: (id: string) => void;
  onCollapse?: (id: string) => void;
};

const RuleCard: FC<RuleCardProps> = ({
  filters,
  value,
  createdOn,
  defaultExpanded = false,
  className,
  onExpand,
  onCollapse,
  id,
}) => {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);
  const filterStatement = getFilterStatementValues(filters);
  const nonExpandedFilterStatement = filterStatement?.slice(0, 1)?.[0];
  const filterStatementLength = filterStatement?.length;

  const handleClickMore = () => {
    setIsExpanded(true);
    onExpand?.(id ?? '');
  };

  const handleClickCollapse = () => {
    setIsExpanded(false);
    onCollapse?.(id ?? '');
  };

  return (
    <div
      className={cn(
        'rounded-md border transition-all duration-500 bg-white',
        isExpanded ? 'border-GRAY_600 h-[270px]' : 'border-GRAY_400 h-[118px]',
        className,
      )}
    >
      <div className='flex items-center justify-between bg-BG_GRAY_2 py-2 px-2.5 rounded-t-md'>
        {!!createdOn && (
          <span className='f-12-400 text-GRAY_700'>
            Created on {format(new Date(createdOn), DATE_FORMATS.ddMMMyyyy)}
          </span>
        )}
        <SvgSpriteLoader
          id='dots-horizontal'
          iconCategory={ICON_SPRITE_TYPES.GENERAL}
          width={14}
          height={14}
          className='hidden'
        />
      </div>
      <div className='px-2.5 py-3 space-y-3'>
        {!!value && <TagChip item={`${getTagLabel(value)}`} />}
        <div
          className={cn('flex f-11-400 gap-1.5', {
            'h-[150px] overflow-auto flex-wrap items-start content-start': isExpanded,
            'items-center': !isExpanded,
          })}
          style={{ scrollbarWidth: 'none' }}
        >
          <Image src={RULE_ICON} alt='rule' width={12} height={13} className='max-h-[22px]' />

          {!isExpanded && (
            <>
              <div className='flex items-center gap-1.5 f-11-400 w-[300px] overflow-hidden'>
                <span className='text-GRAY_1000 pl-1.5 pr-2 py-1'>If</span>
                <RuleStatement
                  index={0}
                  filterStatement={nonExpandedFilterStatement}
                  numberOfFilters={filterStatementLength}
                />
              </div>
              <span className='text-GRAY_1000 f-11-400'>...</span>
              <Button
                id='expand-rule'
                type={BUTTON_TYPES.TEXT_NAV}
                onClick={handleClickMore}
                childrenClassName='text-GRAY_1000 underline'
                size={SIZE_TYPES.SMALL}
              >
                +more
              </Button>
            </>
          )}
          {isExpanded && (
            <>
              <span className='text-GRAY_1000 pl-1.5 pr-2 py-1 h-fit'>If</span>
              {filterStatement?.map((value, index) => (
                <RuleStatement
                  index={index}
                  filterStatement={value}
                  numberOfFilters={filterStatementLength}
                  key={`filter-statement-${index}`}
                />
              ))}
              <span className='text-GRAY_1000 pl-1.5 pr-2 py-1 h-fit'>then apply tag</span>
              <span className='border-BORDER_GRAY_400 border bg-white rounded-md pl-1.5 pr-2 py-1 text-nowrap h-fit'>
                {value}
              </span>
            </>
          )}
        </div>
      </div>
      {isExpanded && (
        <Button
          id='collapse-rule'
          type={BUTTON_TYPES.TEXT_NAV}
          onClick={handleClickCollapse}
          childrenClassName='text-GRAY_1000'
          size={SIZE_TYPES.SMALL}
        >
          Show less
        </Button>
      )}
    </div>
  );
};

export default RuleCard;
