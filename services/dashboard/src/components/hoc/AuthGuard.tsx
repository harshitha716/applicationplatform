import React, { useEffect } from 'react';
import { useWhoAmIQuery } from 'apis/auth';
import { ALLOWED_EMAIL_DOMAINS, ENVIRONMENT } from 'constants/common.constants';
import { useAppDispatch, useAppSelector } from 'hooks/toolkit';
import OrgMembershipPending from 'modules/login/OrgMembershipPending';
import { useRouter } from 'next/router';
import { RootState } from 'store';
import { setDashboardLoader, setRoles, setUser, setWorkspace } from 'store/slices/user';
import { UserRoleIdType } from 'types/api/auth.types';
import NotAuthorized from 'components/NotAuthorized';

type Props = {
  loginRoute: string;
  children: React.ReactNode;
};

export const AuthGuard: React.FC<Props> = (props) => {
  const router = useRouter();
  const dispatch = useAppDispatch();

  const { data: session, isLoading, isError, isSuccess } = useWhoAmIQuery();
  const workspace = useAppSelector((state: RootState) => state.user.workspace);

  useEffect(() => {
    if (session && isSuccess) {
      dispatch(setUser(session));
      const defaultWorkspace = session?.organization_id;
      const user_role = session?.orgs[0]?.resource_audience_policies[0]?.privilege;

      dispatch(setRoles([{ id: UserRoleIdType.USER, name: user_role }]));

      dispatch(setWorkspace(defaultWorkspace));
    }
  }, [session, isSuccess, dispatch]);

  if (isError && router.pathname !== props.loginRoute) {
    let query = {};

    if (router.query.redirect_to) {
      query = { redirect_to: router.query.redirect_to };
    }
    router.push(props.loginRoute, {
      query,
    });

    return null;
  }

  if (isSuccess && session?.user_id && router.pathname === props.loginRoute) {
    router.push((router.query.redirect_to as string) ?? '/');

    return;
  }

  if (isLoading || (session?.user_id && workspace === null)) {
    dispatch(setDashboardLoader(true));

    return null;
  }

  if (!session) {
    if (router.pathname !== props.loginRoute) {
      return <div>Not logged in. Redirecting...</div>;
    } else {
      return props.children;
    }
  }

  if (session?.orgs?.length === 0) {
    return <OrgMembershipPending />;
  }

  if (
    session &&
    ENVIRONMENT === 'staging' &&
    ALLOWED_EMAIL_DOMAINS.every((eachDomain: string) => !session?.user_email?.endsWith(eachDomain))
  ) {
    return <NotAuthorized />;
  }

  return props.children;
};
