import Header from "./components/Header/Header";
import ChatHistory from "./components/ChatHistory/ChatHistory";
import { useEffect, useState } from 'react';
import { connect, sendMsg } from './api/index';

function App() {
  const [chatHistory, setChatHistory] = useState([]);
  useEffect(() => {
    connect((msg) => {
      console.log("New Message")
      setChatHistory([...chatHistory, msg])
    });
  })

  const send = (e) => {
    sendMsg('hello hieu hoc code day')
  }

  return (
    <div className="App">
      <Header />
      <ChatHistory chatHistory={chatHistory} />
      <button onClick={send}>Hit</button>
    </div>
  );
}


export default App;
