import {IHistoryRequest} from "../types/history";
import { v4 as uuidv4 } from 'uuid';
import {BASE_URL, GET_HISTORY, GET_USER_PROFILE, LOGIN, SEND_MESSAGE} from "../const/urls";
import {ILoginRequest} from "../types/login";
import {IUserProfile} from "../types/user";
import {IMessagesResponse, ISendMessage, ISendMessageResponse} from "../types/messages";

const defaultHistoryPageSize = 10;

export class APIClient {

  static async getHistory(cursor: string): Promise<IMessagesResponse> {
    const request: IHistoryRequest = {};
    if (cursor) {
      request.cursor = cursor;
    } else {
      request.pageSize = defaultHistoryPageSize;
    }

    const response = await fetch(BASE_URL + GET_HISTORY, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8',
        'Authorization': 'Bearer ' + this.getToken(),
        'X-Request-ID': uuidv4(),
      },
      body: JSON.stringify(request),
    });
    return await this.extractData(response);
  }

  static async sendMessage(msgBody): Promise<ISendMessageResponse> {
    const response = await fetch(BASE_URL + SEND_MESSAGE, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8',
        'Authorization': 'Bearer ' + this.getToken(),
        'X-Request-ID': uuidv4(),
      },
      body: JSON.stringify({
        messageBody: msgBody,
      }),
    });
    return await this.extractData(response);
  }

  static async login(payload: ILoginRequest) {
    const response = await fetch(BASE_URL + LOGIN, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8',
        'X-Request-ID': uuidv4(),
      },
      body: JSON.stringify(payload),
    });
    return await this.extractData(response);
  }

  static async getUserProfile(): Promise<IUserProfile> {
    const response = await fetch(BASE_URL + GET_USER_PROFILE, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8',
        'Authorization': 'Bearer ' + this.getToken(),
        'X-Request-ID': uuidv4(),
      },
    });
    return await this.extractData(response);
  }

  static getToken(): string {
    return localStorage.getItem('token') || '';
  }

  static async extractData(response) {
    if (!response.ok) {
      throw new Error(`${response.status}`);
    }

    const result = await response.json();
    if (result.error) {
      throw new Error(`${result.error.code}: ${result.error.message}`);
    }
    return result.data;
  }
}
