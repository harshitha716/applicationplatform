import React, { useMemo } from 'react';
import { ICellRendererParams } from 'ag-grid-community';
import TagChip from 'components/common/table/CustomCellEditors/CustomTagEditor/TagChip';
import { getTagLabel } from 'components/filter/filter.utils';

const CustomTagRenderer = (props: ICellRendererParams) => {
  const { value, colDef } = props;
  const tagColorMap = colDef?.cellRendererParams?.tagColorMap;

  const tag = useMemo(() => getTagLabel(value), [value]);

  return tag ? <TagChip item={tag} externalColor={tagColorMap?.[value]} /> : <></>;
};

export default CustomTagRenderer;
