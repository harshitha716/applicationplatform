import React, { FC } from 'react';
import { FormattedJsonPropsType } from 'modules/admin/admin.types';
import { displayConfigType } from 'types/api/admin.types';
import { toast } from 'components/common/toast/Toast';

const highlightDifferences = (
  original: displayConfigType[],
  formatted: displayConfigType[],
  searchQuery: string,
  indent = 2,
) => {
  const originalJsonStr = JSON.stringify(original, null, indent);
  const formattedJsonStr = JSON.stringify(formatted, null, indent);
  const regexLineMatch = /^\s*("[^"]+"):\s*(.*)/;

  const originalLines = originalJsonStr.split('\n');
  const formattedLines = formattedJsonStr.split('\n');

  let searchCount = 0;

  const highlightSearch = (text: string, query: string) => {
    if (!query) return text;

    const lowerText = text.toLowerCase();
    const lowerQuery = query.toLowerCase();

    const parts = lowerText.split(lowerQuery);

    if (parts.length === 1) return text;

    let lastIndex = 0;
    const highlighted: (string | JSX.Element)[] = [];

    parts.forEach((part, index) => {
      if (index > 0) {
        highlighted.push(
          <span key={index} className='bg-RED_300 text-black'>
            {text.slice(lastIndex, lastIndex + query.length)}
          </span>,
        );
        lastIndex += query.length;
      }
      highlighted.push(part);
      lastIndex += part.length;
    });

    searchCount += parts.length - 1;

    return highlighted;
  };

  const highlightedJson: JSX.Element[] = formattedLines.map((line, index) => {
    const originalLine = originalLines[index] || '';
    const match = line.match(regexLineMatch);

    if (match) {
      const [, , valuePart] = match;
      const originalMatch = originalLine.match(regexLineMatch);
      const originalValue = originalMatch ? originalMatch[2] : '';

      let highlightedValue: JSX.Element | string = valuePart;

      if (valuePart !== originalValue) {
        highlightedValue = <span className='bg-yellow-300 text-black'>{valuePart}</span>;
      }

      if (searchQuery) {
        if (searchQuery.length > 50) {
          toast.error('Search query is too long, please try a shorter query.');

          return <span key={index}>{line + '\n'}</span>;
        }

        const searchHighlighted = highlightSearch(valuePart, searchQuery);

        highlightedValue =
          valuePart !== originalValue ? (
            <span className='bg-yellow-300 text-black'>{searchHighlighted}</span>
          ) : (
            <>{searchHighlighted}</>
          );
      }

      return (
        <span key={index}>
          {match[0].replace(valuePart, '')}
          {highlightedValue}
          <br />
        </span>
      );
    }

    return <span key={index}>{line + '\n'}</span>;
  });

  return { highlightedJson, searchCount };
};

const FormattedJson: FC<FormattedJsonPropsType> = ({ originalJson, formattedJson, search }) => {
  const { highlightedJson, searchCount } = highlightDifferences(originalJson, formattedJson, search);

  return (
    <>
      <div className='sticky w-full h-fit rounded-md border border-GRAY_300'>
        <div className='flex w-full justify-end absolute'>
          <span className='f-10-400 bg-GRAY_300 py-0.5 p-4 rounded-bl-md text-GRAY_700'>JSON</span>
        </div>
        <div
          className='flex flex-col w-full max-h-[calc(100vh-180px)] overflow-y-auto rounded-md'
          style={{ scrollbarWidth: 'thin' }}
        >
          <pre className='f-14-400 w-full p-4 whitespace-pre-wrap'>{highlightedJson}</pre>
        </div>
      </div>
      {search && (
        <span className='flex justify-end f-12-400 mb-2 mt-2 gap-1'>
          Found <span className='f-12-600'>{searchCount}</span> search {searchCount > 1 ? 'results' : 'result'}
        </span>
      )}
    </>
  );
};

export default FormattedJson;
