import { API_ENDPOINTS, REQUEST_TYPES } from 'apis/apiEndpoint.constants';
import { APITags } from 'constants/api.constants';
import baseApi from 'services/api';
import {
  AudiencesByOrganisationIdRequest,
  AudiencesByOrganisationIdResponse,
  DeleteAudienceFromOrganizationAccessType,
  GetMembershipRequestsByOrganizationIdRequest,
  GetMembershipRequestsByOrganizationIdResponse,
  GetTeamsByOrganizationIdRequestType,
  GetTeamsByOrganizationIdResponseType,
  InvitedAudiencesByOrganisationIdResponse,
  PatchChangeAudienceRoleInOrganizationType,
  PostAddTeamToAudienceRequestType,
  PostAudiencesInviteData,
  PostTeamsByOrganizationIdRequestType,
  PostTeamsByOrganizationIdResponseType,
  RemoveTeamFromAudienceRequestType,
} from 'types/api/people.types';
import { formRequestUrlWithParams } from 'utils/common';

const People = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAudiencesByOrganisationId: builder.query<AudiencesByOrganisationIdResponse[], AudiencesByOrganisationIdRequest>({
      query: ({ organizationId }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.AUDIENCES_BY_ORGANIZATION_ID_GET, { organizationId }),
      }),
      transformResponse: (data) => data,
      providesTags: [APITags.GET_PEOPLE_TEAM_MEMBERS],
    }),
    getInvitedAudiencesByOrganisationId: builder.query<
      InvitedAudiencesByOrganisationIdResponse[],
      AudiencesByOrganisationIdRequest
    >({
      query: ({ organizationId }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.INVITED_AUDIENCES_BY_ORGANIZATION_ID_GET, { organizationId }),
      }),
      transformResponse: (data) => data,
      providesTags: [APITags.GET_PEOPLE_INVITATIONS],
    }),

    postInviteAudiencesByOrganisationId: builder.mutation<
      void,
      { organizationId: string; body: PostAudiencesInviteData }
    >({
      query: ({ organizationId, body }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.INVITE_AUDIENCES_BY_ORGANIZATION_ID_POST, { organizationId }),
        method: REQUEST_TYPES.POST,
        body: body,
      }),
      invalidatesTags: [APITags.GET_PEOPLE_INVITATIONS],
    }),
    patchChangeAudienceRoleInOrganization: builder.mutation<void, PatchChangeAudienceRoleInOrganizationType>({
      query: ({ organizationId, body }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.CHANGE_AUDIENCE_ROLE_IN_ORGANIZATION_PATCH, { organizationId }),
        method: REQUEST_TYPES.PATCH,
        body: body,
      }),
      invalidatesTags: [APITags.GET_PEOPLE_TEAM_MEMBERS],
    }),
    deleteAudienceFromOrganizationAccess: builder.mutation<void, DeleteAudienceFromOrganizationAccessType>({
      query: ({ organizationId, body }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.DELETE_AUDIENCE_FROM_ORGANIZATION_ACCESS, { organizationId }),
        method: REQUEST_TYPES.DELETE,
        body: body,
      }),
      invalidatesTags: [APITags.GET_PEOPLE_TEAM_MEMBERS],
    }),
    getMembershipRequestsByOrganizationId: builder.query<
      GetMembershipRequestsByOrganizationIdResponse,
      GetMembershipRequestsByOrganizationIdRequest
    >({
      query: ({ organizationId }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.MEMBERSHIP_REQUESTS_BY_ORGANIZATION_ID_GET, { organizationId }),
      }),
    }),
    getOrganizationMembershipRequestsAll: builder.query<GetMembershipRequestsByOrganizationIdResponse, void>({
      query: () => ({ url: API_ENDPOINTS.MEMBERSHIP_REQUESTS_ALL_GET }),
    }),
    getTeamsByOrganizationId: builder.query<GetTeamsByOrganizationIdResponseType, GetTeamsByOrganizationIdRequestType>({
      query: ({ organizationId }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.TEAMS_BY_ORGANIZATION_ID_GET, { organizationId }),
      }),
      providesTags: [APITags.GET_ALL_TEAMS],
    }),
    postAddTeamToOrganization: builder.mutation<
      PostTeamsByOrganizationIdResponseType,
      PostTeamsByOrganizationIdRequestType
    >({
      query: ({ organizationId, payload }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.TEAMS_BY_ORGANIZATION_ID_POST, { organizationId }),
        method: REQUEST_TYPES.POST,
        body: payload,
      }),
      invalidatesTags: [APITags.GET_ALL_TEAMS],
    }),
    postAddTeamToAudience: builder.mutation<void, PostAddTeamToAudienceRequestType>({
      query: ({ organizationId, teamId, payload }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.ADD_TEAMS_TO_AUDIENCE_POST, { organizationId, teamId }),
        method: REQUEST_TYPES.POST,
        body: payload,
      }),
      invalidatesTags: [APITags.GET_ALL_TEAMS],
    }),
    removeTeamFromAudience: builder.mutation<void, RemoveTeamFromAudienceRequestType>({
      query: ({ organizationId, teamId, payload }) => ({
        url: formRequestUrlWithParams(API_ENDPOINTS.REMOVE_TEAMS_FROM_AUDIENCE_POST, { organizationId, teamId }),
        method: REQUEST_TYPES.POST,
        body: payload,
      }),
      invalidatesTags: [APITags.GET_ALL_TEAMS],
    }),
  }),
});

export const {
  useGetAudiencesByOrganisationIdQuery,
  useGetInvitedAudiencesByOrganisationIdQuery,
  usePostInviteAudiencesByOrganisationIdMutation,
  usePatchChangeAudienceRoleInOrganizationMutation,
  useDeleteAudienceFromOrganizationAccessMutation,
  useGetMembershipRequestsByOrganizationIdQuery,
  useGetOrganizationMembershipRequestsAllQuery,
  useGetTeamsByOrganizationIdQuery,
  usePostAddTeamToOrganizationMutation,
  usePostAddTeamToAudienceMutation,
  useRemoveTeamFromAudienceMutation,
} = People;
