import { FC, useState } from 'react';
import { copyToClipBoard } from 'utils/common';

type CopyToClipboardBrowserUrlPropsType = {
  initialText?: string;
  copiedText?: string;
};

const CopyToClipboardBrowserUrl: FC<CopyToClipboardBrowserUrlPropsType> = ({
  initialText = 'Copy link',
  copiedText = 'Copied!',
}) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    const currentUrl = window.location.href;

    copyToClipBoard(currentUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <span className='f-11-500 text-GRAY_1000 cursor-pointer select-none' onClick={handleCopy}>
      {copied ? copiedText : initialText}
    </span>
  );
};

export default CopyToClipboardBrowserUrl;
