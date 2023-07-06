import { FC } from 'react';
import { green } from '@mui/material/colors';
import { CircularProgress, Container, Grid } from '@mui/material';

const Loader: FC = () => (
  <Container>
    <Grid
      justifyContent="center"
      alignItems="center"
    >
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
    </Grid>
  </Container>
);

export default Loader;
