import { FC, useEffect } from 'react';
import { useGetDatasetListingQuery } from 'apis/dataset';
import { useGetPagesQuery } from 'apis/pages';
import { ENVIRONMENT, ENVIRONMENT_TYPES } from 'constants/common.constants';
import { FEATURE_FLAGS } from 'constants/featureFlags';
import { ROUTES_PATH } from 'constants/routeConfig';
import { useFeatureFlags } from 'hooks/useFeatureFlags';
import { useWindowDimensions } from 'hooks/useWindowDimensions';
import ScreenSupport from 'modules/cards/ScreenSupport';
import { useRouter } from 'next/router';
import { checkScreenBreakpoint, getLeadingPathFromURL } from 'utils/common';
import { PAGE_SIZE } from 'components/common/table/table.constants';

type AuthGuardPropsType = {
  children: React.ReactNode;
};

export const RouteGuard: FC<AuthGuardPropsType> = (props) => {
  const router = useRouter();
  const { pathname } = router;
  const { id } = router.query;
  const currentPathName = getLeadingPathFromURL(pathname);
  const PAGES = getLeadingPathFromURL(ROUTES_PATH.PAGES);
  const DATASETS = getLeadingPathFromURL(ROUTES_PATH.DATASET);
  const isAdminRoute = router.pathname.startsWith(ROUTES_PATH.ADMIN);
  const { evaluate } = useFeatureFlags();
  const { width, height } = useWindowDimensions();
  const { data: datasetListingData, isLoading: isDatasetListingLoading } = useGetDatasetListingQuery(
    { page: 1, pageSize: PAGE_SIZE },
    {
      skip: currentPathName !== DATASETS || id === undefined,
      refetchOnMountOrArgChange: false,
    },
  );
  const { data: pages, isLoading: isPagesLoading } = useGetPagesQuery(undefined, {
    refetchOnMountOrArgChange: false,
  });

  useEffect(() => {
    if (currentPathName === DATASETS && id !== undefined && !isDatasetListingLoading) {
      const pageExists = datasetListingData?.datasets?.some((dataset) => dataset?.id === id);

      if (!pageExists) {
        router.push(ROUTES_PATH.NO_ACCESS);
      }
    }
  }, [currentPathName, id, isDatasetListingLoading, datasetListingData]);

  useEffect(() => {
    if (currentPathName === PAGES && !isPagesLoading && id && pages) {
      const pageExists = pages.some((page) => page?.page_id === id);

      if (!pageExists) {
        router.push(ROUTES_PATH.NO_ACCESS);
      }
    }
  }, [currentPathName, id, pages, isPagesLoading, router]);

  useEffect(() => {
    const checkAdminAccess = async () => {
      if (isAdminRoute) {
        const isAdminFeatureEnabled = await evaluate(FEATURE_FLAGS.ADMIN_PAGE);

        if (!isAdminFeatureEnabled) {
          router.push(ROUTES_PATH.NO_ACCESS);
        }
      }
    };

    checkAdminAccess();
  }, [isAdminRoute]);

  const breakpoint = checkScreenBreakpoint(width, height);

  if (breakpoint && ENVIRONMENT === ENVIRONMENT_TYPES.PRODUCTION) return <ScreenSupport />;

  return props.children;
};
