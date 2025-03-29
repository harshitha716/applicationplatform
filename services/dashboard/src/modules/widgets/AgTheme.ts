import { AgChartTheme } from 'ag-charts-community';
import { CHART_PALETTE, COLORS } from 'constants/colors';

export const AG_CHART_THEME: AgChartTheme = {
  palette: CHART_PALETTE.palette,
  overrides: {
    common: {
      legend: {
        item: {
          label: {
            fontSize: 12,
            fontWeight: 450,
            fontFamily: 'Inter',
            color: COLORS.GRAY_900,
          },
        },
      },
      axes: {
        category: {
          line: {
            stroke: COLORS.GRAY_400,
          },
          label: {
            enabled: true,
            fontSize: 11,
            fontWeight: 450,
            fontFamily: 'Inter',
            color: COLORS.GRAY_700,
          },
          tick: {
            stroke: COLORS.GRAY_700,
          },
        },
        number: {
          line: {
            stroke: COLORS.GRAY_400,
          },
          label: {
            enabled: true,
            fontSize: 11,
            fontWeight: 450,
            fontFamily: 'Inter',
            color: COLORS.GRAY_700,
          },
          tick: {
            stroke: COLORS.GRAY_700,
          },
          gridLine: {
            style: [
              {
                stroke: COLORS.GRAY_400,
              },
            ],
          },
        },
      },
    },
  },
};
