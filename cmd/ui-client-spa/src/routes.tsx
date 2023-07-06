import { JSX } from "react";
import Login from "./components/Login/Login";
import Chat from "./components/Chat/Chat";

export interface IRoutes {
  path: string;
  Component: JSX;
}

export const publicRoutes: IRoutes[] = [
  {
    path: '/login',
    Component: Login,
  },
];

export const privateRoutes: IRoutes[] = [
  {
    path: '/',
    Component: Chat,
  },
];
