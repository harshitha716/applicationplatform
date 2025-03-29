import React, { ReactNode } from 'react';
import { captureException } from '@sentry/browser';
import { ErrorCardTypes } from 'components/commonWrapper/commonWrapper.types';
import ErrorCard from 'components/commonWrapper/ErrorCard';

class ErrorBoundary extends React.Component<{ children: ReactNode }, { hasError: boolean }> {
  constructor(props: any) {
    super(props);

    // Define a state variable to track whether is an error or not
    this.state = { hasError: false };
  }
  static getDerivedStateFromError() {
    // Update state so the next render will show the fallback UI

    return { hasError: true };
  }
  componentDidCatch(error: Error, errorInfo: any) {
    // You can use your own error logging service here
    captureException(error, errorInfo);
    console.log({ error, errorInfo });
  }
  render() {
    // Check if the error is thrown
    if (this.state.hasError)
      // You can render any custom fallback UI
      return (
        <div className='flex justify-center items-center h-screen w-full'>
          <ErrorCard
            title='Something went wrong'
            className='w-full'
            subtitle='Please try again later'
            type={ErrorCardTypes.GENERAL_API_FAIL}
            onClose={() => window.location.reload()}
          />
        </div>
      );
    // Return children components in case of no error

    return this.props.children;
  }
}

export default ErrorBoundary;
