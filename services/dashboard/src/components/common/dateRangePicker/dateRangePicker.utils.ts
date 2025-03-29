import {
  DATE_FILTER_CATEGORIES,
  DATE_RANGE_TYPES,
  dateFilterValueType,
  DateRangeValue,
  DaysConfig,
  MonthsConfig,
  QuartersConfig,
  RangeType,
} from 'constants/date.constants';
import { endOfMonth, endOfQuarter, endOfYear, startOfMonth, startOfQuarter, startOfYear, sub } from 'date-fns';
import { OptionsType } from 'types/commonTypes';
import { getStartOfYear, isUserInUS } from 'utils/common';

export const getYearList = () => {
  const currentYear = new Date().getFullYear();
  const yearList = [];

  // Add five years before the current year
  for (let i = currentYear - 10; i < currentYear; i++) {
    yearList.push(i);
  }

  // Add the current year
  yearList.push(currentYear);

  return yearList;
};

export function searchDateRange(input: string): DateRangeValue | null {
  const today = new Date(); // Get the current date

  const { direction, unit, offset } = parseRelativeDatePhrase(input);
  const negativeDirectionPatters = /(last|previous|back|before)/i;

  if (unit) {
    const parsedOffset = offset !== null ? offset : 1; // Treat undefined offset as 1

    const date = new Date(today);

    const isNegativeDirection = direction ? negativeDirectionPatters.test(direction) : true;

    switch (unit.toLowerCase()) {
      case 'day':
      case 'days':
        date.setDate(today.getDate() + (isNegativeDirection ? -parsedOffset : parsedOffset));

        return {
          type: DATE_RANGE_TYPES.DAY,
          value: date,
          year: date.getFullYear(),
          label: date.toDateString(),
        };
      case 'month':
      case 'months':
        date.setMonth(today.getMonth() + (isNegativeDirection ? -parsedOffset : parsedOffset));

        return {
          type: DATE_RANGE_TYPES.MONTH,
          value: date.getMonth(),
          year: date.getFullYear(),
          label: MonthsConfig[date.getMonth()].label,
        };
      case 'quarter':
      case 'quarters':
        date.setMonth(today.getMonth() + (isNegativeDirection ? -parsedOffset * 3 : parsedOffset * 3));

        return {
          type: DATE_RANGE_TYPES.QUARTER,
          value: Math.floor(date.getMonth() / 3),
          year: date.getFullYear(),
          label: QuartersConfig[Math.floor(date.getMonth() / 3)].short,
        };
      case 'year':
      case 'years':
        date.setFullYear(today.getFullYear() + (isNegativeDirection ? -parsedOffset : parsedOffset));

        return {
          type: DATE_RANGE_TYPES.YEAR,
          value: getStartOfYear(date.getFullYear()),
          year: date.getFullYear(),
          label: date.getFullYear().toString(),
        };
      default:
        return null;
    }
  }

  return null;
}

