import * as React from 'react';
import { ReactCookieProps } from 'react-cookie';
import { useDispatch, useSelector } from 'react-redux';
import { Route, Routes, useLocation } from 'react-router-dom';
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import { Button, Header, Icon, Modal } from 'semantic-ui-react';
import { AppHeader } from './components/Header';
import { PermissionLock, PermissionLockDenied } from './components/Locks';
import { GotdPage } from './components/pages/GameOfTheDay';
import { Home } from './components/pages/Home';
import { ContentReportsPage } from './components/pages/Moderation';
import { NotFound } from './components/pages/NotFound';
import { PlaylistPage } from './components/pages/Playlist';
import { PlaylistsPage } from './components/pages/Playlists';
import { PostPage } from './components/pages/Post';
import { PostCreator } from './components/pages/PostCreator';
import { ProfilePage } from './components/pages/Profile';
import { TestPage } from './components/pages/Test';
import { closeAgeVerification, finishedLoading, setAdult } from './redux/mainSlice';
import { setPlaylistSearchQuery } from './redux/playlistSearchSlice';
import { logout, setUser } from './redux/userSlice';
import { RootState } from './store';
import { User } from './types';
import { UserPerm, getUserPermissions } from './utils/permissions';

type AppProps = ReactCookieProps;

export function App(props: AppProps) {
  const dispatch = useDispatch();
  const { loading, ageVerificationOpen } = useSelector((state: RootState) => state.mainState);
  const playlistSearchState = useSelector((state: RootState) => state.searchPlaylistsState);
  const location = useLocation();

  React.useEffect(() => {
    // Load user info from cookies
    const idCookie = props.cookies.get('uid', { doNotParse: true });
    const usernameCookie = props.cookies.get('username');
    const avatarCookie = props.cookies.get('avatar_url');
    const rolesCookie = props.cookies.get('roles');
    const adult = localStorage.getItem('adult');
    dispatch(setAdult(adult === 'true'));

    if (usernameCookie && avatarCookie && rolesCookie) {
      // User is logged in
      const userObject: User = {
        id: idCookie,
        authed: true,
        username: usernameCookie,
        avatarUrl: avatarCookie,
        roles: rolesCookie.split(','),
        perms: [],
      };
      userObject.perms = getUserPermissions(userObject.roles);
      dispatch(setUser(userObject));
    } else {
      // User is not logged in
      dispatch(logout());
    }

    // Handle loading initial redux states if the location matches a route that has a redux store
    if (location.pathname === '/playlists') {
      // Load playlist search query
      const params = new URLSearchParams(location.search);
      const dictedParams = {};
      params.forEach((value, key) => {
        dictedParams[key] = value;
      });
      const newQuery = { ...playlistSearchState.query };
      newQuery.page = parseInt(dictedParams['page'] || playlistSearchState.query.page);
      newQuery.pageSize = parseInt(dictedParams['pageSize'] || playlistSearchState.query.pageSize);
      newQuery.order = dictedParams['order'] || playlistSearchState.query.order;
      newQuery.orderReverse = dictedParams['orderReverse'] === 'true' || playlistSearchState.query.orderReverse;
      newQuery.library = dictedParams['library'] || playlistSearchState.query.library;
      newQuery.extreme = false;
      // Respect age restriction
      if (adult === 'true') {
        newQuery.extreme = dictedParams['extreme'] === 'true' || playlistSearchState.query.extreme;
      }
      dispatch(setPlaylistSearchQuery(newQuery));
    }

    dispatch(finishedLoading());
  }, []);

  return (
    <div className="App">
      <Modal open={ageVerificationOpen} size="tiny">
        <Modal.Header>Age Verification</Modal.Header>
        <Modal.Content className='centered-modal'>
          <Header icon>
            <Icon name='warning sign' />
            You must be over 18 to access this content
          </Header>
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={() => {
            localStorage.setItem('adult', 'true');
            dispatch(closeAgeVerification(true));
          }}>Yes, I am over 18</Button>
          <Button onClick={() => {
            localStorage.setItem('adult', 'false');
            dispatch(closeAgeVerification(false));
          }}>No, I am not over 18</Button>
        </Modal.Actions>
      </Modal>
      <AppHeader />
      <PermissionLock perm={UserPerm.STAFF} fallback={<PermissionLockDenied/>}>
        <TransitionGroup>
          <CSSTransition key={location.key}
            timeout={350}
            classNames="fade"
            unmountOnExit>
            <div className="content">
              {!loading && (
                <Routes location={location}>
                  <Route path="*" element={<NotFound />} />
                  <Route path="/" element={<Home/>}/>
                  <Route path="/gotd/" element={<GotdPage/>}/>
                  <Route path="/post/:id" element={<PostPage/>}/>
                  <Route path="/playlists" element={<PlaylistsPage/>}/>
                  <Route path="/playlist/:id" element={<PlaylistPage/>}/>
                  <Route path="/test" element={<TestPage/>}/>
                  <Route path="/profile" element={<ProfilePage/>}/>
                  <Route path="/profile/:id" element={<ProfilePage preview/>}/>
                  <Route path="/post/create" element={
                    <PermissionLock perm={UserPerm.CREATE_NEWS_POST} fallback={<PermissionLockDenied/>}>
                      <PostCreator />
                    </PermissionLock>
                  }/>
                  <Route path="/moderation" element={
                    <PermissionLock perm={UserPerm.MODERATE} fallback={<PermissionLockDenied/>}>
                      <ContentReportsPage />
                    </PermissionLock>
                  }/>
                </Routes>
              )}
            </div>
          </CSSTransition>
        </TransitionGroup>
      </PermissionLock>
    </div>
  );
}
