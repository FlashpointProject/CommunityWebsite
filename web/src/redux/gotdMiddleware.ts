import { isAnyOf } from '@reduxjs/toolkit';
import { startAppListening } from '../listenerMiddleware';
import { ResponseGotdSuggestions } from '../resTypes';
import { mapRawGotdSuggestion } from '../utils/mappers';
import { forceGotdSuggestionsLoad, setGotdSuggestionsQuery, setGotdSuggestionsResults, setGotdSuggestionsResultsPage } from './gotdSlice';
import { setContentReportsResults } from './moderationSlice';

export function addGotdSuggestionsSearchMiddleware() {
  startAppListening({
    matcher: isAnyOf(setGotdSuggestionsQuery, setGotdSuggestionsResultsPage, forceGotdSuggestionsLoad),
    effect: async (action, listenerApi) => {
      const state = listenerApi.getState();
      if (action.type === forceGotdSuggestionsLoad.type && state.gotdState.initialLoadState === 2) {
        return; // Another action already used the force load
      }
      const query = state.gotdState.query;

      listenerApi.dispatch(setContentReportsResults({ results: [], totalResults: -1 })); // Clear results while loading
      // Fetch gotd suggestions from the server based on filter options and pagination
      const url = new URL('/api/gotd/suggestions', window.location.origin);
      url.searchParams.append('page', query.page.toString());
      url.searchParams.append('page_size', query.pageSize.toString());
      url.searchParams.append('order_by', query.order);
      url.searchParams.append('order_direction', query.orderReverse ? 'desc' : 'asc');
      url.searchParams.append('include_total', 'true');
      fetch(url) // Add filter parameters
      .then((response) => {
        if (!response.ok) {
          throw new Error(response.status + ' ' + response.statusText);
        }
        return response.json();
      })
      .then((data: ResponseGotdSuggestions) => {
        listenerApi.dispatch(setGotdSuggestionsResults({ results: data.suggestions.map(mapRawGotdSuggestion), totalResults: data.total }));
        // Update the total number of pages based on the response
      }).catch((error) => {
        alert('Failed to fetch gotd suggestions - ' + error);
      });
    }
  });
}
