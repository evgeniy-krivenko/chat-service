import React, {FC, JSX} from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import useAuth from "../../hook/useAuth";

interface WithAuthProps {
  children: JSX,
}

const WithAuth: FC<WithAuthProps> = ({ children }) => {
  const { user } = useAuth();
  const location = useLocation();

  if (!user) {
    return <Navigate to="/login" state={location} />
  }

  return children;
}

export default WithAuth;
