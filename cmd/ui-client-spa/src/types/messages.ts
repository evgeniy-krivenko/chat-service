export interface IMessagesResponse {
  readonly messages: IMessage[];
  readonly next: string;
}

export interface IMessage {
  readonly id: string;
  readonly body: string;
  readonly createdAt: Date;
  readonly isReceived: boolean;
  readonly isBlocked: boolean;
  readonly isService: boolean;
  readonly authorId?: string;
  userIsAuthor?: boolean;
  readonly authorName?: string;
}

export interface ISendMessage {
  readonly messageBody: string;
}

export interface ISendMessageResponse {
  readonly id: string;
  readonly createdAt: Date;
  readonly authorId: string;
}
