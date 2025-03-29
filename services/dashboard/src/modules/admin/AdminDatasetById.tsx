import React, { FC, useEffect, useState } from 'react';
import { useGetDatasetDisplayConfigQuery } from 'apis/admin';
import { ZAMP_LOGO_LOADER } from 'constants/lottie/zamp-logo-loader';
import { DisplayConfigHeadersList } from 'modules/admin/admin.constants';
import { AdminDatasetByIdPropsType, DISPLAY_CONFIG_HEADERS } from 'modules/admin/admin.types';
import AdminHeader from 'modules/admin/AdminHeader';
import EditableConfigField from 'modules/admin/components/previewSidebar/EditConfig';
import { cn } from 'utils/common';
import ToggleSwitch from 'components/common/toggleSwitch';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import DynamicLottiePlayer from 'components/DynamicLottiePlayer';

const AdminDatasetById: FC<AdminDatasetByIdPropsType> = ({ id }) => {
  const { data, isLoading, isError } = useGetDatasetDisplayConfigQuery({ datasetId: id });
  const displayConfigData = data?.display_config;
  const [editMode, setEditMode] = useState<{ [key: string]: boolean }>({});
  const [displayConfigUpdatedData, setDisplayConfigUpdatedData] = useState(displayConfigData);

  const handleEditToggle = (index: number, columnName: string) => {
    setEditMode((prev) => ({
      ...prev,
      [`${index}-${columnName}`]: !prev[`${index}-${columnName}`],
    }));
  };

  const handleChange = (index: number, field: string, value: string | boolean | number) => {
    const updatedData = [...(displayConfigUpdatedData || [])];

    updatedData[index] = { ...updatedData[index], [field]: value };
    setDisplayConfigUpdatedData(updatedData);
  };

  useEffect(() => {
    if ((displayConfigData?.length ?? 0) > 0) {
      setDisplayConfigUpdatedData(displayConfigData);
    }
  }, [displayConfigData]);

  return (
    <CommonWrapper
      className={cn('h-full', {
        'flex flex-col items-center justify-center': isLoading,
      })}
      isLoading={isLoading}
      isError={isError}
      skeletonType={SkeletonTypes.CUSTOM}
      loader={
        <div className='flex justify-center items-center h-full overflow-y-auto w-full z-1000 bg-white'>
          <DynamicLottiePlayer src={ZAMP_LOGO_LOADER} className='lottie-player h-[140px]' autoplay loop keepLastFrame />
        </div>
      }
    >
      <AdminHeader
        displayConfigInitialData={displayConfigData ?? []}
        displayConfigFinalData={displayConfigUpdatedData ?? []}
        datasetId={id}
      />

      <div className='flex flex-row overflow-y-auto'>
        <div className='flex flex-col w-full f-14-400'>
          <div className='grid grid-cols-6 border-GRAY_400 px-10'>
            {DisplayConfigHeadersList.map((header, index: number) => (
              <div key={index} className='border border-GRAY_400 px-2 py-5'>
                {header?.value}
              </div>
            ))}
          </div>
          <div className='flex flex-col p-10 pt-0 h-[calc(100vh-130px)] overflow-y-auto [&::-webkit-scrollbar]:hidden'>
            {displayConfigUpdatedData?.map((config, index) => (
              <div key={index} className='grid grid-cols-6 border-b border-GRAY_400'>
                <EditableConfigField
                  value={config?.column || ''}
                  isEditing={editMode[`${index}-${DISPLAY_CONFIG_HEADERS.COLUMN}`]}
                  onEditToggle={() => handleEditToggle(index, DISPLAY_CONFIG_HEADERS.COLUMN)}
                  onChange={(e) => handleChange(index, DISPLAY_CONFIG_HEADERS.COLUMN, e.target.value)}
                  firstColumn
                />

                <div className='flex gap-2 border-r border-GRAY_400 p-2'>
                  <ToggleSwitch
                    id='toggle-is_hidden'
                    toggleClassName='relative w-10 h-5 rounded-full border-none'
                    sliderClassName='absolute top-0.5 rounded-full w-4 h-4 transition-all duration-200'
                    checked={Boolean(config?.is_hidden)}
                    onChange={(state: boolean) => handleChange(index, DISPLAY_CONFIG_HEADERS.IS_HIDDEN, state)}
                  />
                  {config?.is_hidden ? 'True' : 'False'}
                </div>

                <div className='flex gap-2 border-r border-GRAY_400 p-2'>
                  <ToggleSwitch
                    id='toggle-is_editable'
                    toggleClassName='relative w-10 h-5 rounded-full border-none'
                    sliderClassName='absolute top-0.5 rounded-full w-4 h-4 transition-all duration-200'
                    checked={Boolean(config?.is_editable)}
                    onChange={(state: boolean) => handleChange(index, DISPLAY_CONFIG_HEADERS.IS_EDITABLE, state)}
                  />
                  {config?.is_editable ? 'True' : 'False'}
                </div>

                <span className=' border-r border-GRAY_400 p-2 overflow-hidden'>{config?.type || '-'}</span>

                <EditableConfigField
                  value={config?.config?.amount_column || ''}
                  isEditing={editMode[`${index}-${DISPLAY_CONFIG_HEADERS.AMOUNT_COLUMN}`]}
                  onEditToggle={() => handleEditToggle(index, DISPLAY_CONFIG_HEADERS.AMOUNT_COLUMN)}
                  onChange={(e) => handleChange(index, DISPLAY_CONFIG_HEADERS.AMOUNT_COLUMN, e.target.value)}
                />

                <EditableConfigField
                  value={config?.config?.currency_column || ''}
                  isEditing={editMode[`${index}-${DISPLAY_CONFIG_HEADERS.CURRENCY_COLUMN}`]}
                  onEditToggle={() => handleEditToggle(index, DISPLAY_CONFIG_HEADERS.CURRENCY_COLUMN)}
                  onChange={(e) => handleChange(index, DISPLAY_CONFIG_HEADERS.CURRENCY_COLUMN, e.target.value)}
                />
              </div>
            ))}
          </div>
        </div>
      </div>
    </CommonWrapper>
  );
};

export default AdminDatasetById;
