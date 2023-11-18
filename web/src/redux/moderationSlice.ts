import { PayloadAction, createSlice } from '@reduxjs/toolkit';
import { ContentReport, FilterContentReports } from '../types';

export type ModerationState = {
  query: FilterContentReports;
  results: ContentReport[];
  totalResults: number;
  searching: boolean;
  initialLoadState: number;
  openReport: ContentReport | null;
};

const initialState: ModerationState = {
  query: {
    page: 1,
    pageSize: 10,
    order: 'updated_at',
    orderReverse: true,
    content: '',
    reportState: '',
  },
  results: [],
  totalResults: -1,
  searching: true,
  initialLoadState: 0,
  openReport: null,
};

export type ContentReportsResultsAction = {
  results: ContentReport[];
  totalResults?: number;
};

export const moderationSlice = createSlice({
  name: 'moderation',
  initialState,
  reducers: {
    setContentReportsQuery: (state: ModerationState, action: PayloadAction<FilterContentReports>) => {
      return {
        ...state,
        searching: true,
        query: action.payload
      };
    },
    setContentReportsResultsPage: (state: ModerationState, action: PayloadAction<number>) => {
      return {
        ...state,
        query: {
          ...state.query,
          page: action.payload
        }
      };
    },
    setContentReportsResults: (state: ModerationState, action: PayloadAction<ContentReportsResultsAction>) => {
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
    forceContentReportsLoad: (state: ModerationState) => {
      return {
        ...state,
        initialLoadState: Math.min(state.initialLoadState + 1, 2)
      };
    },
    setOpenReport: (state: ModerationState, action: PayloadAction<ContentReport | null>) => {
      return {
        ...state,
        openReport: action.payload
      };
    }
  }
});

export const { setContentReportsQuery, setContentReportsResultsPage, setContentReportsResults, forceContentReportsLoad, setOpenReport } = moderationSlice.actions;
export default moderationSlice.reducer;
