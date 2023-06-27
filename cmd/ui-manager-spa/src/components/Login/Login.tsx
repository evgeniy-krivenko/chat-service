import React, { ChangeEvent, FC } from 'react';
import { useLocation, Navigate } from 'react-router-dom';
import { green } from '@mui/material/colors';
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Container,
  Grid,
  TextField,
} from '@mui/material';

import { shallow } from 'zustand/shallow';
import { useManagersStore } from '../../store/manager';
import Loader from '../Loader';
import { ILoginForm } from './types';

export interface LoginForm {
  login: { value: string };
  password: { value: string };
}

const Login: FC = () => {
  const location = useLocation();

  const {
    manager, error, login: signIn, loading, resetError,
  } = useManagersStore((state) => ({
    manager: state.manager,
    error: state.error,
    login: state.login,
    loading: state.loading,
    resetError: state.resetError,
  }), shallow);

  const buttonSx = {
    ...(manager && {
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

  const handleSubmit: React.FormEventHandler<ILoginForm> = (event: ChangeEvent<ILoginForm>) => {
    event.preventDefault();
    // const form = event.target as LoginForm;

    // const login = form.login?.value;
    // const password = form.password?.value;
    const login = event.currentTarget.login.value;
    const password = event.currentTarget.login.value;

    signIn(login, password);
  };

  if (loading) {
    return <Loader />;
  }

  if (manager) {
    return <Navigate to={fromPage} state={location} />;
  }

  return (
    <div>
      <Container>
        <Grid
          container
          style={{ marginTop: '15vh' }}
          alignItems="center"
          justifyContent="center"
        >
          <Grid
            container
            alignItems="center"
            justifyContent="center"
            direction="column"
          >
            <h1>Bank Manager Login</h1>
            <Box
              p={4}
              onSubmit={handleSubmit}
              flexDirection="column"
              component="form"
              sx={{
                '& .MuiTextField-root': { m: 2, width: '250px' },
                display: 'block',
              }}
              noValidate
              autoComplete="off"
            >
              <Grid xs={12}>
                <TextField
                  onChange={() => resetError()}
                  name="login"
                  required
                  id="login"
                  label="Login"
                  disabled={loading}
                />
              </Grid>
              <Grid xs={12}>
                <TextField
                  onChange={() => resetError()}
                  name="password"
                  required
                  type="password"
                  id="password"
                  label="Password"
                  disabled={loading}
                />
              </Grid>
              <Button
                sx={buttonSx}
                type="submit"
                variant="contained"
                disabled={loading}
              >
                Login
                {loading && (
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
              {error && (
                <Alert style={{ margin: '24px 0' }} severity="error">
                  {error}
                </Alert>
              )}
            </Box>
          </Grid>
        </Grid>
      </Container>
    </div>
  );
};

export default Login;
