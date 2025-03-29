import { ReactNode } from 'react';
import { defaultFnType } from 'types/commonTypes';

export enum SkeletonTypes {
  DEFAULT,
  PROGRESS_BAR,
  CUSTOM,
}

export enum ErrorCardTypes {
  GENERAL_API_FAIL = 'GENERAL_API_FAIL',
  KPI_API_FAIL = 'KPI_API_FAIL',
}

export interface CommonWrapperPropsTypes {
  children: React.ReactNode;
  isNoData?: boolean;
  isLoading?: boolean;
  isSuccess?: boolean;
  isError?: boolean;
  errorCardTitle?: string;
  errorCardSubTitle?: string;
  height?: number;
  errorCardType?: ErrorCardTypes;
  refetchFunction?: defaultFnType;
  errorCardStyle?: string;
  skeletonType?: SkeletonTypes;
  noDataBanner?: React.ReactNode;
  successCard?: React.ReactNode;
  loaderClassName?: string;
  errorCardProps?: any;
  skeleton?: any;
  renderError?: ReactNode;
  className?: string;
  skeletonItemCount?: number;
  skeletonClassName?: string;
  loader?: ReactNode;
}

export interface ErrorCardPropTypes {
  className?: string;
  onClose?: defaultFnType;
  title?: string;
  type?: ErrorCardTypes;
  isLoading?: boolean;
  height?: number;
  subtitle?: string;
  customWrapperClassName?: string;
  customImageDimensions?: { width?: number; height?: number };
  customTitleClassName?: string;
  customSubtitleClassName?: string;
  refetchButtonId?: string;
  contentClassName?: string;
}
