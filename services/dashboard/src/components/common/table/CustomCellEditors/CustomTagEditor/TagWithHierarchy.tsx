import { FC } from 'react';
import TagChip from 'components/common/table/CustomCellEditors/CustomTagEditor/TagChip';
import { getTagLabel, getTagParents } from 'components/filter/filter.utils';

type TagWithHierarchyProps = { tag: string; labelColor?: string };

const TagWithHierarchy: FC<TagWithHierarchyProps> = ({ tag, labelColor }) => {
  return (
    <div className='p-1 mx-1 hover:bg-GRAY_100 rounded-md space-y-1'>
      <div className='flex items-center'>
        <TagChip item={getTagLabel(tag)} externalColor={labelColor} />
      </div>
      <div className='f-11-400 text-GRAY_700 ml-1'>{getTagParents(tag)}</div>
    </div>
  );
};

export default TagWithHierarchy;
