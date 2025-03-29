export type GetDatasetDisplayConfigRequestType = {
  datasetId: string;
};

export type displayConfigType = {
  column: string;
  is_hidden: string;
  is_editable: string;
  type: string;
  config: {
    amount_column: string;
    currency_column: string;
  };
};

export type GetDatasetDisplayConfigResponseType = {
  display_config: displayConfigType[];
};

export type PostDatasetDisplayConfigRequestType = {
  datasetId: string;
  body: { display_config: displayConfigType[] };
};

export type PostDatasetDisplayConfigResponseType = {
  action_id: string;
};
