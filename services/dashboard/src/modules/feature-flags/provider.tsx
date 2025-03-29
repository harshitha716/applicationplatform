import { useSelector } from 'react-redux';
import { ENVIRONMENT } from 'constants/common.constants';
import { LAUNCH_DARKLY_CLIENT_SIDE_ID } from 'constants/featureFlags';
import { LDProvider } from 'launchdarkly-react-client-sdk';
import { RootState } from 'store';

type Props = {
  children: React.ReactNode;
};

export const FeatureFlagsProvider = ({ children }: Props) => {
  const user = useSelector((state: RootState) => state.user.user);

  if (ENVIRONMENT === 'local') return children;

  return (
    <LDProvider
      clientSideID={LAUNCH_DARKLY_CLIENT_SIDE_ID}
      context={{
        kind: 'user',
        key: user?.user_id || '',
        email: user?.user_email || '',
        organizationIds: user?.orgs?.map((org) => org.organization_id) || [],
      }}
    >
      {children}
    </LDProvider>
  );
};
