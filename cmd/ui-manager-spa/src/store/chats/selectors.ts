import { IChatState } from './index';

export const selectChats = (state: IChatState): Omit<IChatState, 'resetError'> => ({
  loading: state.loading,
  chats: state.chats,
  error: state.error,
  getChats: state.getChats,
  getFreeHandsBtnAvailability: state.getFreeHandsBtnAvailability,
  addChat: state.addChat,
  canTakeMoreProblems: state.canTakeMoreProblems,
  removeChat: state.removeChat,
  setCanTakeMoreProblems: state.setCanTakeMoreProblems,
  freeHands: state.freeHands,
  freeHandsLoading: state.freeHandsLoading,
  freeHandsError: state.freeHandsError,
  closeChat: state.closeChat,
});
