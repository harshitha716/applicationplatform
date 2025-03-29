import { store } from 'store';
import { UserRoleIdType } from 'types/api/auth.types';
import { acceptableRolesForAdminPurpose } from 'utils/accessPermission/accessPermission.constants';
import { PERMISSION_ROLES } from 'utils/accessPermission/accessPermission.types';

export const accessPermissionForPage = (userRole: string) => {
  const hasAccess = acceptableRolesForAdminPurpose.find((role) => role === userRole);

  if (hasAccess) return true;

  return false;
};

export const accessPermissionForDataset = (userRole: string) => {
  const hasAccess = acceptableRolesForAdminPurpose.find((role) => role === userRole);

  if (hasAccess) return true;

  return false;
};

export const accessPermissionForPeople = () => {
  const userRole = store.getState()?.user?.roles?.find((role) => role.id === UserRoleIdType.USER)?.name;
  const hasAccess = userRole && userRole === PERMISSION_ROLES.SYSTEM_ADMIN;

  if (hasAccess) return true;

  return false;
};
