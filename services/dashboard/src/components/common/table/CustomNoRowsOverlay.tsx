import SvgSpriteLoader from 'components/SvgSpriteLoader';

const CustomNoRowsOverlay = () => {
  return (
    <div
      role='presentation'
      className='flex flex-col items-center gap-2.5 h-full justify-center text-GRAY_700 f-12-450'
    >
      <SvgSpriteLoader id='coins-stacked-03' width={24} height={24} />
      <div>No data available, try again with different filters</div>
    </div>
  );
};

export default CustomNoRowsOverlay;
