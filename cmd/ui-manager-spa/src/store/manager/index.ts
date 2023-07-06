import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { AxiosResponse } from 'axios';
import { devtools } from 'zustand/middleware';
import { IManager } from '../../types/user';
import { GET_MANAGER_PROFILE, LOGIN } from '../../const/urls';
import { IBackApiResponse } from '../../types/messages';
import { ILoginResponse } from '../../types/login';
import $api from '../../api/index';

export interface IManagerStore {
  manager: IManager;
  loading: boolean;
  error: string;
  getManagerProfile: () => void;
  login: (login: string, password: string) => Promise<void>;
  resetError: () => void;
}

export const useManagersStore = create<IManagerStore>()(immer(devtools((set) => ({
  manager: {} as IManager,
  loading: false,
  error: '',
  getManagerProfile: async () => {
    set({ loading: true });
    try {
      const { data }: AxiosResponse<IBackApiResponse<IManager>> = await $api.post(
        GET_MANAGER_PROFILE,
      );
      set({ manager: data.data });
    } catch (e) {
      set({ manager: null, error: e.message });
    } finally {
      set({ loading: false });
    }
  },
  login: async (login: string, password: string) => {
    set({ loading: true });
    try {
      const { data }: AxiosResponse<IBackApiResponse<ILoginResponse>> = await $api.post(LOGIN, {
        login,
        password,
      });
      localStorage.setItem('token', data.data.token);
      set({ manager: data.data.user });
    } catch (e) {
      set({ manager: null, error: e.message });
    } finally {
      set({ loading: false });
    }
  },
  resetError: () => set({ error: null }),
}))));
