import { FC, Fragment, useMemo } from 'react';
import TagChip from 'components/common/table/CustomCellEditors/CustomTagEditor/TagChip';

type CreateTagProps = {
  value: string;
  handleCreateTag: () => void;
  existingList: string[];
};

const CreateTag: FC<CreateTagProps> = ({ value, handleCreateTag, existingList }) => {
  const tags = useMemo(() => value?.split('/')?.map((item: string) => item?.trim()), [value]);

  return (
    <Fragment>
      {!!value && (
        <div
          className='mx-1 mb-1 px-2.5 py-2 space-x-1.5 text-GRAY_700 f-11-500 hover:bg-GRAY_100 rounded-md cursor-pointer flex items-center'
          onClick={handleCreateTag}
        >
          <span>Create</span>
          {tags.map((item: string, index: number) => (
            <Fragment key={item}>
              {!!item && <TagChip existingList={existingList} item={item} />}
              {index !== tags?.length - 1 && <span>/</span>}
            </Fragment>
          ))}
        </div>
      )}
    </Fragment>
  );
};

export default CreateTag;
