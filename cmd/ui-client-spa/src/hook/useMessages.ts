import {IMessage} from "../types/messages";
import {MutableRefObject, useEffect, useRef, useState} from "react";
import {APIClient} from "../api";
import useWebSocket from "react-use-websocket";
import useAuth from "./useAuth";
import {BLOCK_MESSAGE_TEXT} from "../const/blockMessage";

type TUseMessage = { messages: IMessage[], onScroll: () => void, addMessage: (message: IMessage) => void }

const set = new Set<string>();

enum MessageEvent {
  NewMessageEvent = 'NewMessageEvent',
  MessageSentEvent = 'MessageSentEvent',
  MessageBlockedEvent = 'MessageBlockedEvent',
}

interface NewMessageEvent extends Omit<IMessage, 'id'> {
  eventType: MessageEvent.NewMessageEvent;
  messageId: string;
}

interface MessageSentEvent {
  eventType: MessageEvent.MessageSentEvent;
  messageId: string;
}

interface MessageBlockEvent {
  eventType: MessageEvent.MessageBlockedEvent;
  messageId: string;
}

type Events = NewMessageEvent | MessageSentEvent | MessageBlockEvent

export const useMessages = (ref?: MutableRefObject<HTMLDivElement>): TUseMessage => {
  const {user} = useAuth();

  const [cursor, setCursor] = useState('');
  const innerCursor = useRef<string>();
  const [messages, setMessages] = useState<IMessage[]>([]);

  const token = localStorage.getItem('token');
  const {lastMessage} = useWebSocket(
    import.meta.env.VITE_WS_ENDPOINT,
    {protocols: [import.meta.env.VITE_WS_PROTOCOL, token]},
  )

  useEffect(() => {

    if (lastMessage) {
      setTimeout(() => {
        ref.current?.lastElementChild?.scrollIntoView();
      }, 200)
    }
  }, [lastMessage])

  const addNewMessages = (newMessages: IMessage[]): void => {
    setMessages((messages) => {
        const result: IMessage[] = [];

        for (const msg of newMessages) {
          if (!set.has(msg.id)) {
            set.add(msg.id);

            result.push({
              ...msg,
              userIsAuthor: msg.authorId === user?.id,
              body: msg.isBlocked ? BLOCK_MESSAGE_TEXT : msg.body,
            });
          }
        }

        return [...result, ...messages];
      }
    )
  }

  const addOneMessage = (message: IMessage): void => {
    if (!set.has(message.id)) {
      set.add(message.id);

      const newMessage = {...message, userIsAuthor: message.authorId === user?.id}

      setMessages((messages) => [...messages, newMessage])
    }
  }

  const setSentMessage = (messageId: string): void => {
    setMessages((messages) => {
      return messages.map((m) => {
        if (m.id === messageId) {
          return {
            ...m,
            isReceived: true,
          }
        }

        return m;
      })
    })
  }

  const setBlockMessage = (messageId: string): void => {
    setMessages((messages) => {
      return messages.map((m) => {
        if (m.id === messageId) {
          return {
            ...m,
            isBlocked: true,
            body: BLOCK_MESSAGE_TEXT,
          }
        }

        return m;
      })
    })
  }

  useEffect(() => {
    if (lastMessage) {
      const msg: Events = JSON.parse(lastMessage.data);

      switch (msg.eventType) {
        case "NewMessageEvent":
          addOneMessage({
            id: msg.messageId,
            isService: msg.isService,
            body: msg.body,
            authorId: msg.authorId,
            isBlocked: false,
            isReceived: false,
            userIsAuthor: msg.authorId === user?.id,
            createdAt: msg.createdAt,
          } as IMessage)
          return;
        case "MessageSentEvent":
          setSentMessage(msg.messageId);
          return;

        case "MessageBlockedEvent":
          setBlockMessage(msg.messageId);
          return;
      }
    }
  }, [lastMessage])

  useEffect(() => {
    APIClient.getHistory(cursor)
      .then((result) => {
        result.messages.reverse();
        addNewMessages(result.messages);

        if (result.next) {
          innerCursor.current = result.next;
        }
      })
  }, [cursor])

  return {
    messages,
    onScroll: () => {
      if (innerCursor.current) {
        setCursor(innerCursor.current as string)
      }
    },
    addMessage: addOneMessage,
  }
}
