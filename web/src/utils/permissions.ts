import { DiscordRole, DiscordRoleIds } from '../types';

export enum UserPerm {
  STAFF = 'staff',
  MODERATE = 'moderate',
  CREATE_NEWS_POST = 'create:news_post',
  TEST = 'test',
}

const DEFAULT_PERMS: UserPerm[] = [];

const MOD_ROLES = [
  DiscordRoleIds.Administrator,
  DiscordRoleIds.Moderator,
];

const STAFF_ROLES = [
  DiscordRoleIds.Administrator,
  DiscordRoleIds.Moderator,
  DiscordRoleIds.Archivist,
  DiscordRoleIds.Developer,
  DiscordRoleIds.Mechanic,
  DiscordRoleIds.Hacker,
  DiscordRoleIds.Editor,
  DiscordRoleIds.Tester,
  DiscordRoleIds.Curator,
];

// Client side check for valid permissions to control visibility of UI elements
export function getUserPermissions(roleIds: string[]): UserPerm[] {
  if (!roleIds) {
    return DEFAULT_PERMS;
  }

  const perms: UserPerm[] = [...DEFAULT_PERMS];
  if (containsAnyElement(roleIds, MOD_ROLES)) {
    perms.push(UserPerm.CREATE_NEWS_POST);
    perms.push(UserPerm.MODERATE);
  }
  if (containsAnyElement(roleIds, STAFF_ROLES)) {
    perms.push(UserPerm.STAFF);
  }

  return perms;
}

export function loadEncodedRoles(roles: string): string[] {
  if (!roles) {
    return [];
  }

  return roles.split(',');
}

export function decodeBase64(encodedStr: string): string {
  try {
    return atob(encodedStr);
  } catch (e) {
    return '';
  }
}

function containsAnyElement<T>(arrA: Array<T>, arrB: Array<T>): boolean {
  for (let i = 0; i < arrA.length; i++) {
    if (arrB.includes(arrA[i])) {
      return true;
    }
  }
  return false;
}

export function convertRoleIdsToDiscordRoles(roleIds: string[], roles: DiscordRole[]): DiscordRole[] {
  const discordRoles: DiscordRole[] = [];
  for (const r of roles) {
    if (roleIds.includes(r.id)) {
      discordRoles.push(r);
    }
  }
  return discordRoles;
}

export function throwOn404(res: Response): Response {
  if (res.status === 404) {
    throw new Error('Not Found');
  }
  return res;
}

export function getRoleName(roleId: string): string | undefined {
  // Find the key in the enum that matches the given role ID
  const roleName = (Object.keys(DiscordRoleIds) as Array<keyof typeof DiscordRoleIds>).find(key => DiscordRoleIds[key] === roleId);

  return roleName;
}
