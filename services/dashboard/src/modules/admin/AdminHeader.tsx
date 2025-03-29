import React, { FC, useState } from 'react';
import { usePostDatasetDisplayConfigMutation } from 'apis/admin';
import { AdminHeaderPropsType } from 'modules/admin/admin.types';
import JsonPreviewSidebar from 'modules/admin/components/previewSidebar';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';
import { toast } from 'components/common/toast/Toast';

const AdminHeader: FC<AdminHeaderPropsType> = ({ displayConfigInitialData, displayConfigFinalData, datasetId }) => {
  const [isJsonPreviewSidebarOpen, setIsJsonPreviewSidebarOpen] = useState(false);
  const [postDatasetDisplayConfig, { isLoading }] = usePostDatasetDisplayConfigMutation();
  const disableUpdateButton = JSON.stringify(displayConfigInitialData) === JSON.stringify(displayConfigFinalData);

  const handleOpenJsonPreviewSidebar = () => setIsJsonPreviewSidebarOpen(true);
  const handleCloseJsonPreviewSidebar = () => setIsJsonPreviewSidebarOpen(false);

  const handleUpdateJson = () => {
    postDatasetDisplayConfig({ datasetId: datasetId, body: { display_config: displayConfigFinalData } })
      .unwrap()
      .then(() => {
        toast.success('Updated successfully.');
      })
      .catch(() => {
        toast.error('Failed to update.');
      });
  };

  return (
    <div className='py-3 px-10 flex justify-end gap-2'>
      {isJsonPreviewSidebarOpen && (
        <JsonPreviewSidebar
          originalJson={displayConfigInitialData}
          formattedJson={displayConfigFinalData}
          isOpen={isJsonPreviewSidebarOpen}
          onClose={handleCloseJsonPreviewSidebar}
        />
      )}
      <Button
        type={BUTTON_TYPES.PRIMARY}
        id='preview-json'
        size={SIZE_TYPES.SMALL}
        onClick={handleOpenJsonPreviewSidebar}
      >
        Preview Json
      </Button>
      <Button
        type={BUTTON_TYPES.PRIMARY}
        id='update-preview-json'
        size={SIZE_TYPES.SMALL}
        onClick={handleUpdateJson}
        isLoading={isLoading}
        disabled={disableUpdateButton}
      >
        Update
      </Button>
    </div>
  );
};

export default AdminHeader;
