import { PayloadAction, createSlice } from '@reduxjs/toolkit';
import { FilterPlaylists, PlaylistInfo } from '../types';

export type PlaylistSearchState = {
  query: FilterPlaylists;
  results: PlaylistInfo[];
  totalResults: number;
  searching: boolean;
  initialLoadState: number;
};

const initialState: PlaylistSearchState = {
  query: {
    page: 1,
    pageSize: 10,
    order: 'updated_at',
    orderReverse: true,
    library: '',
    extreme: localStorage.getItem('playlistIncludeAdult') === 'true',
  },
  results: [],
  totalResults: -1,
  searching: true,
  initialLoadState: 0,
};

export type PlaylistResultsAction = {
  results: PlaylistInfo[];
  totalResults?: number;
};

export const playlistSearchSlice = createSlice({
  name: 'playlist',
  initialState,
  reducers: {
    setPlaylistSearchQuery: (state: PlaylistSearchState, action: PayloadAction<FilterPlaylists>) => {
      return {
        ...state,
        searching: true,
        query: action.payload
      };
    },
    setPlaylistSearchPage: (state: PlaylistSearchState, action: PayloadAction<number>) => {
      return {
        ...state,
        query: {
          ...state.query,
          page: action.payload
        }
      };
    },
    setPlaylistSearchResults: (state: PlaylistSearchState, action: PayloadAction<PlaylistResultsAction>) => {
      const newState = { ...state };
      if (action.payload.totalResults) {
        newState.totalResults = action.payload.totalResults;
      }
      newState.results = action.payload.results;
      if (action.payload.totalResults === -1 ) {
        newState.searching = true;
      } else {
        newState.searching = false;
      }
      return newState;
    },
    forcePlaylistSearchLoad: (state: PlaylistSearchState) => {
      return {
        ...state,
        initialLoadState: Math.min(state.initialLoadState + 1, 2)
      };
    },
  }
});

export const { setPlaylistSearchQuery, setPlaylistSearchResults, setPlaylistSearchPage, forcePlaylistSearchLoad } = playlistSearchSlice.actions;
export default playlistSearchSlice.reducer;
