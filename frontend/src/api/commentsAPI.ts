import axiosInstance from './axiosInstance';
import { Comment } from '../types/comment';


export const createComment = async (comment: Omit<Comment, 'id'>) => {
    try {
        const response = await axiosInstance.post('/comments', comment);
        return response.data.payload.data; 
    } catch (error) {
        console.error('Error creating comment:', error);
        throw error;
    }
};


export const updateComment = async (commentId: number, comment: Partial<Comment>) => {
    const response = await axiosInstance.put(`/comments/${commentId}`, comment);
    return response.data.payload.data;
};


export const deleteComment = async (commentId: number) => {
    await axiosInstance.delete(`/comments/${commentId}`);
};


export const getCommentsByThreadId = async (threadId: number): Promise<Comment[]> => {
    try {
        const response = await axiosInstance.get(`/threads/${threadId}/comments`);
        return response.data.payload.data || [];
    } catch (error) {
        console.error('Error fetching comments:', error);
        return [];
    }
};

export const getCommentsByUserId = async (userId: number): Promise<Comment[]> => {
    try {
      const response = await axiosInstance.get(`/users/${userId}/comments`);
      return response.data.payload.data || [];
    } catch (error) {
      console.error('Error fetching comments by user ID:', error);
      throw error;
    }
};