import { API_ENDPOINTS, REQUEST_TYPES } from 'apis/apiEndpoint.constants';
import baseApi from 'services/api';
import {
  GetDatasetDisplayConfigRequestType,
  GetDatasetDisplayConfigResponseType,
  PostDatasetDisplayConfigRequestType,
  PostDatasetDisplayConfigResponseType,
} from 'types/api/admin.types';
import { formRequestUrlWithParams } from 'utils/common';

const Admin = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getDatasetDisplayConfig: builder.query<GetDatasetDisplayConfigResponseType, GetDatasetDisplayConfigRequestType>({
      query: ({ datasetId }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.ADMIN_DATASET_DISPLAY_CONFIG_GET, { datasetId }),
      }),
    }),
    postDatasetDisplayConfig: builder.mutation<
      PostDatasetDisplayConfigResponseType,
      PostDatasetDisplayConfigRequestType
    >({
      query: ({ datasetId, body }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.ADMIN_DATASET_DISPLAY_CONFIG_POST, { datasetId }),
        method: REQUEST_TYPES.POST,
        body: body,
      }),
    }),
  }),
});

export const { useGetDatasetDisplayConfigQuery, usePostDatasetDisplayConfigMutation } = Admin;
