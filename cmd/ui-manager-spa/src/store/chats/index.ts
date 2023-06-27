import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { toast } from 'react-toastify';
import { devtools } from 'zustand/middleware';
import { AxiosResponse } from 'axios';
import { IBackApiResponse } from '../../types/messages';
import { FREE_HANDS, GET_CHATS, GET_FREE_HANDS_BTN_AVAILABILITY } from '../../const/urls';
import $api from '../../api/index';
import { IGetFreeHandsBtn } from '../../types/chats';
import { TOAST } from '../../const/toast.config';

export interface IChat {
  chatId: string;
  clientId: string;
  firstName?: string;
  lastName?: string;
}

export interface IChatState {
  loading: boolean;
  chats: IChat[];
  error: string;
  getChats: () => void;
  addChat: (chat: IChat, canTakeMoreProblems: boolean) => void;
  removeChat: (chatId: string, canTakeMoreProblems: boolean) => void;
  resetError: () => void;
  canTakeMoreProblems: boolean;
  setCanTakeMoreProblems: (value: boolean) => void;
  getFreeHandsBtnAvailability: () => void;
  freeHands: () => void;
  freeHandsLoading: boolean;
  freeHandsError: string;
}

export const useChats = create<IChatState>()(immer(devtools((set) => ({
  loading: false,
  error: null,
  chats: [],
  canTakeMoreProblems: false,
  freeHandsLoading: false,
  freeHandsError: '',
  addChat: (chat: IChat, canTakeMoreProblems: boolean) => {
    set((state) => {
      state.chats.push(chat);
      state.canTakeMoreProblems = canTakeMoreProblems;
    });
  },
  removeChat: (chatId: string, canTakeMoreProblems: boolean) => {
    set((state) => {
      state.chats = state.chats.filter((c) => c.chatId !== chatId);
      state.canTakeMoreProblems = canTakeMoreProblems;
    });
  },
  setCanTakeMoreProblems: (value: boolean) => {
    set((state) => {
      state.canTakeMoreProblems = value;
    });
  },
  getChats: async () => {
    set({ loading: true });
    try {
      const { data }: AxiosResponse<IBackApiResponse<Record<'chats', IChat[]>>> = await $api.post(GET_CHATS);
      set({ chats: data.data?.chats });
    } catch (e) {
      set({ chats: [], error: e.message });
    } finally {
      set({ loading: false });
    }
  },
  getFreeHandsBtnAvailability: () => {
    $api.post<IBackApiResponse<IGetFreeHandsBtn>>(GET_FREE_HANDS_BTN_AVAILABILITY)
      .then(({ data } :AxiosResponse) => {
        set({ canTakeMoreProblems: data.data.available });
      })
      .catch((e) => toast.error(e.message, TOAST));
  },
  resetError: () => set({ error: null }),
  freeHands: () => {
    set({ freeHandsLoading: true });
    $api.post(FREE_HANDS)
      .then(() => set({ freeHandsError: '' }))
      .catch((e) => {
        set({ freeHandsError: e.message });
        toast.error(e.message, TOAST);
      })
      .finally(() => set({ freeHandsLoading: false }));
  },
}))));
