const PivotTableLoader = () => {
  return (
    <div className='overflow-x-auto w-full h-full border border-GRAY_400 rounded-xl animate-pulse overflow-hidden'>
      <table className='w-full text-left border-collapse'>
        <thead className='h-[84px]'>
          <tr className='border-b border-GRAY_400'>
            {Array.from({ length: 6 }).map((_, i) => (
              <th key={i} className='py-6 px-4 border-r-0.5 border-GRAY_400 last:border-r-0 first:w-[380px] w-[170px]'>
                <div className='w-24 h-4 bg-GRAY_50 rounded'></div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: 12 }).map((_, rowIndex) => (
            <tr key={rowIndex} className='border-b-0.5 border-GRAY_400 last:border-b-0'>
              {Array.from({ length: 6 }).map((_, colIndex) => (
                <td key={colIndex} className='py-3 px-4 border-r-0.5 border-GRAY_400 last:border-r-0 h-[42px]'>
                  <div className='h-4 bg-GRAY_50 rounded w-full'></div>
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default PivotTableLoader;
