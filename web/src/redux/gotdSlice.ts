import { PayloadAction, createSlice } from '@reduxjs/toolkit';
import { FilterGotdSuggestions, GotdSuggestion } from '../types';

export type GotdState = {
  query: FilterGotdSuggestions;
  results: GotdSuggestion[];
  totalResults: number;
  searching: boolean;
  initialLoadState: number;
};

const initialState: GotdState = {
  query: {
    page: 1,
    pageSize: 10,
    order: 'created_at',
    orderReverse: true,
  },
  results: [],
  totalResults: -1,
  searching: true,
  initialLoadState: 0,
};

export type GotdResultsAction = {
  results: GotdSuggestion[];
  totalResults?: number;
};

export const gotdSlice = createSlice({
  name: 'gotd',
  initialState,
  reducers: {
    setGotdSuggestionsQuery: (state: GotdState, action: PayloadAction<FilterGotdSuggestions>) => {
      return {
        ...state,
        searching: true,
        query: action.payload
      };
    },
    setGotdSuggestionsResultsPage: (state: GotdState, action: PayloadAction<number>) => {
      return {
        ...state,
        query: {
          ...state.query,
          page: action.payload
        }
      };
    },
    setGotdSuggestionsResults: (state: GotdState, action: PayloadAction<GotdResultsAction>) => {
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
    forceGotdSuggestionsLoad: (state: GotdState) => {
      return {
        ...state,
        initialLoadState: Math.min(state.initialLoadState + 1, 2)
      };
    },
  }
});

export const { setGotdSuggestionsQuery, setGotdSuggestionsResults, setGotdSuggestionsResultsPage, forceGotdSuggestionsLoad } = gotdSlice.actions;
export default gotdSlice.reducer;
