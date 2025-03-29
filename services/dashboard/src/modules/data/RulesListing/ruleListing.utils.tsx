import { RuleCardProps } from 'modules/data/RulesListing/RuleCard';

export const searchRules = (rules: RuleCardProps[], searchTerm: string): RuleCardProps[] => {
  if (!searchTerm) return rules;
  const filteredRules: RuleCardProps[] = [];
  const searchTermLower = searchTerm.toLowerCase();

  rules.forEach((rule) => {
    const isSearchTermInValue = rule?.value?.toLowerCase().includes(searchTermLower);

    const isSearchTermInFilters = rule?.filters?.conditions?.some((condition) => {
      const checkValue = (value: string | string[]) => {
        if (Array.isArray(value)) {
          return value.some((item) => item.toLowerCase().includes(searchTermLower));
        }

        return String(value)?.toLowerCase().includes(searchTermLower);
      };

      return (
        condition.column?.column?.toLowerCase().includes(searchTermLower) ||
        (condition.column?.alias?.toLowerCase() || '').includes(searchTermLower) ||
        checkValue(condition.value)
      );
    });

    if (isSearchTermInValue || isSearchTermInFilters) {
      filteredRules.push(rule);
    }
  });

  return filteredRules;
};
