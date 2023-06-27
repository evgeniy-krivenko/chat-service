import { FC, useEffect, useState } from 'react';
import {
  Alert,
  Button, CircularProgress, Container, Grid,
} from '@mui/material';
import useWebSocket from 'react-use-websocket';
import { shallow } from 'zustand/shallow';
import { green } from '@mui/material/colors';
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
  const [currentChat, setCurrentChat] = useState<IChat>(null);
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
  } = useChats(selectChats, shallow);

  const token = localStorage.getItem('token');
  const { lastMessage } = useWebSocket(WS_ENDPOINT, { protocols: [WS_PROTOCOL, token] });

  useEffect(() => {
    getChats();
    getFreeHandsBtnAvailability();
  }, []);

  useEffect(() => {
    if (lastMessage) {
      const msg: Events = JSON.parse(lastMessage.data);

      if (msg.eventType === 'NewChatEvent') {
        addChat({
          chatId: msg.chatId,
          firstName: msg.firstName,
          lastName: msg.lastName,
          clientId: msg.clientId,
        }, msg.canTakeMoreProblems);
      }

      if (msg.eventType === 'ChatClosedEvent') {
        removeChat(msg.chatId, msg.canTakeMoreProblems);
      }
    }
  }, [lastMessage]);

  return (
    <Container>
      <h1 className="chat-list__title">Chats</h1>
      <Grid justifyContent="center" container>
        <Grid xs={4} item>
          <h4 className="chat-list__subtitle">Open Problems</h4>
          <div className="chat-list__chats">
            {chats.map((chat) => (
              <div
                className="chat-list__item"
                role="button"
                tabIndex={0}
                key={chat.chatId}
                onClick={() => setCurrentChat(chat)}
              >
                {`${chat.firstName} ${chat.lastName}`}
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
          {currentChat
            ? <Chat chatId={currentChat.chatId} lastMessage={lastMessage} />
            : <div>Выберите активный чат 👈</div>}
        </Grid>
      </Grid>
    </Container>
  );
};

export default ChatList;
