import { IMessage } from './messages';

// eslint-disable-next-line no-shadow
export enum MessageEvent {
  NewMessageEvent = 'NewMessageEvent',
  NewChatEvent = 'NewChatEvent',
  ChatClosedEvent = 'ChatClosedEvent',
}

export interface NewMessageEvent extends Omit<IMessage, 'id'> {
  eventType: MessageEvent.NewMessageEvent;
  messageId: string;
  chatId: string;
}

export interface NewChatEvent {
  eventType: MessageEvent.NewChatEvent;
  chatId: string;
  clientId: string;
  firstName: string;
  lastName: string;
  canTakeMoreProblems: boolean;
}

export interface ChatClosedEvent {
  eventType: MessageEvent.ChatClosedEvent;
  chatId: string;
  canTakeMoreProblems: boolean;
}

export type Events = NewMessageEvent | NewChatEvent |
  ChatClosedEvent
