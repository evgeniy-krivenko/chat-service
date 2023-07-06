export interface IMessagesResponse {
  readonly messages: IMessage[];
  readonly next: string;
}

export interface IMessage {
  readonly id: string;
  readonly body: string;
  readonly createdAt: Date;
  readonly authorId: string;
  userIsAuthor?: boolean;
  authorName?: string;
}

export interface IChat {
  readonly clientId: string;
  readonly clientName: string;
}

export interface ISendMessage {
  readonly messageBody: string;
}

export interface ISendMessageResponse {
  readonly id: string;
  readonly createdAt: Date;
  readonly authorId: string;
}

export interface IBackApiResponse<T> {
  readonly data?: T
  readonly error?: {
    message: string;
    code: number;
  }
}

export interface IGetMessagesReq {
  chatId: string;
  cursor?: string;
  pageSize?: number;
}
