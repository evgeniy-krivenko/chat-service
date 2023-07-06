import { useContext } from 'react';
import { AuthContext, IAuthProvider } from '../hoc/AuthProvider/AuthProvider';

export default function useAuth(): IAuthProvider {
  return useContext(AuthContext);
}
