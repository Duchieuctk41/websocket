import React from 'react'

function Message() {
    return (
        <>
            <div className="message">
                <div className="photo" style={{ backgroundImage: `url(https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1050&q=80)` }}>
                    <div className="online"></div>
                </div>
                <p className="text"> Hi, how are you ? </p>
            </div>
            <div className="message text-only">
                <p className="text"> What are you doing tonight ? Want to go take a drink ?</p>
            </div>
            <p className="time">14h58</p>

            <div className="discussion message-active">
                <div className="photo" style={{ backgroundImage: `url(https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1050&q=80)` }}>
                    <div className="online"></div>
                </div>
                <div className="desc-contact">
                    <p className="name">Megan Leib</p>
                    <p className="message">9 pm at the bar if possible 😳</p>
                </div>
                <div className="timer">12 sec</div>
            </div>

            <div className="message text-only">
                <div className="response">
                    <p className="text"> Hey Megan ! It's been a while 😃</p>
                </div>
            </div>
            <div className="message text-only">
                <div className="response">
                    <p className="text"> When can we meet ?</p>
                </div>
            </div>
            <p className="response-time time"> 15h04</p>
        </>
    )
}

export default Message