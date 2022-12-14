
import { React, ReactDOM, useEffect } from 'react'
// Sending message to specific user with GoLang WebSocket
// @author Shashank Tiwari

const domElement = document.querySelector(".chat__app-container");

function App() {
  const [chatUserList, setChatUserList] = []
  const [message, setMessage] = null
  const [selectedUserID, setSelectedUserID] = null
  const [userID, setUserID] = null
  this.webSocketConnection = null;

  useEffect(() => {
    setWebSocketConnection();
    subscribeToSocketMessage();
  })

  const setWebSocketConnection = () => {
    const username = prompt("What's Your name");
    if (window["WebSocket"]) {
      const socketConnection = new WebSocket("ws://" + document.location.host + "/ws/" + username);
      this.webSocketConnection = socketConnection;
    }
  }

  const subscribeToSocketMessage = () => {
    if (this.webSocketConnection === null) {
      return;
    }

    this.webSocketConnection.onclose = (evt) => {
      this.setState({
        message: 'Your Connection is closed.',
        chatUserList: []
      });
    };

    this.webSocketConnection.onmessage = (event) => {
      try {
        const socketPayload = JSON.parse(event.data);
        switch (socketPayload.eventName) {
          case 'join':
          case 'disconnect':
            if (!socketPayload.eventPayload) {
              return
            }

            const userInitPayload = socketPayload.eventPayload;

            this.setState({
              chatUserList: userInitPayload.users,
              userID: this.state.userID === null ? userInitPayload.userID : this.state.userID
            });

            break;

          case 'message response':

            if (!socketPayload.eventPayload) {
              return
            }

            const messageContent = socketPayload.eventPayload;
            const sentBy = messageContent.username ? messageContent.username : 'An unnamed fellow'
            const actualMessage = messageContent.message;

            this.setState({
              message: `${sentBy} says: ${actualMessage}`
            });

            break;

          default:
            break;
        }
      } catch (error) {
        console.log(error)
        console.warn('Something went wrong while decoding the Message Payload')
      }
    };
  }

  const setNewUserToChat = (event) => {
    if (event.target && event.target.value) {
      if (event.target.value === "select-user") {
        alert("Select a user to chat");
        return;
      }
      this.setState({
        selectedUserID: event.target.value
      })
    }
  }

  const getChatList = () => {
    if (this.state.chatUserList.length === 0) {
      return (
        <h3>No one has joined yet</h3>
      )
    }
    return (
      <div className="chat__list-container">
        <p>Select a user to chat</p>
        <select onChange={this.setNewUserToChat}>
          <option value={'select-user'} className="username-list">Select User</option>
          {
            this.state.chatUserList.map(user => {
              if (user.userID !== this.state.userID) {
                return (
                  <option value={user.userID} className="username-list">
                    {user.username}
                  </option>
                )
              }
            })
          }
        </select>
      </div>
    );
  }

  const handleKeyPress = (event) => {
    try {
      if (event.key === 'Enter') {
        if (!this.webSocketConnection) {
          return false;
        }
        if (!event.target.value) {
          return false;
        }

        this.webSocketConnection.send(JSON.stringify({
          EventName: 'message',
          EventPayload: {
            userID: this.state.selectedUserID,
            message: event.target.value
          },
        }));

        event.target.value = '';
      }
    } catch (error) {
      console.log(error)
      console.warn('Something went wrong while decoding the Message Payload')
    }
  }

  return (
    <div>
      <div class="chat__message-container">
        <div class="message-container">
          {this.state.message}
        </div>
        <input type="text" id="message-text" size="64" autofocus placeholder="Type Your message" onKeyPress={this.handleKeyPress} />
      </div>
      {this.getChatList()}
      {this.getChatContainer()}
    </div>
  );
}

export default App;