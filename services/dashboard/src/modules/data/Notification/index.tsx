import { FC, useRef, useState } from 'react';
import { COLORS } from 'constants/colors';
import ProgressBar from 'components/common/RingProgress';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

type NotificationProps = {
  isPolling: boolean;
};

const Notification: FC<NotificationProps> = ({ isPolling }: NotificationProps) => {
  const [showNotificationPanel, setShowNotificationPanel] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const toggleNotificationPanel = () => {
    setShowNotificationPanel((prev) => !prev);
  };

  return (
    <>
      {isPolling ? (
        <div className='relative  cursor-pointer w-5.5 h-5.5 rounded' onClick={toggleNotificationPanel}>
          <div className='hover:bg-GRAY_100 h-full w-full rounded flex items-center justify-center'>
            <ProgressBar
              trackColor={COLORS.GRAY_400}
              indicatorColor={'#22A356'}
              indicatorWidth={2}
              trackWidth={2}
              size={16}
              className='animate-spin'
              progress={20}
            />
          </div>
          {showNotificationPanel && (
            <div
              ref={dropdownRef}
              className='p-5 absolute top-7 -right-[86px] h-[55px] f-13-500 bg-white rounded-[10px] text-GRAY_1000 f-12-450 z-1000 flex items-center w-[308px] border-0.5 border-GRAY_500 gap-3'
            >
              <ProgressBar
                trackColor={COLORS.GRAY_400}
                indicatorColor={'#22A356'}
                indicatorWidth={2}
                trackWidth={2}
                size={16}
                className='animate-spin'
                progress={20}
              />
              <div className='grow'>Tagging in progress</div>
              <SvgSpriteLoader
                id='x-close'
                width={16}
                height={16}
                onClick={(e) => {
                  e.stopPropagation();
                  setShowNotificationPanel(false);
                }}
                className='text-GRAY_800 hover:text-GRAY_1000'
              />
            </div>
          )}
        </div>
      ) : null}
    </>
  );
};

export default Notification;
