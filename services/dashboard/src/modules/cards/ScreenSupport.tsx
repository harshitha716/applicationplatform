import { SCREEN_SUPPORT, ZAMP_LOGO } from 'constants/icons';
import Image from 'next/image';

const ScreenSupport = () => {
  return (
    <div className=' flex p-6 fixed w-screen h-screen bg-white z-1000 items-center justify-center'>
      <Image
        width={115}
        height={28}
        alt={'zamp logo'}
        className='absolute top-8 left-8'
        src={ZAMP_LOGO}
        draggable='false'
        priority
      />
      <div className='flex flex-col items-center text-center '>
        <Image width={260} height={112} alt='small screen banner' src={SCREEN_SUPPORT} draggable='false' priority />
        <div className='mt-8 text-GRAY_600'>We are live on desktop only. See you there!</div>
      </div>
    </div>
  );
};

export default ScreenSupport;
