import { ROUTES_PATH } from 'constants/routeConfig';

export const KEYS_DELIMITER = '::';

export const MONITOR_KEYS = [
  'KeyA',
  'KeyB',
  'KeyC',
  'KeyD',
  'KeyG',
  'KeyH',
  'KeyI',
  'KeyK',
  'KeyM',
  'KeyO',
  'KeyP',
  'KeyR',
  'KeyT',
  'KeyK',
  'ShiftLeft',
  'MetaLeft',
  'Escape',
  'MetaLeft',
];

export const KEYBOARD_FUNCTION_KEYS = {
  SHIFT_KEY: 'shiftKey',
  META_KEY: 'metaKey',
  OPTION: 'altKey',
};

export enum KEYBOARD_KEYS {
  A = 'KeyA',
  B = 'KeyB',
  C = 'KeyC',
  D = 'KeyD',
  E = 'KeyE',
  F = 'KeyF',
  G = 'KeyG',
  H = 'KeyH',
  I = 'KeyI',
  J = 'KeyJ',
  K = 'KeyK',
  L = 'KeyL',
  M = 'KeyM',
  N = 'KeyN',
  O = 'KeyO',
  P = 'KeyP',
  Q = 'KeyQ',
  R = 'KeyR',
  S = 'KeyS',
  T = 'KeyT',
  U = 'KeyU',
  V = 'KeyV',
  W = 'KeyW',
  X = 'KeyX',
  Y = 'KeyY',
  Z = 'KeyZ',
  ENTER = 'Enter',
  ESCAPE = 'Escape',
  SPACE = 'Space',
  ARROW_UP = 'ArrowUp',
  ARROW_DOWN = 'ArrowDown',
  ARROW_LEFT = 'ArrowLeft',
  ARROW_RIGHT = 'ArrowRight',
  BACKSPACE = 'Backspace',
  DELETE = 'Delete',
  SLASH = 'Slash',
  DOT = 'Period',
  COMMA = 'Comma',
  SHIFT = 'Shift',
  META = 'Meta',
  CONTROL = 'Control',
}

export const FUNCTION_KEYS_ICON = {
  SHIFT_KEY: '⇧',
  META_KEY: '⌘',
  OPTION: '⌥',
};

export const SHORTCUTS_TABS = [
  {
    label: 'Data',
    value: 'data',
    id: ROUTES_PATH.DATA,
  },
];
