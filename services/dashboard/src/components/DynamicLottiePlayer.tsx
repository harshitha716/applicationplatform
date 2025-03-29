import React, { FC } from 'react';
import dynamic from 'next/dynamic';
import { MapAny } from 'types/commonTypes';

const Player = dynamic(() => import('@lottiefiles/react-lottie-player').then((mod) => mod.Player), {
  ssr: false,
});

interface DynamicLottiePlayerPropsType {
  src: MapAny;
  loop?: boolean;
  autoplay?: boolean;
  style?: React.CSSProperties;
  keepLastFrame?: boolean;
  className?: string;
}

const DynamicLottiePlayer: FC<DynamicLottiePlayerPropsType> = ({
  src,
  autoplay,
  style = {},
  keepLastFrame,
  className = '',
  loop = false,
}) => {
  return src ? (
    <Player
      src={src}
      style={style}
      className={className}
      autoplay={autoplay}
      keepLastFrame={keepLastFrame}
      loop={loop}
    />
  ) : null;
};

export default DynamicLottiePlayer;
