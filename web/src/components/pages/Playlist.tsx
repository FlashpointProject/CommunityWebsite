import { ReactNode, useCallback, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Button, Dimmer, Loader, Segment } from 'semantic-ui-react';
import { Playlist } from '../../types';
import { fetchPlaylist } from '../../utils/fetchers';
import { BoxPlaylistInfo } from '../BoxPlaylistInfo';
import { Comment } from '../Comment';
import { GameListPlaylist } from '../GameList';
import { SegmentWithHeader } from '../SegmentWithHeader';
import { UserMini } from '../UserBanners';
import { useSelector } from 'react-redux';
import { RootState } from '../../store';
import { UserPerm } from '../../utils/permissions';
import { useConfirm } from '../../hooks/useConfirm';
import { useSmartNavigate } from '../../hooks/useSmartNavigate';

export function PlaylistPage() {
  const { id } = useParams();
  const [playlist, setPlaylist] = useState<Playlist | null>(null);
  const [forceRefresh, setForceRefresh] = useState(false);

  // Cheeky way to force the playlist to reload
  const forceRefreshFunc = useCallback(() => {
    setForceRefresh(!forceRefresh);
  }, [forceRefresh]);

  useEffect(() => {
    // Fetch playlist from the server
    fetchPlaylist(id)
    .then((playlist) => {
      setPlaylist(playlist);
    })
    .catch(console.error);
  }, [id, forceRefresh]);

  const header = playlist ? (
    <>
      <div className="left-flex">
        {playlist.name}
      </div>
      <div className="right-flex">
        <UserMini user={playlist.author} />
      </div>
    </>
  ) : '';

  return (
    <PlaylistFullBox forceRefresh={forceRefreshFunc} playlist={playlist} header={header} />
  );
}

export type PlaylistFullBoxProps = {
  playlist: Playlist | null;
  header: ReactNode;
  forceRefresh: () => void;
};

export function PlaylistFullBox(props: PlaylistFullBoxProps) {
  const [uploading, setUploading] = useState(false);
  const { user } = useSelector((state: RootState) => state.userState);
  const { confirm, modal } = useConfirm();
  const navigate = useSmartNavigate();

  const onDeletePlaylist = useCallback(() => {
    if (props.playlist) {
      fetch(`/api/playlist/${props.playlist.id}`, {
        method: 'DELETE',
      }).then((response) => {
        if (response.status === 200) {
          // Force reload the playlist
          navigate('/playlists');
        } else {
          alert('Failed to delete playlist.');
        }
      })
      .catch(() => {
        alert('Failed to delete playlist. Server connection failed.');
      });
    }
  }, [props.playlist]);

  const onUploadNewVersion = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.item(0);
    if (file && props.playlist) {
      setUploading(true);
      // Upload file
      const reader = new FileReader();
      reader.readAsText(file, 'UTF-8');
      reader.onload = (event) => {
        try {
          const data = JSON.parse(event.target.result as string);
          if (data.title && data.library && data.games) {
            // Valid file, upload
            // Upload json text in body
            const reader = new FileReader();
            reader.readAsText(file, 'UTF-8');
            reader.onload = (event) => {
              fetch(`/api/playlist/${props.playlist.id}`, {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json',
                },
                body: event.target?.result as string,
              }).then((response) => {
                if (response.status === 200) {
                  response.json().then((data) => {
                    // Force reload the playlist
                    props.forceRefresh();
                  });
                } else {
                  alert('Failed to upload playlist.');
                }
              })
              .finally(() => setUploading(false));
            };
          } else {
            alert('Invalid playlist file structure. Missing one of the following: title, library, games');
          }
        } catch (e) {
          setUploading(false);
          alert('Failed to load as JSON file.');
        }
      };
    }
  }, [props.playlist]);

  const onClickUploadNewVersion = useCallback(() => {
    // Open a file select and immediately upload if given a file
    const fileBox = document.getElementById('new-playlist-file');
    if (fileBox) {
      fileBox.click();
    }
  }, []);

  return (
    <div>
      {uploading && (
        <Dimmer active>
          <Loader>Uploading</Loader>
        </Dimmer>
      )}
      {modal}
      <h3>Playlist</h3>
      {props.playlist === null ? (
        <Loader active inline='centered'>Loading</Loader>
      ) : (
        <SegmentWithHeader header={props.header}>
          { (user && (user.id === props.playlist.author.id || user.perms.includes(UserPerm.MODERATE))) && (
            <div className='two-column-responsive-grid'>
              { user.id === props.playlist.author.id && (
                <Segment>
                  <h4>Owner Actions</h4>
                  <div className='button-container'>
                    <Button
                      onClick={onClickUploadNewVersion}
                      positive>Upload New Version</Button>
                    <input
                      onChange={onUploadNewVersion}
                      id='new-playlist-file' type='file' name='name' style={{display: 'none'}} />
                    <Button
                      onClick={confirm('Are you sure', 'Deleting a playlist is irreversable', onDeletePlaylist)}
                      negative>Delete</Button>
                  </div>
                </Segment>
              )}
              { user.perms.includes(UserPerm.MODERATE) && (
                <Segment>
                  <h4>Moderator Actions</h4>
                  <div className='button-container'>
                    <Button
                      onClick={confirm('Are you sure', 'Deleting a playlist is irreversable', onDeletePlaylist)}
                      negative>Delete</Button>
                  </div>
                </Segment>
              )}
            </div>
          )}
          <BoxPlaylistInfo playlist={props.playlist} />
          <Segment>
            <h3>Games</h3>
            <GameListPlaylist games={props.playlist.games} />
          </Segment>
          <h3>Comments</h3>
          <div className='playlist-comments'>
            Not Implemented
          </div>
        </SegmentWithHeader>
      )}
    </div>
  );
}
