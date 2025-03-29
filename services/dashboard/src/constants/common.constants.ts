import { Inter } from 'next/font/google';

// Configure the Inter font
export const inter = Inter({
  subsets: ['latin'], // Specify subsets you need (e.g., 'latin', 'latin-ext').
  variable: '--font-inter', // Define a CSS variable to use in your styles.
  display: 'swap', // Controls font-display behavior.
});

export const SENTRY_DSN =
  'https://3129cf83b7bf9bd6c715ba81823cd0db@o4504767438520320.ingest.us.sentry.io/4508794285129728';

export enum ENVIRONMENT_TYPES {
  PRODUCTION = 'production',
  DEVELOPMENT = 'development',
  STAGING = 'staging',
}

export enum SIZE {
  XSMALL = 'xsmall',
  SMALL = 'small',
  MEDIUM = 'medium',
  LARGE = 'large',
  XLARGE = 'xlarge',
}

export enum BUTTON_TYPE {
  PRIMARY = 'primary',
  SECONDARY = 'secondary',
  TERTIARY = 'tertiary',
  ICON = 'icon',
  DESTRUCTIVE_HC = 'destructive-hc',
  DESTRUCTIVE_LC = 'destructive-lc',
  LINK = 'link',
}

export enum POSITION {
  LEFT = 'left',
  RIGHT = 'right',
  TOP = 'top',
  BOTTOM = 'bottom',
}

export const ALLOWED_EMAIL_DOMAINS = ['zamp.ai', 'zamp.finance'];
export const ENVIRONMENT = process.env.NEXT_PUBLIC_ENVIRONMENT || 'production';

export enum STORAGE_TYPES {
  SESSION = 'session',
  LOCAL = 'local',
}

export const SCREEN_BREAKPOINTS = {
  MIN_WIDTH: 854,
  MIN_HEIGHT: 300,
  XL_WIDTH: 1280,
  LG_WIDTH: 1024,
};
