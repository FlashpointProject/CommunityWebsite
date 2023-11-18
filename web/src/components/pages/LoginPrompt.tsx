import { useCallback } from 'react';
import { Button } from 'semantic-ui-react';
import { RootState } from '../../store';
import { useSelector } from 'react-redux';
import { NotFound } from './NotFound';

export function LoginPrompt() {
  const { user } = useSelector((state: RootState) => state.userState);

  const onLogin = useCallback(() => {
    window.location.href = '/auth';
  }, []);

  return !user ? (
    <div>
      <h2>You must be logged in to view this page.</h2>
      <p><Button onClick={onLogin}>Login</Button></p>
    </div>
  ) : (
    <NotFound />
  );
}
