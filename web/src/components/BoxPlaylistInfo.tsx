import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Button, Grid, Icon, Label, Segment, Image } from 'semantic-ui-react';
import { PlaylistInfo } from '../types';
import { UserMini } from './UserBanners';
import { useSmartNavigate } from '../hooks/useSmartNavigate';

export type BoxPlaylistInfoProps = {
  playlist: PlaylistInfo;
};

/**
 * A medium sized box for a playlist with buttons to download, view, share, and add to Flashpoint
 * @param playlist Playlist to render
 */
export function BoxPlaylistInfo({ playlist }: BoxPlaylistInfoProps) {
  const navigate = useSmartNavigate();
  const [copied, setCopied] = useState(false);

  const copyLink = () => {
    navigator.clipboard.writeText(window.location.origin + `/playlist/${playlist.id}`)
    .then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  const viewPlaylist = () => {
    navigate(`/playlist/${playlist.id}`);
  };

  return (
    <Segment>
      <div className='box-playlist-container'>
        <Grid divided className='box-playlist-grid'>
          <Grid.Row>
            <Grid.Column width={2}>
              <div className='box-playlist-info-header-container'>
                <Link to={`/playlist/${playlist.id}`}>
                  <div className='box-playlist-info-header-id'>{playlist.id}</div>
                </Link>
                <div className='box-playlist-info-header-title'>Title</div>
              </div>
            </Grid.Column>
            <Grid.Column width={14}>
              <div className='box-playlist-info-title'>{
                playlist.extreme ? (
                  <div className='box-playlist-title-container'>
                    <Image size='mini' alt='Adult Content Warning' src={'/images/icons/Extreme.png'}/>
                    <div>{playlist.name}</div>
                  </div>
                ) : playlist.name
              }</div>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column width={2}>
              <div className='box-playlist-info-header-title'>Created By</div>
            </Grid.Column>
            <Grid.Column width={14}>
              <UserMini user={playlist.author} popup/>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column width={2}>
              <div className='box-playlist-info-header-title'>Total Games</div>
            </Grid.Column>
            <Grid.Column width={14}>
              <div className='box-playlist-info-title'>{playlist.totalGames}</div>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column width={2}>
              <div className='box-playlist-info-header-title'>Description</div>
            </Grid.Column>
            <Grid.Column width={14}>
              <div>{playlist.description}</div>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column width={2}>
              <div className='box-playlist-info-header-title'>Content Warnings</div>
            </Grid.Column>
            <Grid.Column width={14}>
              { playlist.filterGroups.length > 0 ?
                playlist.filterGroups.map((filterGroup) => (
                  <Label icon='warning sign' content={filterGroup}/>
                ))
                : 'None'}
            </Grid.Column>
          </Grid.Row>
        </Grid>
        <div className='box-playlist-buttons'>
          <Button as='a' href={'flashpoint://playlist/add/?url=' + encodeURIComponent(window.location.origin + `/api/playlist/${playlist.id}/download`)}>
            <Icon name='plus' />
            Add to Flashpoint Launcher
          </Button>
          <Button as='a' href={window.location.origin + `/api/playlist/${playlist.id}/download`} download={playlist.name + '.json'}>
            <Icon name='download' />
            Download
          </Button>
          <Button onClick={viewPlaylist}>
            <Icon name='magnify' />
            View
          </Button>
          <Button onClick={copyLink}>
            <Icon name='linkify' />
            { copied ? 'Copied!' : 'Copy Link' }
          </Button>
        </div>
      </div>
    </Segment>
  );
}
