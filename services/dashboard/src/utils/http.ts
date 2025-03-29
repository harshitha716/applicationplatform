export const getCookieByName = (name: string): string => {
  const cookie: any = {};

  document.cookie.split(';').forEach(function (el) {
    const [k, v] = el.split('=');

    cookie[k.trim()] = v;
  });

  return cookie[name];
};
