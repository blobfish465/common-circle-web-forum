import axiosInstance from './axiosInstance';
import { Thread } from '../types/thread';



export const createThread = async (thread: Omit<Thread, 'id'>) => {
    const response = await axiosInstance.post('/threads', thread);
    return response.data;
};


export const updateThread = async (threadId: number, thread: Partial<Thread>) => {
    const response = await axiosInstance.put(`/threads/${threadId}`, thread);
    return response.data;
};


export const deleteThread = async (threadId: number) => {
    const response = await axiosInstance.delete(`/threads/${threadId}`);
    return response.data;
};

export const getAllThreads = async (): Promise<Thread[]> => {
    try {
        const response = await axiosInstance.get('/threads');
        console.log('API Response:', response.data);
        const threads = response.data.payload?.data;

        if (Array.isArray(threads)) {
            return threads;
        } else {
            console.error('Invalid data format:', threads);
            return [];
        }
    } catch (error) {
        console.error('Error fetching threads:', error);
        return [];
    }
};

export const getThreadsByUserId = async (userId: number): Promise<Thread[]> => {
    try {
    console.log(`Fetching threads for userId: ${userId}`);
      const response = await axiosInstance.get(`/users/${userId}/threads`);
      
      // Extract the threads from the response payload
      const threads = response.data.payload?.data;
      
      if (Array.isArray(threads)) {
        return threads;
      } else {
        console.error('Invalid data format:', threads);
        return [];
      }
    } catch (error) {
      console.error('Error fetching threads:', error);
      return [];
    }
};

export const getThreadsByCategory = async (categoryId: number): Promise<Thread[]> => {
    try {
        const response = await axiosInstance.get(`/categories/${categoryId}/threads`);
        return response.data.payload.data;
    } catch (error) {
        console.error('Error fetching threads by category:', error);
        return [];
    }
};
