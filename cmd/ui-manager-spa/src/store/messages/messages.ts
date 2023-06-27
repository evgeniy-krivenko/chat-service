import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { devtools } from 'zustand/middleware';
import { AxiosResponse } from 'axios';
import {
  IBackApiResponse, IGetMessagesReq, IMessage, IMessagesResponse, ISendMessageResponse,
} from '../../types/messages';
import $api from '../../api/index';
import { GET_CHAT_HISTORY, SEND_MESSAGE } from '../../const/urls';

const DEFAULT_PAGE_SIZE = 10;

export interface IMessagesState {
  loading: boolean;
  messages: Map<string, IMessage>;
  error: string;
  cursor: string;
  getMessages: (chatId: string, authorId: string) => Promise<void>;
  addMessage: (msg: IMessage) => void;
  sendMessage: (chatId: string, messageBody: string) => Promise<void>;
  resetMessages: () => void;
}

const useMessages = create<IMessagesState>()(immer(devtools((set, get) => ({
  loading: false,
  messages: new Map(),
  error: null,
  cursor: null,
  resetMessages: () => {
    set({ messages: new Map<string, IMessage>() });
  },
  addMessage: (msg: IMessage): void => {
    set((state) => {
      state.messages.set(msg.id, msg);
    });
  },
  getMessages: async (chatId: string, authorId: string) => {
    try {
      const req: IGetMessagesReq = { chatId };
      const { cursor } = get();

      if (cursor) {
        req.cursor = cursor;
      } else {
        req.pageSize = DEFAULT_PAGE_SIZE;
      }

      const {
        data: { data },
      }: AxiosResponse<IBackApiResponse<IMessagesResponse>> = await $api.post(GET_CHAT_HISTORY, req);

      const newMessages = data?.messages
        .map((m) => ({ ...m, userIsAuthor: m.authorId === authorId }))
        .reverse();

      const { messages } = get();

      // eslint-disable-next-line no-restricted-syntax
      for (const msg of newMessages) {
        messages.set(msg.id, msg);
      }

      set((state) => {
        state.messages = messages;
        state.cursor = data.next;
      });
    } catch (e) {
      set({ error: e.message });
    } finally {
      set({ loading: false });
    }
  },
  sendMessage: async (chatId: string, messageBody: string) => {
    try {
      set({ loading: true });
      const { data: { data } }: AxiosResponse<IBackApiResponse<ISendMessageResponse>> = await $api.post(
        SEND_MESSAGE,
        { chatId, messageBody },
      );

      set((state) => {
        state.messages.set(data.id, {
          id: data.id,
          createdAt: data.createdAt,
          authorId: data.authorId,
          userIsAuthor: true,
          body: messageBody,
        });
      });
    } catch (e) {
      set({ error: e.message });
    } finally {
      set({ loading: false });
    }
  },
}))));

export default useMessages;
