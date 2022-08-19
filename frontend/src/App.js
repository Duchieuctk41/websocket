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
      console.log("New Message", msg)
      setChatHistory([...chatHistory, msg])
    });
  })

  const send = (e) => {
    sendMsg(msg)
    // console.log(msg)
    // is_online -> data: {"event_name":"is_online","event_payload":{"body":"im here","from":"ca248164-c87b-4f29-8f6a-fe78cedb9718"}}
    // chat -> data: {"event_name":"chat","event_payload":{"body":"im here","from":"ca248164-c87b-4f29-8f6a-fe78cedb9718"}}
    // is_typing -> -> data: {"event_name":"is_typing","event_payload":{"body":"im here","from":"ca248164-c87b-4f29-8f6a-fe78cedb9718"}}
  }

  const handleChange = (e) => {
    setMsg(JSON.stringify({"topic": "is_online", "from": e.target.value}));
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
