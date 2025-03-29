import { PERMISSION_ROLES, PERMISSION_TYPES } from 'utils/accessPermission/accessPermission.types';

export const PERMISSION_MESSAGES = {
  [PERMISSION_TYPES.ROLE_CHANGE]: 'You do not have permission to change role',
  [PERMISSION_TYPES.DELETE]: 'You do not have permission to delete',
  [PERMISSION_TYPES.INVITE]: 'You do not have permission to invite',
};

export enum VALIDATION_ERROR_MESSAGES {
  INVALID_EMAIL = 'Invalid email address',
  DUPLICATE_EMAIL = 'Duplicate email address',
  USER_ALREADY_IN_ORG = 'This user is already part of the organization.',
  USER_NOT_IN_ORG = 'This user is not part of the organization.',
  ORG_ALREADY_HAS_ACCESS = 'Organization already has access.',
  USER_ALREADY_HAS_ACCESS = 'This user already has access.',
  USER_ALREADY_INVITED = 'This user is already invited.',
  CANNOT_ADD_SELF = 'You cannot add yourself.',
  CANNOT_INVITE_SELF = 'You cannot invite yourself.',
}

export const acceptableRolesForAdminPurpose = [
  PERMISSION_ROLES.ADMIN,
  PERMISSION_ROLES.SYSTEM_ADMIN,
  PERMISSION_ROLES.MEMBER,
];
