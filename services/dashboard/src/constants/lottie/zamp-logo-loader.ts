export const ZAMP_LOGO_LOADER = {
  v: '5.12.1',
  fr: 60,
  ip: 0,
  op: 60,
  w: 1080,
  h: 1080,
  nm: 'Zamp Logo Loader Lottie',
  ddd: 0,
  assets: [],
  layers: [
    {
      ddd: 0,
      ind: 1,
      ty: 3,
      nm: 'NULL 1',
      sr: 1,
      ks: {
        o: { a: 0, k: 0, ix: 11 },
        r: { a: 0, k: 0, ix: 10 },
        p: { a: 0, k: [540, 540, 0], ix: 2, l: 2 },
        a: { a: 0, k: [50, 50, 0], ix: 1, l: 2 },
        s: { a: 0, k: [100, 100, 100], ix: 6, l: 2 },
      },
      ao: 0,
      ip: 0,
      op: 121,
      st: 0,
      bm: 0,
    },
    {
      ddd: 0,
      ind: 2,
      ty: 4,
      nm: 'Shape Layer 1',
      parent: 1,
      sr: 1,
      ks: {
        o: { a: 0, k: 100, ix: 11 },
        r: { a: 0, k: 0, ix: 10 },
        p: {
          a: 1,
          k: [
            {
              i: { x: 0.833, y: 0.586 },
              o: { x: 0.167, y: 0.115 },
              t: 0,
              s: [39.25, -25.25, 0],
              to: [-0.583, 1.667, 0],
              ti: [1.5, -8.667, 0],
            },
            {
              i: { x: 0.833, y: 0.754 },
              o: { x: 0.167, y: 0.104 },
              t: 10,
              s: [35.75, -15.25, 0],
              to: [-1.5, 8.667, 0],
              ti: [0.083, -20.75, 0],
            },
            {
              i: { x: 0.833, y: 0.86 },
              o: { x: 0.167, y: 0.126 },
              t: 20,
              s: [30.25, 26.75, 0],
              to: [-0.083, 20.75, 0],
              ti: [-3.333, -18.75, 0],
            },
            {
              i: { x: 0.833, y: 0.88 },
              o: { x: 0.167, y: 0.243 },
              t: 30,
              s: [35.25, 109.25, 0],
              to: [2.428, 13.66, 0],
              ti: [-3.943, -4.602, 0],
            },
            {
              i: { x: 0.833, y: 0.829 },
              o: { x: 0.167, y: 0.408 },
              t: 36,
              s: [46.306, 135.289, 0],
              to: [1.469, 1.715, 0],
              ti: [-1.109, -0.973, 0],
            },
            {
              i: { x: 0.833, y: 0.879 },
              o: { x: 0.167, y: 0.165 },
              t: 40,
              s: [50.25, 139.25, 0],
              to: [4.083, 3.583, 0],
              ti: [-1.833, 2.467, 0],
            },
            {
              i: { x: 0.833, y: 0.833 },
              o: { x: 0.167, y: 0.266 },
              t: 50,
              s: [59.75, 130.75, 0],
              to: [1.833, -2.467, 0],
              ti: [-0.25, 1.05, 0],
            },
            { t: 60, s: [61.25, 124.45, 0] },
          ],
          ix: 2,
          l: 2,
        },
        a: { a: 0, k: [-13.75, -85.25, 0], ix: 1, l: 2 },
        s: { a: 0, k: [100, 100, 100], ix: 6, l: 2 },
      },
      ao: 0,
      ef: [
        {
          ty: 5,
          nm: 'Slider Control',
          np: 3,
          mn: 'ADBE Slider Control',
          ix: 1,
          en: 1,
          ef: [
            {
              ty: 0,
              nm: 'Slider',
              mn: 'ADBE Slider Control-0001',
              ix: 1,
              v: {
                a: 1,
                k: [
                  { i: { x: [0.679], y: [0.803] }, o: { x: [0.54], y: [0.148] }, t: 0, s: [100] },
                  { i: { x: [0.833], y: [0.452] }, o: { x: [0.485], y: [-0.27] }, t: 20, s: [174.1] },
                  { i: { x: [0.925], y: [1] }, o: { x: [0.167], y: [0.169] }, t: 30, s: [133.1] },
                  { i: { x: [0.145], y: [0.628] }, o: { x: [0.098], y: [0] }, t: 40, s: [0] },
                  { i: { x: [0.833], y: [0.833] }, o: { x: [0.209], y: [0.233] }, t: 50, s: [72] },
                  { t: 60, s: [100] },
                ],
                ix: 1,
              },
            },
          ],
        },
      ],
      shapes: [
        {
          ty: 'gr',
          it: [
            {
              ty: 'rc',
              d: 1,
              s: {
                a: 0,
                k: [253.5, 0],
                ix: 2,
                x: "var $bm_rt;\nvar temp;\ntemp = effect('Slider Control')('Slider');\n$bm_rt = $bm_sum(value, [\n    0,\n    temp\n]);",
              },
              p: { a: 0, k: [0, 0], ix: 3 },
              r: { a: 0, k: 0, ix: 4 },
              nm: 'Rectangle Path 1',
              mn: 'ADBE Vector Shape - Rect',
              hd: false,
            },
            {
              ty: 'fl',
              c: { a: 0, k: [0.813333429075, 0.813333429075, 0.813333429075, 1], ix: 4 },
              o: { a: 0, k: 100, ix: 5 },
              r: 1,
              bm: 0,
              nm: 'Fill 1',
              mn: 'ADBE Vector Graphic - Fill',
              hd: false,
            },
            {
              ty: 'tr',
              p: { a: 0, k: [-13.75, -85.25], ix: 2 },
              a: { a: 0, k: [0, 0], ix: 1 },
              s: { a: 0, k: [100, 100], ix: 3 },
              r: { a: 0, k: 0, ix: 6 },
              o: { a: 0, k: 100, ix: 7 },
              sk: {
                a: 1,
                k: [
                  { i: { x: [0.8], y: [0.484] }, o: { x: [0.2], y: [-0.017] }, t: 0, s: [20] },
                  { i: { x: [0.933], y: [0.881] }, o: { x: [0.607], y: [0.23] }, t: 30, s: [-12.5] },
                  { i: { x: [0.681], y: [0.993] }, o: { x: [0.123], y: [-0.01] }, t: 39, s: [-79] },
                  { i: { x: [0.254], y: [0.874] }, o: { x: [0.06], y: [-0.04] }, t: 40, s: [83] },
                  { i: { x: [0.833], y: [0.833] }, o: { x: [0.206], y: [0.236] }, t: 50, s: [28.1] },
                  { t: 60, s: [20] },
                ],
                ix: 4,
              },
              sa: { a: 0, k: 0, ix: 5 },
              nm: 'Transform',
            },
          ],
          nm: 'Rectangle 1',
          np: 2,
          cix: 2,
          bm: 0,
          ix: 1,
          mn: 'ADBE Vector Group',
          hd: false,
        },
        {
          ty: 'st',
          c: { a: 0, k: [0.811764765721, 0.811764765721, 0.811764765721, 1], ix: 3 },
          o: { a: 0, k: 100, ix: 4 },
          w: { a: 0, k: 4, ix: 5 },
          lc: 2,
          lj: 2,
          bm: 0,
          nm: 'Stroke 1',
          mn: 'ADBE Vector Graphic - Stroke',
          hd: false,
        },
      ],
      ip: 0,
      op: 121,
      st: 0,
      ct: 1,
      bm: 0,
    },
    {
      ddd: 0,
      ind: 3,
      ty: 3,
      nm: 'NULL 2',
      sr: 1,
      ks: {
        o: { a: 0, k: 0, ix: 11 },
        r: { a: 0, k: 0, ix: 10 },
        p: { a: 0, k: [540, 540, 0], ix: 2, l: 2 },
        a: { a: 0, k: [50, 50, 0], ix: 1, l: 2 },
        s: { a: 0, k: [-100, -100, 100], ix: 6, l: 2 },
      },
      ao: 0,
      ip: 0,
      op: 121,
      st: 0,
      bm: 0,
    },
    {
      ddd: 0,
      ind: 4,
      ty: 4,
      nm: 'Shape Layer 2',
      parent: 3,
      sr: 1,
      ks: {
        o: {
          a: 1,
          k: [
            { i: { x: [0.833], y: [0.833] }, o: { x: [0.167], y: [0.167] }, t: 6, s: [100] },
            { i: { x: [0.833], y: [0.833] }, o: { x: [0.167], y: [0.167] }, t: 21, s: [50] },
            { t: 38, s: [100] },
          ],
          ix: 11,
        },
        r: { a: 0, k: 0, ix: 10 },
        p: {
          a: 1,
          k: [
            {
              i: { x: 0.833, y: 0.586 },
              o: { x: 0.167, y: 0.115 },
              t: 0,
              s: [39.25, -25.25, 0],
              to: [-0.583, 1.667, 0],
              ti: [1.5, -8.667, 0],
            },
            {
              i: { x: 0.833, y: 0.754 },
              o: { x: 0.167, y: 0.104 },
              t: 10,
              s: [35.75, -15.25, 0],
              to: [-1.5, 8.667, 0],
              ti: [0.083, -20.75, 0],
            },
            {
              i: { x: 0.833, y: 0.86 },
              o: { x: 0.167, y: 0.126 },
              t: 20,
              s: [30.25, 26.75, 0],
              to: [-0.083, 20.75, 0],
              ti: [-3.333, -18.75, 0],
            },
            {
              i: { x: 0.833, y: 0.88 },
              o: { x: 0.167, y: 0.243 },
              t: 30,
              s: [35.25, 109.25, 0],
              to: [2.428, 13.66, 0],
              ti: [-3.943, -4.602, 0],
            },
            {
              i: { x: 0.833, y: 0.829 },
              o: { x: 0.167, y: 0.408 },
              t: 36,
              s: [46.306, 135.289, 0],
              to: [1.469, 1.715, 0],
              ti: [-1.109, -0.973, 0],
            },
            {
              i: { x: 0.833, y: 0.879 },
              o: { x: 0.167, y: 0.165 },
              t: 40,
              s: [50.25, 139.25, 0],
              to: [4.083, 3.583, 0],
              ti: [-1.833, 2.467, 0],
            },
            {
              i: { x: 0.833, y: 0.833 },
              o: { x: 0.167, y: 0.266 },
              t: 50,
              s: [59.75, 130.75, 0],
              to: [1.833, -2.467, 0],
              ti: [-0.25, 1.05, 0],
            },
            { t: 60, s: [61.25, 124.45, 0] },
          ],
          ix: 2,
          l: 2,
        },
        a: { a: 0, k: [-13.75, -85.25, 0], ix: 1, l: 2 },
        s: { a: 0, k: [100, 100, 100], ix: 6, l: 2 },
      },
      ao: 0,
      ef: [
        {
          ty: 5,
          nm: 'Slider Control',
          np: 3,
          mn: 'ADBE Slider Control',
          ix: 1,
          en: 1,
          ef: [
            {
              ty: 0,
              nm: 'Slider',
              mn: 'ADBE Slider Control-0001',
              ix: 1,
              v: {
                a: 1,
                k: [
                  { i: { x: [0.679], y: [0.803] }, o: { x: [0.54], y: [0.148] }, t: 0, s: [100] },
                  { i: { x: [0.833], y: [0.452] }, o: { x: [0.485], y: [-0.27] }, t: 20, s: [174.1] },
                  { i: { x: [0.925], y: [1] }, o: { x: [0.167], y: [0.174] }, t: 30, s: [133.1] },
                  { i: { x: [0.145], y: [0.607] }, o: { x: [0.098], y: [0] }, t: 40, s: [4] },
                  { i: { x: [0.833], y: [0.833] }, o: { x: [0.209], y: [0.233] }, t: 50, s: [72] },
                  { t: 60, s: [100] },
                ],
                ix: 1,
              },
            },
          ],
        },
      ],
      shapes: [
        {
          ty: 'gr',
          it: [
            {
              ty: 'rc',
              d: 1,
              s: {
                a: 0,
                k: [253.5, 0],
                ix: 2,
                x: "var $bm_rt;\nvar temp;\ntemp = effect('Slider Control')('Slider');\n$bm_rt = $bm_sum(value, [\n    0,\n    temp\n]);",
              },
              p: { a: 0, k: [0, 0], ix: 3 },
              r: { a: 0, k: 0, ix: 4 },
              nm: 'Rectangle Path 1',
              mn: 'ADBE Vector Shape - Rect',
              hd: false,
            },
            {
              ty: 'fl',
              c: { a: 0, k: [0.813333429075, 0.813333429075, 0.813333429075, 1], ix: 4 },
              o: { a: 0, k: 100, ix: 5 },
              r: 1,
              bm: 0,
              nm: 'Fill 1',
              mn: 'ADBE Vector Graphic - Fill',
              hd: false,
            },
            {
              ty: 'tr',
              p: { a: 0, k: [-13.75, -85.25], ix: 2 },
              a: { a: 0, k: [0, 0], ix: 1 },
              s: { a: 0, k: [100, 100], ix: 3 },
              r: { a: 0, k: 0, ix: 6 },
              o: { a: 0, k: 100, ix: 7 },
              sk: {
                a: 1,
                k: [
                  { i: { x: [0.8], y: [0.484] }, o: { x: [0.2], y: [-0.017] }, t: 0, s: [20] },
                  { i: { x: [0.933], y: [0.881] }, o: { x: [0.607], y: [0.23] }, t: 30, s: [-12.5] },
                  { i: { x: [0.681], y: [0.993] }, o: { x: [0.123], y: [-0.01] }, t: 39, s: [-79] },
                  { i: { x: [0.254], y: [0.874] }, o: { x: [0.06], y: [-0.04] }, t: 40, s: [83] },
                  { i: { x: [0.833], y: [0.833] }, o: { x: [0.206], y: [0.236] }, t: 50, s: [28.1] },
                  { t: 60, s: [20] },
                ],
                ix: 4,
              },
              sa: { a: 0, k: 0, ix: 5 },
              nm: 'Transform',
            },
          ],
          nm: 'Rectangle 1',
          np: 2,
          cix: 2,
          bm: 0,
          ix: 1,
          mn: 'ADBE Vector Group',
          hd: false,
        },
        {
          ty: 'st',
          c: { a: 0, k: [0.811764765721, 0.811764765721, 0.811764765721, 1], ix: 3 },
          o: { a: 0, k: 100, ix: 4 },
          w: { a: 0, k: 4, ix: 5 },
          lc: 2,
          lj: 2,
          bm: 0,
          nm: 'Stroke 1',
          mn: 'ADBE Vector Graphic - Stroke',
          hd: false,
        },
      ],
      ip: 0,
      op: 121,
      st: 0,
      ct: 1,
      bm: 0,
    },
  ],
  markers: [{ tm: 60, cm: '', dr: 0 }],
  props: {},
};
