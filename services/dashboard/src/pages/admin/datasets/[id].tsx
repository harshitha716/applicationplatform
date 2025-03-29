import React, { ReactElement } from 'react';
import AdminDatasetById from 'modules/admin/AdminDatasetById';
import { useParams } from 'next/navigation';
import DashboardLayout from 'components/layouts/dashboard-layout';

const AdminDatasetId = () => {
  const { id } = useParams();

  return <AdminDatasetById id={id as string} />;
};

AdminDatasetId.getLayout = function getLayout(page: ReactElement) {
  return (
    <div>
      <DashboardLayout>{page}</DashboardLayout>
    </div>
  );
};

export default AdminDatasetId;
