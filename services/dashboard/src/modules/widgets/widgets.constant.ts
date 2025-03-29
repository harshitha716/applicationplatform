import { AgCartesianAxisOptions, time } from 'ag-charts-community';
import { CHART_PALETTE_COLORS, COLORS } from 'constants/colors';
import { DATE_FORMATS, PERIODICITY_TYPES } from 'constants/date.constants';
import { endOfWeek, format, startOfWeek } from 'date-fns';
import { WIDGET_TYPES } from 'types/api/widgets.types';
import { MapAny } from 'types/commonTypes';
import { formatNumber, isValidDate, snakeCaseToSentenceCase, trimString } from 'utils/common';

export enum SCREEN_BREAKPOINTS_NAMES {
  SM = 'sm',
  MD = 'md',
  LG = 'lg',
  XL = 'xl',
  DEFAULT = 'default',
}

export const SCREEN_BREAKPOINTS = { lg: 1200, md: 996, sm: 768, xs: 480, xxs: 0 };

export const ROW_HEIGHT = 56; // Height of a single row in px
export const WIDGETS_LAYOUT_MARGIN = [20, 20]; // Space between components (20px)
export const MAX_DONUT_CHART_SLICE_COUNT = 5;

export enum WidgetDataValueType {
  STRING = 'STRING',
  DECIMAL = 'DECIMAL',
  NUMBER = 'NUMBER',
  BIGINT = 'BIGINT',
  DOUBLE = 'DOUBLE',
  BOOLEAN = 'BOOLEAN',
  FLOAT = 'FLOAT',
  SMALLINT = 'SMALLINT',
  TINYINT = 'TINYINT',
  INT = 'INT',
  DATE = 'DATE',
  TIMESTAMP = 'TIMESTAMP',
  TIME = 'TIME',
  LONG = 'LONG',
  DATETIME = 'DATETIME',
}

export const AG_CHART_TYPES = {
  [WIDGET_TYPES.BAR_CHART]: 'bar',
  [WIDGET_TYPES.LINE_CHART]: 'line',
  [WIDGET_TYPES.PIE_CHART]: 'pie',
  [WIDGET_TYPES.DONUT_CHART]: 'donut',
};

export enum CHART_SLICE_TYPES {
  OTHERS = 'others',
}

export const getFormattedDateWithPeriodicity = (periodicity: PERIODICITY_TYPES, date: string) => {
  switch (periodicity) {
    case PERIODICITY_TYPES.DAILY: {
      return format(new Date(date), DATE_FORMATS.ddMMMyyyy);
    }
    case PERIODICITY_TYPES.WEEKLY: {
      const start = startOfWeek(new Date(date), { weekStartsOn: 1 }); // Monday as start of week
      const end = endOfWeek(new Date(date), { weekStartsOn: 1 });

      return `${format(start, DATE_FORMATS.DD)}-${format(end, DATE_FORMATS.d_MMM_yyyy)}`;
    }
    case PERIODICITY_TYPES.MONTHLY:
      return format(new Date(date), DATE_FORMATS.MMM_yyyy);
    case PERIODICITY_TYPES.QUARTERLY:
      return format(new Date(date), DATE_FORMATS.QQ_yyyy);
    case PERIODICITY_TYPES.YEARLY:
      return format(new Date(date), DATE_FORMATS.YYYY);
  }
};

export const getCategoryAxis = (periodicity: PERIODICITY_TYPES) => {
  return {
    type: 'category' as const,
    position: 'bottom',
    crosshair: {
      enabled: false,
    },
    label: {
      minSpacing: 20,
      autoRotate: false,
      formatter: function (params: MapAny) {
        if (isValidDate(params.value)) {
          return getFormattedDateWithPeriodicity(periodicity, params.value);
        }

        return params.value;
      },
    },
    tick: {
      size: 10, // Changed from length to size
      width: 0.75,
    },
    line: {
      width: 1,
      stroke: COLORS.GRAY_400,
    },
  };
};

export const CHART_NUMBER_AXES: AgCartesianAxisOptions = {
  type: 'number' as const,
  position: 'right',
  crosshair: {
    enabled: false,
  },
  gridLine: {
    enabled: true,
    width: 0.5,
    style: [
      {
        stroke: COLORS.GRAY_100,
      },
    ],
  },
  label: {
    formatter: ({ value }) => {
      return formatNumber(value, 2, false);
    },
  },
};

export const AG_CHART_TIME_AXES: AgCartesianAxisOptions = {
  type: 'time',
  nice: false,
  position: 'bottom',
  interval: { step: time.month },
  label: {
    format: '%d %b',
  },
  tick: {
    size: 10, // Changed from length to size
    width: 0.75,
  },
  line: {
    width: 1,
    stroke: COLORS.GRAY_400,
  },
};

export const AG_CHART_LEGEND_CONFIG = {
  enabled: true,
  item: {
    showSeriesStroke: false,
    paddingX: 16,
    marker: {
      size: 8,
      shape: 'square' as const,
      strokeWidth: 0,
      padding: 6,
    },
    label: {
      fontFamily: 'Inter',
      fontWeight: 450,
      fontSize: 12,
      color: COLORS.GRAY_900,
      formatter: ({ value = '' }) => trimString(snakeCaseToSentenceCase(value), 20),
    },
  },
};

export const getDonutChartSeriesConfig = (dataLength: number) => {
  return {
    innerRadiusRatio: 0.75,
    sectorSpacing: dataLength > 1 ? 3 : 0,
    cornerRadius: dataLength > 1 ? 2 : 0,
    label: {
      color: COLORS.GRAY_950,
      fontSize: 26,
    },
    calloutLine: {
      length: 18,
      strokeWidth: 2,
      colors: [COLORS.GRAY_400],
    },
    fills: CHART_PALETTE_COLORS,
    innerCircle: {
      fill: COLORS.WHITE,
    },
  };
};

export const DEFAULT_TRANSFORMED_DATA = {
  transformedData: [],
  stackedValues: [],
  donutOthersData: [],
  yAxisTitle: '',
  maxValueLength: 0,
  showCurrency: false,
};
