import React, { FC, createContext, useState, JSX } from 'react';

export interface IAuthProvider {
  user?: IUser;
  signIn: ISignIn;
  signOut: ISignOut;
}

const initialValue: IAuthProvider = {
  user: undefined,
  signOut: () => null,
  signIn: () => null,
}

export type ISignIn = (user: IUser, cb: () => void) => void

export type ISignOut = (cb: () => void) => void

export interface IUser {
  username: string;
}

export const AuthContext = createContext<IAuthProvider>(initialValue);

export interface AuthProviderProps {
  children: JSX,
}

const AuthProvider: FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<IUser>(null);

  const signIn = (newUser: IUser, cb) => {
    setUser(newUser);
    cb();
  }

  const signOut = (cb) => {
    setUser(null);
    cb();
  }

  const value = { user, signIn, signOut, };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
