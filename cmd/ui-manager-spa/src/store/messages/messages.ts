import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { devtools } from 'zustand/middleware';
import { AxiosResponse } from 'axios';
import {
  IBackApiResponse, IGetMessagesReq, IMessage, IMessagesResponse, ISendMessageResponse,
} from '../../types/messages';
import $api from '../../api/index';
import { GET_CHAT_HISTORY, SEND_MESSAGE } from '../../const/urls';

const DEFAULT_PAGE_SIZE = 20;

export interface IMessagesState {
  loading: boolean;
  messages: IMessage[];
  error: string;
  cursor: string;
  getMessages: (chatId: string, authorId: string) => Promise<void>;
  addMessage: (msg: IMessage) => void;
  sendMessage: (chatId: string, messageBody: string) => Promise<void>;
  resetMessages: () => void;
}

const useMessages = create<IMessagesState>()(immer(devtools((set, get) => ({
  loading: false,
  messages: [],
  error: null,
  cursor: null,
  resetMessages: () => {
    set({ messages: [] });
  },
  addMessage: (msg: IMessage): void => {
    set((state) => {
      if (state.messages.findIndex((m) => m.id === msg.id) === -1) {
        state.messages.push(msg);
      }
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

      const existingMsgIds = new Set<string>();

      // eslint-disable-next-line no-restricted-syntax
      for (const msg of messages) {
        existingMsgIds.add(msg.id);
      }

      set((state) => {
        state.messages = [...messages, ...newMessages.filter((m) => !existingMsgIds.has(m.id))];
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
        if (state.messages.findIndex((m) => m.id === data.id) === -1) {
          state.messages.push({
            id: data.id,
            createdAt: data.createdAt,
            authorId: data.authorId,
            userIsAuthor: true,
            body: messageBody,
          });
        }
      });
    } catch (e) {
      set({ error: e.message });
    } finally {
      set({ loading: false });
    }
  },
}))));

export default useMessages;
