import axiosInstance from './axiosInstance';
import { Category } from '../types/category';

// get all the categories, for category dropdown in ThreadForm
export const getCategories = async (): Promise<Category[]> => {
    const response = await axiosInstance.get('/categories');
    return response.data.payload.data;
};
