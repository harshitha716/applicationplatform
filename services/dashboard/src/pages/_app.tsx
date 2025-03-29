import { Provider } from 'react-redux';
import { ToastContainer } from 'react-toastify';
import { LicenseManager as LicenseManagerCharts } from 'ag-charts-enterprise';
import { LicenseManager } from 'ag-grid-enterprise';
import { inter } from 'constants/common.constants';
import { FAVICON } from 'constants/icons';
import { FeatureFlagsProvider } from 'modules/feature-flags/provider';
import type { AppProps } from 'next/app';
import Head from 'next/head';
import ErrorBoundary from 'pages/ErrorBoundary';
import { store } from 'store';
import { NextPageWithLayout } from 'types/commonTypes';
import { AG_CHART_KEY, AG_GRID_KEY } from 'components/common/agGridTable/agGridTable.constants';
import { AuthGuard } from 'components/hoc/AuthGuard';
import { RouteGuard } from 'components/hoc/RouteGuard';
import NetworkStatus from 'components/NetWorkStatus';
import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';
import 'react-toastify/dist/ReactToastify.css';
import 'styles/colors.css';
import 'styles/fonts.css';
import 'styles/globals.css';
import 'styles/react-datepicker.css';
import 'styles/react-dates.css';

LicenseManager.setLicenseKey(AG_GRID_KEY);
LicenseManagerCharts.setLicenseKey(AG_CHART_KEY);

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

export default function App({ Component, pageProps }: AppPropsWithLayout) {
  const getLayout = Component.getLayout ?? ((page: any) => page);

  const getComponent = () => {
    return <div>{getLayout(<Component {...pageProps} />)}</div>;
  };

  return (
    <>
      <Head>
        <meta name='viewport' content='width=device-width, initial-scale=1.0' />
        <title>Zamp</title>
        <link rel='icon' type='image/x-icon' href={FAVICON} />
      </Head>
      <div className={inter.className}>
        <NetworkStatus />
        <ErrorBoundary>
          <Provider store={store}>
            <AuthGuard loginRoute='/login'>
              <FeatureFlagsProvider>
                <ToastContainer />
                <RouteGuard>
                  <div className={'h-screen light-mode'}>{getComponent()}</div>
                </RouteGuard>
              </FeatureFlagsProvider>
            </AuthGuard>
          </Provider>
        </ErrorBoundary>
      </div>
    </>
  );
}
