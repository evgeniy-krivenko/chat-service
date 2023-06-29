import React, {
  ChangeEvent,
  FC,
  useEffect,
  useRef,
} from 'react';
import {
  Box, Button, Container, Grid, TextField,
} from '@mui/material';
import { shallow } from 'zustand/shallow';
import SendIcon from '@mui/icons-material/Send';
import Message from '../Message';
import './Chat.css';
import useMessages from '../../store/messages/messages';
import { Events } from '../../types/events';
import useAuth from '../../hook/useAuth';
import { MessageForm } from './types';

interface ChatProps {
 chatId: string;
 lastMessage: MessageEvent<any>;
 authorName?: string;
}

const Chat: FC<ChatProps> = ({ chatId, lastMessage, authorName }) => {
  const listRef = useRef<HTMLDivElement>(null);
  const { manager } = useAuth();
  const {
    messages, getMessages, addMessage, sendMessage, loading, error, resetMessages, cursor,
  } = useMessages((state) => ({
    loading: state.loading,
    messages: state.messages,
    getMessages: state.getMessages,
    addMessage: state.addMessage,
    sendMessage: state.sendMessage,
    error: state.error,
    resetMessages: state.resetMessages,
    cursor: state.cursor,
  }), shallow);

  useEffect(() => {
    getMessages(chatId, manager.id);
    return () => {
      resetMessages();
    };
  }, [chatId]);

  const observerTarget = useRef<HTMLDivElement | null >(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          if (cursor) { getMessages(chatId, manager.id); }
        }
      },
      { threshold: 1 },
    );

    if (observerTarget.current) {
      observer.observe(observerTarget.current as HTMLDivElement);
    }

    return () => {
      if (observerTarget.current) {
        observer.unobserve(observerTarget.current as HTMLDivElement);
      }
    };
  }, []);

  useEffect(() => {
    setTimeout(() => {
      listRef.current?.lastElementChild?.scrollIntoView();
    }, 50);
  }, [lastMessage, chatId, loading]);

  useEffect(() => {
    if (lastMessage) {
      const msg: Events = JSON.parse(lastMessage.data);

      if (msg.eventType === 'NewMessageEvent' && msg.chatId === chatId) {
        addMessage({
          id: msg.messageId,
          body: msg.body,
          authorId: msg.authorId,
          createdAt: msg.createdAt,
          userIsAuthor: msg.authorId === manager.id,
        });
      }
    }
  }, [lastMessage]);

  const handleSubmit: React.FormEventHandler = (event: ChangeEvent<MessageForm>) => {
    event.preventDefault();

    const message = event.currentTarget.message.value;

    sendMessage(chatId, message);

    if (!error) {
      event.target.message.value = '';
    }
  };

  return (
    <Container>
      <Grid container justifyContent="center" direction="column">
        <div className="chat__wrapper">
          <div className="chat" ref={listRef}>
            <div style={{ flex: '1 1 auto' }} ref={observerTarget} />
            {messages.map((m) => (
              <Message key={m.id} authorName={authorName} message={m} />
            ))}
          </div>
          <Grid
            container
            justifyContent="center"
            justifySelf="center"
            alignItems="center"
            style={{ width: '100%', marginLeft: '-9px' }}
          >
            <Grid xs={12} item>
              <Box
                onSubmit={handleSubmit}
                component="form"
                noValidate
                autoComplete="off"
              >
                <Grid xs={12} item container alignItems="center">
                  <Grid xs={10} item>
                    <TextField
                      name="message"
                      fullWidth
                      maxRows={3}
                      rows={2}
                      variant="outlined"
                      style={{ padding: '4px' }}
                    />
                  </Grid>
                  <Grid xs={2} item alignContent="center">
                    <Button type="submit" disabled={loading}>
                      <SendIcon />
                    </Button>
                  </Grid>
                </Grid>
              </Box>
            </Grid>
          </Grid>
        </div>
      </Grid>
    </Container>
  );
};

export default Chat;
