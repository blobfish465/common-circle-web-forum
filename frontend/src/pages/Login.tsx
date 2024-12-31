import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

const Login: React.FC = () => {
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
    try {
        await login(username, password);
        console.log("Logged in successfully");
        navigate('/');
    } catch (error) {
        if (error instanceof Error) {
            alert(`Login failed: ${error.message}`);
        } else {
            console.error("An unexpected error occurred:", error);
            alert("Login failed due to an unexpected error.");
        }
    }
    };

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label>
                    Username:
                    <input type="text" value={username} onChange={e => setUsername(e.target.value)} />
                </label>
            </div>
            <div>
                <label>
                    Password:
                    <input type="password" value={password} onChange={e => setPassword(e.target.value)} />
                </label>
            </div>
            <div>
                <button type="submit">Login</button>
            </div>
        </form>
    );
};

export default Login;