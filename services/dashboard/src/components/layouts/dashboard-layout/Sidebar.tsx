import React, { useEffect } from 'react';
import { useGetPagesQuery } from 'apis/pages';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { SIDEBAR_ITEMS } from 'constants/routeConfig';
import { useAppSelector } from 'hooks/toolkit';
import { usePersistedPageNavigation } from 'hooks/useLastVisitedPage';
import { useLogout } from 'hooks/useLogout';
import { useRouter } from 'next/router';
import { RootState } from 'store';
import { cn } from 'utils/common';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import PageNavTab from 'components/layouts/dashboard-layout/components/PageNavTab';
import SidebarTab from 'components/layouts/dashboard-layout/components/SidebarTab';
import SkeletonLoaderSidebarPages from 'components/layouts/dashboard-layout/components/SkeletonLoaderSidebarPages';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const Sidebar = () => {
  const { isSidebarOpen } = useAppSelector((state: RootState) => state.layoutConfig);
  const router = useRouter();
  const pathname = router?.pathname;
  const pageId = router?.query?.id;
  const { logout } = useLogout();
  const { data: pages, isLoading: isLoadingPages } = useGetPagesQuery(undefined, {
    refetchOnMountOrArgChange: false,
  });
  const { pushToMostRelevantPage } = usePersistedPageNavigation(pages ?? []);

  useEffect(() => {
    if (pages) {
      pushToMostRelevantPage();
    }
  }, [pages]);

  const filteredSidebarItems = SIDEBAR_ITEMS.filter((item) => !item.isHidden);

  return (
    <div className={cn('relative transition-all', isSidebarOpen ? 'w-60' : 'w-0')}>
      <div className='w-60'>
        <div className='px-2 border-b border-GRAY_400 pb-4'>
          {filteredSidebarItems.map((item) => (
            <SidebarTab
              key={item.label}
              name={item.label}
              path={item.path}
              iconId={item.iconId}
              isSelected={pathname === item?.path}
            />
          ))}
        </div>
        <div className='px-2 py-2.5'>
          <div className='f-11-600 text-GRAY_700 px-1.5 py-2'>Pages</div>
          <CommonWrapper
            isLoading={isLoadingPages}
            skeletonType={SkeletonTypes.CUSTOM}
            loader={<SkeletonLoaderSidebarPages />}
          >
            {pages?.map((item) => (
              <PageNavTab
                key={item?.page_id}
                label={item?.name}
                pageId={item?.page_id}
                isSelected={pageId === item?.page_id}
              />
            ))}
          </CommonWrapper>
        </div>
        <div
          className='border-t border-GRAY_400 px-4 py-3 absolute bottom-0 w-full cursor-pointer h-[57px] flex items-center gap-2.5 text-GRAY_900'
          onClick={logout}
        >
          <SvgSpriteLoader iconCategory={ICON_SPRITE_TYPES.GENERAL} id='log-out-02' height={14} width={14} />
          <div className='f-13-500'>Logout</div>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
