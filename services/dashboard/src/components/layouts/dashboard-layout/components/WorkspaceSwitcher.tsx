import React, { useRef, useState } from 'react';
import { WORKSPACE_ITEMS } from 'constants/dummyData';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useOnClickOutside } from 'hooks';
import { cn } from 'utils/common';
// import PageNavTab from 'components/layouts/dashboard-layout/components/PageNavTab';
import WorkspaceTab from 'components/layouts/dashboard-layout/components/WorkspaceTab';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface WorkspaceType {
  workspace_id: string;
  label: string;
  color: string;
}

const WorkspaceSwitcher = () => {
  const [isWorkspacePopoverOpen, setIsWorkspacePopoverOpen] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState<WorkspaceType>(WORKSPACE_ITEMS[0]);
  const ref = useRef(null);

  useOnClickOutside(ref, () => {
    setIsWorkspacePopoverOpen(false);
  });

  const handleWorkspaceClick = (workspace: WorkspaceType) => {
    setSelectedWorkspace(workspace);
    setIsWorkspacePopoverOpen(false);
  };

  return (
    <div className='px-2' ref={ref}>
      <div className='relative hidden'>
        <div
          className=' flex items-center gap-1 px-2 py-2.5 f-13-500 select-none cursor-pointer'
          onClick={() => setIsWorkspacePopoverOpen((prev) => !prev)}
        >
          <WorkspaceTab label={selectedWorkspace.label} className='pr-0' color={selectedWorkspace.color} />
          <SvgSpriteLoader
            iconCategory={ICON_SPRITE_TYPES.ARROWS}
            id='chevron-down'
            className={cn('transition-transform duration-300 -mb-0.5', isWorkspacePopoverOpen ? '-rotate-180' : '')}
          />
        </div>
        {isWorkspacePopoverOpen && (
          <div className='bg-white absolute rounded-md top-[90%] left-0 w-[264px] border border-GRAY_400 z-10 px-2 py-3'>
            {WORKSPACE_ITEMS.map((workspace) => (
              <WorkspaceTab
                key={workspace.workspace_id}
                label={workspace.label}
                onClick={() => handleWorkspaceClick(workspace)}
                isSelected={selectedWorkspace.workspace_id === workspace.workspace_id}
                color={workspace.color}
                className='gap-0.5'
              />
            ))}
          </div>
        )}
      </div>
      <div className='px-1 py-2.5'>
        <div className='f-11-600 text-GRAY_700 px-1.5 py-2'>Pages</div>
        {/* {PAGES_ITEMS.map((item) => (
          <PageNavTab key={item.label} label={item.label} />
        ))} */}
      </div>
    </div>
  );
};

export default WorkspaceSwitcher;
