import { Route, Routes } from 'react-router-dom';
import { ToastContainer } from 'react-toastify';
import { enableMapSet } from 'immer';
import Login from './components/Login';
import NotFound from './components/NotFound';
import WithAuth from './hoc/WithAuth';
import ChatList from './components/ChatList';
import './App.css';
import AuthProvider from './hoc/AuthProvider';
import 'react-toastify/dist/ReactToastify.css';

enableMapSet();

const App = () => (
  <AuthProvider>
    <Routes>
      <Route
        path="/"
        element={(
          <WithAuth>
            <ChatList />
          </WithAuth>
        )}
      />
      <Route path="/login" element={<Login />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
    <ToastContainer icon={false} />
  </AuthProvider>
);

export default App;
