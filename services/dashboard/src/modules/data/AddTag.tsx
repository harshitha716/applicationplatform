import { useMemo, useState } from 'react';
import { IServerSideGetRowsRequest } from 'ag-grid-community';
import { useUpdateDatasetDataMutation } from 'apis/dataset';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { convertFilterModelToRuleFilters } from 'modules/data/data.utils';
import RuleStatement from 'modules/data/RulesListing/RuleStatement';
import { DatasetUpdateResponseType } from 'types/api/dataset.types';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { Button } from 'components/common/button/Button';
import Input from 'components/common/input';
import { MenuWrapper } from 'components/common/MenuWrapper';
import CreateTag from 'components/common/table/CustomCellEditors/CustomTagEditor/CreateTag';
import TagWithHierarchy from 'components/common/table/CustomCellEditors/CustomTagEditor/TagWithHierarchy';
import { convertToFilterModel, getFilterModelFromGroupAndFilterModel } from 'components/common/table/table.utils';
import ToggleSwitch from 'components/common/toggleSwitch';
import { getFilterStatementValues, getTagLabel } from 'components/filter/filter.utils';
import { useFiltersContextStore } from 'components/filter/filters.context';
import SvgSpriteLoader from 'components/SvgSpriteLoader';
const fieldOperatorClassName = 'text-GRAY_1000 pl-1.5 pr-2 py-1';

const AddTag = ({
  datasetId,
  handleSuccessfulUpdate,
  tagList,
  column,
  onClose,
}: {
  datasetId: string;
  handleSuccessfulUpdate: (data: DatasetUpdateResponseType) => void;
  tagList: string[];
  column: string;
  onClose: defaultFnType;
}) => {
  const [isActive, setIsActive] = useState(false);
  const [searchValue, setSearchValue] = useState<string>('');
  const [searchResults, setSearchResults] = useState<string[]>(tagList);
  const [selectedTag, setSelectedTag] = useState<string>('');
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [updateDatasetData, { isLoading }] = useUpdateDatasetDataMutation();

  const {
    state: { selectedFilters, totalRows },
  } = useFiltersContextStore();

  const filterStatement = useMemo(
    () => getFilterStatementValues(convertFilterModelToRuleFilters(convertToFilterModel(selectedFilters))),
    [selectedFilters],
  );

  const handleClickAddTag = () => {
    updateDatasetData({
      datasetId: datasetId,
      data: {
        filters: getFilterModelFromGroupAndFilterModel({ filterModel: selectedFilters } as IServerSideGetRowsRequest),
        update: {
          column: column,
          value: selectedTag
            ?.trim()
            ?.split('/')
            ?.map((str) => str?.trim())
            ?.join('.'),
        },
        save_as_rule: isActive,
      },
    })
      .unwrap()
      .then((data) => {
        onClose();
        handleSuccessfulUpdate(data);
      });
  };

  const handleTagClick = (tag: string) => {
    setSelectedTag(tag);
    setSearchValue(getTagLabel(tag));
    setIsOpen(false);
  };

  const handleCreateTag = () => {
    setSelectedTag(searchValue);
    setSearchValue(
      searchValue
        .trim()
        .split('/')
        .map((str) => str.trim())
        .pop() || '',
    );
    setIsOpen(false);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;

    if (value?.includes('.')) return;

    setSearchValue(value);
    setSearchResults(tagList?.filter((tag) => tag?.toLowerCase()?.includes(value?.toLowerCase())));
    if (value) {
      setIsOpen(true);
    }
  };

  return (
    <div className='w-[300px]'>
      <div className='py-3'>
        <div className='flex items-center justify-between mb-3.5 px-3'>
          <div className='f-12-500 text-GRAY_1000'>Add Tag</div>
          <SvgSpriteLoader
            id='x-close'
            iconCategory={ICON_SPRITE_TYPES.GENERAL}
            width={12}
            height={12}
            onClick={onClose}
            className='cursor-pointer'
          />
        </div>
        <div className='px-4'>
          <Input placeholder='Search Tag' onChange={handleChange} value={searchValue} onFocus={() => setIsOpen(true)} />
          {isOpen && (
            <MenuWrapper
              id='custom-tag-editor-menu'
              className='!fixed mt-1 w-64 z-10'
              childrenWrapperClassName='!overflow-y-visible !max-h-fit'
            >
              <div className='text-GRAY_700 f-11-500 p-2'>Select an option or create one</div>
              <div className='space-y-1 my-1 overflow-y-auto max-h-[300px]'>
                {searchResults.map((tag: string) => (
                  <div key={tag} onClick={() => handleTagClick(tag)}>
                    <TagWithHierarchy tag={tag} />
                  </div>
                ))}
              </div>
              <CreateTag value={searchValue} handleCreateTag={handleCreateTag} existingList={tagList} />
              <div className='flex items-center p-2 bg-BG_GRAY_2 gap-2 rounded-b-md'>
                <span>💡</span>
                <span className='text-GRAY_900 f-11-400'>Use “ / “ to create hierarchy</span>
              </div>
            </MenuWrapper>
          )}
          {filterStatement.length > 0 && (
            <>
              <div
                className='rounded-md bg-BG_GRAY_2 px-3 py-2.5 f-11-400 text-GRAY_1000 border border-BORDER_GRAY_400 my-2.5 h-[140px] overflow-y-auto flex flex-wrap gap-y-2 items-center'
                style={{ scrollbarWidth: 'none' }}
              >
                <span className={fieldOperatorClassName}>If</span>
                {filterStatement.map((value, index) => (
                  <RuleStatement
                    index={index}
                    filterStatement={value}
                    numberOfFilters={filterStatement.length}
                    key={`filter-statement-${index}`}
                  />
                ))}
              </div>
              <div className='flex items-center gap-1.5 mb-1.5'>
                <ToggleSwitch id='add-tag-make-rule' onChange={setIsActive} checked={isActive} />
                <div className='f-11-400 text-GRAY_1000'>Make this a rule</div>
              </div>
              <div className='f-11-400 text-GRAY_700 text-wrap'>
                Rule applies selected tags to all transactions meeting its criteria, historical & future, replacing any
                existing tags
              </div>
            </>
          )}
        </div>
      </div>
      <div className='flex flex-row-reverse items-center justify-between px-4 py-3 border-t border-BORDER_GRAY_400'>
        <Button
          size={SIZE_TYPES.XSMALL}
          id='add-tag-transactions'
          onClick={handleClickAddTag}
          isLoading={isLoading}
          disabled={!selectedTag}
        >
          Add tag to {totalRows} transactions
        </Button>
      </div>
    </div>
  );
};

export default AddTag;
