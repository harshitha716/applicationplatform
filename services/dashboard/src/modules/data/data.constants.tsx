import { ColDef, ICellRendererParams } from 'ag-grid-community';
import { COLORS } from 'constants/colors';
import { DATASET_ICON } from 'constants/icons';
import { DATASET_ACCESS_PRIVILEGES } from 'modules/data/data.types';
import Image from 'next/image';
import { cn } from 'utils/common';
import CustomTagRenderer from 'components/common/table/CustomCellRenderers/CustomTagRenderer';
import { DATA_TABLE_CONFIG } from 'components/common/table/table.constants';
import { CUSTOM_COLUMNS_TYPE } from 'components/common/table/table.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export const LISTING_COLUMNS: ColDef[] = [
  {
    field: 'title',
    headerName: 'Datasets',
    cellRenderer: (params: ICellRendererParams) => {
      return (
        <div className='flex items-center gap-2.5 f-13-500'>
          <Image src={DATASET_ICON} alt='dataset' width={20} height={20} />
          {params.value}
        </div>
      );
    },
  },
  {
    field: 'description',
    headerName: 'Description',
  },
  {
    field: 'updated_at',
    headerName: 'Last Updated',
  },
  {
    field: '',
    headerName: '',
    cellRenderer: () => {
      return <SvgSpriteLoader id='arrow-narrow-right' width={14} height={14} color={COLORS.GRAY_900} />;
    },
    width: 108,
    flex: 0,
    minWidth: 108,
    cellClass: cn(DATA_TABLE_CONFIG.cellClass, 'hidden-cell'),
  },
];

export const CustomColumnsMapping: Record<CUSTOM_COLUMNS_TYPE, (props: ICellRendererParams) => JSX.Element> = {
  [CUSTOM_COLUMNS_TYPE.TAG]: CustomTagRenderer,
};

export enum TEAM_OPTIONS {
  ENGG = 'engg',
  DESIGN = 'design',
  SALES_MARKETING = 'sales_marketing',
  PRODUCT = 'product',
  HIRING = 'hiring',
}

export const TEAM_OPTIONS_LIST = [
  {
    label: 'Engg',
    value: TEAM_OPTIONS.ENGG,
    color: COLORS.ORANGE_200,
  },
  {
    label: 'Design',
    value: TEAM_OPTIONS.DESIGN,
    color: COLORS.BLUE_150,
  },
  {
    label: 'Sales/Marketing',
    value: TEAM_OPTIONS.SALES_MARKETING,
    color: COLORS.VIOLET_100,
  },
  {
    label: 'Product',
    value: TEAM_OPTIONS.PRODUCT,
    color: COLORS.BLUE_150,
  },
  {
    label: 'Hiring',
    value: TEAM_OPTIONS.HIRING,
    color: COLORS.RED_250,
  },
];

export const DATASET_ACCESS_PRIVILEGES_LIST = [
  {
    label: 'Admin',
    value: DATASET_ACCESS_PRIVILEGES.ADMIN,
  },
  {
    label: 'Viewer',
    value: DATASET_ACCESS_PRIVILEGES.VIEWER,
  },
];

export const CHANGE_ACCESS_PRIVILEGES_LIST = [
  {
    label: 'Admin',
    value: DATASET_ACCESS_PRIVILEGES.ADMIN,
    desc: 'Can manage and share dataset',
  },
  {
    label: 'Viewer',
    value: DATASET_ACCESS_PRIVILEGES.VIEWER,
    desc: 'Can read data only',
  },
];
