/* eslint-disable @typescript-eslint/ban-ts-comment */
import _matrix from 'matrix-js-sdk';
// @ts-ignore
import { LoginPayload, MatrixEvent, Room } from '@types/matrix-js-sdk';
// @ts-ignore
const matrix: typeof import('@types/matrix-js-sdk') = _matrix;

type fuckingTypesAreMissing = unknown;

import { UserCredentials } from './auth';

const client = matrix.createClient('https://matrix.org');

export const loginWithCredentials = (
  credentials: UserCredentials,
): Promise<LoginPayload> =>
  client.loginWithPassword(credentials.username, credentials.password);

export const loginWithToken = (token: string): Promise<LoginPayload> =>
  client.loginWithToken(token);

export const registerWithCredentials = (
  credentials: UserCredentials,
): Promise<LoginPayload> =>
  client.register(credentials.username, credentials.password);

export const startClient = (): Promise<void> => client.startClient();
export const waitForSync = (): Promise<void> =>
  new Promise((resolve) =>
    client.on(
      'sync',
      (
        state: fuckingTypesAreMissing,
        prevState: fuckingTypesAreMissing,
        res: fuckingTypesAreMissing,
      ) => {
        if (state === 'PREPARED' || prevState === 'PREPARED') resolve();
      },
    ),
  );

export const getRooms = (): Room[] => client.getRooms();

export const listenToEvents = (cb: (e: MatrixEvent) => any): void => {
  client.on('event', cb);
};

export const listenToRoomEvents = (
  cb: (
    event: fuckingTypesAreMissing,
    room: Room,
    toStartOfTimeline: fuckingTypesAreMissing,
  ) => unknown,
  room: string,
): void => {
  client.on(`Room.${room}`, cb);
};

export const joinRoom = (roomName: string) => client.joinRoom('');
