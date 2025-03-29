import { API_ENDPOINTS } from 'apis/apiEndpoint.constants';
import baseApi from 'services/api';
import { WidgetDataRequestType, WidgetDataResponseType, WidgetInstanceResponseType } from 'types/api/widgets.types';
import { formRequestUrlWithParams } from 'utils/common';

const Widgets = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getWidgetInstance: builder.query<WidgetInstanceResponseType, string>({
      query: (widgetId) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.WIDGET_INSTANCE_GET, { widgetId }),
      }),
      transformResponse: ({ data }) => data,
    }),
    getWidgetData: builder.query<WidgetDataResponseType, WidgetDataRequestType>({
      query: ({ widgetId, payload }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.WIDGET_DATA_GET, { widgetId }),
        params: payload,
      }),
    }),
  }),
});

export const { useGetWidgetInstanceQuery, useGetWidgetDataQuery } = Widgets;
