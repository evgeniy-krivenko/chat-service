import { FC, ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import useAuth from '../../hook/useAuth';

interface WithAuthProps {
  children: ReactNode,
}

const WithAuth: FC<WithAuthProps> = ({ children }) => {
  const { manager } = useAuth();
  const location = useLocation();

  if (!manager) {
    return <Navigate to="/login" state={location} />;
  }

  return children;
};

export default WithAuth;
