import React, {FC} from 'react';
import {IMessage} from "../../types/messages";
import {Grid} from "@mui/material";
import {formatRelative} from 'date-fns';
import {ru} from 'date-fns/locale';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import DoneIcon from '@mui/icons-material/Done';
import cn from 'classnames';
import './Message.css';

export interface MessageProps {
  message: IMessage;
}

const Message: FC<MessageProps> = ({message}) => {

  return (
      <Grid
        container
        direction="column"
        className="message__container"
      >
          <div
            className={cn('message', {
              'service': message.isService,
              'blocked': message.isBlocked,
              'right': message.userIsAuthor && !message.isBlocked,
              'left': !message.userIsAuthor,
            })}
          >
            {message.body}
            {message.userIsAuthor && (message.isReceived
              ? <DoneAllIcon className="message__icon" style={{ width: '14px', height: '14px' }} />
              : <DoneIcon className="message__icon" style={{ width: '14px', height: '14px' }} />)}
            <div className="message__time">
              {formatRelative(new Date(message.createdAt), new Date(), {locale: ru})}
            </div>
          </div>
      </Grid>
  );
};

export default Message;
