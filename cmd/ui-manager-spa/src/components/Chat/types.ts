import { IMessage } from '../../types/messages';

export interface FormElements extends HTMLFormControlsCollection {
  message: HTMLInputElement;
}

export interface MessageForm extends HTMLFormElement {
  readonly elements: FormElements;
}

export interface NewMessageEvent extends Omit<IMessage, 'id'> {
  eventType: 'NewMessageEvent';
  messageId: string;
  chatId: string;
}
