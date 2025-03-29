'use client';
import React, { useEffect, useState } from 'react';
import { useGetErrorDetailsQuery } from 'apis/auth';
import { LOGIN_METHODS } from 'constants/auth.constants';
import { ICON_SPRITE_TYPES, ZAMP_ICON_BLACK } from 'constants/icons';
import { LOGIN_ERROR_TEXT } from 'modules/login/constants';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { LoginFlow } from 'types/api/auth.types';
import { SIZE_TYPES } from 'types/common/components';
import { getFromLocalStorage, LOCAL_STORAGE_KEYS, setToLocalStorage } from 'utils/localstorage';
import { Button } from 'components/common/button/Button';
import Input from 'components/common/input';

type LoginFormProps = {
  className?: string;
  loginFlow: LoginFlow;
  setLoginFlow: (loginFlow: LoginFlow) => void;
};

const commonFetchConfig = {
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  },
};

const LoginForm: React.FC<LoginFormProps> = ({ className = '', loginFlow, setLoginFlow }) => {
  const cachedUserEmail = JSON.parse(getFromLocalStorage(LOCAL_STORAGE_KEYS.XZAMP_USER) ?? '{}');
  const router = useRouter();
  const errorId = router.query.error?.toString() ?? '';

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [isEmailLogin, setIsEmailLogin] = useState<boolean>(false);

  const { data: userFacingError } = useGetErrorDetailsQuery(errorId, { skip: !errorId });

  const [email, setEmail] = useState<string>(cachedUserEmail.email ?? '');
  const [password, setPassword] = useState<string>(cachedUserEmail.password ?? '');

  const handlePasswordSubmit = (
    e?: React.FormEvent<HTMLFormElement> | React.MouseEvent<HTMLButtonElement, MouseEvent>,
  ) => {
    setIsEmailLogin(true);
    e?.preventDefault?.();
    if (loginFlow) {
      setLoading(true);

      const csrfNode = loginFlow.ui.nodes.find((node) => {
        const nodeAttributes = node.attributes;

        if ('name' in nodeAttributes && nodeAttributes['name'] === 'csrf_token') {
          return true;
        }

        return false;
      });
      let csrfToken = '';

      if (csrfNode && 'value' in csrfNode.attributes) {
        csrfToken = csrfNode.attributes.value;
      }

      fetch(loginFlow.ui.action, {
        ...commonFetchConfig,
        method: loginFlow.ui.method,
        credentials: 'include',
        body: JSON.stringify({
          password: password,
          csrf_token: csrfToken,
          method: LOGIN_METHODS.PASSWORD,
          identifier: email,
        }),
      }).then((response) => {
        return response
          .json()
          .then((responseJson) => {
            setToLocalStorage(LOCAL_STORAGE_KEYS.XZAMP_USER, JSON.stringify({ email, password }));

            if (response.status < 300) {
              window.location.reload();

              return;
            }
            if (response.status === 400) {
              setLoginFlow(responseJson);
              setError(responseJson?.ui?.messages?.[0]?.text ?? LOGIN_ERROR_TEXT);
              setLoading(false);
            }
          })
          ?.catch(() => {
            setError(LOGIN_ERROR_TEXT);
          });
      });
    }
  };

  useEffect(() => {
    const token = localStorage.getItem('token');

    if (token) {
      router.push('/payments');
    }
  }, []);

  const formDisabled = loading || !loginFlow;

  return (
    <div className={`w-96 mx-auto mt-[30vh] items-center flex flex-col gap-10 ${className}`}>
      <Image src={ZAMP_ICON_BLACK} width={48} height={38} alt='Zamp' priority />

      {userFacingError &&
        userFacingError.map((error, index) => (
          <div key={index} className='text-red-600'>
            {error.message}
          </div>
        ))}
      <form className='flex flex-col gap-3 w-full' onSubmit={handlePasswordSubmit}>
        <Input
          id='login-email'
          label='Email'
          required
          placeholder='Enter your email'
          name='email'
          type='email'
          value={email}
          onChange={(e) => {
            if (e?.target?.value !== undefined) {
              setEmail(e.target.value);
            }
          }}
          disabled={formDisabled}
        />
        <Input
          id='login-password'
          label='Password'
          required
          disabled={formDisabled}
          placeholder='Enter your password'
          type='password'
          name='password'
          value={password}
          onChange={(e) => {
            if (e?.target?.value !== undefined) {
              setPassword(e.target.value);
            }
          }}
          error={error ?? ''}
        />
        <Button
          id='login'
          className='w-fit'
          disabled={formDisabled}
          size={SIZE_TYPES.LARGE}
          iconProps={{
            id: 'arrow-right',
            iconCategory: ICON_SPRITE_TYPES.ARROWS,
          }}
          isLoading={isEmailLogin ? loading : false}
        >
          Login
        </Button>
      </form>
    </div>
  );
};

export default LoginForm;
