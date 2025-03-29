import React, { FC, ReactNode, useEffect } from 'react';
import { useGetSheetFilterConfigQuery } from 'apis/pages';
import { getDefaultFilterValues, getFormattedSheetsFiltersConfig } from 'modules/sheets/sheets.utils';
import { filtersContextActions, useFiltersContextStore } from 'components/filter/filters.context';

const InitializeSheetsFilters: FC<{ children: ReactNode; pageId: string; sheetId: string }> = ({
  children,
  pageId,
  sheetId,
}) => {
  const { dispatch } = useFiltersContextStore();
  const {
    data: sheetFilterConfig,
    isSuccess,
    isFetching,
  } = useGetSheetFilterConfigQuery(
    { pageId: pageId as string, sheetId: sheetId as string },
    { skip: !sheetId, refetchOnMountOrArgChange: false },
  );

  useEffect(() => {
    if (!isFetching && sheetFilterConfig?.native_filter_config?.length) {
      const filtersConfig = sheetFilterConfig?.native_filter_config;
      const defaultFilterValues = getDefaultFilterValues(filtersConfig);

      const filters = filtersConfig.map((filter) => {
        dispatch({
          type: filtersContextActions.ADD_EMPTY_STATE_FILTER,
          payload: { filterKey: filter.targets?.[0]?.column },
        });

        return getFormattedSheetsFiltersConfig(filter);
      });

      if (filters.length) {
        if (Object.keys(defaultFilterValues).length)
          dispatch({
            type: filtersContextActions.INITIALIZE_DEFAULT_FILTERS,
            payload: {
              selectedFilters: defaultFilterValues ?? {},
            },
          });

        dispatch({
          type: filtersContextActions.SET_INITIALISED,
        });

        dispatch({
          type: filtersContextActions.SET_FILTERS_CONFIG,
          payload: {
            filtersConfig: filters,
          },
        });
      } else {
        dispatch({
          type: filtersContextActions.SET_INITIALISED,
        });
      }
    } else {
      if (isSuccess)
        dispatch({
          type: filtersContextActions.SET_INITIALISED,
        });
    }
  }, [sheetFilterConfig, sheetId]);

  useEffect(() => {
    if (isFetching) {
      dispatch({
        type: filtersContextActions.RESET_ALL_FILTERS,
        payload: { shouldClearDate: false },
      });
      dispatch({
        type: filtersContextActions.SET_FILTERS_CONFIG,
        payload: { filtersConfig: [] },
      });
    }
    dispatch({
      type: filtersContextActions.SET_FILTER_LOADING,
      payload: { isFilterLoading: isFetching },
    });
  }, [isFetching]);

  return <div>{children}</div>;
};

export default InitializeSheetsFilters;
