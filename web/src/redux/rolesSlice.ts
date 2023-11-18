import { PayloadAction, createSlice } from '@reduxjs/toolkit';
import { DiscordRole } from '../types';

export type RolesState = {
  roles: DiscordRole[];
};

const initialState: RolesState = {
  roles: []
};

export const rolesSlice = createSlice({
  name: 'roles',
  initialState,
  reducers: {
    setRoles: (state: RolesState, action: PayloadAction<DiscordRole[]>) => {
      return {
        ...state,
        roles: action.payload
      };
    },
  }
});


export const { setRoles } = rolesSlice.actions;
export default rolesSlice.reducer;
