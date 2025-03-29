import { FC, Fragment } from 'react';
import { captureException } from '@sentry/browser';
import { COLORS } from 'constants/colors';
import { SIZE } from 'constants/common.constants';
import { cn } from 'utils/common';
import { Loader } from 'components/common/loader/Loader';
import ProgressBar from 'components/common/RingProgress';
import { CommonWrapperPropsTypes, ErrorCardTypes, SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import ErrorCard from 'components/commonWrapper/ErrorCard';

const CommonWrapper: FC<CommonWrapperPropsTypes> = ({
  children,
  isNoData = false,
  isLoading = false,
  isSuccess = false,
  isError = false,
  errorCardSubTitle,
  errorCardTitle = 'Something went wrong',
  height,
  errorCardType = ErrorCardTypes.GENERAL_API_FAIL,
  refetchFunction,
  errorCardStyle = '',
  skeletonType = SkeletonTypes.DEFAULT,
  noDataBanner = null,
  successCard = null,
  loaderClassName = 'h-32 py-5',
  errorCardProps = {},
  renderError,
  className = '',
  skeletonItemCount = 1,
  loader,
}) => {
  const getSkeleton = () => {
    switch (skeletonType) {
      case SkeletonTypes.DEFAULT:
        return (
          <div className={cn(`${loaderClassName} flex justify-center items-center`)} style={{ height }}>
            <Loader size={SIZE.LARGE} />
          </div>
        );
      case SkeletonTypes.PROGRESS_BAR:
        return (
          <div className={cn(`${loaderClassName} flex justify-center items-center`)} style={{ height }}>
            <ProgressBar
              trackColor={COLORS.BLACK}
              indicatorColor={COLORS.WHITE}
              indicatorWidth={10}
              trackWidth={5}
              className='animate-spin'
              size={100}
              progress={30}
            />
          </div>
        );
      case SkeletonTypes.CUSTOM:
        return <>{loader}</>;
    }
  };

  if (isLoading)
    return (
      <div
        key='loader'
        className={cn('animate-opacity', className, skeletonItemCount > 1 ? 'flex flex-col gap-4' : '')}
      >
        {new Array(skeletonItemCount).fill(0).map((_, index) => (
          <Fragment key={index}>{getSkeleton()}</Fragment>
        ))}
      </div>
    );

  if (isError) {
    captureException(new Error('Error in fetching data'));

    return renderError ? (
      renderError
    ) : (
      <ErrorCard
        height={height}
        type={errorCardType}
        onClose={() => refetchFunction?.()}
        isLoading={isLoading}
        className={errorCardStyle}
        title={errorCardTitle}
        subtitle={errorCardSubTitle}
        {...errorCardProps}
      />
    );
  }

  if (isSuccess && successCard)
    return (
      <div key='success-card' className='animate-opacity'>
        {successCard}
      </div>
    );

  if (isNoData)
    return (
      <div key='zero-state' className={cn('animate-opacity', className)}>
        {noDataBanner}
      </div>
    );

  return children ? (
    <div key='main-body' className={cn('animate-opacity', className)}>
      {children}
    </div>
  ) : null;
};

export default CommonWrapper;
