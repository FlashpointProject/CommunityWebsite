import { useSelector } from 'react-redux';
import { LoginPrompt } from './LoginPrompt';
import { Loader, Segment } from 'semantic-ui-react';
import { RootState } from '../../store';
import { UserSmall } from '../UserBanners';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { User } from '../../types';
import { mapResponseUserProfile } from '../../resTypes';
import { throwOn404 } from '../../utils/permissions';

type ProfilePageProps = {
  preview?: boolean;
};

export function ProfilePage({ preview }: ProfilePageProps) {
  const userState = useSelector((state: RootState) => state.userState);
  const { id } = useParams();
  const [fetching, setFetching] = useState(true);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    if (id) {
      fetch(`/api/profile/${id}`, {
        headers: {
          'Content-Type': 'application/json',
        },
      })
      .then(throwOn404)
      .then((response) => response.json())
      .then((data) => {
        const d = mapResponseUserProfile(data);
        const newUser: User = {
          id: d.uid,
          authed: false,
          username: d.username,
          avatarUrl: d.avatar_url,
          roles: d.roles,
          perms: [],
        };
        setUser(newUser);
        setFetching(false);
      })
      .catch(() => {
        setFetching(false);
      });
    } else {
      setUser(userState.user);
      setFetching(false);
    }
  }, [id]);

  return (
    <div>
      {fetching ? (
        <Loader />
      ) : user === null ? (
        <LoginPrompt />
      ) : (
        <>
          <h3>Profile</h3>
          <Segment>
            <UserSmall user={user} />
          </Segment>
        </>
      )}
    </div>
  );
}
