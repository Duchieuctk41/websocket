import Header from "./components/Header/Header";
import ChatHistory from "./components/ChatHistory/ChatHistory";
import { useEffect, useState } from 'react';
import { connect, sendMsg } from './api/index';

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
      <form onSubmit={onSubmit}>
        <label>
          Name:
          <input type="text" onChange={handleChange} />
        </label>
        <button onClick={send}>Hit</button>
      </form>
      <ChatHistory chatHistory={chatHistory} />
    </div>
  );
}


export default App;
