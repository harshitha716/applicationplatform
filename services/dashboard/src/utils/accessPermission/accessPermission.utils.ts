import { store } from 'store';
import { UserRoleIdType } from 'types/api/auth.types';
import { PERMISSION_ROLES } from 'utils/accessPermission/accessPermission.types';

// Get the user privilege from store
// @params none
// @returns string
export const getUserPrivilege = () => {
  const userRole = store.getState()?.user?.roles?.find((role) => role.id === UserRoleIdType.USER)?.name ?? '';

  return userRole;
};

// Get the user email from store
// @params none
// @returns string
export const getUserEmail = () => store.getState()?.user?.user?.user_email ?? '';

// Checks if current user is same as userEmail passed
// @params string
// @returns boolean
export const checkIfCurrentUser = (userEmail: string) => (userEmail === '' ? false : getUserEmail() === userEmail);

// Checks if current user is a member in org
// @params none
// @returns boolean
export const checkIfCurrentUserIsMember = () => {
  const userRole = store.getState()?.user?.roles?.find((role) => role.id === UserRoleIdType.USER)?.name;

  return userRole === PERMISSION_ROLES.MEMBER;
};

/**
 * Get the user id from store
 * @returns string
 */
export const getUserId = () => store.getState()?.user?.user?.user_id ?? '';
