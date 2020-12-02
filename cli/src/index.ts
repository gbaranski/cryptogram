// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { LoginPayload } from '@types/matrix-js-sdk';
import { runApp } from './app';
import { runAuth } from './auth';

(async () => {
  const loginPayload: LoginPayload = await runAuth();
  runApp(loginPayload);
})();
