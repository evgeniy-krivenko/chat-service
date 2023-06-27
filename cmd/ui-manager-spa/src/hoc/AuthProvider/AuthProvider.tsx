import {
  FC, createContext, useEffect, ReactNode,
} from 'react';
import { shallow } from 'zustand/shallow';
import { IManager } from '../../types/user';
import { useManagersStore } from '../../store/manager';
import Loader from '../../components/Loader/Loader';

export interface IAuthProvider {
  manager?: IManager;

}

const initialValue: IAuthProvider = {
  manager: undefined,
};

export type ISignIn = (payload: IManager, cb: () => void) => void

export type ISignOut = (cb: () => void) => void

export const AuthContext = createContext<IAuthProvider>(initialValue);

export interface AuthProviderProps {
  children: ReactNode,
}

const AuthProvider: FC<AuthProviderProps> = ({ children }) => {
  const {
    manager, getManagerProfile, loading,
  } = useManagersStore((state) => ({
    manager: state.manager,
    error: state.error,
    getManagerProfile: state.getManagerProfile,
    loading: state.loading,
  }), shallow);

  useEffect(() => {
    getManagerProfile();
  }, []);

  // const signIn = (newUser: IManager, cb) => {
  //   setUser(newUser);
  //   cb();
  // };
  //
  // const signOut = (cb) => {
  //   setUser(null);
  //   cb();
  // };

  const value = { manager };

  if (loading) {
    return <Loader />;
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
