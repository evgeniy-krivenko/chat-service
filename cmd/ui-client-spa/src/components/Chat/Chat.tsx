import {FC, useEffect, useRef} from 'react';
import {Box, Button, Container, Grid, TextField} from "@mui/material";
import SendIcon from '@mui/icons-material/Send';
import Message from "../Message/Message";
import {useMessages} from "../../hook/useMessages";
import {APIClient} from "../../api";
import './Chat.css';

const Chat: FC = () => {
  const listRef = useRef<HTMLDivElement>(null);
  const {messages, onScroll, addMessage} = useMessages(listRef);

  const observerTarget = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting) {
          onScroll();
        }
      },
      {threshold: 1}
    );

    if (observerTarget.current) {
      observer.observe(observerTarget.current);
    }

    return () => {
      if (observerTarget.current) {
        observer.unobserve(observerTarget.current);
      }
    };
  }, []);

  useEffect(() => {
    setTimeout(() => {

      listRef.current?.lastElementChild?.scrollIntoView();
    }, 200)
  }, [])

  const handleSubmit = (event: unknown) => {
    event.preventDefault();

    const message = event.currentTarget.message.value;

    APIClient.sendMessage(message)
      .then((result) => {
        addMessage({
          ...result,
          body: message,
          isBlocked: false,
          isReceived: false,
          isService: false,
        });

        setTimeout(() => {
          listRef.current?.lastElementChild?.scrollIntoView();
        }, 100)

      }).finally(() => event.target.message.value = '')
  }

  return (
    <Container>
      <Grid
        container
        justifyContent="center"
        direction="column"
      >
        <div className="chat__wrapper">
          <div className="chat"
            ref={listRef}
          >
            <div
              style={{ flex: '1 1 auto' }}
              ref={observerTarget}
            />
            {messages.map((m) => (
              <Message key={m.id} message={m}/>
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
                <Grid
                  xs={12}
                  item
                  container
                  alignItems="center"
                >
                  <Grid xs={10} item>
                    <TextField
                      name="message"
                      fullWidth
                      maxRows={3}
                      rows={2}
                      variant="outlined"
                      style={{padding: '4px'}}
                    />
                  </Grid>
                  <Grid
                    xs={2}
                    item
                    alignContent="center"
                  >
                    <Button type="submit">
                      <SendIcon/>
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
