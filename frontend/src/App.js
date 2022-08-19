import Header from "./components/Header/Header";
import ChatHistory from "./components/ChatHistory/ChatHistory";
import { useEffect, useState } from 'react';
import { connect, sendMsg } from './api/index';
import Login from "./components/Login/Login";

function App() {
  const [chatHistory, setChatHistory] = useState([]);
  const [msg, setMsg] = useState('')
  useEffect(() => {
    connect((msg) => {
      console.log("New Message")
      setChatHistory([...chatHistory, msg])
    });
  })

  const send = (e) => {
    sendMsg(msg)
    console.log(msg)
  }

  const handleChange = (e) => {
    setMsg(e.target.value);
  }

  const onSubmit = (e) => {
    e.preventDefault();
  }

  return (
    <div className="App">
      <Header />
      <Login handleChange={handleChange} send={send} onSubmit={onSubmit} />
      <ChatHistory chatHistory={chatHistory} />
    </div>
  );
}


export default App;
