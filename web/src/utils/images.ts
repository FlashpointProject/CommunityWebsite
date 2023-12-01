import { RolesWithIcon } from '../types';
import { getRoleName } from './permissions';

export function getUserRoleIcon(roleIds: string[]): RoleIcon | undefined {
  for (const roleId of roleIds) {
    if (RolesWithIcon.includes(roleId)) {
      return {
        src: `/images/roles/${roleId}.png`,
        alt: getRoleName(roleId) + ' role'
      };
    }
  }
}

type RoleIcon = {
  src: string;
  alt: string;
};

export function getRoleIcon(roleId: string): RoleIcon | undefined {
  if (RolesWithIcon.includes(roleId)) {
    return {
      src: `/images/roles/${roleId}.png`,
      alt: getRoleName(roleId) + ' role'
    };
  }
}
