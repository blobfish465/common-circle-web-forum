import axios from 'axios';
import Cookies from 'js-cookie';

const axiosInstance = axios.create({
    baseURL: process.env.REACT_APP_API_URL || 'http://54.206.114.24:8080',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add an interceptor to include the token in the Authorization header
axiosInstance.interceptors.request.use(
    (config) => {
      const token = Cookies.get('authToken'); // Retrieve token from cookies
      console.log('Auth token:', Cookies.get('authToken'));
      console.log('Token attached to request:', token);
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
  );

// Handle token expiry or unauthorized access
axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
      if (error.response?.status === 401) {
          console.log('Unauthorized or token expired. Logging out...');
          Cookies.remove('authToken'); // Remove expired token
          window.location.href = '/login'; // Redirect to login page
      }
      return Promise.reject(error);
  }
);

export default axiosInstance;
