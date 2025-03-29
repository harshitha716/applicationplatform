import { FC } from 'react';

type RuleStatementProps = {
  index: number;
  filterStatement: JSX.Element;
  numberOfFilters: number;
};

const RuleStatement: FC<RuleStatementProps> = ({ index, filterStatement, numberOfFilters }) => {
  return (
    <>
      {filterStatement}
      {index !== numberOfFilters - 1 && <span className='text-GRAY_1000 pl-1.5 pr-2 py-1 h-fit'>and</span>}
    </>
  );
};

export default RuleStatement;
