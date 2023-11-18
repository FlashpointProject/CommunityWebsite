import { configureStore } from '@reduxjs/toolkit';
import mainReducer from './redux/mainSlice';
import rolesReducer from './redux/rolesSlice';
import userReducer from './redux/userSlice';
import playlistsReducer from './redux/playlistSearchSlice';
import moderationReducer from './redux/moderationSlice';
import gotdReducer from './redux/gotdSlice';
import { listenerMiddleware } from './listenerMiddleware';
import { addPlaylistMiddleware } from './redux/playlistSearchMiddleware';
import { addContentReportSearchMiddleware } from './redux/moderationMiddleware';
import { addGotdSuggestionsSearchMiddleware } from './redux/gotdMiddleware';

addPlaylistMiddleware();
addContentReportSearchMiddleware();
addGotdSuggestionsSearchMiddleware();

export const store = configureStore({
  reducer: {
    userState: userReducer,
    rolesState: rolesReducer,
    searchPlaylistsState: playlistsReducer,
    mainState: mainReducer,
    moderationState: moderationReducer,
    gotdState: gotdReducer,
  },
  middleware: (getDefaultMiddleware) => {
    return getDefaultMiddleware().prepend(listenerMiddleware.middleware);
  }
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export default store;
