import { act, renderHook } from '@testing-library/react';
import { ROUTES_PATH } from 'constants/routeConfig';
import { usePersistedPageNavigation } from 'hooks/useLastVisitedPage';
import { useRouter } from 'next/router';
import { PageResponseType } from 'types/api/pagesApi.types';
import { getFromLocalStorage, removeFromLocalStorage, setToLocalStorage } from 'utils/localstorage';

jest.mock('next/router', () => ({
  useRouter: jest.fn(),
}));

jest.mock('utils/localstorage', () => ({
  getFromLocalStorage: jest.fn(),
  setToLocalStorage: jest.fn(),
  removeFromLocalStorage: jest.fn(),
  LOCAL_STORAGE_KEYS: {
    LAST_VISITED_PAGE_ID: 'LAST_VISITED_PAGE_ID',
  },
}));

describe('usePersistedPageNavigation', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  const samplePages: PageResponseType[] = [
    {
      page_id: 'page1',
      name: 'Page 1',
      description: 'Page 1 description',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      fractional_index: 0,
      organization_id: 'org1',
    },
    {
      page_id: 'page2',
      name: 'Page 2',
      description: 'Page 2 description',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      fractional_index: 0,
      organization_id: 'org1',
    },
  ];

  const testCases = [
    {
      name: 'should navigate to last visited page if it exists in pages list',
      pathname: ROUTES_PATH.HOME,
      lastVisitedPageId: 'page1',
      pages: samplePages,
      expectedPath: '/pages/page1',
      expectedPushCalls: 1,
      shouldRemoveFromStorage: false,
    },
    {
      name: 'should clear last visited page and navigate to first page if last visited page not found',
      pathname: ROUTES_PATH.HOME,
      lastVisitedPageId: 'nonexistent',
      pages: samplePages,
      expectedPath: '/pages/page1',
      expectedPushCalls: 1,
      shouldRemoveFromStorage: true,
    },
    {
      name: 'should navigate to first page if no last visited page exists',
      pathname: ROUTES_PATH.HOME,
      lastVisitedPageId: 'page3',
      pages: samplePages,
      expectedPath: '/pages/page1',
      expectedPushCalls: 1,
      shouldRemoveFromStorage: true,
    },
    {
      name: 'should not navigate if not on home page',
      pathname: '/some-other-page',
      lastVisitedPageId: 'page1',
      pages: samplePages,
      expectedPath: null,
      expectedPushCalls: 0,
      shouldRemoveFromStorage: false,
    },
    {
      name: 'should not navigate if first navigation already done',
      pathname: ROUTES_PATH.HOME,
      lastVisitedPageId: 'page1',
      pages: [samplePages[0]],
      expectedPath: '/pages/page1',
      expectedPushCalls: 1,
      shouldRemoveFromStorage: false,
    },
  ];

  test.each(testCases)(
    '$name',
    ({ pathname, lastVisitedPageId, pages, expectedPath, expectedPushCalls, shouldRemoveFromStorage }) => {
      const pushFn = jest.fn();

      (useRouter as jest.Mock).mockReturnValue({
        pathname,
        push: pushFn,
      });

      (getFromLocalStorage as jest.Mock).mockReturnValue(lastVisitedPageId);

      const { result } = renderHook(() => usePersistedPageNavigation(pages));

      act(() => result.current.pushToMostRelevantPage());

      // For the "already done" test case, call it twice
      if (expectedPushCalls === 1) {
        act(() => result.current.pushToMostRelevantPage());
      }

      expect(pushFn).toHaveBeenCalledTimes(expectedPushCalls);
      if (expectedPath) {
        expect(pushFn).toHaveBeenCalledWith(expectedPath);
      }
      if (shouldRemoveFromStorage) {
        expect(removeFromLocalStorage).toHaveBeenCalled();
      }
    },
  );

  it('should persist last visited page', () => {
    const pagesList: PageResponseType[] = [];
    const pageId = 'test-page';

    const { result } = renderHook(() => usePersistedPageNavigation(pagesList));

    act(() => result.current.setLastVisitedPage(pageId));

    expect(setToLocalStorage).toHaveBeenCalledWith('LAST_VISITED_PAGE_ID', pageId);
  });
});
