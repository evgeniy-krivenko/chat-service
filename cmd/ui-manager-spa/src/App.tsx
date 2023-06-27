import { Route, Routes } from 'react-router-dom';
import Login from './components/Login';
import NotFound from './components/NotFound';
import WithAuth from './hoc/WithAuth';
import ChatList from './components/ChatList';
import './App.css';
import AuthProvider from './hoc/AuthProvider';

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
  </AuthProvider>

);

export default App;
