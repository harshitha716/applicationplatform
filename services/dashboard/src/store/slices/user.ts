import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { LOADER_STATUS } from 'modules/data/data.types';
import { Session, Workspace } from 'types/api/auth.types';
import { MapAny } from 'types/commonTypes';

export type DatasetBulkLoadersType = { id: string; status: LOADER_STATUS; title: string; description: string };

export type UserState = {
  user: Session | null;
  userAccessFlags: any;
  userSessionExpired?: boolean;
  isGodMode?: boolean;
  workspace: Workspace | null;
  configuration?: MapAny;
  merchantDetails?: MapAny;
  roles?: { id: string; name: string }[];
  dashboardLoader: boolean;
  datasetBulkLoaders?: DatasetBulkLoadersType[];
};

const initialState: UserState = {
  user: null,
  userAccessFlags: {},
  userSessionExpired: false,
  isGodMode: false,
  configuration: undefined,
  merchantDetails: {},
  workspace: null,
  dashboardLoader: false,
  datasetBulkLoaders: [],
};

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUserInfo: (state, action) => {
      state.user = action.payload;
      state.userSessionExpired = false;
    },
    setUserAccessFlags: (state, action) => {
      state.userAccessFlags = { ...state.userAccessFlags, ...action.payload };
    },
    setConfiguration: (state, action) => {
      state.configuration = action.payload;
    },
    setIsGodMode: (state, action) => {
      state.isGodMode = action.payload;
    },
    setMerchantDetails: (state, action) => {
      state.merchantDetails = action.payload;
    },
    setRoles: (state, action) => {
      state.roles = action.payload;
    },
    setUser: (state, action: PayloadAction<Session>) => {
      state.user = action.payload;

      return state;
    },
    setDashboardLoader: (state, action: PayloadAction<boolean>) => {
      state.dashboardLoader = action.payload;
    },
    addDatasetBulkLoaders: (
      state,
      action: PayloadAction<{ id: string; status: LOADER_STATUS; title: string; description: string }>,
    ) => {
      if (!state.datasetBulkLoaders) {
        state.datasetBulkLoaders = [];
      }
      state.datasetBulkLoaders.push(action.payload);
    },
    removeDatasetBulkLoader: (state, action: PayloadAction<string>) => {
      if (state.datasetBulkLoaders) {
        state.datasetBulkLoaders = state.datasetBulkLoaders.filter((loader) => loader.id !== action.payload);
      }
    },
    setWorkspace: (state, action: PayloadAction<Workspace>) => {
      state.workspace = action.payload;

      return state;
    },
    resetUser: () => {
      return initialState;
    },
  },
});

export const {
  setUserInfo,
  setUserAccessFlags,
  resetUser,
  setIsGodMode,
  setConfiguration,
  setMerchantDetails,
  setRoles,
  setUser,
  addDatasetBulkLoaders,
  removeDatasetBulkLoader,
  setWorkspace,
  setDashboardLoader,
} = userSlice.actions;

export default userSlice.reducer;
