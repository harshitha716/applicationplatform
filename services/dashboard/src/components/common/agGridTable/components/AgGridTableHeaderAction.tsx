import React, { FC, useRef, useState } from 'react';
import { MapAny } from 'types/commonTypes';
import { DataAlign } from 'components/common/agGridTable/agGrid.types';
import AgGridTableActionPortal from 'components/common/agGridTable/components/AgGridTableActionPortal';
import SvgSpriteLoader, { SvgSpriteLoaderProps } from 'components/SvgSpriteLoader';

export interface AgGridTableHeaderActionProps {
  actionIcon?: SvgSpriteLoaderProps;
  ActionComponent?: FC<any>;
  actionComponentProps?: MapAny;
  actionAlign?: DataAlign;
  actionWrapperClassName?: string;
  id?: string;
  columnIndex?: number;
}

const AgGridTableHeaderAction: FC<AgGridTableHeaderActionProps> = ({
  actionIcon,
  ActionComponent,
  actionComponentProps,
  actionAlign,
  actionWrapperClassName = '',
  id = '',
  columnIndex = -1,
}) => {
  const headerRef = useRef(null);
  const [isOpen, setIsOpen] = useState(false);
  const [clickedOutside, setClickedOutside] = useState(false);

  const handleVisibility = () => {
    if (!isOpen && !clickedOutside) setIsOpen(true);
  };

  const handleClose = () => setIsOpen(false);

  const handleClickOutside = () => {
    setClickedOutside(true);
    setTimeout(() => {
      setIsOpen(false);
      setClickedOutside(false);
    }, 100);
  };

  const onClickAction = (data: MapAny) => {
    actionComponentProps?.onClick?.(data);
  };

  return (
    <div
      ref={headerRef}
      className='relative'
      data-testid={`ag-grid-table-header-row-cell-action-wrapper-${columnIndex}-${id}`}
    >
      {!!actionIcon?.id && <SvgSpriteLoader className='cursor-pointer' onClick={handleVisibility} {...actionIcon} />}
      {isOpen && !clickedOutside && (
        <AgGridTableActionPortal
          parentEl={headerRef?.current ?? null}
          actionAlign={actionAlign}
          actionWrapperClassName={actionWrapperClassName}
          onClose={handleClickOutside}
        >
          {!!ActionComponent && (
            <ActionComponent {...actionComponentProps} onClick={onClickAction} onClose={handleClose} isOpen={isOpen} />
          )}
        </AgGridTableActionPortal>
      )}
    </div>
  );
};

export default AgGridTableHeaderAction;
