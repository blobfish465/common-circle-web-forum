import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { TextField, Button, Box, Typography } from '@mui/material';

const Login: React.FC = () => {
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
    try {
        await login(username, password);
        console.log("Logged in successfully");
        setError(null);
        navigate('/');
    } catch (error) {
        setError('Invalid username or password');
        if (error instanceof Error) {
            alert(`Login failed: ${error.message}`);
        } else {
            console.error("An unexpected error occurred:", error);
            alert("Login failed due to an unexpected error.");
        }
    }
    };

    return (
        <Box
            sx={{
                maxWidth: 400,
                margin: '0 auto',
                marginTop: 4, // Add margin at the top
                padding: 2,
                border: '1px solid #ccc',
                borderRadius: 4,
            }}
        >
        <Typography variant="h4" gutterBottom>
            Login
        </Typography>
        {error && (
            <Typography color="error" gutterBottom>
            {error}
            </Typography>
        )}
            <form onSubmit={handleSubmit}>
                <TextField
                    fullWidth
                    label="Username"
                    margin="normal"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                />
                <TextField
                    fullWidth
                    label="Password"
                    margin="normal"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <Button type="submit" variant="contained" fullWidth sx={{ marginTop: 2 }}>
                    Login
                </Button>
            </form>
        </Box>

    );
};

export default Login;