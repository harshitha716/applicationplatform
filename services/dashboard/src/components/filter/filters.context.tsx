import React, { createContext, Dispatch, FC, ReactElement, useContext, useReducer } from 'react';
import { PERSISTENT_FILTER_ID, usePersistFilters } from 'hooks/usePersistFilters';
import { MapAny } from 'types/commonTypes';
import { FilterConfigType, FilterEntityMenuType } from 'components/filter/filter.types';
import { FILTER_PERIODICITIES } from 'components/filter/filters.constants';

enum filtersContextActions {
  INITIALIZE_DEFAULT_FILTERS = 'INITIALIZE_DEFAULT_FILTERS',
  SET_SELECTED_FILTERS = 'SET_SELECTED_FILTERS',
  SET_FILTERS = 'SET_FILTERS',
  RESET_ALL_FILTERS = 'RESET_ALL_FILTERS',
  ADD_EMPTY_STATE_FILTER = 'ADD_EMPTY_STATE_FILTER',
  REMOVE_FILTER = 'REMOVE_FILTER',
  GET_FILTERS_FROM_LOCAL_STORAGE = 'GET_FILTERS_FROM_LOCAL_STORAGE',
  SET_PERSIST_ID = 'SET_PERSIST_ID',
  SET_INITIALISED = 'SET_INITIALISED',
  RESET_INITIALISED = 'RESET_INITIALISED',
  REPLACE_FILTERS = 'REPLACE_FILTERS',
  INCREMENT_SELECTED_FILTERS_CHANGE_COUNT = 'INCREMENT_SELECTED_FILTERS_CHANGE_COUNT',
  SET_FILTERS_CONFIG = 'SET_FILTERS_CONFIG',
  SET_TAG_SUGGESTIONS = 'SET_TAG_SUGGESTIONS',
  SET_PERIODICITY = 'SET_PERIODICITY',
  SET_TOTAL_ROWS = 'SET_TOTAL_ROWS',
  SET_FILTER_LOADING = 'SET_FILTER_LOADING',
  SET_STATUS_BAR = 'SET_STATUS_BAR',
}

interface InitialStateType {
  filters: MapAny;
  filtersConfig?: FilterConfigType[] | null;
  selectedFilters: MapAny; // Powers the filters to be sent to the api calls
  selectedFiltersInUI: MapAny; // Powers the filters to be displayed in the UI - We need this as we need to show the filters in the UI even if their value is empty
  selectedEntity?: FilterEntityMenuType;
  search?: string;
  selectedFiltersChangeCount: number;
  persistId: PERSISTENT_FILTER_ID | null;
  isFilterInitialized?: boolean;
  isFilterLoading?: boolean;
  periodicity?: FILTER_PERIODICITIES;
  currentPageFilters: string[];
  totalRows: number;
}

export interface ActionType {
  type: keyof typeof filtersContextActions;
  payload?: MapAny;
}

const initialState: InitialStateType = {
  filters: {},
  filtersConfig: null,
  selectedFilters: {},
  selectedFiltersChangeCount: 0,
  search: '',
  persistId: null,
  isFilterInitialized: false,
  isFilterLoading: false,
  selectedFiltersInUI: {},
  currentPageFilters: [],
  totalRows: 0,
};

const context = createContext<{
  state: InitialStateType;
  dispatch: Dispatch<ActionType>;
}>({
  state: initialState,
  dispatch: () => null,
});

const { Provider } = context;

