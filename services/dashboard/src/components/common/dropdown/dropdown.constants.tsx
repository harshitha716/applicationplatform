import { SIZE_TYPES } from 'types/common/components';

export const SELECT_ALL_OPTION = { label: 'Select All', value: 'select_all' };
export const ACTIONS = {
  SELECT_OPTION: 'select-option',
  DESELECT_OPTION: 'deselect-option',
  REMOVE_VALUE: 'remove-value',
  POP_VALUE: 'pop-value',
};

export const DROPDOWN_SIZE_STYLES = {
  [SIZE_TYPES.XLARGE]: {
    customStyles: {
      control: {
        borderRadius: '6px',
        padding: '21px 24px',
        minHeight: '72px',
        width: 'fit-content',
        cursor: 'pointer',
      },
      option: {
        height: '64px',
      },
      menu: {
        borderRadius: '0px',
      },
      input: {},
      valueContainer: {},
      noOptionsMessage: {},
    },
    dropdownIndicatorProps: {
      width: 24,
      height: 24,
    },
    menuOptionClasses: {
      wrapperClass: 'px-2 h-16',
      labelOverrideClassName: 'f-16-400',
      contentWrapper: 'pl-2 py-3 w-full',
    },
    customClassNames: {
      placeholder: 'f-16-300',
    },
  },
  [SIZE_TYPES.LARGE]: {
    customStyles: {
      control: {
        borderRadius: '6px',
        padding: '21px 24px',
        minHeight: '72px',
        width: 'fit-content',
        cursor: 'pointer',
      },
      option: {
        height: '64px',
      },
      menu: {
        borderRadius: '0px',
      },
      input: {},
      valueContainer: {},
      noOptionsMessage: {},
    },
    dropdownIndicatorProps: {
      width: 24,
      height: 24,
    },
    menuOptionClasses: {
      wrapperClass: 'px-1 h-16',
      labelOverrideClassName: 'f-16-400',
      contentWrapper: 'pl-2 py-3 w-full',
    },
    customClassNames: {
      placeholder: 'f-16-300',
    },
  },
  [SIZE_TYPES.MEDIUM]: {
    customStyles: {
      control: {
        borderRadius: '6px',
        padding: '6px 12px',
        minHeight: '40px',
        width: 'fit-content',
        cursor: 'pointer',
      },
      option: {
        fontSize: '13px',
      },
      menu: {
        borderRadius: '6px',
        padding: '4px',
      },
      input: {
        fontSize: '13px',
        fontWeight: '300',
      },
      valueContainer: {},
      noOptionsMessage: {
        fontSize: '13px',
      },
    },
    dropdownIndicatorProps: {
      width: 20,
      height: 20,
    },
    menuOptionClasses: {
      wrapperClass: 'h-full overflow-y-scroll',
      labelOverrideClassName: 'f-12-500',
      contentWrapper: 'pl-2 py-3 w-full',
    },
    customClassNames: {
      placeholder: 'f-13-300 text-GRAY_700',
    },
  },
  [SIZE_TYPES.SMALL]: {
    customStyles: {
      control: {
        borderRadius: '0px',
        padding: '0px 20px',
        minHeight: '44px',
      },
      option: {
        height: '40px',
      },
      menu: {
        borderRadius: '0px',
      },
      input: {},
      valueContainer: {},
      noOptionsMessage: {},
    },
    dropdownIndicatorProps: {
      width: 16,
      height: 16,
    },
    menuOptionClasses: {
      wrapperClass: 'h-10 overflow-clip',
      labelOverrideClassName: 'f-16-400',
      contentWrapper: '!px-4 !py-3 w-full',
    },
    customClassNames: {
      placeholder: 'f-16-300',
    },
  },
  [SIZE_TYPES.XSMALL]: {
    customStyles: {
      control: {
        borderRadius: '6px',
        padding: '0px 8px',
        minHeight: '28px',
        width: 'fit-content',
        cursor: 'pointer',
      },
      option: {
        fontSize: '13px',
      },
      menu: {
        borderRadius: '6px',
        padding: '4px',
      },
      input: {
        fontSize: '13px',
        fontWeight: '300',
      },
      valueContainer: {},
      noOptionsMessage: {
        fontSize: '13px',
      },
    },
    dropdownIndicatorProps: {
      width: 20,
      height: 20,
    },
    menuOptionClasses: {
      wrapperClass: 'h-full overflow-y-scroll',
      labelOverrideClassName: 'f-12-500',
      contentWrapper: 'pl-2 py-3 w-full',
    },
    customClassNames: {
      placeholder: 'f-13-300 text-GRAY_700',
    },
  },
};
