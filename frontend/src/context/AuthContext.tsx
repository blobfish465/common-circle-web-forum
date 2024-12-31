import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import { jwtDecode } from 'jwt-decode';
import Cookies from 'js-cookie';

interface DecodedToken {
    user_id: string;
    exp: number;
}

interface AuthContextType {
    userId: string | null;
    token: string | null;
    loading: boolean;
    login: (username: string, password: string) => Promise<void>;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [auth, setAuth] = useState<{ userId: string | null, token: string | null }>({ userId: null, token: null });
    const [loading, setLoading] = useState(true);

    // Check if the auth token is expired
    const isTokenExpired = (token: string): boolean => {
        const decoded: DecodedToken = jwtDecode(token);
        const currentTime = Math.floor(Date.now() / 1000);
        return decoded.exp < currentTime; 
    };

    useEffect(() => {
        // Get token from cookies
        const token = Cookies.get('authToken');
        if (token) {
            if (isTokenExpired(token)) {
                console.log('Token expired. Logging out...');
                Cookies.remove('authToken');
                setAuth({ userId: null, token: null });
            } else {
                const decoded: DecodedToken = jwtDecode(token);
                setAuth({ userId: decoded.user_id, token });
            }
        }
        // When Authentication status is determined
        setLoading(false); 
    }, []);

    const login = async (username: string, password: string) => {
        const response = await fetch('https://common-circle-web-forum.onrender.com/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        });
        const data = await response.json();
        if (response.ok) {
            Cookies.set('authToken', data.token, { expires: 7 }); // Token stored for 7 days
            const decoded: DecodedToken = jwtDecode(data.token);
            console.log('Decoded token:', decoded); 
            setAuth({ userId: decoded.user_id, token: data.token });
        } else {
            throw new Error(data.message || 'Login failed');
        }
    };

    const logout = () => {
        Cookies.remove('authToken');
        setAuth({ userId: null, token: null });
    };

    return (
        <AuthContext.Provider value={{ ...auth, loading, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};
