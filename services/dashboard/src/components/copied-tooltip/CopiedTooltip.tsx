import { CSSProperties, FC } from 'react';
import { cn } from 'utils/common';
import { Button } from 'components/common/button/Button';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

type CopiedTooltipProps = {
  style?: CSSProperties;
  show: boolean;
  className?: string;
  buttonId?: string;
  wrapperOverrideClassName?: string;
};

const CopiedTooltip: FC<CopiedTooltipProps> = ({
  show = false,
  className = '',
  buttonId = 'COPY_CONTENT_BUTTON',
  wrapperOverrideClassName = '!rounded-[5px] !py-1 !px-6 !h-6 top-14',
}) =>
  show && (
    <Button className={cn('absolute flex', wrapperOverrideClassName, className)} id={buttonId}>
      <div className='flex f-12-300'>
        <SvgSpriteLoader id='check' className='mr-1 min-w-[15px]' width={15} height={15} />
        Copied!
      </div>
    </Button>
  );

export default CopiedTooltip;
