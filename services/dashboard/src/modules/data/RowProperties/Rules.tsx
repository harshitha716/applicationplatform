import { FC, useMemo } from 'react';
import { useGetRulesByRuleIdsQuery } from 'apis/dataset';
import RuleCard, { RuleCardProps } from 'modules/data/RulesListing/RuleCard';
import { cn } from 'utils/common';
import CommonWrapper from 'components/commonWrapper';
import { getTagLabel } from 'components/filter/filter.utils';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

type RulesProps = {
  ruleIds: string[];
  selectedRuleId: string;
};

const Rules: FC<RulesProps> = ({ ruleIds, selectedRuleId }) => {
  const {
    data: rulesData,
    isLoading,
    isError,
  } = useGetRulesByRuleIdsQuery({ rule_ids: ruleIds }, { skip: !ruleIds.length });

  const listOfFilters: RuleCardProps[] = useMemo(
    () =>
      rulesData?.map((rule) => {
        return {
          filters: rule?.filter_config?.query_config?.filters,
          value: rule?.value,
          createdOn: rule?.created_at,
          defaultExpanded: selectedRuleId === rule?.rule_id,
        };
      }) ?? [],
    [rulesData, selectedRuleId],
  );

  return (
    <CommonWrapper
      isLoading={isLoading}
      isError={isError}
      isNoData={!ruleIds.length}
      className={cn({ 'h-full': !ruleIds.length })}
      noDataBanner={
        <div className='flex items-center gap-2.5 h-full justify-center text-GRAY_700 f-12-450'>
          <SvgSpriteLoader id='lightning-01' width={24} height={24} />
          <div>No rules found</div>
        </div>
      }
    >
      <div className='space-y-3.5'>
        {listOfFilters?.map((filter, index) => (
          <RuleCard
            filters={filter?.filters}
            key={index}
            value={getTagLabel(filter?.value ?? '')}
            createdOn={filter?.createdOn}
            defaultExpanded={filter?.defaultExpanded}
          />
        ))}
      </div>
    </CommonWrapper>
  );
};

export default Rules;
