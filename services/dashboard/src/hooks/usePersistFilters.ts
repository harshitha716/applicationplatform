import { STORAGE_TYPES } from 'constants/common.constants';
import { useAppSelector } from 'hooks/toolkit';
import { RootState } from 'store';
import { defaultFn, MapAny } from 'types/commonTypes';
import { FILTER_KEYS } from 'components/filter/filters.constants';

export enum PERSISTENT_FILTER_ID {
  TRANSACTIONS = 'transactions',
  CREATE_RULES = 'createRules',
}

export const usePersistFilters = (storageType: STORAGE_TYPES = STORAGE_TYPES.SESSION) => {
  const { user } = useAppSelector((state: RootState) => state.user);

  if (!user?.user_id)
    return {
      setFiltersToStorageForPage: defaultFn,
      clearFiltersFromStorage: defaultFn,
      getFiltersFromStorageForPage: defaultFn,
    };

  const getFiltersFromStorageForPage = (page: string) => {
    try {
      let filters = null;
      const key = `${page}_filters`;

      if (storageType === STORAGE_TYPES.LOCAL) {
        filters = localStorage.getItem(key);
      } else {
        filters = sessionStorage.getItem(key);
      }

      const value = filters ? JSON.parse(filters) : {};

      if (value && !value[user?.user_id]) {
        clearFiltersFromStorage(page);

        return {};
      }

      const selectedFilters = { ...value[user?.user_id] };
      const dateKey = FILTER_KEYS.DATE_RANGE;

      if (
        selectedFilters?.[dateKey] &&
        selectedFilters?.[dateKey]?.start_date &&
        selectedFilters?.[dateKey]?.end_date
      ) {
        selectedFilters[dateKey].start_date = new Date(
          value[user?.user_id][dateKey].start_date || value[user?.user_id][dateKey].start,
        );
        selectedFilters[dateKey].end_date = new Date(
          value[user?.user_id][dateKey].end_date || value[user?.user_id][dateKey].end,
        );
      }

      return selectedFilters ?? {};
    } catch (e) {
      console.error('Error while getting filters from session storage', e);

      return {};
    }
  };

  const setFiltersToStorageForPage = (page: string, selectedFilters: MapAny) => {
    try {
      const value = {
        [user?.user_id]: selectedFilters,
      };
      const key = `${page}_filters`;
      const storageValue = JSON.stringify(value);

      if (storageType === STORAGE_TYPES.LOCAL) {
        localStorage.setItem(key, storageValue);

        return;
      }

      sessionStorage.setItem(key, storageValue);
    } catch (e) {
      console.error('Error while setting filters from session storage', e);

      return {};
    }
  };

  const clearFiltersFromStorage = (page: string) => {
    const key = `${page}_filters`;

    if (storageType === STORAGE_TYPES.LOCAL) {
      localStorage.removeItem(key);

      return;
    }

    sessionStorage.removeItem(key);
  };

  return {
    setFiltersToStorageForPage,
    clearFiltersFromStorage,
    getFiltersFromStorageForPage,
  };
};
