import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { captureException } from '@sentry/browser';
import { Mutex } from 'async-mutex';
import { ABORT_ERROR, API_DOMAIN, API_TAGS } from 'constants/api.constants';
import { ROUTES_PATH } from 'constants/routeConfig';
import { getFromLocalStorage, LOCAL_STORAGE_KEYS } from 'utils/localstorage';

const mutex = new Mutex();

const baseQuery = fetchBaseQuery({
  baseUrl: `${API_DOMAIN}/`,
  credentials: 'include',
  prepareHeaders: (headers) => {
    headers.set('Accept', 'application/json');
    headers.set(
      LOCAL_STORAGE_KEYS.XZAMP_ORGANIZATION_ID,
      getFromLocalStorage(LOCAL_STORAGE_KEYS.XZAMP_ORGANIZATION_ID) || '',
    );

    return headers;
  },
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const baseQueryWithAuth: any = async (args: any, api: any, extraOptions: any) => {
  await mutex.waitForUnlock();

  const result = await baseQuery(args, api, extraOptions);
  const path = window.location.pathname;

  const isLoginRoute = path === ROUTES_PATH.LOGIN;

  const error = result?.error;

  if (error) {
    const status = error?.status;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const data: any = error?.data;

    if (status === 401 && !isLoginRoute) {
      let loginUrl = ROUTES_PATH.LOGIN;

      if (window.location.pathname && window.location.pathname !== '/') {
        loginUrl += '?redirect_to=' + window.location.pathname;
      }

      window.location.href = loginUrl;
    }

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore status can be string on abort
    if (status !== 401 && error?.error !== ABORT_ERROR) {
      const errorObj = new Error(JSON.stringify(`${status}=${data?.error?.code ?? 'NA'}`));

      captureException(errorObj, {
        extra: {
          error: JSON.stringify(error),
          args: JSON.stringify(args),
          rtkEndpoint: api?.endpoint,
        },
      });
    }
  }

  return result;
};

const baseApi = createApi({
  reducerPath: 'api',
  tagTypes: Object.values(API_TAGS),
  baseQuery: baseQueryWithAuth,
  endpoints: () => ({}),
  refetchOnMountOrArgChange: true,
});

export default baseApi;
