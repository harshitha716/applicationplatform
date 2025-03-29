import { ROUTES_PATH } from 'constants/routeConfig';
import Listing from 'modules/data';
import { useRouter } from 'next/router';
import DashboardLayout from 'components/layouts/dashboard-layout';

const AdminDataset = () => {
  const router = useRouter();
  const DATASETS = 'datasets';
  const onRowClicked = (event: any) => {
    router.push(`${ROUTES_PATH.ADMIN}/${DATASETS}/${event?.data?.id}`);
  };

  return <Listing onRowClicked={onRowClicked} />;
};

AdminDataset.getLayout = function getLayout(page: React.ReactElement) {
  return <DashboardLayout>{page}</DashboardLayout>;
};

export default AdminDataset;
