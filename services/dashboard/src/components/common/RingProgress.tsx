import React, { FC, useEffect, useRef, useState } from 'react';

interface ProgressBarProps {
  size?: number;
  progress?: number;
  trackWidth?: number;
  trackColor?: string;
  indicatorWidth?: number;
  indicatorColor?: string;
  spinnerMode?: boolean;
  spinnerSpeed?: number;
  className?: string;
}

const ProgressBar: FC<ProgressBarProps> = ({
  size = 48,
  progress = 0,
  trackWidth = 2,
  trackColor = `#E9E9E0`,
  indicatorWidth = 2,
  indicatorColor = `#111619`,
  className = '',
}) => {
  const timerId = useRef<ReturnType<typeof setInterval>>();
  const [currentProgress, setCurrentProgress] = useState(0);

  const center = size / 2,
    radius = center - (trackWidth > indicatorWidth ? trackWidth : indicatorWidth),
    dashArray = 2 * Math.PI * radius,
    dashOffset = dashArray * ((100 - currentProgress) / 100);

  useEffect(() => {
    timerId.current = setInterval(() => {
      setCurrentProgress((prev) => {
        if (prev === progress) {
          clearInterval(timerId.current);

          return prev;
        } else {
          if (progress >= currentProgress) return prev + 1;
          else return prev - 1;
        }
      });
    }, 20);

    return () => {
      clearInterval(timerId.current);
    };
  }, [progress]);

  return (
    <div className={`${className}`} style={{ width: size, height: size }}>
      <svg className='-rotate-90' style={{ width: size, height: size }}>
        <circle cx={center} cy={center} fill='transparent' r={radius} stroke={trackColor} strokeWidth={trackWidth} />
        <circle
          cx={center}
          cy={center}
          fill='transparent'
          r={radius}
          stroke={indicatorColor}
          strokeWidth={indicatorWidth}
          strokeDasharray={dashArray}
          strokeDashoffset={dashOffset}
        />
      </svg>
    </div>
  );
};

export default ProgressBar;
