import axiosInstance from './axiosInstance';
import { User } from '../types/user';


// Get all users
export const getUsers = async (): Promise<User[]> => {
    try {
        const response = await axiosInstance.get('/users');
        return response.data;
    } catch (error) {
        console.error('Error fetching users:', error);
        throw error;
    }
};
  
// Create a new user
export const createUser = async (user: Omit<User, 'id'>): Promise<User> => {
    try {
        const response = await axiosInstance.post('/users', user);
        return response.data;
    } catch (error) {
        console.error('Error creating user:', error);
        throw error;
    }
};
  
// Get a user by user ID
export const getUserById = async (userId: number): Promise<User> => {
    try {
        const response = await axiosInstance.get(`/users/${userId}`);
        console.log('getUserById response:', response);
        return response.data.payload.data;
    } catch (error) {
        console.error('Error fetching user by ID:', error);
        throw error;
    }
};
  
// Delete a user by user ID
export const deleteUser = async (userId: number): Promise<void> => {
    try {
        await axiosInstance.delete(`/users/${userId}`);
    } catch (error) {
        console.error('Error deleting user:', error);
        throw error;
    }
};