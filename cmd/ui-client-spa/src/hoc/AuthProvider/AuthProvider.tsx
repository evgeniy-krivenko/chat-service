import React, {FC, createContext, useState, JSX, useEffect} from 'react';
import {IUserProfile} from "../../types/user";
import {APIClient} from "../../api";

export interface IAuthProvider {
  user?: IUserProfile;
  signIn: ISignIn;
  signOut: ISignOut;
}

const initialValue: IAuthProvider = {
  user: undefined,
  signOut: () => null,
  signIn: () => null,
}

export type ISignIn = (payload: IUserProfile, cb: () => void) => void

export type ISignOut = (cb: () => void) => void


export const AuthContext = createContext<IAuthProvider>(initialValue);

export interface AuthProviderProps {
  children: JSX,
}

const AuthProvider: FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<IUserProfile>(null);

  useEffect(() => {
    APIClient.getUserProfile()
      .then((user) => setUser(user))
      .catch(() => setUser(null))
  }, [])

  const signIn = (newUser: IUserProfile, cb) => {
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
