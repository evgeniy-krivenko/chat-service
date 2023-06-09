import {FC, FormEvent} from 'react';
import {useNavigate, useLocation} from "react-router-dom";
import useAuth from "../../hook/useAuth";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface LoginProps {

}

export interface LoginForm {
  username: { value: string };
}

const Login: FC<LoginProps> = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { signIn } = useAuth();

  const fromPage = location.state?.from?.pathname || '/';

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    const form = event.target as LoginForm;

    const username = form.username.value;

    signIn({ username }, () => navigate(fromPage, { replace: true } ));
  }

  return (
    <div>
      <h1>Login</h1>
      <form onSubmit={handleSubmit}>
        <label >
          Name:
          <input name="username" />
        </label>
        <button type="submit">Login</button>
      </form>
    </div>
  );
};

export default Login;
