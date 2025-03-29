import { AgChartOptions } from 'ag-charts-community';

export const chartData = [
  {
    quarter: "Q1'18",
    iphone: 140,
    mac: 16,
    ipad: 14,
    wearables: 12,
    services: 20,
  },
  {
    quarter: "Q2'18",
    iphone: 124,
    mac: 20,
    ipad: 14,
    wearables: 12,
    services: 30,
  },
  {
    quarter: "Q3'18",
    iphone: 112,
    mac: 20,
    ipad: 18,
    wearables: 14,
    services: 36,
  },
  {
    quarter: "Q4'18",
    iphone: 118,
    mac: 24,
    ipad: 14,
    wearables: 14,
    services: 36,
  },
];

export const chartDataJson = {
  title: {
    text: "Apple's Revenue by Product Category",
  },
  subtitle: {
    text: 'In Billion U.S. Dollars',
  },
  series: [
    {
      type: 'bar',
      xKey: 'quarter',
      yKey: 'iphone',
      yName: 'iPhone',
      stacked: true,
    },
  ],
};

export const barChartOptions = {
  data: chartData,
  title: {
    text: "Apple's Revenue by Product Category",
  },
  subtitle: {
    text: 'In Billion U.S. Dollars',
  },
  listeners: {
    seriesNodeClick: (event: any) => {
      console.log('seriesNodeClick', event);
    },
  },
  animation: {
    enabled: true,
  },
};

export const barGraph: AgChartOptions = {
  ...barChartOptions,
  series: [
    {
      type: 'bar',
      xKey: 'quarter',
      yKey: 'iphone',
      yName: 'iPhone',
      stacked: true,
    },
    {
      type: 'bar',
      xKey: 'quarter',
      yKey: 'mac',
      yName: 'Mac',
      stacked: true,
    },
    {
      type: 'bar',
      xKey: 'quarter',
      yKey: 'ipad',
      yName: 'iPad',
      stacked: true,
    },
    {
      type: 'bar',
      xKey: 'quarter',
      yKey: 'wearables',
      yName: 'Wearables',
      stacked: true,
    },
    {
      type: 'bar',
      xKey: 'quarter',
      yKey: 'services',
      yName: 'Services',
      stacked: true,
    },
  ],
  legend: {
    position: 'top',
  },
};
