/**
 * get the extension of the file
 * @param fileName string demo.csv
 * @returns string csv
 */
export const getFileExtension = (fileName: string) => {
  return fileName.split('.').pop() ?? '';
};

/**
 * get formatted date
 * @param date string 2025-02-13 10:43:51
 * @returns string Thu, 13 Feb 2025
 */
export const formattedDate = (date: string) => {
  if (!date) return '';

  return new Date(date).toLocaleDateString('en-GB', {
    weekday: 'short',
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });
};

/**
 *  get formatted file name
 * @param fileName string LsRwCYzZQZFPvLetDo6T5B_01_29_1732008234522_qx4SidxmzEoABnvxFvWyUb_11_19.csv
 * @returns string LsRw...Ub_11_19.csv
 */
export const maskString = (str: string, start: number, end: number, limit?: number) => {
  if (!str) return '';

  const parts = str.split('.');
  const extension = parts.pop();
  const name = parts.join('.');

  limit = limit ?? 16;
  if (name.length > limit) {
    return `${name.slice(0, start)}...${name.slice(-end)}.${extension}`;
  }

  return str;
};

/**
 * generate random string
 * @param length number
 * @returns string
 */
export const generateUniqueId = (length: number) => {
  if (length <= 0) return '';

  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const array = new Uint32Array(length);

  crypto.getRandomValues(array);

  return Array.from(array, (x) => chars[x % chars.length]).join('');
};

/**
 * Get file type from file object
 * @param file
 * @returns file type
 */

export const getFileType = (file: File) => {
  if (file && !file?.type) {
    const fileName = file?.name;
    const fileExtension = fileName?.split('.')?.pop()?.toLowerCase();

    if (fileExtension) {
      return fileExtension;
    }
  }

  return file.type;
};
