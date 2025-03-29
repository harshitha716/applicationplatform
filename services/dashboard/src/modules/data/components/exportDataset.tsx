import { useRef, useState } from 'react';
import { toast } from 'react-toastify';
import {
  useLazyGetActionStatusQuery,
  useLazyGetDatasetExportQuery,
  useLazyGetDatasetExportsSignedUrlQuery,
} from 'apis/dataset';
import { COLORS } from 'constants/colors';
import { useOnClickOutside } from 'hooks';
import usePolling from 'hooks/usePolling';
import LoadingWidthAnimation from 'modules/data/components/LoadingWidthAnimation';
import { useRouter } from 'next/router';
import { DatasetActionStatusResponseType } from 'types/api/dataset.types';
import ProgressBar from 'components/common/RingProgress';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface ExportDatasetProps {
  query: string;
  datasetId: string;
  hasFilters: boolean;
}

const ExportDataset = ({ query, datasetId, hasFilters }: ExportDatasetProps) => {
  const router = useRouter();
  const dropdownRef = useRef<HTMLDivElement>(null);
  const { startPolling } = usePolling();
  const [getActionStatus] = useLazyGetActionStatusQuery();
  const [getDatasetExport] = useLazyGetDatasetExportQuery();
  const [getDatasetExportsSignedUrl] = useLazyGetDatasetExportsSignedUrlQuery();
  const [showExportStatus, setShowExportStatus] = useState(false);
  const [isPolling, setIsPolling] = useState<boolean>(false);

  useOnClickOutside(dropdownRef, () => setShowExportStatus(false));

  const handleDownloadCsvErr = () => {
    setIsPolling(false);
    toast.error('Unable to download CSV. Please try again.');
  };

  const handleSuccessfullyUpdate = (workflowId: string) => {
    setIsPolling(true);
    startPolling({
      fn: () => getActionStatus({ datasetId, params: { action_ids: [workflowId] } }),
      validate: (data: DatasetActionStatusResponseType[]) => {
        return data.filter((item) => !item.is_completed).length === 0;
      },
      interval: 3000,
      maxAttempts: 50,
    })
      .then(() => {
        getDatasetExportsSignedUrl({ datasetId, workflowId })
          .unwrap()
          .then((data) => {
            setIsPolling(false);
            router.push(data?.signed_url);
          });
      })
      .catch(handleDownloadCsvErr);
  };

  const downloadCsv = async () => {
    setShowExportStatus(true);

    if (isPolling) return;
    getDatasetExport({ datasetId, query_config: query })
      .unwrap()
      .then((data) => {
        if (data?.workflow_id) {
          handleSuccessfullyUpdate(data?.workflow_id);
        }
      })
      .catch(handleDownloadCsvErr);
  };

  return (
    <div className='relative z-1000 cursor-pointer w-5.5 h-5.5 rounded' onClick={downloadCsv}>
      <Tooltip
        tooltipBody={hasFilters ? 'Export filtered' : 'Export all'}
        tooltipBodyClassName='f-12-300 rounded-md whitespace-nowrap z-[1000] bg-black text-GRAY_200'
        className='z-1 h-full w-full'
        tooltipBodystyle='f-10-400'
        disabled={isPolling}
        position={TooltipPositions.BOTTOM}
      >
        <div className='hover:bg-GRAY_100 h-full w-full rounded flex items-center justify-center'>
          <SvgSpriteLoader id='download-02' width={14} height={14} className='text-GRAY_900' />
          {isPolling && (
            <div className='absolute bottom-px left-[3px]'>
              <LoadingWidthAnimation />
            </div>
          )}
        </div>
      </Tooltip>
      {isPolling && showExportStatus && (
        <div
          ref={dropdownRef}
          className='p-5 absolute top-7 -right-[86px] h-[55px] f-13-500 bg-white rounded-[10px] text-GRAY_1000 f-12-450 z-1000 flex items-center w-[308px] border-0.5 border-GRAY_500 gap-3'
        >
          <ProgressBar
            trackColor={COLORS.GRAY_400}
            indicatorColor={'#22A356'}
            indicatorWidth={2}
            trackWidth={2}
            size={16}
            className='animate-spin'
            progress={20}
          />
          <div className='grow'>Export in progress</div>
          <SvgSpriteLoader
            id='x-close'
            width={16}
            height={16}
            onClick={(e) => {
              e.stopPropagation();
              setShowExportStatus(false);
            }}
            className='text-GRAY_800 hover:text-GRAY_1000'
          />
        </div>
      )}
    </div>
  );
};

export default ExportDataset;
