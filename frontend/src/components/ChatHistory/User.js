import React from 'react'

function User() {
    return (
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
    )
}

export default User