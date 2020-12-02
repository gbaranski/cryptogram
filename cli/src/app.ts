// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { LoginPayload } from '@types/matrix-js-sdk';
import inquirer from 'inquirer';
import { getRooms, listenToEvents, startClient, waitForSync } from './matrix';

enum Action {
  send_message,
  join_room,
  get_rooms,
}

const sendMessage = async () => {
  const message = await inquirer.prompt<{ roomID: string; content: string }>([
    {
      type: 'input',
      message: 'RoomID',
      name: 'roomID',
    },
    {
      type: 'input',
      message: 'Content of message',
      name: 'content',
    },
  ]);
};

export const joinRoom = async () => {
  console.log('To be implemented');
};

export const askForAction = async (): Promise<Action> => {
  return (
    await inquirer.prompt<{ action: Action }>([
      {
        type: 'rawlist',
        message: 'What to do?',
        name: 'action',
        choices: [
          {
            name: 'Send message',
            value: Action.send_message,
          },
          {
            name: 'Join room',
            value: Action.join_room,
          },
          {
            name: 'Get available rooms',
            value: Action.get_rooms,
          },
        ],
      },
    ])
  ).action;
};

const handlerInterval = async (): Promise<void> => {
  const action = await askForAction();
  switch (action) {
    case Action.send_message: {
      await sendMessage();
      break;
    }
    case Action.get_rooms: {
      console.log(`Rooms: ${getRooms()}`);
      break;
    }
  }
  console.log();
  handlerInterval();
};

export const runApp = async (loginPayload: LoginPayload): Promise<void> => {
  console.log('Initializing...');
  await startClient();
  await waitForSync();
  listenToEvents((e) => console.log(e));
  console.log('Successfully initialized');
  handlerInterval();
};
