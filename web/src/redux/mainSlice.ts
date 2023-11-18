import { PayloadAction, createSlice } from '@reduxjs/toolkit';

type ListenerFunction = (adult: boolean) => void;

export type MainState = {
  loading: boolean;
  adult: boolean | null;
  ageVerificationOpen: boolean;
  ageVerificationListeners: ListenerFunction[];
};

const initialState: MainState = {
  loading: true,
  adult: null,
  ageVerificationOpen: false,
  ageVerificationListeners: []
};

export const mainSlice = createSlice({
  name: 'main',
  initialState,
  reducers: {
    finishedLoading: (state: MainState) => {
      return {
        ...state,
        loading: false
      };
    },
    setAdult: (state: MainState, action: PayloadAction<boolean>) => {
      return {
        ...state,
        adult: action.payload
      };
    },
    openAgeVerification: (state: MainState, action: PayloadAction<ListenerFunction>) => {
      return {
        ...state,
        ageVerificationOpen: true,
        ageVerificationListeners: [...state.ageVerificationListeners, action.payload]
      };
    },
    closeAgeVerification: (state: MainState, action: PayloadAction<boolean>) => {
      action.payload ?
        localStorage.setItem('adult', 'true') :
        localStorage.setItem('adult', 'false');
      for (const listener of state.ageVerificationListeners) {
        listener(action.payload);
      }
      return {
        ...state,
        adult: action.payload,
        ageVerificationOpen: false,
        ageVerificationListeners: []
      };
    },
  }
});


export const { setAdult, finishedLoading, openAgeVerification, closeAgeVerification} = mainSlice.actions;
export default mainSlice.reducer;
