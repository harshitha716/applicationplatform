import React from 'react';
import { API_ENDPOINTS, REQUEST_TYPES } from 'apis/apiEndpoint.constants';
import { API_DOMAIN } from 'constants/api.constants';
import { LOGIN_PROVIDERS } from 'constants/auth.constants';
import { ZAMP_FULL_LOGO, ZAMP_LOGIN_BG } from 'constants/icons';
import { LOGIN_ERROR_TEXT } from 'modules/login/constants';
import LocaldevEmailPasswordLogin from 'modules/login/LocaldevEmailPasswordLogin';
import { LOGIN_GROUPS, VALID_SESSION_DETECTED_ERROR_MSG } from 'modules/login/login.constants';
import LoginButton from 'modules/login/LoginButton';
import Image from 'next/image';
import { LoginFlow } from 'types/api/auth.types';
import { SIZE_TYPES } from 'types/common/components';
import { getDomainFromEmail, isValidEmail } from 'utils/common';
import { getFromLocalStorage, LOCAL_STORAGE_KEYS, removeFromLocalStorage, setToLocalStorage } from 'utils/localstorage';
import Input from 'components/common/input';

export const LoginForm = () => {
  const [email, setEmail] = React.useState(getFromLocalStorage(LOCAL_STORAGE_KEYS.LAST_LOGGED_IN_OIDC_EMAIL) ?? '');
  const [loginFlow, setLoginFlow] = React.useState<LoginFlow | null>(null);
  const [error, setError] = React.useState<string | null>(null);
  const [loading, setLoading] = React.useState<boolean>(false);
  const [hasError, setHasError] = React.useState<boolean>(false);

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e?.target?.value !== undefined) {
      setEmail(e.target.value);
    }
  };

  const initiateOidcLogin = async (url: string, method: string, providerId: LOGIN_PROVIDERS) => {
    setLoading(true);
    try {
      const resp = await fetch(url, {
        method: method,
        body: JSON.stringify({
          provider: providerId,
        }),
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
      });
      const respJson = await resp.json();

      const validSessionMsg = respJson?.ui?.messages?.[0]?.text.includes(VALID_SESSION_DETECTED_ERROR_MSG);

      if (resp.status === 422 || resp.status === 200 || (resp.status === 400 && validSessionMsg)) {
        setToLocalStorage(LOCAL_STORAGE_KEYS.LAST_LOGGED_IN_OIDC_EMAIL, email);

        const redirectUrl = respJson.redirect_browser_to;

        try {
          const url = new URL(redirectUrl);
          const emailDomain = getDomainFromEmail(email);

          setHasError(false);
          url.searchParams.set('hd', emailDomain);
          window.location.href = url.toString();
        } catch (error) {
          setLoading(false);
          console.error(error);
          setHasError(true);
        }
      } else if (resp.status === 400) {
        setError(respJson?.ui?.messages?.[0]?.text ?? LOGIN_ERROR_TEXT);
        setHasError(true);
      } else {
        setError(respJson?.error?.message ?? LOGIN_ERROR_TEXT);
        setHasError(true);
        setLoading(false);
        setLoginFlow(respJson);
      }
    } catch (error) {
      setLoading(false);
      console.error(error);
      setHasError(true);
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e?.preventDefault();
    setError(null);
    setLoading(true);
    if (!isValidEmail(email)) {
      setError('Please enter a valid email address');
      setLoading(false);

      return;
    }

    try {
      const response = await fetch(`${API_DOMAIN}/${API_ENDPOINTS.AUTH_INITIAL_LOGIN_FLOW_BY_EMAIL_POST}`, {
        method: REQUEST_TYPES.POST,
        body: JSON.stringify({
          email,
        }),
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
        credentials: 'include',
      });

      const respJson = await response.json();

      setHasError(false);

      if (response.status !== 200) {
        setError(respJson.error);
        removeFromLocalStorage(LOCAL_STORAGE_KEYS.LAST_LOGGED_IN_OIDC_EMAIL);
        setHasError(true);
        setLoading(false);

        return;
      }

      setLoginFlow(respJson);

      // if the number of login methods is 1 and it is OIDC, we can directly login
      if (respJson?.ui?.nodes?.length == 1) {
        const loginNode = respJson.ui.nodes[0];

        if (loginNode?.group === LOGIN_GROUPS.OIDC) {
          await initiateOidcLogin(
            respJson.ui.action,
            respJson.ui.method,
            loginNode.attributes.value as LOGIN_PROVIDERS,
          );
        }
      } else {
        setLoading(false);
      }
    } catch {
      setHasError(true);
      setLoading(false);
    }
  };

  const inputDisabled = loading;

  if (!hasError && loginFlow && loginFlow?.ui?.nodes?.length > 1) {
    return <LocaldevEmailPasswordLogin loginFlow={loginFlow} setLoginFlow={setLoginFlow} />;
  }

  return (
    <div className='relative flex items-center justify-center w-screen h-screen bg-BG_GRAY_5'>
      <video autoPlay muted loop className='absolute z-0 w-full h-full object-cover'>
        <source src={ZAMP_LOGIN_BG} type='video/mp4' />
        <span className='f-14-400 text-GRAY_1000'>Your browser does not support the video tag.</span>
      </video>
      <div className='bg-white z-50 w-[580px] rounded-4.5 shadow-tableFilterMenu px-16 py-[82px] border border-GRAY_100'>
        <Image src={ZAMP_FULL_LOGO} alt='ZAMP' width={98} height={24} />
        <form onSubmit={handleSubmit}>
          <Input
            id='login-email'
            placeholder='Enter your email address'
            className='mt-10'
            name='email'
            type='email'
            value={email}
            error={error ? error : ''}
            autoFocus
            onChange={handleEmailChange}
            disabled={inputDisabled}
            size={SIZE_TYPES.LARGE}
          />
          <LoginButton loading={loading} onClick={() => handleSubmit} />
        </form>
      </div>
    </div>
  );
};
