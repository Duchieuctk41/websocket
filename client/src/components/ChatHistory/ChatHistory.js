import React from "react";
import "./ChatHistory.scss";

function ChatHistory(props) {
  const { chatHistory } = props
  const messages = chatHistory.map((msg, index) => (
    <p key={index}>{msg.data}</p>
  ));

  return (
    <div className="ChatHistory">
      <h2>Chat History</h2>
      <div>{messages}</div>
    </div>
  );
}

export default ChatHistory;