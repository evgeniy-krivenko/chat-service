import { IManager } from './user';

export interface ILoginRequest {
  readonly login: string;
  readonly password: string;
}

export interface ILoginResponse {
  user: IManager;
  token: string;
}
