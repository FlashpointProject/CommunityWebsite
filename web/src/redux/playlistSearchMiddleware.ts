import { isAnyOf } from '@reduxjs/toolkit';
import { startAppListening } from '../listenerMiddleware';
import { PlaylistInfo, RawPlaylistInfo } from '../types';
import { forcePlaylistSearchLoad, setPlaylistSearchPage, setPlaylistSearchQuery, setPlaylistSearchResults } from './playlistSearchSlice';

export function addPlaylistMiddleware() {
  startAppListening({
    matcher: isAnyOf(setPlaylistSearchQuery),
    effect: async (action, listenerApi) => {
      const state = listenerApi.getState();
      // Save adult state to local storage
      localStorage.setItem('playlistIncludeAdult', state.searchPlaylistsState.query.extreme ? 'true' : 'false');
    }
  });

  startAppListening({
    matcher: isAnyOf(setPlaylistSearchQuery, setPlaylistSearchPage, forcePlaylistSearchLoad),
    effect: async (action, listenerApi) => {
      const state = listenerApi.getState();
      if (action.type === forcePlaylistSearchLoad.type && state.searchPlaylistsState.initialLoadState === 2) {
        return; // Another action already used the force load
      }
      const query = state.searchPlaylistsState.query;

      listenerApi.dispatch(setPlaylistSearchResults({ results: [], totalResults: -1 })); // Clear results while loading
      // Fetch playlists from the server based on filter options and pagination
      const url = new URL('/api/playlists', window.location.origin);
      url.searchParams.append('page', query.page.toString());
      url.searchParams.append('page_size', query.pageSize.toString());
      url.searchParams.append('order_by', query.order);
      url.searchParams.append('order_direction', query.orderReverse ? 'desc' : 'asc');
      if (query.library !== '') {
        url.searchParams.append('library', query.library);
      }
      if (query.extreme) {
        url.searchParams.append('extreme', 'true');
      }
      url.searchParams.append('include_total', 'true');
      fetch(url) // Add filter parameters
      .then((response) => {
        if (!response.ok) {
          throw new Error(response.status + ' ' + response.statusText);
        }
        return response.json();
      })
      .then((data) => {
        const rawPlaylists = data['playlists'] as RawPlaylistInfo[];
        // Map authors to users
        const playlists = rawPlaylists.map<PlaylistInfo>((rawPlaylist) => {
          return {
            id: rawPlaylist.id,
            name: rawPlaylist.name,
            description: rawPlaylist.description,
            totalGames: rawPlaylist.total_games,
            library: rawPlaylist.library,
            extreme: rawPlaylist.extreme,
            filterGroups: rawPlaylist.filter_groups,
            author: {
              authed: false,
              id: rawPlaylist.author.uid,
              username: rawPlaylist.author.username,
              avatarUrl: rawPlaylist.author.avatar_url,
              roles: rawPlaylist.author.roles,
              perms: [],
            }
          };
        });
        listenerApi.dispatch(setPlaylistSearchResults({ results: playlists, totalResults: data['total'] }));
        // Update the total number of pages based on the response
      }).catch((error) => {
        alert('Failed to fetch playlists - ' + error);
      });
    }
  });
}
