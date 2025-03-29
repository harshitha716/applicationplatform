import React, { FC, useRef, useState } from 'react';
import { useOnClickOutside } from 'hooks';
import { JsonPreviewSidebarPropsType } from 'modules/admin/admin.types';
import FormattedJson from 'modules/admin/components/previewSidebar/FormattedJson';
import Input from 'components/common/input';
import SideDrawer from 'components/common/SideDrawer/SideDrawer';

const JsonPreviewSidebar: FC<JsonPreviewSidebarPropsType> = ({ formattedJson, originalJson, onClose, isOpen }) => {
  const jsonPreviewRef = useRef<HTMLDivElement>(null);
  const [search, setSearch] = useState('');

  useOnClickOutside(jsonPreviewRef, onClose);

  return (
    <SideDrawer
      id='json-preview-sidebar'
      isOpen={isOpen}
      onClose={onClose}
      title='JSON Preview'
      className='h-screen overflow-hidden'
    >
      <>
        <Input
          type='text'
          placeholder='search keyword ...'
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className=' mb-4'
        />
        <FormattedJson originalJson={originalJson} formattedJson={formattedJson} search={search} />
      </>
    </SideDrawer>
  );
};

export default JsonPreviewSidebar;
