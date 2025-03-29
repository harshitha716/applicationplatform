export const PAGES_ITEMS = [
  {
    label: 'Daily Liquidity Summary',
    iconId: 'notebook',
  },
  {
    label: 'Cash Summary',
    iconId: 'notebook',
  },
  {
    label: 'Bank Account Balances',
    iconId: 'notebook',
  },
  {
    label: 'Cash Positioning',
    iconId: 'notebook',
  },
];

export const WORKSPACE_ITEMS = [
  {
    label: 'Reconciliation',
    workspace_id: 'reconciliation',
    color: '#40A97F',
  },
  {
    label: 'Cash Management',
    workspace_id: 'cash-management',
    color: '#0052D6',
  },
  {
    label: 'Financial Forecasting',
    workspace_id: 'financial-forecasting',
    color: '#BF0000',
  },
];

export const barGraphInstance = {
  instance_id: 'currency_volume_analysis',
  widget_id: 1,
  type: 'bar',
  title: 'Transaction Volume by Currency',
  data_mappings: {
    datasets: [
      {
        id: 'CashOpsBankTransactions',
      },
    ],
    mappings: {
      x_axis: {
        field: 'CurrencyCode',
      },
      y_axis: {
        field: 'IntegerAmount',
        aggregation: 'sum',
      },
    },
  },
  visual_config: {},
};

export const barGraphData = {
  result: [
    {
      status: 'success',
      error: null,
      rowcount: 5,
      columns: [
        {
          column_name: 'CurrencyCode',
          column_type: 'STRING',
        },
        {
          column_name: 'SUM(IntegerAmount)',
          column_type: 'NUMBER',
        },
      ],
      data: [
        {
          CurrencyCode: 'USD',
          IntegerAmount: 1543437.0,
          IntegerAmountV2: 1543437.0,
        },
        {
          CurrencyCode: 'EUR',
          IntegerAmount: 756909.0,
          IntegerAmountV2: 756909.0,
        },
        {
          CurrencyCode: 'GBP',
          IntegerAmount: 432224.0,
          IntegerAmountV2: 432224.0,
        },
        {
          CurrencyCode: 'JPY',
          IntegerAmount: 234567.0,
          IntegerAmountV2: 234567.0,
        },
        {
          CurrencyCode: 'SGD',
          IntegerAmount: 123456.0,
          IntegerAmountV2: 123456.0,
        },
      ],
    },
  ],
};
