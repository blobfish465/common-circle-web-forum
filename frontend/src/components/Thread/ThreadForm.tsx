import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { createThread } from '../../api/threadsAPI';
import { Button, TextField, MenuItem, Select, FormControl, InputLabel } from '@mui/material';
import { getCategories } from '../../api/categoriesAPI';
import { Category } from 'types/category';
import { Thread } from 'types/thread';

interface ThreadFormProps {
  userId: number;
}

const ThreadForm: React.FC<ThreadFormProps> = ({ userId }) => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [categories, setCategories] = useState<Category[]>([]);
    const [categoryId, setCategoryId] = useState<number | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const categoriesData = await getCategories();
                setCategories(categoriesData);
            } catch (error) {
                console.error('Error fetching categories:', error);
            }
        };

        fetchCategories();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!categoryId) {
            alert('Please select a category.');
            return;
        }

        try {
            const singaporeTime = new Date();
            singaporeTime.setHours(singaporeTime.getHours() + 8); // UTC +8

            const newThread: Omit<Thread, 'id'> = {
                user_id: userId,  
                title,
                content,
                created_at: singaporeTime.toISOString(), 
                category_id: categoryId,
            };
            
            await createThread(newThread);  
            alert('Thread created successfully!');
            navigate('/');
        } catch (error) {
            alert('Error creating thread');
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <TextField
                label="Title"
                fullWidth
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />
            <TextField
                label="Content"
                fullWidth
                multiline
                rows={4}
                value={content}
                onChange={(e) => setContent(e.target.value)}
            />
            <FormControl fullWidth>
                <InputLabel>Category</InputLabel>
                <Select
                    value={categoryId || ''}
                    onChange={(e) => setCategoryId(Number(e.target.value))}
                >
                    {categories.map((category) => (
                        <MenuItem key={category.id} value={category.id}>
                            {category.name}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>
            <Button type="submit" variant="contained">
                Create Thread
            </Button>
        </form>
    );
};

export default ThreadForm;
