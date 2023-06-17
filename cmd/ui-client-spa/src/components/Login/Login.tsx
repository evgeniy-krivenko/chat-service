import React, {FC, FormEvent, useState} from 'react';
import {useNavigate, useLocation, Navigate} from "react-router-dom";
import useAuth from "../../hook/useAuth";
import {green} from '@mui/material/colors';
import {Alert, Box, Button, CircularProgress, Container, Grid, TextField} from "@mui/material";
import {APIClient} from "../../api";

export interface LoginForm {
  login: { value: string };
  password: { value: string };
}

const Login: FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const {user, signIn} = useAuth();
  const [err, setErr] = useState<string>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const buttonSx = {
    ...(user && {
      bgcolor: green[500],
      '&:hover': {
        bgcolor: green[700],
      },
    }),
    marginTop: '8px',
    width: '160px',
    height: '48px',
  };

  const fromPage = location.state?.from?.pathname || '/';

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    const form = event.target as LoginForm;

    const login = form.login?.value;
    const password = form.password?.value;

    setIsLoading(true);

    APIClient.login({login, password})
      .then((res) => {
        signIn(res.user, () => navigate(fromPage, {replace: true}));

        localStorage.setItem('token', res.token);
      })
      .catch((e) => setErr(e.message))
      .finally(() => setIsLoading(false));
  }

  if (user) {
    return <Navigate to={fromPage} state={location} />
  }

  return (
    <div>
      <Container>
        <Grid
          container
          style={{marginTop: '15vh'}}
          alignItems="center"
          justifyContent="center"
        >
          <Grid
            container
            alignItems="center"
            justifyContent="center"
            direction="column"
          >
            <h1>Login</h1>
            <Box
              p={4}
              onSubmit={handleSubmit}
              flexDirection="column"
              component="form"
              sx={{
                '& .MuiTextField-root': {m: 2, width: '250px'},
                display: 'block',
              }}
              noValidate
              autoComplete="off"
            >
              <Grid xs={12}>
                <TextField
                  onChange={() => setErr(null)}
                  name="login"
                  required
                  id="login"
                  label="Login"
                  disabled={isLoading}
                />
              </Grid>
              <Grid xs={12}>
                <TextField
                  onChange={() => setErr(null)}
                  name="password"
                  required
                  type="password"
                  id="password"
                  label="Password"
                  disabled={isLoading}
                />
              </Grid>
              <Button
                sx={buttonSx}
                type="submit"
                variant="contained"
                disabled={isLoading}
              >
                Login
                {isLoading && (
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
              {err && <Alert style={{margin: '24px 0'}} severity="error">{err}</Alert>}
            </Box>
          </Grid>
        </Grid>
      </Container>
    </div>
  );
};

export default Login;