/* eslint-disable react/display-name */
export const StateProvider: FC<{ children: ReactElement }> = ({ children }) => {
  const { setFiltersToStorageForPage, getFiltersFromStorageForPage, clearFiltersFromStorage } = usePersistFilters();

  const onSetFiltersToStorage = (persistId: PERSISTENT_FILTER_ID, selectedFilters: MapAny) => {
    const filters: MapAny = {};

    for (const key in selectedFilters) {
      // if (isFilterValueEmpty(key, selectedFilters[key])) {
      //   continue;
      // }

      filters[key] = selectedFilters[key];
    }

    setFiltersToStorageForPage(persistId, filters);
  };

  const [state, dispatch] = useReducer((state: InitialStateType, action: ActionType): InitialStateType => {
    switch (action.type) {
      case filtersContextActions.SET_PERSIST_ID:
        return { ...state, persistId: action?.payload?.persistId };
      case filtersContextActions.SET_INITIALISED:
        return { ...state, isFilterInitialized: true };
      case filtersContextActions.RESET_INITIALISED:
        return { ...state, isFilterInitialized: false };
      case filtersContextActions.INITIALIZE_DEFAULT_FILTERS:
        return {
          ...state,
          selectedFilters: { ...state?.selectedFilters, ...action?.payload?.selectedFilters },
          selectedFiltersInUI: { ...state?.selectedFiltersInUI, ...action?.payload?.selectedFilters },
          currentPageFilters: Object.keys(action?.payload?.selectedFilters),
          isFilterInitialized: true,
        };
      case filtersContextActions.GET_FILTERS_FROM_LOCAL_STORAGE: {
        const selectedFilters = getFiltersFromStorageForPage(action?.payload?.persistId);

        // Needed to migrate to the new accounts filter that add an array of account ids as opposed to an object with accounts data
        // removeAccountsFilterObject(selectedFilters);

        return { ...state, selectedFilters, selectedFiltersInUI: { ...selectedFilters }, isFilterInitialized: true };
      }
      case filtersContextActions.SET_FILTERS:
        return { ...state, filters: { ...state?.filters, ...action?.payload?.filters } };
      case filtersContextActions.SET_FILTERS_CONFIG:
        return { ...state, filtersConfig: action?.payload?.filtersConfig };
      case filtersContextActions.SET_PERIODICITY:
        return { ...state, periodicity: action?.payload?.periodicity };
      case filtersContextActions.SET_SELECTED_FILTERS: {
        const selectedFilters = {
          ...state?.selectedFilters,
          ...action?.payload?.selectedFilters,
        };

        const selectedFiltersInUI = {
          ...state?.selectedFiltersInUI,
          ...action?.payload?.selectedFilters,
        };

        // Needed to migrate to the new accounts filter that add an array of account ids as opposed to an object with accounts data

        return {
          ...state,
          selectedFilters,
          selectedFiltersInUI,
          currentPageFilters: Object.keys(selectedFilters),
        };
      }

      case filtersContextActions.REPLACE_FILTERS: {
        const selectedFilters = { ...action?.payload?.selectedFilters };

        return {
          ...state,
          selectedFilters,
          selectedFiltersInUI: { ...selectedFilters },
          selectedFiltersChangeCount: (state?.selectedFiltersChangeCount ?? 0) + 1,
        };
      }

      case filtersContextActions.RESET_ALL_FILTERS: {
        const date_range = state?.selectedFilters?.date_range;

        if (!action?.payload?.shouldClearDate && date_range) {
          if (state?.persistId) {
            onSetFiltersToStorage(state?.persistId, { date_range });
          }

          return {
            ...state,
            search: '',
            selectedFilters: { date_range },
            selectedFiltersInUI: { date_range },
            selectedFiltersChangeCount: state.selectedFiltersChangeCount + 1,
          };
        }

        clearFiltersFromStorage(state?.persistId as any);

        return {
          ...state,
          search: '',
          selectedFilters: {},
          selectedFiltersInUI: {},
          selectedFiltersChangeCount: 0,
        };
      }

      case filtersContextActions.ADD_EMPTY_STATE_FILTER: {
        const selectedStateValue = state?.selectedFilters[action?.payload?.filterKey];
        const currentPageFilters = state?.currentPageFilters;

        // if (selectedStateValue && !isFilterValueEmpty(action?.payload?.filterKey, selectedStateValue)) {
        if (selectedStateValue) {
          return { ...state };
        }

        const selectedFiltersCurrentState = { ...state?.selectedFilters };
        const selectedFiltersInUICurrentState = { ...state?.selectedFiltersInUI };

        return {
          ...state,
          selectedFilters: { ...selectedFiltersCurrentState },
          currentPageFilters: [...currentPageFilters, action?.payload?.filterKey],
          selectedFiltersInUI: {
            ...selectedFiltersInUICurrentState,
            [action?.payload?.filterKey]: null,
          },
        };
      }

      case filtersContextActions.REMOVE_FILTER: {
        const selectedFilters = { ...state?.selectedFilters };
        const selectedFiltersInUI = { ...state?.selectedFiltersInUI };
        const activeCurrentPageFilters =
          state.currentPageFilters.filter((filter) => filter !== action?.payload?.filterKey) ?? [];

        if (action?.payload?.filterKey in selectedFilters) {
          delete selectedFilters[action?.payload?.filterKey];
          if (state?.persistId) {
            onSetFiltersToStorage(state?.persistId, selectedFilters);
          }
        }

        delete selectedFiltersInUI[action?.payload?.filterKey];

        return {
          ...state,
          selectedFilters,
          selectedFiltersChangeCount: (state?.selectedFiltersChangeCount ?? 0) + 1,
          selectedFiltersInUI,
          currentPageFilters: activeCurrentPageFilters,
        };
      }

      case filtersContextActions.INCREMENT_SELECTED_FILTERS_CHANGE_COUNT: {
        return { ...state, selectedFiltersChangeCount: state.selectedFiltersChangeCount + 1 };
      }

      case filtersContextActions.SET_TOTAL_ROWS: {
        return { ...state, totalRows: action?.payload?.totalRows };
      }

      case filtersContextActions.SET_FILTER_LOADING: {
        return { ...state, isFilterLoading: action?.payload?.isFilterLoading };
      }

      default:
        return state;
    }
  }, initialState);

  return <Provider value={{ state, dispatch }}>{children}</Provider>;
};

const withFiltersContext = (WrappedComponent: FC<any>) => {
  return (props: MapAny) => (
    <StateProvider>
      <WrappedComponent {...props} />
    </StateProvider>
  );
};

const useFiltersContextStore = () => useContext(context);

export { filtersContextActions, useFiltersContextStore, withFiltersContext };
