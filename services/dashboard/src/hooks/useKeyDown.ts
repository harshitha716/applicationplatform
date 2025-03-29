import { useEffect } from 'react';
import { NavigationItemSchema } from 'types/config';
const useKeyDown = (
  handleShortcuts: (val: KeyboardEvent) => void,
  filterKey?: string | string[],
  functionKey?: string | string[],
  ignoreInputs = true,
  navigation?: NavigationItemSchema[],
) => {
  useEffect(() => {
    const handleKeyDown = async (event: KeyboardEvent) => {
      const ignoreTags = ['INPUT', 'TEXTAREA'];
      const isIgnoreInput =
        !ignoreInputs || (ignoreInputs && !ignoreTags.includes((event?.target as HTMLElement)?.tagName?.toUpperCase()));
      const isFilterEvents =
        !filterKey || (typeof filterKey === 'string' && event.code === filterKey) || filterKey.includes(event.code);
      const isFunctionKey =
        !functionKey ||
        (typeof functionKey === 'string' && event[functionKey as keyof KeyboardEvent]) ||
        (Array.isArray(functionKey) &&
          functionKey.some((key: string) => event[key as keyof KeyboardEvent] || event.code === key));

      if (isIgnoreInput && !event.repeat && isFilterEvents && isFunctionKey) {
        handleShortcuts(event);
      }
    };

    window?.addEventListener('keydown', handleKeyDown, { passive: true });

    return () => {
      window?.removeEventListener('keydown', handleKeyDown);
    };
  }, [navigation]);
};

export default useKeyDown;
