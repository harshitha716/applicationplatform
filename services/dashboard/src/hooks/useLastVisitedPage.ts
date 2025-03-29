import React from 'react';
import { getPageRouteById, ROUTES_PATH } from 'constants/routeConfig';
import { useRouter } from 'next/router';
import { PageResponseType } from 'types/api/pagesApi.types';
import { getFromLocalStorage, LOCAL_STORAGE_KEYS, removeFromLocalStorage, setToLocalStorage } from 'utils/localstorage';

const getLastVisitedPage = (): string => {
  return getFromLocalStorage(LOCAL_STORAGE_KEYS.LAST_VISITED_PAGE_ID) || '';
};

export const persistLastVisitedPage = (pageId: string) => {
  setToLocalStorage(LOCAL_STORAGE_KEYS.LAST_VISITED_PAGE_ID, pageId);
};
const clearLastVisitedPage = () => {
  removeFromLocalStorage(LOCAL_STORAGE_KEYS.LAST_VISITED_PAGE_ID);
};

export const usePersistedPageNavigation = (pagesList: PageResponseType[]) => {
  const [firstNavigationDone, setFirstNavigationDone] = React.useState(false);

  const router = useRouter();

  const pushToMostRelevantPage = () => {
    // navigate to the most relevant page when the user visits the home page
    if (router.pathname === ROUTES_PATH.HOME) {
      // if the user has a last visited page, check if it exists in the pages list and navigate to it
      // if it doesn't exist, clear the last visited page
      const lastVisitedPageId = getLastVisitedPage();

      if (lastVisitedPageId) {
        if (pagesList.find((page) => page.page_id === lastVisitedPageId)) {
          router.push(getPageRouteById(lastVisitedPageId));

          return;
        } else {
          clearLastVisitedPage();
        }
      }

      // if the user has no last visited page, navigate to the first page in the list if it exists
      if (pagesList.length > 0) {
        router.push(getPageRouteById(pagesList[0].page_id));
      }
    }
  };

  return {
    pushToMostRelevantPage: () => {
      if (firstNavigationDone) return;
      setFirstNavigationDone(true);
      pushToMostRelevantPage();
    },
    setLastVisitedPage: (pageId: string) => {
      persistLastVisitedPage(pageId);
    },
  };
};
