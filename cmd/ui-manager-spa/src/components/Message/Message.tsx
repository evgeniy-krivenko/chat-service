import { FC } from 'react';
import { Grid } from '@mui/material';
import { formatRelative } from 'date-fns';
import { ru } from 'date-fns/locale';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import cn from 'classnames';
import { IMessage } from '../../types/messages';
import './Message.css';

export interface MessageProps {
  message: IMessage;
}

const Message: FC<MessageProps> = ({ message }) => (
  <Grid
    container
    direction="column"
    className="message__container"
  >
    <div
      className={cn('message', {
        right: message.userIsAuthor,
        left: !message.userIsAuthor,
      })}
    >
      {message.authorName && <p className="message__author-name">{message.authorName}</p>}
      {message.body}
      {message.userIsAuthor
          && <DoneAllIcon className="message__icon" style={{ width: '14px', height: '14px' }} />}
      <div className="message__time">
        {formatRelative(new Date(message.createdAt), new Date(), { locale: ru })}
      </div>
    </div>
  </Grid>
);

export default Message;
