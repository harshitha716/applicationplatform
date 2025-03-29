import { FC } from 'react';
import { MembersEmailPropsType } from 'modules/team/people.types';

const MembersEmail: FC<MembersEmailPropsType> = ({ value = '' }) => {
  return (
    <div className='f-12-400 text-GRAY_1000 h-full flex items-start justify-start text-left py-3 px-2'>{value}</div>
  );
};

export default MembersEmail;
