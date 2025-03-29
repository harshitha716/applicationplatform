//Reference : https://levelup.gitconnected.com/polling-in-javascript-ab2d6378705a

import { useEffect, useRef } from 'react';
import { MapAny } from 'types/commonTypes';

export const POLLING_ERROR = 'Exceeded max attempts';

/**
 * Hook to facilitate polling with support for fixed interval and exponential backoff strategies
 * @returns {Object} Object containing startPolling and stopPolling functions
 */

const usePolling = () => {
  const pollingTimer = useRef<ReturnType<typeof setTimeout>>();
  const pollingFlag = useRef<boolean>();

  pollingFlag.current = true;

  const clearTimer = () => {
    if (pollingTimer) clearTimeout(pollingTimer.current);
  };

  /**
   * A higher-order function that returns a function, executePoll
   * @param {Function} fn - function that will be executed over a given interval. Typically this will be an API request
   * @param {Function} validate - function where we define a test/check to see if the data matches what we want, which will end the poll
   * @param {number} interval - time to wait between poll requests in milliseconds
   * @param {number} maxAttempts - upper bound for the number of poll requests, to prevent it from running infinitely
   * @param {boolean} [pollTillSuccess=false] - whether to continue polling until success even if errors occur
   * @param {boolean} [isExponential=false] - whether to use exponential backoff for polling intervals
   * @param {number} [backoffFactor=2] - factor by which to increase the interval when using exponential backoff
   * @param {number} [maxInterval=300000] - maximum interval time in milliseconds when using exponential backoff (default: 5 minutes)
   * @returns {Promise} Promise that resolves with the result of the polling or rejects with an error
   */
  const startPolling = ({
    fn,
    validate,
    interval = 0,
    maxAttempts = 0,
    pollTillSuccess = false,
    isExponential = false,
    backoffFactor = 2,
    maxInterval = 300000, // 5 minutes as default max interval
  }: {
    fn: any;
    validate: (res: any) => boolean;
    interval: number;
    maxAttempts: number;
    pollTillSuccess?: boolean;
    isExponential?: boolean;
    backoffFactor?: number;
    maxInterval?: number;
  }) => {
    pollingFlag.current = true;

    let attempts = 0;
    let currentInterval = interval;

    /**
     * A function that will run recursively until a stopping condition is met
     *  stopping conditions : validation is truthy or maxAttempts exceeded or pollingFlag is unset
     * @param resolve
     * @param reject
     * @returns promise
     */
    const executePoll = async (resolve: (val: MapAny) => void, reject: (val: MapAny) => void) => {
      try {
        const result = await fn()?.unwrap();

        attempts++;

        if (validate?.(result)) {
          return resolve(result);
        } else if (maxAttempts && attempts === maxAttempts) {
          return reject({ error: POLLING_ERROR });
        } else {
          clearTimer();
          if (isExponential) {
            // Calculate new interval with exponential backoff
            currentInterval = Math.min(currentInterval * backoffFactor, maxInterval);
          }
          if (pollingFlag.current) pollingTimer.current = setTimeout(executePoll, currentInterval, resolve, reject);
        }
      } catch (err) {
        if (err && pollTillSuccess) {
          clearTimer();
          if (isExponential) {
            // Calculate new interval with exponential backoff
            currentInterval = Math.min(currentInterval * backoffFactor, maxInterval);
          }
          if (pollingFlag.current) pollingTimer.current = setTimeout(executePoll, currentInterval, resolve, reject);
        } else return reject({ error: POLLING_ERROR });
      }
    };

    return new Promise(executePoll);
  };

  /**
   * Function to stop polling
   */
  const stopPolling = () => {
    clearTimer();
    pollingFlag.current = false;
  };

  useEffect(() => {
    return () => stopPolling();
  }, []);

  return { startPolling, stopPolling };
};

export default usePolling;
