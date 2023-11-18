import { Image, Popup } from 'semantic-ui-react';
import { User } from '../types';
import { getUserRoleIcon } from '../utils/images';

type UserSmallProps = {
  user: User;
  skipText?: boolean;
  popup?: boolean;
};

/**
 * A user banner with a mini avatar and username
 * @param user User to render
 * @param skipText Whether to skip rendering the username
 */
export function UserMini({ user, skipText }: UserSmallProps) {
  const roleImage = getUserRoleIcon(user.roles);
  return (
    <div className='user-banner user-small'>
      <Image alt='User avatar' src={user.avatarUrl} avatar size='mini'/>
      {roleImage && (
        <Image className='user-role-icon absolutely-tiny' title={roleImage.alt} alt={roleImage.alt} src={roleImage.src} />
      )}
      {!skipText && <span>{user.username}</span>}
    </div>
  );
}

/**
 * A user banner with a small avatar and username
 * @param user User to render
 * @param skipText Whether to skip rendering the username
 */
export function UserSmall({ user, skipText }: UserSmallProps) {
  const roleImage = getUserRoleIcon(user.roles);
  return (
    <div className='user-banner user-large'>
      <Image alt='User avatar' src={user.avatarUrl} avatar size='tiny'/>
      {roleImage && (
        <Image className='user-role-icon' size='mini' title={roleImage.alt} alt={roleImage.alt} src={roleImage.src} />
      )}
      {!skipText && <span>{user.username}</span>}
    </div>
  );
}
