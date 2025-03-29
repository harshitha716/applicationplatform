import React, { FC, ReactNode, useEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { useOnClickOutside } from 'hooks';
import { useWindowDimensions } from 'hooks/useWindowDimensions';
import { defaultFn, defaultFnType } from 'types/commonTypes';
import { checkIsObjectEmpty } from 'utils/common';
import { DataAlign } from 'components/common/agGridTable/agGrid.types';
import { COLUMN_HEADER_HEIGHT } from 'components/common/agGridTable/agGridTable.constants';

interface AgGridTableActionPortalProps {
  parentEl: HTMLDivElement | null;
  actionWrapperClassName?: string;
  actionAlign?: DataAlign;
  children: ReactNode;
  onClose?: defaultFnType;
}

const AgGridTableActionPortal: FC<AgGridTableActionPortalProps> = ({
  parentEl,
  actionWrapperClassName = '',
  actionAlign = DataAlign.CENTER,
  children,
  onClose = defaultFn,
}) => {
  const ref = useRef(null);
  const { width: windowWidth } = useWindowDimensions();
  const [positionStyles, setPositionStyles] = useState({});

  useOnClickOutside(ref, onClose);

  const { top, left, right, width } = parentEl?.getBoundingClientRect?.() || {};

  useEffect(() => {
    if (top !== undefined && left !== undefined && right !== undefined && width !== undefined) {
      const leftPosition = actionAlign === DataAlign.CENTER ? left + width / 2 : left;
      const rightPosition = windowWidth - right;

      setPositionStyles({
        top: top + COLUMN_HEADER_HEIGHT,
        left: actionAlign === DataAlign.LEFT || actionAlign === DataAlign.CENTER ? leftPosition : 'auto',
        right: actionAlign === DataAlign.RIGHT ? rightPosition : 'auto',
        transform: actionAlign === DataAlign.CENTER ? 'translateX(-50%)' : 'none',
      });
    }
  }, [windowWidth, top, left, right, width]);

  return (
    !checkIsObjectEmpty(positionStyles) &&
    createPortal(
      <div
        style={{ ...positionStyles }}
        className={`z-[998] fixed top-0 ${actionWrapperClassName}`}
        ref={ref}
        onClick={(e) => {
          e?.stopPropagation();
        }}
      >
        {children}
      </div>,
      document.body,
    )
  );
};

export default AgGridTableActionPortal;
