import React, { ReactElement, useEffect } from 'react';
import { useAppDispatch } from 'hooks/toolkit';
import PeoplePage from 'modules/team/PeoplePage';
import { resetBreadcrumb } from 'store/slices/layout-configs';
import DashboardLayout from 'components/layouts/dashboard-layout';

const Team = () => {
  const appDispatch = useAppDispatch();

  useEffect(() => {
    appDispatch(resetBreadcrumb(['Team']));
  }, []);

  return <PeoplePage />;
};

Team.getLayout = function getLayout(page: ReactElement) {
  return (
    <div>
      <DashboardLayout>{page}</DashboardLayout>
    </div>
  );
};

export default Team;
