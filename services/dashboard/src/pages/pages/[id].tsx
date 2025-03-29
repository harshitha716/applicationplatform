import React, { ReactElement, useEffect, useMemo } from 'react';
import { useGetPageDetailsQuery, useGetPagesQuery } from 'apis/pages';
import { useAppDispatch } from 'hooks/toolkit';
import { persistLastVisitedPage } from 'hooks/useLastVisitedPage';
import Sheets from 'modules/sheets';
import SheetsTabs from 'modules/sheets/SheetsTabs';
import { getSheetIdFromPath } from 'modules/widgets/widgets.utils';
import { useRouter } from 'next/router';
import { resetBreadcrumb } from 'store/slices/layout-configs';
import CommonWrapper from 'components/commonWrapper';
import DashboardLayout from 'components/layouts/dashboard-layout';
import 'ag-charts-enterprise';

const Page = () => {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const { id } = router.query;
  const {
    data: pageDetails,
    isLoading,
    isFetching,
    isError,
    refetch,
  } = useGetPageDetailsQuery(id as string, { refetchOnMountOrArgChange: false, skip: !id });
  const { data: pages } = useGetPagesQuery(undefined, {
    refetchOnMountOrArgChange: false,
  });
  const currentSheetId = useMemo(
    () => getSheetIdFromPath(router.asPath, id as string) ?? pageDetails?.sheets?.[0]?.sheet_id,
    [pageDetails, router.asPath],
  );

  useEffect(() => {
    if (pageDetails) {
      persistLastVisitedPage(pageDetails.page_id);
    }
  }, [pageDetails]);

  const tabs = useMemo(
    () =>
      pageDetails?.sheets
        ?.map((sheet) => ({
          value: sheet?.sheet_id,
          label: sheet?.name,
          fractionalIndex: sheet?.fractional_index,
        }))
        .sort((sheet1, sheet2) => sheet1?.fractionalIndex - sheet2?.fractionalIndex) ?? [],
    [pageDetails],
  );

  useEffect(() => {
    const currentPageTitle = pages?.find((page) => page.page_id === id)?.name ?? 'Loading...';

    persistLastVisitedPage(id as string);
    dispatch(resetBreadcrumb([currentPageTitle]));
  }, [id, pages]);

  return (
    <CommonWrapper isError={isError} refetchFunction={refetch}>
      <div className='relative h-full rounded-tl-md w-full'>
        <Sheets pageId={id as string} sheetId={currentSheetId as string} isPageLoading={isLoading} />
        <SheetsTabs tabs={tabs} currentSheetId={currentSheetId as string} isPageLoading={isFetching} />
      </div>
    </CommonWrapper>
  );
};

Page.getLayout = function getLayout(page: ReactElement) {
  return (
    <div>
      <DashboardLayout>{page}</DashboardLayout>
    </div>
  );
};

export default Page;
