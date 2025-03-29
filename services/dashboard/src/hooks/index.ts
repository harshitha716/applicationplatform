import { useEffect, useRef, useState } from 'react';
import { defaultFnType } from 'types/commonTypes';

const useGetCountdown = (time: number, callback?: defaultFnType) => {
  const [seconds, setSeconds] = useState(time);
  const interval = useRef<ReturnType<typeof setInterval> | null>(null);

  const getFormattedCount = (count: number) => {
    if (count > 9) return count;

    return '0' + count;
  };

  const resetCountdown = () => {
    setSeconds(time);
    startCountdown();
  };

  const stopTimer = () => {
    if (interval.current) clearInterval(interval.current);
  };

  const startCountdown = () => {
    stopTimer();
    interval.current = setInterval(() => {
      setSeconds((prev) => {
        if (prev) {
          return prev - 1;
        } else {
          stopTimer();
          if (callback) callback();

          return 0;
        }
      });
    }, 1000);
  };

  useEffect(() => {
    if (time) resetCountdown();

    return () => {
      stopTimer();
    };
  }, [time]);

  return {
    totalSeconds: getFormattedCount(seconds),
    hours: getFormattedCount(Math.floor(seconds / 3600)),
    minutes: getFormattedCount(Math.floor((seconds % 3600) / 60)),
    seconds: getFormattedCount(seconds % 60),
    startCountdown,
    resetCountdown,
  };
};

function useOnClickOutside(ref: any, handler: (event?: MouseEvent) => void, ignoreRefs?: any[]) {
  useEffect(() => {
    const listener = (event: any) => {
      event?.stopPropagation();
      //Pass a list of refs you want to ignore
      if (ignoreRefs?.length) {
        for (let i = 0; i < ignoreRefs.length; i++) {
          if (ignoreRefs[i]?.current?.contains(event.target)) {
            return;
          }
        }
      }

      // Do nothing if clicking ref's element or descendent elements
      if (!ref.current || ref.current.contains(event.target)) return;

      handler(event);
    };

    document.addEventListener('mousedown', listener);
    document.addEventListener('touchstart', listener);

    return () => {
      document.removeEventListener('mousedown', listener);
      document.removeEventListener('touchstart', listener);
    };
  }, [ref, handler]);
}

/**
 *
 * Debounce: doesn't work on using directly on functional components
 * because it creates new instance on each re-render.
 * Use with useCallback
 *
 * TODO : Create hook
 */
function debounce<T extends (...args: any[]) => any>(func: T, wait: number) {
  let timeout: NodeJS.Timeout | null;

  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const context = this;

    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(() => {
      timeout = null;
      func.apply(context, args);
    }, wait);
  };
}

export { debounce, useGetCountdown, useOnClickOutside };
