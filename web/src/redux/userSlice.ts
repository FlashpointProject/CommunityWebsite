import { PayloadAction, createSlice } from '@reduxjs/toolkit';
import { User } from '../types';

export type UserState = {
  user: User | null;
};

const initialState: UserState = {
  user: null,
};

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUser: (state: UserState, action: PayloadAction<User>) => {
      return {
        ...state,
        user: action.payload
      };
    },
    logout: (state: UserState) => {
      return {
        ...state,
        user: null
      };
    }
  }
});


export const { setUser, logout } = userSlice.actions;
export default userSlice.reducer;
