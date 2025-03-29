import React from 'react';
import { NOTEBOOK_ICON } from 'constants/icons';
import { getPageRouteById } from 'constants/routeConfig';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { cn } from 'utils/common';

interface PageNavTabProps {
  label: string;
  pageId: string;
  isSelected?: boolean;
}

const PageNavTab = ({ label, pageId, isSelected }: PageNavTabProps) => {
  const router = useRouter();

  return (
    <div
      className={cn(
        'flex items-center gap-3 text-GRAY_900 px-2 py-2 f-13-500 hover:bg-GRAY_20 rounded-md cursor-pointer select-none',
        isSelected ? 'bg-GRAY_100 text-GRAY_1000' : '',
      )}
      onClick={() => router.push(getPageRouteById(pageId))}
    >
      <Image
        width={16}
        height={16}
        alt='page file'
        className='w-[14px] align-middle cursor-pointer'
        src={NOTEBOOK_ICON}
        priority={true}
      />
      <div>{label}</div>
    </div>
  );
};

export default PageNavTab;
