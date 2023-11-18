import { ReactNode, useMemo } from 'react';
import { UserPerm } from '../utils/permissions';
import { useSelector } from 'react-redux';
import { RootState } from '../store';

export type PermissionLockProps = {
  children: ReactNode;
  perm: UserPerm[] | UserPerm;
  any?: boolean;
  fallback?: ReactNode;
};

export type LoginLockProps = {
  children: ReactNode;
  fallback?: ReactNode;
};

/**
 * Renders child elements if the user has the required permissions
 * @param children Child elements to render if the user has the required permissions
 * @param perm Permission(s) required to view the child elements
 * @param fallback (Optional) Fallback component to render if the user does not have the required permissions
 * @returns
 */
export function PermissionLock({ children, perm, fallback, any }: PermissionLockProps) {
  const { user } = useSelector((state: RootState) => state.userState);

  const valid = useMemo(() => {
    if (!user) return false;

    if (typeof perm === 'string') {
      return user.perms.includes(perm);
    } else {
      if (any) {
        return perm.some((perm) => user.perms.includes(perm));
      } else {
        return perm.every((perm) => user.perms.includes(perm));
      }
    }
  }, [user]);

  return valid ? (
    children
  ) : (
    fallback || <></>
  );
}

/**
 * Renders child elements if the user is logged in
 * @param children Child elements to render if the user is logged in
 * @param fallback (Optional) Fallback component to render if the user is not logged in
 * @returns
 */
export function LoginLock({ children, fallback }: LoginLockProps) {
  const { user } = useSelector((state: RootState) => state.userState);

  return !!user ? (
    children
  ) : (
    fallback || <></>
  );
}
export function PermissionLockDenied() {
  return (
    <h1>You are not authorized to view this content</h1>
  );
}
