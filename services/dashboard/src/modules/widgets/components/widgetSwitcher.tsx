import { FC, useMemo, useState } from 'react';
import WidgetsWrapper from 'modules/widgets/WidgetsWrapper';
import { WidgetType } from 'types/api/pagesApi.types';
import { WidgetInstanceType } from 'types/api/widgets.types';

interface WidgetSwitcherProps {
  widgetConfig: WidgetType;
  widgetInstances: WidgetInstanceType[];
  currency: string[];
  defaultCurrency: string;
  handleWidgetHeightChange: (height: number, isSingleHeader: boolean) => void;
}

const WidgetSwitcher: FC<WidgetSwitcherProps> = ({
  widgetConfig,
  widgetInstances,
  currency,
  defaultCurrency,
  handleWidgetHeightChange,
}) => {
  const [activeWidget, setActiveWidget] = useState<string>(widgetConfig?.default_widget);

  const onWidgetChange = (widgetId: string) => {
    setActiveWidget(widgetId);
  };

  const groupWidgetsOptions = useMemo(() => {
    return widgetInstances
      ?.filter((widget) => widgetConfig?.widget_group?.includes(widget.widget_instance_id))
      ?.map((widget) => ({ label: widget?.title, value: widget?.widget_instance_id }));
  }, [widgetInstances, widgetConfig]);

  const widgetDetails = useMemo(
    () => widgetInstances?.find((widget) => widget?.widget_instance_id === activeWidget),
    [widgetInstances, activeWidget],
  );

  return widgetDetails ? (
    <WidgetsWrapper
      widgetDetails={widgetDetails}
      groupWidgetsOptions={groupWidgetsOptions}
      onWidgetChange={onWidgetChange}
      currency={currency}
      defaultCurrency={defaultCurrency}
      activeWidget={activeWidget}
      handleWidgetHeightChange={handleWidgetHeightChange}
    />
  ) : (
    <div>No widget found</div>
  );
};

export default WidgetSwitcher;
