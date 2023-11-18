import { ChangeEvent, useCallback, useEffect, useMemo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Button, Dimmer, Input, Loader, PaginationProps, Tab } from 'semantic-ui-react';
import { forcePlaylistSearchLoad, setPlaylistSearchPage, setPlaylistSearchQuery } from '../../redux/playlistSearchSlice';
import { RootState } from '../../store';
import { FilterPlaylists } from '../../types';
import { BoxPlaylistInfo } from '../BoxPlaylistInfo';
import { FilterBarPlaylists } from '../FilterBarPlaylists';
import { PaginationBar } from '../PaginationBar';
import { useSmartNavigate } from '../../hooks/useSmartNavigate';

export function PlaylistsPage() {
  const panes = [{
    menuItem: 'Browse',
    render: () => <Tab.Pane className='undo-segment'>
      <PlaylistsBrowse />
    </Tab.Pane>,
  }, {
    menuItem: 'Submit',
    render: () => <Tab.Pane>
      <PlaylistsSubmit />
    </Tab.Pane>,
  }];

  return (
    <>
      <h3>Playlists</h3>
      <Tab menu={{ pointing: true }}  panes={panes}/>
    </>
  );
}

function PlaylistsBrowse() {
  const { query, results, totalResults, searching } = useSelector((state: RootState) => state.searchPlaylistsState);
  const dispatch = useDispatch();

  const onQueryChange = useCallback((query: FilterPlaylists) => {
    dispatch(setPlaylistSearchQuery(query));
  }, []);

  const onPageChange = useCallback((event: React.MouseEvent<HTMLAnchorElement, MouseEvent>, data: PaginationProps) => {
    dispatch(setPlaylistSearchPage(data.activePage as number));
  }, [query]);

  const headerStyle = results.length > 0 ? {
    marginBottom: '0',
  } : {
    marginBottom: '1rem',
  };

  const totalPages = useMemo(() => {
    return Math.ceil(totalResults / query.pageSize);
  }, [totalResults]);

  useEffect(() => {
    dispatch(forcePlaylistSearchLoad());
  }, []);

  return (
    <>
      <FilterBarPlaylists filter={query} onChange={onQueryChange} style={headerStyle}/>
      { totalPages <= 0 ? searching ? (
        <Loader active inline='centered'>Searching</Loader>
      ) : (
        <h2>No Playlists Found</h2>
      ) : (
        <>
          {results.map((playlist) => <BoxPlaylistInfo playlist={playlist}/>)}
          <PaginationBar query={query} totalResults={totalResults} onPageChange={onPageChange} />
        </>
      )}
    </>
  );
}

function PlaylistsSubmit() {
  const [file, setFile] = useState<File | null>(null);
  const [validFile, setValidFile] = useState(false);
  const [uploading, setUploading] = useState(false);
  const navigate = useSmartNavigate();

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files[0];
    if (file && file.type === 'application/json') {
      setFile(file);
      setValidFile(false);
      // Validate file
      const reader = new FileReader();
      reader.readAsText(file, 'UTF-8');
      reader.onload = (event) => {
        try {
          const data = JSON.parse(event.target.result as string);
          if (data.title && data.library && data.games) {
            setValidFile(true);
          } else {
            alert('Invalid playlist file structure. Missing one of the following: title, library, games');
          }
        } catch (e) {
          alert('Failed to load as JSON file.');
        }
      };
    } else {
      alert('Please select a JSON file.');
    }
  };

  const handleUpload = useCallback(() => {
    if (file && validFile) {
      setUploading(true);
      // Upload json text in body
      const reader = new FileReader();
      reader.readAsText(file, 'UTF-8');
      reader.onload = (event) => {
        fetch('/api/playlists', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: event.target?.result as string,
        }).then((response) => {
          if (response.status === 200) {
            // Read playlist from data
            response.json().then((data) => {
              // Navigate to uploaded playlist
              navigate(`/playlist/${data.id}`);
            });
          } else {
            alert('Failed to upload playlist.');
          }
        })
        .finally(() => setUploading(false));
      };
    }
  }, [file, validFile]);

  return (
    <>
      <h4>Upload a playlist</h4>
      { uploading && (
        <Dimmer active>
          <Loader>Uploading...</Loader>
        </Dimmer>
      )}
      <h4>Upload a playlist</h4>
      <Input type="file" onChange={handleFileChange} accept=".json" />
      <Button disabled={!validFile} onClick={handleUpload}>Upload</Button>
    </>
  );
}
