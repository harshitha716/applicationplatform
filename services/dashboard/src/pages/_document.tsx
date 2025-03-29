import React from 'react';
import { Head, Html, Main, NextScript } from 'next/document';

export default function Document() {
  return (
    <Html lang='en'>
      <Head>
        {/* <meta
          httpEquiv='Content-Security-Policy'
          content="
              default-src 'self';
              style-src 'unsafe-inline' 'self';
              font-src 'self' data:;
              connect-src 'self' http://localhost:8080;
            "
        /> */}
      </Head>
      <body className='antialiased light-mode'>
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
