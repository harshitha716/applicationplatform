import { useEffect, useState } from 'react';

export interface WindowDimensionsType {
  width: number;
  height: number;
}

export const useWindowDimensions = (): WindowDimensionsType => {
  const getWindowDimensions = () => {
    return {
      width: window?.innerWidth,
      height: window?.innerHeight,
    };
  };

  const [windowDimensions, setWindowDimensions] = useState({ width: 0, height: 0 });

  useEffect(() => {
    const handleResize = () => setWindowDimensions(getWindowDimensions());

    window?.addEventListener('resize', handleResize);
    handleResize();

    return () => window?.removeEventListener('resize', handleResize);
  }, []);

  return windowDimensions;
};
