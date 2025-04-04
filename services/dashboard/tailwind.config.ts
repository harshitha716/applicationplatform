/* eslint @typescript-eslint/no-require-imports: 0 */

module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx}',
    './src/components/**/*.{js,ts,jsx,tsx}',
    './src/layout/**/*.{js,ts,jsx,tsx}',
    './src/modules/**/*.{js,ts,jsx,tsx}',
    './node_modules/destiny/dist/components/**/*.{js,ts}',
  ],
  theme: {
    extend: {
      colors: {
        // ZAMP PLATFORM COLORS (_D for dark mode)
        GRAY_20: 'var(--GRAY_20)',
        GRAY_50: 'var(--GRAY_50)',
        GRAY_70: 'var(--GRAY_70)',
        GRAY_100: 'var(--GRAY_100)',
        GRAY_200: 'var(--GRAY_200)',
        GRAY_300: 'var(--GRAY_300)',
        GRAY_400: 'var(--GRAY_400)',
        GRAY_500: 'var(--GRAY_500)',
        GRAY_600: 'var(--GRAY_600)',
        GRAY_700: 'var(--GRAY_700)',
        GRAY_800: 'var(--GRAY_800)',
        GRAY_900: 'var(--GRAY_900)',
        GRAY_950: 'var(--GRAY_950)',
        GRAY_1000: 'var(--GRAY_1000)',

        BG_GRAY_1: 'var(--BG_GRAY_1)',
        BG_GRAY_2: 'var(--BG_GRAY_2)',
        BG_GRAY_3: 'var(--BG_GRAY_3)',
        BG_GRAY_4: 'var(--BG_GRAY_4)',
        BG_GRAY_5: 'var(--BG_GRAY_5)',
        BORDER_GRAY_400: 'var(--BORDER_GRAY_400)',

        BLUE_50: 'var(--BLUE_50)',
        BLUE_100: 'var(--BLUE_100)',
        BLUE_200: 'var(--BLUE_200)',
        BLUE_300: 'var(--BLUE_300)',
        BLUE_400: 'var(--BLUE_400)',
        BLUE_500: 'var(--BLUE_500)',
        BLUE_600: 'var(--BLUE_600)',
        BLUE_700: 'var(--BLUE_700)',
        BLUE_800: 'var(--BLUE_800)',
        BLUE_900: 'var(--BLUE_900)',
        BLUE_1000: 'var(--BLUE_1000)',

        GREEN_100: 'var(--GREEN_100)',
        GREEN_200: 'var(--GREEN_200)',
        GREEN_300: 'var(--GREEN_300)',
        GREEN_400: 'var(--GREEN_400)',
        GREEN_500: 'var(--GREEN_500)',
        GREEN_600: 'var(--GREEN_600)',
        GREEN_700: 'var(--GREEN_700)',
        GREEN_800: 'var(--GREEN_800)',
        GREEN_900: 'var(--GREEN_900)',
        GREEN_1000: 'var(--GREEN_1000)',

        ORANGE_100: 'var(--ORANGE_100)',
        ORANGE_200: 'var(--ORANGE_200)',
        ORANGE_300: 'var(--ORANGE_300)',
        ORANGE_400: 'var(--ORANGE_400)',
        ORANGE_500: 'var(--ORANGE_500)',
        ORANGE_600: 'var(--ORANGE_600)',
        ORANGE_700: 'var(--ORANGE_700)',
        ORANGE_800: 'var(--ORANGE_800)',
        ORANGE_900: 'var(--ORANGE_900)',
        ORANGE_1000: 'var(--ORANGE_1000)',

        RED_100: 'var(--RED_100)',
        RED_200: 'var(--RED_200)',
        RED_300: 'var(--RED_300)',
        RED_400: 'var(--RED_400)',
        RED_500: 'var(--RED_500)',
        RED_600: 'var(--RED_600)',
        RED_700: 'var(--RED_700)',
        RED_800: 'var(--RED_800)',
        RED_900: 'var(--RED_900)',
        RED_1000: 'var(--RED_1000)',

        BACKGROUND_GRAY_1: 'var(--BG_GRAY_1)',
        BACKGROUND_GRAY_2: 'var(--BG_GRAY_2)',
      },
      height: {
        0.25: '1px',
        3.5: '14px',
        4.5: '18px',
        5.5: '22px',
        7.5: '30px',
        8.5: '34px',
        11.5: '42px',
        15: '60px',
        18: '72px',
        62.5: '250px',
        105: '420px',
        107.5: '430px',
        112.5: '450px',
        115: '460px',
        topbar: '72px',
        body: 'calc(100vh - 72px)',
      },
      maxHeight: {
        10.5: '42px',
        100: '400px',
        105: '420px',
        125: '500px',
      },
      width: {
        0.25: '1px',
        3.5: '14px',
        4.5: '18px',
        5.5: '22px',
        7.5: '30px',
        12.5: '50px',
        13.5: '54px',
        18: '72px',
        25: '100px',
        30: '120px',
        34: '136px',
        34.5: '138px',
        36: '144px',
        42.5: '170px',
        50: '200px',
        55: '220px',
        60: '240px',
        64.5: '258px',
        65: '260px',
        69.5: '278px',
        87.5: '350px',
        100: '400px',
        104: '416px',
        content: 'calc(100vw - 272px)',
        sideDrawer: '480px',
      },
      maxWidth: {
        formLayout: '438px',
        40: '160px',
        55: '220px',
        60: '240px',
        75: '300px',
        94: '376px',
        95.5: '382px',
        104: '416px',
        145: '580px',
        360: '1440px',
      },
      minWidth: {
        4: '16px',
        4.5: '18px',
        25: '100px',
        sidebar: '200px',
        sidebarmini: '60px',
        5: '20px',
        6: '24px',
        14: '56px',
        17.5: '70px',
        32: '128px',
        40: '160px',
        50: '200px',
        64: '256px',
        75: '300px',
        83: '332px',
        86: '344px',
      },
      minHeight: {
        4: '16px',
        5: '20px',
        6: '24px',
        8: '32px',
        9: '36px',
        12: '48px',
      },
      lineHeight: {
        3.5: '14px',
        4.5: '18px',
      },
      margin: {
        0.25: '1px',
        1.5: '6px',
        4.5: '18px',
        5.5: '22px',
        2.5: '10px',
        12.5: '50px',
        13.5: '54px',
        17.5: '70px',
        19: '76px',
        38: '152px',
        45: '180px',
      },
      inset: {
        12.5: '50px',
        13.5: '54px',
        37.5: '150px',
      },
      padding: {
        0.5: '2.5px',
        2.5: '10px',
        4.5: '18px',
        5.5: '22px',
        6.5: '26px',
        8.5: '34px',
        18: '72px',
      },
      marginTop: {
        5.5: '22px',
      },
      boxShadow: {
        overlay: '1px 2px 10px rgba(197, 220, 255, 0.54)',
        inputOutlineShadow: '0px 0px 0px 3px var(--GRAY_400)',
        inputErrorOutlineShadow: '0px 0px 0px 3px var(--RED_100)',
        tableFilterMenu: '1px 2px 10px 0px #A6A6A61A',
        pageBottomBar: '0px -4px 0px 0px #00000005',
        sideDrawer: '-3px 0px 0px 0px #00000005',
        sideDrawerInner: '10px 0px 50px 0px #0000000d',
        menuList: '1px 2px 20px 0px #0000001A',
      },
      borderRadius: {
        2.5: '10px',
        3.5: '14px',
        4.5: '18px',
      },

      screens: {
        '2xl_custom': { max: '1440px' },
      },
      borderWidth: {
        0.5: '0.5px',
      },
      fontFamily: {
        outfit: 'Outfit',
      },
      zIndex: {
        1000: 1000,
      },
      backgroundImage: {
        'faded-white':
          'linear-gradient(to right, transparent 0%, rgba(255,255,255,0.8) 5%, white 10%, white 80%, white 90%, rgba(255,255,255,0.8) 95%, transparent 100%)',
      },
      animation: {
        opacity: 'opacity 0.3s ease-in-out',
        'file-upload': 'file-upload 0.5s linear ',
        'reverse-spin': 'reverse-spin 1.5s linear infinite',
        'rightSideDrawer-mount': 'rightSideDrawerTransition 0.4s normal forwards ease-out',
        'bottomSideDrawer-mount': 'bottomSideDrawerTransition 0.4s normal forwards ease-out',
        'rightSideDrawer-unMount': 'rightSideDrawerUnMountTransition 0.4s normal forwards ease-out',
        'bottomSideDrawer-unMount': 'bottomSideDrawerUnMountTransition 0.4s normal forwards ease-out',
        'shimmer-round': 'shimmer-round 1.5s infinite linear',
        width: 'position 1.5s linear infinite',
        slide: 'slide 1.5s linear infinite',
        slideInOut: 'slideInOut 5s cubic-bezier(0.85, 0, 0.15, 1) forwards',
        'slide-in': 'slideIn 0.5s ease-in-out',
      },
      keyframes: {
        'reverse-spin': {
          from: {
            transform: 'rotate(360deg)',
          },
        },
        'file-upload': {
          from: { opacity: '0', marginTop: '-56px', zIndex: '-1' },
          to: { opacity: '1', marginTop: '0px', zIndex: '-1' },
        },
        opacity: {
          '0%': { opacity: 0 },
          '100%': { opacity: 1 },
        },
        rightSideDrawerTransition: {
          '0%': {
            right: '-50vw',
          },

          '100%': {
            right: '0px',
          },
        },
        bottomSideDrawerTransition: {
          '0%': {
            bottom: '-50vw',
          },

          '100%': {
            bottom: '0px',
          },
        },
        rightSideDrawerUnMountTransition: {
          '0%': {
            right: '0px',
          },

          '100%': {
            right: '-50vw',
          },
        },
        bottomSideDrawerUnMountTransition: {
          '0%': {
            bottom: '0px',
          },

          '100%': {
            bottom: '-50vw',
          },
        },
        'shimmer-round': {
          '0%': { transform: 'rotate(0deg)' },
          '100%': { transform: 'rotate(360deg)' },
        },
        position: {
          '0%': { left: '1px' },
          '50%': { left: '6px' },
          '100%': { left: '1px' },
        },
        'shimmer-skeleton': {
          '100%': { transform: 'translateX(100%)' },
        },
        slide: {
          '0%': { transform: 'translateX(-100%)' },
          '100%': { transform: 'translateX(200%)' },
        },
        slideInOut: {
          '0%': { transform: 'translateX(100%)', opacity: '0' },
          '10%': { transform: 'translateX(0)', opacity: '1' },
          '90%': { transform: 'translateX(0)', opacity: '1' },
          '100%': { transform: 'translateX(100%)', opacity: '0' },
        },
        slideIn: {
          '0%': { transform: 'translateY(-100%)', opactiy: 0 },
          '100%': { transform: 'translateY(0%)', opactiy: 1 },
        },
      },
      transitionProperty: {
        height: 'height',
      },
    },
  },
  variants: {
    textColor: ['group-hover'],
  },
  // eslint-disable-next-line import/no-unresolved
  plugins: [import('@tailwindcss/typography')],
};
