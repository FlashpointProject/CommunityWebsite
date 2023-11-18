import { useCallback } from 'react';
import { useCookies } from 'react-cookie';
import { useSelector } from 'react-redux';
import { Button, Dropdown } from 'semantic-ui-react';
import { RootState } from '../store';
import { UserPerm } from '../utils/permissions';
import { PermissionLock } from './Locks';
import { UserMini } from './UserBanners';
import { useSmartNavigate } from '../hooks/useSmartNavigate';

/**
 * Universal website page header
 */
export function AppHeader() {
  const userState = useSelector((state: RootState) => state.userState);
  const [,, removeCookie] = useCookies();
  const navigate = useSmartNavigate();

  const onLogout = useCallback(() => {
    // Clear cookies if present
    console.log('removing login cookies');
    removeCookie('login');
    removeCookie('username');
    removeCookie('avatar_url');
    removeCookie('roles');
    removeCookie('uid');
    window.location.reload();
  }, []);

  const onLogin = useCallback(() => {
    window.location.href = '/auth';
  }, []);

  return (
    <div className="header block-header">
      <div className="header__left">
        <div className="header__logo">
          <img className='clickable' src="/images/logo.png" alt="Flashpoint Community Logo" onClick={() => navigate('/')} />
        </div>
        <div className="header__title">Flashpoint Community</div>
      </div>
      <div className='header__right'>
        <div className="header__links">
          <Button onClick={() => navigate('/')}>Home</Button>
          <Button onClick={() => navigate('/gotd')}>Game of the Day</Button>
          <Button onClick={() => navigate('/playlists')}>Playlists</Button>
          <PermissionLock perm={UserPerm.MODERATE}>
            <Button onClick={() => navigate('/moderation')}>Mod Tools</Button>
          </PermissionLock>
        </div>
        <div className="header__user">
          { userState.user !== null ? (
            <>
              <UserMini skipText user={userState.user} />
              <Dropdown text={userState.user.username} pointing className='link item'>
                <Dropdown.Menu className='right-dropdown'>
                  <Dropdown.Item onClick={() => navigate('/profile')}>My Profile</Dropdown.Item>
                  <Dropdown.Divider />
                  <Dropdown.Item onClick={onLogout}>Logout</Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            </>
          ) : (
            <Button onClick={onLogin}>Login</Button>
          )}
        </div>
      </div>
    </div>
  );
}
