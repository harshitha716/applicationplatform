import React from 'react';

const NotAuthorized = () => {
  return (
    <div className='flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-100 to-gray-200'>
      <div className='text-center p-8 bg-white shadow-xl rounded-lg max-w-md w-full mx-4'>
        <div className='mb-6'></div>
        <h1 className='text-3xl sm:text-4xl  text-gray-900 mb-2'>Restricted Access</h1>
        <p className='text-xl text-gray-600 mb-8'>You don&apos;t have access to this page.</p>
      </div>
    </div>
  );
};

export default NotAuthorized;
