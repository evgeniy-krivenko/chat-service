import './App.css';
import {Route, Routes} from "react-router-dom";
import Chat from "./components/Chat";
import Login from "./components/Login";
import NotFound from "./components/NotFound";
import WithAuth from "./hoc/WithAuth";
import AuthProvider from "./hoc/AuthProvider";

function App() {

  return (
    <AuthProvider>
      <Routes>
        <Route path="/" element={
          <WithAuth>
            <Chat />
          </WithAuth>
        } />
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </AuthProvider>
  )
}

export default App
