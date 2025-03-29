export type LogoutFlow = {
  logout_url: string;
  logout_token: string;
};

export type LoginFlow = {
  id: string;
  organization_id: null;
  type: string;
  expires_at: string;
  issued_at: string;
  request_url: string;
  ui: {
    action: string;
    method: string;
    nodes: {
      type: string;
      group: string;
      attributes: {
        name: string;
        type: string;
        value: string;
        disabled: boolean;
        node_type: string;
      };
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      messages: any[];
      meta: {
        label: {
          id: number;
          text: string;
          type: string;
          context: {
            provider: string;
          };
        };
      };
    }[];
  };
  created_at: string;
  updated_at: string;
  refresh: boolean;
  requested_aal: string;
  state: string;
};

// TODO: check if type is correct
export type ErrorDetails = {
  message: string;
  id: string;
  error: {
    code: string;
    status: number;
    reason: string;
    message: string;
  };
  created_at: string;
  updated_at: string;
};

export type Workspace = {
  workspace_id: string;
  name: string;
  description: string;
};

export type Organization = {
  organization_id: string;
  name: string;
  resource_audience_policies: {
    privilege: string;
    resource_audience_type: string;
    resource_audience_id: string;
  }[];
};

export enum ResourceAudienceType {
  ORGANIZATION = 'organization',
  USER = 'user',
  TEAM = 'team',
}

export type Session = {
  user_id: string;
  workspaces: Workspace[];
  organization_id: Workspace;
  user_email: string;
  orgs: Organization[];
};

export type loginPayloadType = {
  url: string;
  body: string;
};

export enum UserRoleIdType {
  USER = 'user',
}
