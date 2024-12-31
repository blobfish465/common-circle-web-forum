import React, { useState } from 'react';
import { Card, CardContent, Typography, Button, Box, TextField } from '@mui/material';
import { Link } from 'react-router-dom';
import { Thread } from '../../types/thread';
import { useAuth } from '../../context/AuthContext';
import axiosInstance from '../../api/axiosInstance';

interface ThreadCardProps {
  thread: Thread;
  onThreadUpdated: (updatedThread: Thread) => void;
  onThreadDeleted: (threadId: number) => void;
}

const ThreadCard: React.FC<ThreadCardProps> = ({ thread, onThreadUpdated, onThreadDeleted }) => {
  const { userId } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [editedTitle, setEditedTitle] = useState(thread.title);
  const [editedContent, setEditedContent] = useState(thread.content);

  const handleEdit = async () => {
    try {
      const updatedThread = {
        ...thread,
        title: editedTitle,
        content: editedContent,
      };
      await axiosInstance.put(`/threads/${thread.id}`, updatedThread);
      onThreadUpdated(updatedThread);
      setIsEditing(false);
    } catch (err) {
      console.error('Error updating thread:', err);
    }
  };

  const handleDelete = async () => {
    try {
      await axiosInstance.delete(`/threads/${thread.id}`);
      onThreadDeleted(thread.id);
    } catch (err) {
      console.error('Error deleting thread:', err);
    }
  };

  return (
    <Card sx={{ marginBottom: 2 }}>
      <CardContent>
        {isEditing ? (
          <>
            <TextField
              fullWidth
              label="Edit Title"
              value={editedTitle}
              onChange={(e) => setEditedTitle(e.target.value)}
              sx={{ marginBottom: 1 }}
            />
            <TextField
              fullWidth
              label="Edit Content"
              multiline
              rows={4}
              value={editedContent}
              onChange={(e) => setEditedContent(e.target.value)}
              sx={{ marginBottom: 1 }}
            />
            <Box>
              <Button variant="contained" color="primary" onClick={handleEdit} sx={{ marginRight: 1 }}>
                Save
              </Button>
              <Button variant="outlined" color="secondary" onClick={() => setIsEditing(false)}>
                Cancel
              </Button>
            </Box>
          </>
        ) : (
          <>
            <Typography variant="h6">
              <Link to={`/thread/${thread.id}`}>{thread.title}</Link>
            </Typography>
            <Typography variant="body2" color="textSecondary">
              {thread.content.substring(0, 100)}...
            </Typography>
            {Number(userId) === thread.user_id && ( 
              <Box sx={{ marginTop: 1 }}>
                <Button variant="contained" color="primary" onClick={() => setIsEditing(true)} sx={{ marginRight: 1 }}>
                  Edit
                </Button>
                <Button variant="outlined" color="secondary" onClick={handleDelete}>
                  Delete
                </Button>
              </Box>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
};

export default ThreadCard;
