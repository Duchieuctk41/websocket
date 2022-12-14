let socket = new WebSocket("ws://127.0.0.1:8080/ws");

let connect = cb => {
  // console.log("connecting");

  socket.onopen = () => {
    // console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    console.log(msg);
    cb(msg);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {
  if (socket.readyState === WebSocket.OPEN) {
    // console.log("2")
    socket.send(msg);

  } else {
    setTimeout(() => {
      sendMsg()

    }, 1000);
  }
};

export { connect, sendMsg };