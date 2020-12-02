import inquirer from 'inquirer';
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { LoginPayload } from 'matrix-js-sdk';
import {
  loginWithCredentials,
  loginWithToken,
  registerWithCredentials,
} from './matrix';

export interface UserCredentials {
  username: string;
  password: string;
}

const getCredentials = (): Promise<UserCredentials> => {
  return inquirer.prompt<UserCredentials>([
    {
      type: 'input',
      message: 'Enter username:',
      name: 'username',
    },
    {
      type: 'password',
      message: 'Enter password, at least 8 characters long:',
      name: 'password',
      mask: '*',
      validate: (val) =>
        val.length >= 8 ? true : 'Password must be at least 8 characters long',
    },
  ]);
};

const getToken = async (): Promise<string> =>
  (
    await inquirer.prompt<{ token: string }>([
      {
        type: 'input',
        message: 'Token:',
        name: 'token',
      },
    ])
  ).token;

export const runAuth = async (): Promise<LoginPayload> => {
  enum AuthChoice {
    login_credentials,
    login_token,
    register_credentials,
  }
  const { choice } = await inquirer.prompt<{ choice: AuthChoice }>([
    {
      type: 'list',
      message: 'Authorization required',
      name: 'choice',
      choices: [
        {
          name: 'Log in to existing account using credentials',
          value: AuthChoice.login_credentials,
        },
        {
          name: 'Log in to existing account using token',
          value: AuthChoice.login_token,
        },
        { name: 'Create new account', value: AuthChoice.register_credentials },
      ],
    },
  ]);

  switch (choice) {
    case AuthChoice.login_credentials: {
      const creds = await getCredentials();
      return loginWithCredentials(creds);
    }
    case AuthChoice.register_credentials: {
      const creds = await getCredentials();
      return registerWithCredentials(creds);
    }
    case AuthChoice.login_token: {
      const token = await getToken();
      return loginWithToken(token);
    }
  }
};
