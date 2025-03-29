import React, { ReactElement } from 'react';
import DrilldownByDatasetAndRowId from 'modules/data/DrilldownByDatasetAndRowId';
import DashboardLayout from 'components/layouts/dashboard-layout';

const Drilldown = () => {
  return <DrilldownByDatasetAndRowId />;
};

Drilldown.getLayout = function getLayout(page: ReactElement) {
  return (
    <div>
      <DashboardLayout>{page}</DashboardLayout>
    </div>
  );
};

export default Drilldown;
