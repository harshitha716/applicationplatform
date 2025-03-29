import { PERIODICITY_TYPES } from 'constants/date.constants';
import { widgetData, widgetInstanceDetails } from 'modules/widgets/Pivot/__tests__/pivot.utils.mock';
import {
  AGGridPivotNode,
  concatTagFilters,
  flattenChildrenAfterGroup,
  getPivotColDefs,
  getPivotColumns,
  getPivotData,
  unwrapTagColumn,
} from 'modules/widgets/Pivot/pivot.utils';
import '@testing-library/jest-dom';

jest.mock('next/font/google', () => ({
  Inter: jest.fn(() => ({
    className: 'mocked-inter',
    variable: '--font-inter',
  })),
}));

describe('getPivotColDefs', () => {
  it('works correctly', () => {
    const pivotColumns = getPivotColumns(widgetInstanceDetails, widgetData);
    const result = getPivotColDefs(pivotColumns);

    expect(result).toBeDefined();
  });
});

describe('getPivotData', () => {
  it('works correctly', () => {
    const pivotColumns = getPivotColumns(widgetInstanceDetails, widgetData);

    const periodicity = PERIODICITY_TYPES.DAILY;
    const result = getPivotData(pivotColumns, widgetData, periodicity);

    expect(result).toBeDefined();
  });
});

describe('getPivotColDefs', () => {
  it('works correctly', () => {
    const pivotColumns = getPivotColumns(widgetInstanceDetails, widgetData);
    const result = getPivotColDefs(pivotColumns);

    expect(result).toBeDefined();
  });
});

describe('flattenChildrenAfterGroup', () => {
  it('flattens a simple tree structure', () => {
    const simpleTree: AGGridPivotNode<any> = {
      key: 'root',
      childrenAfterGroup: [
        {
          key: 'child1',
          childrenAfterGroup: [],
        },
        {
          key: 'child2',
          childrenAfterGroup: [],
        },
      ],
    };

    const result = flattenChildrenAfterGroup(simpleTree);

    expect(result).toHaveLength(2);
    expect(result.map((node) => node.key)).toEqual(['child1', 'child2']);
  });

  it('flattens a nested tree structure', () => {
    const nestedTree: AGGridPivotNode<any> = {
      key: 'root',
      childrenAfterGroup: [
        {
          key: 'child1',
          childrenAfterGroup: [
            {
              key: 'grandchild1',
              childrenAfterGroup: [],
            },
          ],
        },
        {
          key: 'child2',
          childrenAfterGroup: [
            {
              key: 'grandchild2',
              childrenAfterGroup: [],
            },
          ],
        },
      ],
    };

    const result = flattenChildrenAfterGroup(nestedTree);

    expect(result).toHaveLength(4);
    expect(result.map((node) => node.key)).toEqual(['child1', 'child2', 'grandchild1', 'grandchild2']);
  });

  it('returns empty array for leaf node', () => {
    const leafNode: AGGridPivotNode<any> = {
      key: 'leaf',
      childrenAfterGroup: [],
    };

    const result = flattenChildrenAfterGroup(leafNode);

    expect(result).toHaveLength(0);
  });

  it('returns empty array if childrenAfterGroup is undefined', () => {
    const node: AGGridPivotNode<any> = {
      key: 'root',
    };

    const result = flattenChildrenAfterGroup(node);

    expect(result).toHaveLength(0);
  });
});

describe('concatTagFilters', () => {
  it('should concatenate tag filters with correct hierarchy', () => {
    const filters = {
      __tag_LEVEL_1: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud'],
      },
      __tag_LEVEL_2: {
        filterType: 'search',
        type: 'startswith',
        values: ['aws'],
      },
      __tag_LEVEL_3: {
        filterType: 'search',
        type: 'startswith',
        values: ['ec2'],
      },
      otherFilter: {
        filterType: 'search',
        values: ['something'],
      },
    };

    const result = concatTagFilters(filters);

    expect(result).toEqual({
      tag: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud.aws.ec2'],
      },
      otherFilter: {
        filterType: 'search',
        values: ['something'],
      },
    });
  });

  it('should handle single level tag filter', () => {
    const filters = {
      __tag_LEVEL_1: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud'],
      },
    };

    const result = concatTagFilters(filters);

    expect(result).toEqual({
      tag: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud'],
      },
    });
  });

  it('should handle non-sequential tag levels', () => {
    const filters = {
      __tag_LEVEL_1: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud'],
      },
      __tag_LEVEL_3: {
        filterType: 'search',
        type: 'startswith',
        values: ['ec2'],
      },
    };

    const result = concatTagFilters(filters);

    expect(result).toEqual({
      tag: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud.ec2'],
      },
    });
  });

  it('should handle empty filters object', () => {
    const filters = {};
    const result = concatTagFilters(filters);

    expect(result).toEqual({});
  });

  it('should preserve non-tag filters', () => {
    const filters = {
      status: {
        filterType: 'search',
        values: ['active'],
      },
      __tag_LEVEL_1: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud'],
      },
    };

    const result = concatTagFilters(filters);

    expect(result).toEqual({
      tag: {
        filterType: 'search',
        type: 'startswith',
        values: ['cloud'],
      },
      status: {
        filterType: 'search',
        values: ['active'],
      },
    });
  });
});

describe('unwrapTagColumn', () => {
  it('works correctly', () => {
    const result = unwrapTagColumn('__tag_LEVEL_1');

    expect(result).toEqual({ name: 'tag', hierarchy: 1 });
  });

  it('works correctly', () => {
    const result = unwrapTagColumn('__tag_name_test_LEVEL_5');

    expect(result).toEqual({ name: 'tag_name_test', hierarchy: 5 });
  });
});