function parseRelativeDatePhrase(phrase: string): {
  direction: string | null;
  unit: string | null;
  offset: number | null;
} {
  const directionPattern = /(last|previous|next|later|after|back|before|this|prev)/i;
  const lastPattern = /(last|previous|back|before|prev)/i;
  const thisPattern = /\bthis\b/i;
  const dayPattern = /\b(sun|mon|tue|wed|thu|fri|sat)(?:day)?\b/i;

  const monthPattern =
    /(\b(?:January|February|March|April|May|June|July|August|September|October|November|December|Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\b)(?:\s+(\d{4}))?/i;

  const unitPattern = /(quarters?|months?|days?|years?|weeks?)/i;
  const offsetPattern = /\b\d+(st|nd|rd|th)?\b/;

  let direction: string | null = null;
  let unit: string | null = null;
  let offset: number | null = null;

  // Match direction
  const directionMatch = phrase.match(directionPattern);

  if (directionMatch) {
    direction = lastPattern.test(directionMatch[0]) ? 'last' : thisPattern.test(directionMatch[0]) ? 'this' : 'next';
    phrase = phrase.replace(directionPattern, '').trim();
  }

  if (phrase.toLowerCase() === 'tomorrow') {
    return { direction: 'next', unit: 'day', offset: 1 };
  } else if (phrase.toLowerCase() === 'yesterday') {
    return { direction: 'last', unit: 'day', offset: 1 };
  } else if (phrase.toLowerCase() === 'today') {
    return { direction: 'next', unit: 'day', offset: 0 };
  }

  // Handle "this" phrases
  if (direction && direction.toLowerCase() === 'this') {
    // Check for "this year", "this month", or "this quarter"
    const thisYearMatch = phrase.match(/\b(?:this\s+)?year\b/i);
    const thisMonthMatch = phrase.match(/\b(?:this\s+)?month\b/i);
    const thisQuarterMatch = phrase.match(/\b(?:this\s+)?quarter\b/i);

    if (thisYearMatch) {
      unit = 'year';
      offset = 0;
      direction = null;

      return { direction, unit, offset };
    } else if (thisMonthMatch) {
      unit = 'month';
      offset = 0;
      direction = null;

      return { direction, unit, offset };
    } else if (thisQuarterMatch) {
      unit = 'quarter';
      offset = 0;
      direction = null;

      return { direction, unit, offset };
    }
  }

  // Match specific dates in the format "DD/MM/YYYY" or "MM/DD/YYYY"
  const specificDateMatch = phrase.match(/(\d{1,2})\/(\d{1,2})\/(\d{4})/);

  if (specificDateMatch) {
    let day = parseInt(specificDateMatch[1]);
    let month = parseInt(specificDateMatch[2]) - 1;
    let year = parseInt(specificDateMatch[3]);

    if (isUserInUS()) {
      day = parseInt(specificDateMatch[2]);
      month = parseInt(specificDateMatch[1]) - 1;
      year = parseInt(specificDateMatch[3]);
    }

    const parsedDate = new Date(year, month, day);
    const currentDate = new Date();

    offset = Math.floor((parsedDate.getTime() - currentDate.getTime()) / (1000 * 60 * 60 * 24));

    direction = offset >= 0 ? 'next' : 'last';

    unit = 'day';

    return { direction, unit, offset: Math.abs(offset) - 1 };
  }

  // Match quarters in the format "Q1", "Q2", "Q3", or "Q4" followed by a year
  const quarterYearMatch = phrase.match(/Q([1-4])(?:\s+(\d{4}))?/i);

  if (quarterYearMatch) {
    const quarter = parseInt(quarterYearMatch[1]) - 1; // Convert quarter to 0-based index
    let year: number | null = null;

    if (quarterYearMatch[2]) {
      year = parseInt(quarterYearMatch[2]);
    } else {
      // If year is not provided, default to the current year
      year = new Date().getFullYear();
    }
    // Calculate offset relative to the current quarter
    const currentQuarter = Math.floor((new Date().getMonth() + 3) / 3) - 1;
    const currentYear = new Date().getFullYear();

    offset = (year - currentYear) * 4 + (quarter - currentQuarter);
    // Set appropriate direction based on offset
    direction = offset >= 0 ? 'next' : 'last';

    return { direction, unit: 'quarter', offset: Math.abs(offset) };
  }

  // Match months in the format "MonthName" followed by an optional year
  const monthMatch = phrase.match(monthPattern);
  const offsetMatch = phrase.match(offsetPattern);
  const digits = offsetMatch && offsetMatch[0].match(/\d+/);

  if (monthMatch) {
    const monthIndex = MonthsConfig.findIndex(
      (month) =>
        month.label.toLowerCase() === monthMatch[1].toLowerCase() ||
        month?.short?.toLowerCase() === monthMatch[1].toLowerCase(),
    );
    let year: number | null = null;

    if (monthMatch[2]) {
      year = parseInt(monthMatch[2]);
    } else {
      // If year is not provided, default to the current year
      year = new Date().getFullYear();
    }
    const currentMonth = new Date().getMonth();
    const currentYear = new Date().getFullYear();

    let offset = (year - currentYear) * 12 + (monthIndex - currentMonth);

    // Set appropriate direction based on offset
    direction = offset >= 0 ? 'next' : 'last';

    if (!digits || digits[0].length > 2) {
      return { direction, unit: 'month', offset: Math.abs(offset) };
    } else {
      const parsedDate = new Date(year, monthIndex, Number(digits[0]));
      const currentDate = new Date();

      offset = Math.floor((parsedDate.getTime() - currentDate.getTime()) / (1000 * 60 * 60 * 24));

      direction = offset >= 0 ? 'next' : 'last';

      unit = 'day';

      return { direction, unit, offset: Math.abs(offset) - (direction === 'next' ? -1 : 1) };
    }
  }

  // Match years in the format "YYYY"
  const yearMatch = phrase.match(/\b\d{4}\b/);

  if (yearMatch) {
    unit = 'year';
    offset = parseInt(yearMatch[0]) - new Date().getFullYear();

    direction = offset >= 0 ? 'next' : 'last';

    return { direction, unit, offset: Math.abs(offset) };
  }

  // Match day of the week
  const dayMatch = phrase.match(dayPattern);

  if (dayMatch) {
    const targetDay = dayMatch[0].toLowerCase();
    const today = new Date();
    let currentDayIndex = today.getDay(); // Sunday is 0, Monday is 1, ..., Saturday is 6
    const targetDayIndex = DaysConfig.findIndex(
      (day) => day.short.toLowerCase() === targetDay || day.label.toLowerCase() === targetDay,
    );

    let finalOffset = 0; // Start with 0 offset

    direction = !direction ? 'last' : direction;
    if (direction && direction.toLowerCase() === 'last') {
      if (currentDayIndex === targetDayIndex) {
        finalOffset = 7;
      } else {
        while (currentDayIndex !== targetDayIndex) {
          currentDayIndex--;
          if (currentDayIndex < 0) currentDayIndex = 6;
          finalOffset++;
        }
      }
    } else if (direction && direction.toLowerCase() === 'next') {
      if (currentDayIndex === targetDayIndex) {
        finalOffset = 7;
      } else {
        while (currentDayIndex !== targetDayIndex) {
          currentDayIndex++;
          if (currentDayIndex > 6) currentDayIndex = 0;
          finalOffset++;
        }
      }
    }

    unit = 'day';

    return { direction, unit, offset: finalOffset };
  }

  // Match unit
  const unitMatch = phrase.match(unitPattern);

  if (unitMatch) {
    unit = unitMatch[0];
    phrase = phrase.replace(unitPattern, '').trim();
  }

  // Match offset

  if (digits) {
    offset = parseInt(digits[0]);
    const isWeekUnit = unit && /weeks?/i.test(unit);

    if (unit && (unit.toLowerCase() === 'year' || isWeekUnit)) {
      if (unit.toLowerCase() === 'year') {
        offset *= 365; // Assuming a year has 365 days
      } else if (isWeekUnit) {
        // If unit is 'week' and offset is provided, treat it as days
        offset *= 7;
      }
      unit = 'day';
    }
  } else if (phrase.toLowerCase().includes('later') || phrase.toLowerCase().includes('after')) {
    offset = 1;
  } else if (unit && /weeks?/i.test(unit)) {
    // If unit is 'week' and offset is not provided, set offset to 1 week (7 days)
    offset = 7;
    unit = 'day';
  }

  // If unit is 'year' and offset is not provided, set offset to 1
  if (unit && unit.toLowerCase() === 'year' && offset === null) {
    offset = 1;
  }

  // Set appropriate direction based on offset
  if (offset !== null) {
    direction = direction ? direction : offset >= 0 ? 'next' : 'last';
  }

  return { direction, unit, offset };
}

export const getDateRangeFromCategory = (value: OptionsType, finalDateRange: RangeType) => {
  const typedValue = value?.value as DATE_FILTER_CATEGORIES;

  const today = new Date();
  const todayEnd = new Date();

  today.setHours(0, 0, 0.0);
  todayEnd.setHours(23, 59, 59, 999);

  let updatedDateRange: dateFilterValueType = { start: finalDateRange?.startDate, end: finalDateRange?.endDate };

  switch (typedValue) {
    case DATE_FILTER_CATEGORIES.ALL_TIME:
      updatedDateRange = {
        start: undefined,
        end: undefined,
      };
      break;
    case DATE_FILTER_CATEGORIES.TODAY:
      updatedDateRange = {
        start: today,
        end: todayEnd,
      };
      break;
    case DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE:
      break;
    case DATE_FILTER_CATEGORIES.LAST_30_DAYS:
      updatedDateRange = {
        start: sub(today, { days: 29 }),
        end: today,
      };
      break;
    case DATE_FILTER_CATEGORIES.THIS_MONTH: {
      const end = endOfMonth(today);

      updatedDateRange = {
        start: startOfMonth(today),
        end: end,
      };
      break;
    }
    case DATE_FILTER_CATEGORIES.LAST_MONTH: {
      const lastMonth = sub(today, { months: 1 });

      updatedDateRange = {
        start: startOfMonth(lastMonth),
        end: endOfMonth(lastMonth),
      };
      break;
    }
    case DATE_FILTER_CATEGORIES.THIS_QUARTER: {
      const end = endOfQuarter(today);

      updatedDateRange = {
        start: startOfQuarter(today),
        end: end,
      };
      break;
    }
    case DATE_FILTER_CATEGORIES.LAST_QUARTER: {
      const dayInLastQuarter = sub(today, { months: 3 });

      updatedDateRange = {
        start: startOfQuarter(dayInLastQuarter),
        end: endOfQuarter(dayInLastQuarter),
      };
      break;
    }
    case DATE_FILTER_CATEGORIES.THIS_YEAR:
      {
        const end = endOfYear(today);

        updatedDateRange = {
          start: startOfYear(today),
          end: end,
        };
      }
      break;
    case DATE_FILTER_CATEGORIES.LAST_YEAR: {
      const lastYear = sub(today, { years: 1 });

      updatedDateRange = {
        start: startOfYear(lastYear),
        end: endOfYear(lastYear),
      };
      break;
    }
  }

  return updatedDateRange;
};

export const getPlacehoderDate = () => {
  return isUserInUS() ? '01/24/2019' : '11/01/2019';
};
