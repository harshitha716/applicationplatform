import React, { FC, memo } from 'react';
import { InputProps } from 'types/common/components/input/input.types';
import InputField from 'components/common/input/InputField';
import { Label } from 'components/common/Label';
import { SupporterInfo } from 'components/common/SupporterInfo';

const Input: FC<InputProps> = ({
  supporterInfoProps = {},
  inputWrapperClassName = 'w-full ',
  className = 'w-full',
  label = '',
  description = '',
  labelClassName = '',
  labelOverrideClassName = 'f-12-500 text-GRAY_900 mb-2 select-none px-1.5',
  required,
  ...rest
}) => (
  <div className={className}>
    {(label || description) && (
      <Label titleClassName={`${labelOverrideClassName} ${labelClassName}`} title={label} description={description} />
    )}
    <div className={inputWrapperClassName}>
      <InputField {...rest} />
    </div>

    {supporterInfoProps.showSupportInfo && <SupporterInfo {...supporterInfoProps} />}
  </div>
);

export default memo(Input);
