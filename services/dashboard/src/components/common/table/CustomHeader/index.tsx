import { FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { ColDef, ColumnHeaderClickedEvent, ColumnResizedEvent } from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import AddTag from 'modules/data/AddTag';
import { getColumnOrderingVisibilityForCurrentDataset, updateLocalStorage } from 'modules/data/data.utils';
import { DatasetFilterConfigMetadataType, DatasetUpdateResponseType } from 'types/api/dataset.types';
import { SIZE_TYPES } from 'types/common/components';
import { MapAny } from 'types/commonTypes';
import { ICON_POSITION_TYPES } from 'types/components/button.type';
import { OrderType } from 'types/components/table.type';
import { cn } from 'utils/common';
import { Button } from 'components/common/button/Button';
import PositionedMenuWrapper from 'components/common/PositionedMenuWrapper';
import { CustomHeaderMenuOptions } from 'components/common/table/CustomHeader/customHeader.constants';
import { CustomHeaderMenuOptionTypes } from 'components/common/table/CustomHeader/customHeader.types';
import { CUSTOM_COLUMNS_TYPE } from 'components/common/table/table.types';
import { FILTER_TYPES } from 'components/filter/filter.types';
import FilterDropdownMenu from 'components/filter/filterMenu/FilterDropdownMenu';
import { useFiltersContextStore } from 'components/filter/filters.context';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

type CustomHeaderProps = {
  metadata: DatasetFilterConfigMetadataType;
  handleRulesListingSideDrawerOpen: (colId: string) => void;
  handleSuccessfulUpdate: (data: DatasetUpdateResponseType) => void;
  datasetId: string;
  tableRef: React.RefObject<AgGridReact>;
  filterType: FILTER_TYPES;
  options: string[];
  column: {
    colId: string;
    colDef: ColDef;
  };
  filterComponentProps?: MapAny;
};
const CustomHeader: FC<CustomHeaderProps> = ({
  metadata,
  handleRulesListingSideDrawerOpen,
  handleSuccessfulUpdate,
  datasetId,
  tableRef,
  filterType,
  options,
  column,
  filterComponentProps,
}) => {
  const { colId, colDef } = column;

  const menuRef = useRef<HTMLDivElement>(null);

  const {
    state: { selectedFilters },
  } = useFiltersContextStore();

  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isAddTagOpen, setIsAddTagOpen] = useState(false);
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [menuPosition, setMenuPosition] = useState<{ top: number; left: number }>({ top: 0, left: 0 });
  const lastResizedTimeRef = useRef<number | null>(null); // Track last resize time

  const filtersCount = selectedFilters ? Object.keys(selectedFilters)?.length : 0;
  const isTagColumn = metadata?.custom_type === CUSTOM_COLUMNS_TYPE.TAG;
  const sortState = tableRef?.current?.api?.getColumn(colId)?.getSort();
  const isFilterActive = tableRef?.current?.api?.getColumn(colId)?.isFilterActive();

  const filteredMenuOptions = useMemo(
    () =>
      CustomHeaderMenuOptions.filter((option) => {
        if (option.value === CustomHeaderMenuOptionTypes.REMOVE_SORT) {
          return !!sortState;
        }

        return option.value === CustomHeaderMenuOptionTypes.RULES ? isTagColumn : true;
      }),
    [isTagColumn, sortState],
  );

  const handleMenuClose = () => setIsMenuOpen(false);

  const handleMenuOptionClick = (option: CustomHeaderMenuOptionTypes) => {
    handleMenuClose();

    switch (option) {
      case CustomHeaderMenuOptionTypes.RULES:
        handleRulesListingSideDrawerOpen(colId);
        break;
      case CustomHeaderMenuOptionTypes.ADD_TAG:
        setIsAddTagOpen(true);
        break;
      case CustomHeaderMenuOptionTypes.SORT_ASC:
        tableRef?.current?.api?.applyColumnState({
          state: [{ colId: colId, sort: OrderType.ASC }],
        });
        break;
      case CustomHeaderMenuOptionTypes.SORT_DESC:
        tableRef?.current?.api?.applyColumnState({
          state: [{ colId: colId, sort: OrderType.DESC }],
        });
        break;
      case CustomHeaderMenuOptionTypes.FILTER:
        setIsFilterOpen(true);
        break;
      case CustomHeaderMenuOptionTypes.REMOVE_SORT:
        tableRef?.current?.api?.applyColumnState({
          state: [{ colId: colId, sort: null }],
        });
        break;
    }
  };

  const handleAddTagClose = () => {
    setIsAddTagOpen(false);
  };

  // Function to calculate and update menu position
  const updateMenuPosition = () => {
    if (!menuRef.current) return;

    const rect = menuRef.current.getBoundingClientRect();

    setMenuPosition({
      top: rect.bottom + window.scrollY, // Stick below the header
      left: rect.left, // Adjust for AG Grid's horizontal scroll
    });
  };

  // Function to open menu and set position
  const toggleMenu = useCallback(
    (event: ColumnHeaderClickedEvent) => {
      if (event.column?.getId() !== colId) return;
      const currentTime = Date.now();
      const lastResizedTime = lastResizedTimeRef.current;

      if (lastResizedTime !== null) {
        const timeDifference = currentTime - lastResizedTime;

        if (timeDifference <= 100) {
          return; // Suppress further handling
        }
      }

      tableRef?.current?.api?.clearCellSelection();
      tableRef?.current?.api?.clearFocusedCell();

      updateMenuPosition();
      setIsMenuOpen((prev) => !prev);
    },
    [colId, filterType],
  );

  const handleFilterClose = () => {
    setIsFilterOpen(false);
  };

  const handleColumnResizing = useCallback(
    (event: ColumnResizedEvent) => {
      if (event.column?.getId() !== colId) return;
      const columnOrderingVisibility = getColumnOrderingVisibilityForCurrentDataset(datasetId);
      const columnOrderingVisibilityIndex = columnOrderingVisibility.findIndex((column) => column.colId === colId);

      columnOrderingVisibility[columnOrderingVisibilityIndex].width = event.column?.getActualWidth();
      updateLocalStorage(columnOrderingVisibility, datasetId);
      lastResizedTimeRef.current = Date.now(); // Update the timestamp
    },
    [colId],
  );

  // Track column resize
  useEffect(() => {
    tableRef?.current?.api?.addEventListener('columnResized', handleColumnResizing);

    return () => {
      tableRef?.current?.api?.removeEventListener('columnResized', handleColumnResizing);
    };
  }, [colId, tableRef, handleColumnResizing]);

  // Track column header clicked
  useEffect(() => {
    tableRef?.current?.api?.addEventListener('columnHeaderClicked', toggleMenu);

    return () => {
      tableRef?.current?.api?.removeEventListener('columnHeaderClicked', toggleMenu);
    };
  }, [colId, tableRef, toggleMenu]);

  return (
    <div ref={menuRef} className='w-full h-full -mx-4 flex-1 relative'>
      <div
        className={cn(
          'w-full h-full flex-1 hover:bg-BACKGROUND_GRAY_1 cursor-pointer flex items-center justify-between px-2 group pt-5 pb-1',
          { 'bg-BACKGROUND_GRAY_1': isMenuOpen },
        )}
      >
        <div className='flex items-center gap-1 truncate self-stretch flex-auto'>
          <span className='truncate'>{colDef?.headerName ?? colId}</span>
          {!!sortState && (
            <span>
              <SvgSpriteLoader
                id={sortState === OrderType.ASC ? 'arrow-narrow-up' : 'arrow-narrow-down'}
                width={12}
                height={12}
                color={COLORS.BLUE_700}
              />
            </span>
          )}
          {isFilterActive && (
            <span>
              <SvgSpriteLoader id='filter-lines' width={12} height={12} color={COLORS.BLUE_700} />
            </span>
          )}
        </div>
        <SvgSpriteLoader id='chevron-down' width={12} height={12} className='ml-2.5' />
      </div>
      {isMenuOpen && (
        <PositionedMenuWrapper
          id='custom-header-menu'
          className='w-52 p-1'
          childrenWrapperClassName='!overflow-auto'
          menuPosition={menuPosition}
          onClose={handleMenuClose}
        >
          {filteredMenuOptions.map((option) => (
            <div
              key={option.value}
              className='flex items-center gap-1.5 px-2.5 py-2 hover:bg-GRAY_100 cursor-pointer rounded-md'
              onClick={(e) => {
                e.stopPropagation();
                handleMenuOptionClick(option.value);
              }}
            >
              <SvgSpriteLoader id={option.iconId} width={12} height={12} />
              <div className='f-12-500'>{option.label}</div>
            </div>
          ))}
          {isTagColumn && (
            <div className='px-2.5 py-3'>
              <Button
                id='add-tag-button'
                iconProps={{ id: 'tag-01', iconCategory: ICON_SPRITE_TYPES.FINANCE_AND_ECOMMERCE }}
                size={SIZE_TYPES.SMALL}
                className='w-full'
                iconPosition={ICON_POSITION_TYPES.LEFT}
                onClick={() => handleMenuOptionClick(CustomHeaderMenuOptionTypes.ADD_TAG)}
              >
                Add Tag
              </Button>
              {!!filtersCount && (
                <div className='f-11-400 text-GRAY_700 mt-1.5'>
                  {Object.keys(selectedFilters)?.length} filters applied
                </div>
              )}
            </div>
          )}
        </PositionedMenuWrapper>
      )}
      {isAddTagOpen && (
        <PositionedMenuWrapper
          id='custom-header-add-tag-menu'
          childrenWrapperClassName='!overflow-visible !max-h-[380px]'
          menuPosition={menuPosition}
          onClose={handleAddTagClose}
        >
          <AddTag
            tagList={options?.filter((option) => !!option)}
            datasetId={datasetId}
            handleSuccessfulUpdate={handleSuccessfulUpdate}
            column={colId}
            onClose={handleAddTagClose}
          />
        </PositionedMenuWrapper>
      )}
      {isFilterOpen && (
        <PositionedMenuWrapper
          id='custom-header-filter-menu'
          className='border-none'
          childrenWrapperClassName='!overflow-visible'
          menuPosition={menuPosition}
          onClose={handleFilterClose}
        >
          <FilterDropdownMenu
            filterKey={colId}
            label={colDef?.headerName}
            filterType={filterType}
            {...(filterType === FILTER_TYPES.TAGS
              ? {
                  filterComponentProps,
                }
              : {})}
          />
        </PositionedMenuWrapper>
      )}
    </div>
  );
};

export default CustomHeader;
