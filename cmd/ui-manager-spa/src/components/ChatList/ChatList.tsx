import {
  FC, SyntheticEvent, useEffect, useState,
} from 'react';
import {
  Alert,
  Button,
  CircularProgress,
  Container,
  Grid,
  Tooltip,
} from '@mui/material';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';
import useWebSocket from 'react-use-websocket';
import { shallow } from 'zustand/shallow';
import { green } from '@mui/material/colors';
import cn from 'classnames';
import { IChat, useChats } from '../../store/chats/index';
import { selectChats } from '../../store/chats/selectors';
import Chat from '../Chat';
import './ChatList.css';
import { WS_ENDPOINT, WS_PROTOCOL } from '../../config';
import { Events } from '../../types/events';

const buttonSx = {
  bgcolor: green[500],
  '&:hover': {
    bgcolor: green[700],
  },
  marginTop: '8px',
  width: '100%',
  height: '48px',
};

const ChatList: FC = () => {
  const [currentChat, setCurrentChat] = useState<IChat | null>(null);
  const {
    chats,
    getChats,
    addChat,
    removeChat,
    canTakeMoreProblems,
    getFreeHandsBtnAvailability,
    freeHands,
    freeHandsLoading,
    freeHandsError,
    closeChat,
  } = useChats(selectChats, shallow);

  const token = localStorage.getItem('token');
  const { lastMessage } = useWebSocket(WS_ENDPOINT, {
    protocols: [WS_PROTOCOL, token],
  });

  useEffect(() => {
    getChats();
    getFreeHandsBtnAvailability();
  }, []);

  useEffect(() => {
    if (!chats.find((c) => c.chatId === currentChat?.chatId)) {
      setCurrentChat(null);
    }
  }, [chats]);

  useEffect(() => {
    if (lastMessage) {
      const msg: Events = JSON.parse(lastMessage.data);

      if (msg.eventType === 'NewChatEvent') {
        addChat(
          {
            chatId: msg.chatId,
            firstName: msg.firstName,
            lastName: msg.lastName,
            clientId: msg.clientId,
          },
          msg.canTakeMoreProblems,
        );
      }

      if (msg.eventType === 'ChatClosedEvent') {
        removeChat(msg.chatId, msg.canTakeMoreProblems);
      }
    }
  }, [lastMessage]);

  const closeChatHandler = (chatId: string) => (e: SyntheticEvent) => {
    e.stopPropagation();

    closeChat(chatId);
  };

  return (
    <Container>
      <h1 className="chat-list__title">Chats</h1>
      <Grid justifyContent="center" container>
        <Grid xs={4} item>
          <h4 className="chat-list__subtitle">Open Problems</h4>
          <div className="chat-list__chats">
            {chats.map((chat) => (
              <div
                className={cn('chat-list__item', {
                  active: currentChat?.chatId === chat.chatId,
                })}
                role="button"
                tabIndex={0}
                key={chat.chatId}
                onClick={() => setCurrentChat(chat)}
              >
                <span>{`${chat.firstName} ${chat.lastName}`}</span>
                <Tooltip title="Close Problem">
                  <HighlightOffIcon
                    className="chat-list__close"
                    role="button"
                    tabIndex={0}
                    onClick={closeChatHandler(chat.chatId)}
                  />
                </Tooltip>
              </div>
            ))}
          </div>
          <Button
            sx={buttonSx}
            onClick={() => freeHands()}
            type="submit"
            variant="contained"
            disabled={!canTakeMoreProblems || freeHandsLoading}
          >
            Ready to work 🫡
            {freeHandsLoading && (
              <CircularProgress
                size={24}
                sx={{
                  color: green[500],
                  position: 'absolute',
                  top: '50%',
                  left: '50%',
                  marginTop: '-12px',
                  marginLeft: '-12px',
                }}
              />
            )}
          </Button>
          {freeHandsError && (
            <Alert style={{ margin: '24px 0' }} severity="error">
              {freeHandsError}
            </Alert>
          )}
        </Grid>
        <Grid xs={8} item>
          {currentChat ? (
            <Chat
              chatId={currentChat.chatId}
              authorName={`${currentChat.firstName || ''} ${currentChat.lastName || ''}`}
              lastMessage={lastMessage}
            />
          ) : (
            <div className="chat-list__no-active-text">Select open problem or push on Ready to work 👈</div>
          )}
        </Grid>
      </Grid>
    </Container>
  );
};

export default ChatList;
