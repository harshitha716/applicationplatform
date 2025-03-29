import { FEATURE_FLAGS } from 'constants/featureFlags';
import { useLDClient } from 'launchdarkly-react-client-sdk';

export const useFeatureFlags = () => {
  const ldClient = useLDClient();

  return {
    evaluate: (flag: FEATURE_FLAGS) => ldClient?.variation(flag, false),
  };
};
