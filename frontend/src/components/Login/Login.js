import React from 'react'
import './login.scss'

function Login({ onSubmit, handleChange, send }) {
    return (
        <div className="form" onSubmit={onSubmit}>
            <div className="form-toggle"></div>
            <div className="form-panel one">
                <div className="form-header">
                    <h1>Account Login</h1>
                </div>
                <div className="form-content">
                    <form>
                        <div className="form-group">
                            <label htmlFor="username">ID Username</label>
                            <input
                                type="text"
                                id="username"
                                name="username"
                                required="required"
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <button type="submit" onClick={send}>Log In</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

export default Login