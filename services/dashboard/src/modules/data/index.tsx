import React, { FC, useMemo } from 'react';
import { IServerSideDatasource, IServerSideGetRowsParams } from 'ag-grid-community';
import { useLazyGetDatasetListingQuery } from 'apis/dataset';
import { LISTING_COLUMNS } from 'modules/data/data.constants';
import { ListingPropsType } from 'modules/data/data.types';
import { formatData } from 'modules/data/data.utils';
import { OrderType } from 'types/components/table.type';
import DataTable from 'components/common/table/DataTable';
import { PAGE_SIZE } from 'components/common/table/table.constants';

const Listing: FC<ListingPropsType> = ({ onRowClicked }) => {
  const [getDatasetListing] = useLazyGetDatasetListingQuery();
  const columns = useMemo(() => LISTING_COLUMNS, []);

  const serverSideDatasource: IServerSideDatasource = useMemo(() => {
    return {
      getRows: (parameters: IServerSideGetRowsParams): void => {
        const sortModel =
          parameters.request.sortModel?.map((item) => ({
            column: item.colId,
            desc: item.sort === OrderType.DESC,
          })) ?? [];

        getDatasetListing(
          {
            page: Math.floor(parameters.request.endRow ?? 0) / PAGE_SIZE,
            pageSize: PAGE_SIZE,
            sort: JSON.stringify(sortModel),
          },
          true,
        )
          .unwrap()
          .then((data) => {
            parameters.success({
              rowData: formatData(data?.datasets ?? []),
              ...(parameters.request.startRow === 0 ? { rowCount: data?.total_count } : {}),
            });
          })
          .catch(() => {
            parameters.fail();
          });
      },
    };
  }, [getDatasetListing]);

  return (
    <div className='rounded-tl-xl overflow-hidden'>
      <DataTable columns={columns} onRowClicked={onRowClicked} serverSideDatasource={serverSideDatasource} />
    </div>
  );
};

export default Listing;
