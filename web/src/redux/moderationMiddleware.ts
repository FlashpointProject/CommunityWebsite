import { isAnyOf } from '@reduxjs/toolkit';
import { startAppListening } from '../listenerMiddleware';
import { ResponseContentReports } from '../resTypes';
import { mapRawContentReport } from '../utils/mappers';
import { forceContentReportsLoad, setContentReportsQuery, setContentReportsResults, setContentReportsResultsPage } from './moderationSlice';

export function addContentReportSearchMiddleware() {
  startAppListening({
    matcher: isAnyOf(setContentReportsQuery, setContentReportsResultsPage, forceContentReportsLoad),
    effect: async (action, listenerApi) => {
      const state = listenerApi.getState();
      if (action.type === forceContentReportsLoad.type && state.moderationState.initialLoadState === 2) {
        return; // Another action already used the force load
      }
      const query = state.moderationState.query;

      listenerApi.dispatch(setContentReportsResults({ results: [], totalResults: -1 })); // Clear results while loading
      // Fetch playlists from the server based on filter options and pagination
      const url = new URL('/api/reports', window.location.origin);
      url.searchParams.append('page', query.page.toString());
      url.searchParams.append('page_size', query.pageSize.toString());
      url.searchParams.append('order_by', query.order);
      url.searchParams.append('order_direction', query.orderReverse ? 'desc' : 'asc');
      if (query.content !== '') {
        url.searchParams.append('content_type', query.content);
      }
      url.searchParams.append('include_total', 'true');
      fetch(url) // Add filter parameters
      .then((response) => {
        if (!response.ok) {
          throw new Error(response.status + ' ' + response.statusText);
        }
        return response.json();
      })
      .then((data: ResponseContentReports) => {
        listenerApi.dispatch(setContentReportsResults({ results: data.reports.map(mapRawContentReport), totalResults: data.total }));
        // Update the total number of pages based on the response
      }).catch((error) => {
        alert('Failed to fetch playlists - ' + error);
      });
    }
  });
}
