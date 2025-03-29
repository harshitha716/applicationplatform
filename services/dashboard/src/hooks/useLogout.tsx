import { useCallback } from 'react';
import { useInitiateLogoutFlowQuery, useLazyLogoutQuery, useLazyWhoAmIQuery } from 'apis/auth';
import { ROUTES_PATH } from 'constants/routeConfig';
import { useRouter } from 'next/router';

export const useLogout = () => {
  const router = useRouter();
  const { data: logoutFlow, refetch: refetchLogoutFlow } = useInitiateLogoutFlowQuery();
  const [logOut] = useLazyLogoutQuery();
  const [whoAmI] = useLazyWhoAmIQuery();

  const handleLogout = useCallback(async () => {
    logOut(logoutFlow?.logout_url ?? '')
      .then(() => {
        whoAmI()
          .then(() => {
            router.push(ROUTES_PATH.LOGIN);
          })
          .catch(() => {
            router.push(ROUTES_PATH.LOGIN);
          });
      })
      .catch(() => {
        refetchLogoutFlow();
      });
  }, [logoutFlow, logOut, router, refetchLogoutFlow]);

  return {
    logout: handleLogout,
  };
};
