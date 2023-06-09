import { FC } from 'react';
import useAuth from "../../hook/useAuth";
import {useNavigate} from "react-router-dom";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface ChatProps {

}

const Chat: FC<ChatProps> = (props) => {
  const {signOut} = useAuth();
  const navigate = useNavigate();

  const logout = () => {
    signOut(() => navigate('/login'));
  }
  return (
    <div>
      Chat
     <button onClick={logout}>Logout</button>
    </div>
  );
};

export default Chat;
