import { MapAny } from 'types/commonTypes';

export type ArrayListOption = {
  value: string;
  resource_audience_type?: string;
  resource_audience_id?: string;
  validationMessage?: string;
  team_membership_id?: string;
  team_id?: string;
  label: string;
  valid: boolean;
  role?: string;
  color?: string;
  isNew?: boolean;
};

export type MultiSelectInputPropsType = {
  id: string;
  checkAudiencePresentInOrg?: boolean;
  search: string;
  setSearch: (value: string) => void;
  selectedRoleRef?: React.MutableRefObject<any>;
  isOpen: boolean;
  placeholderText: string;
  roleOptions?: Array<{ label: string; value: string }>;
  inputArrayList: ArrayListOption[];
  setInputArrayList: React.Dispatch<React.SetStateAction<ArrayListOption[]>>;
  showValidationError?: boolean;
  setShowValidationError?: React.Dispatch<React.SetStateAction<boolean>>;
  validationErrorText?: string;
  onValidateAndAdd: ({
    value,
    label,
    color,
    type,
    team_id,
  }: {
    value: string;
    label: string;
    color?: string;
    type?: string;
    team_id?: string;
  }) => void;
  optionsList?: { value: string; label: string; color?: string; type?: string; team_id?: string }[];
  onSelectOption?: (option: { value: string; label: string; color?: string; type?: string; team_id?: string }) => void;
  selectOnlyFromList?: boolean;
  transformLabel?: (label: string) => string;
  isLoadingOptionsList?: boolean;
  wrapperClassName?: string;
  inputWrapperClassName?: string;
  multiSelectInputClassName?: string;
  setIsCustomInputFocused?: React.Dispatch<React.SetStateAction<boolean>>;
  customOptionsListDropdown?: React.ElementType;
  selectedRole?: string;
  setSelectedRole?: React.Dispatch<React.SetStateAction<Record<number, string> | string>>;
  onCustomDeleteFn?: (item: MapAny) => void;
  optionalOpenDropdownOptions?: boolean;
};

export const KEY_CODES = {
  BACKSPACE: 'Backspace',
  ENTER: 'Enter',
  ESCAPE: 'Escape',
  SPACE: ' ',
  COMMA: ',',
  ARROW_UP: 'ArrowUp',
  ARROW_DOWN: 'ArrowDown',
};
