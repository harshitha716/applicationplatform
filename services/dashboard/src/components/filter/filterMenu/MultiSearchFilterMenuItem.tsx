import { FC, useCallback, useEffect, useRef, useState } from 'react';
import { KEYBOARD_KEYS } from 'constants/shortcuts';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType, MapAny } from 'types/commonTypes';
import { camelCaseToNormalText, cn, debounce } from 'utils/common';
import Input from 'components/common/input';
import { FILTER_TYPES } from 'components/filter/filter.types';
import DescriptionOperatorsDropdown from 'components/filter/filterMenu/components/DescriptionOperatorsDropdown';
import SearchTags, { DESCRIPTION_TAGS } from 'components/filter/filterMenu/components/SearchTags';
import { CONDITION_OPERATOR_TYPE, OPERATOR } from 'components/filter/filters.constants';
import { filtersContextActions, useFiltersContextStore } from 'components/filter/filters.context';

interface MultiSearchFilterMenuItemProps {
  column: { colId: string };
  handleClose?: defaultFnType;
  id?: string;
  isOpen?: boolean;
  forView?: string;
  label?: string;
  placeholder?: string;
  className?: string;
}

const MultiSearchFilterMenuItem: FC<MultiSearchFilterMenuItemProps> = ({
  column,
  isOpen,
  forView = 'table_header',
  className = '',
}) => {
  const columnId = column?.colId;
  const inputWrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const {
    state: { selectedFilters },
    dispatch,
  } = useFiltersContextStore();

  const [selectedOperator, setSelectedOperator] = useState<MapAny>(OPERATOR.ArrayContains);

  const ArraySearchFilter = selectedFilters[columnId];

  const isContainsOperator =
    selectedOperator?.value === CONDITION_OPERATOR_TYPE.ARRAY_IN ||
    selectedOperator?.value === CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS ||
    ArraySearchFilter?.operator?.value === CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS ||
    ArraySearchFilter?.operator?.value === CONDITION_OPERATOR_TYPE.ARRAY_IN;

  const [searchTags, setSearchTags] = useState<MapAny[]>(ArraySearchFilter?.descriptionTags ?? []);
  const [inputValue, setInputValue] = useState(ArraySearchFilter?.values ?? '');
  const [descriptionPropertySearch, setDescriptionPropertySearch] = useState('');

  const setFilter = (operator: string, searchValue: string, descriptionTags: MapAny[]) => {
    dispatch({
      type: filtersContextActions.SET_SELECTED_FILTERS,
      payload: {
        selectedFilters: {
          [columnId]: searchValue
            ? {
                filterType: FILTER_TYPES.ARRAY_SEARCH,
                type: operator,
                value: searchValue,
                descriptionTags: descriptionTags,
              }
            : {},
        },
      },
    });
  };

  const debouncedHandleSetFilters = useCallback(debounce(setFilter, 800), []);

  const updateConditionForValue = (value: string) => {
    const update = {
      operator: selectedOperator ?? OPERATOR.ContainsOperator,
      values: value,
    };

    debouncedHandleSetFilters(
      selectedOperator?.value ?? CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS,
      update.values,
      searchTags,
    );
  };

  const updateCondition = (newDescriptionTags: MapAny[]) => {
    const update = {
      operator: selectedOperator ?? OPERATOR.ContainsOperator,
      values: newDescriptionTags
        ?.filter((tag) => tag?.type === DESCRIPTION_TAGS.DESCRIPTION_VALUE)
        ?.map((tag) => tag?.label)
        ?.join(','),
      descriptionTags: newDescriptionTags,
    };

    debouncedHandleSetFilters(
      selectedOperator?.value ?? CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS,
      update?.values,
      newDescriptionTags,
    );
  };

  const updateOperator = (operator: MapAny) => {
    setSelectedOperator(operator);
    setInputValue('');
    setDescriptionPropertySearch('');
    setSearchTags([]);

    if (!inputValue && !searchTags?.length && forView === 'table_header') {
      return;
    }

    if (ArraySearchFilter?.values?.length) {
      const update = {
        operator,
        values: '',
        descriptionTags: [],
      };

      debouncedHandleSetFilters(operator?.value ?? CONDITION_OPERATOR_TYPE.ARRAY_CONTAINS, update.values, []);
    }
  };

  const handleDescriptionInputKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    event.stopPropagation();

    if (event.key !== 'Backspace' || !isContainsOperator) {
      return;
    }

    if (inputValue === '' && descriptionPropertySearch === '' && !!searchTags?.length) {
      onDeleteDescriptionInputTag(searchTags?.length - 1);
    }
  };

  const handleDescriptionInputKeyPress = (event: React.KeyboardEvent<HTMLInputElement>) => {
    const currentValue = (event?.target as HTMLInputElement)?.value.trim();

    if (
      ![KEYBOARD_KEYS.COMMA, KEYBOARD_KEYS.ENTER].includes(event?.code as KEYBOARD_KEYS) ||
      !isContainsOperator ||
      (!currentValue && searchTags?.length)
    ) {
      return;
    }

    event?.stopPropagation();
    event?.preventDefault();

    if (currentValue) {
      const newDescriptionTags = [...searchTags, { label: currentValue, type: DESCRIPTION_TAGS.DESCRIPTION_VALUE }];

      setSearchTags(newDescriptionTags);
      updateCondition(newDescriptionTags);
    }

    setInputValue('');
    handleResetDescriptionPropertySearch();

    return;
  };

  const descriptionValue: string = inputValue;

  useEffect(() => {
    if (!isOpen) {
      inputRef?.current?.blur();

      return;
    }

    if (!isContainsOperator) {
      setInputValue((ArraySearchFilter?.values as string) ?? '');
    }

    let newDescriptionTags = ArraySearchFilter?.descriptionTags ?? [];

    if (ArraySearchFilter?.values && isContainsOperator && !ArraySearchFilter?.descriptionTags?.length) {
      newDescriptionTags = [
        ...searchTags,
        { label: ArraySearchFilter?.values, type: DESCRIPTION_TAGS.DESCRIPTION_VALUE },
      ];

      updateCondition(newDescriptionTags);
    }

    setSearchTags(newDescriptionTags);
    setInputValue('');
    inputRef?.current?.focus();
    if (ArraySearchFilter?.operator) setSelectedOperator(ArraySearchFilter?.operator);
  }, [isOpen]);

  const handleResetDescriptionPropertySearch = () => {
    setDescriptionPropertySearch('');
  };

  const handleDescriptionValueSearch = (value: string) => {
    setInputValue(value);
    updateConditionForValue(value);
  };

  const onDescriptionInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();

    handleDescriptionValueSearch(e?.target?.value);
  };

  const onDeleteDescriptionInputTag = (index: number) => {
    if (!searchTags?.length) return;

    const newTags: MapAny[] = searchTags?.filter((_, i) => i !== index);

    setSearchTags(newTags);
    updateCondition(newTags);
  };

  return (
    <div
      className={cn(
        'px-2.5 py-2 w-[218px] min-w-[300px] max-w-[360px] border-0.5 border-GRAY_400 rounded-md bg-white shadow-tableFilterMenu',
        className,
      )}
    >
      <DescriptionOperatorsDropdown
        operator={selectedOperator}
        updateOperator={updateOperator}
        label={camelCaseToNormalText(columnId)}
      />

      <div
        className='relative mt-2 w-full rounded-md border border-BORDER_GRAY_400 overflow-hidden shadow-tableFilterMenu focus:shadow-inputOutlineShadow focus:border-GRAY_600'
        ref={inputWrapperRef}
      >
        <Input
          inputRef={inputRef}
          size={SIZE_TYPES.SMALL}
          name='description'
          placeholder='type here....'
          customTags={
            searchTags?.length ? <SearchTags tags={searchTags} onDeleteTag={onDeleteDescriptionInputTag} /> : undefined
          }
          id='description-filter-input'
          inputPillsWrapperClasses={cn('px-2 gap-3 py-3')}
          isMulti
          onKeyPress={handleDescriptionInputKeyPress}
          onKeyDown={handleDescriptionInputKeyDown}
          onDeleteTag={onDeleteDescriptionInputTag}
          overrideInputBgClassName='!bg-white'
          value={(descriptionValue ? descriptionValue : descriptionPropertySearch) as string}
          onChange={onDescriptionInputChange}
          inputClassName='w-full !min-w-[160px] flex-1 outline-none border-none focus:!shadow-none !shadow-none'
          inputSizeClassName='p-0'
          autoFocus
        />
      </div>
    </div>
    // </MenuWrapper>
  );
};

export default MultiSearchFilterMenuItem;
