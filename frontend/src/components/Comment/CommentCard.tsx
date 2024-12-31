import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Typography, Box, Button, TextField } from '@mui/material';
import { Comment } from '../../types/comment';
import { useAuth } from '../../context/AuthContext';
import axiosInstance from '../../api/axiosInstance';

interface CommentCardProps {
  comment: Comment;
  onCommentUpdated: (updatedComment: Comment) => void;
  onCommentDeleted: (commentId: number) => void;
}

const CommentCard: React.FC<CommentCardProps> = ({ comment, onCommentUpdated, onCommentDeleted }) => {
  const [username, setUsername] = useState<string | null>(null);
  const { userId } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [editedContent, setEditedContent] = useState(comment.content);

  useEffect(() => {
    const fetchUsername = async () => {
      try {
        const userRes = await axiosInstance.get(`/users/${comment.user_id}`);
        const user = userRes.data.payload.data;
        setUsername(user.username);
      } catch (error) {
        console.error('Error fetching username:', error);
        setUsername('Unknown User');
      }
    };

    fetchUsername();
  }, [comment.user_id]);

  const handleEdit = async () => {
    try {
      const updatedComment = { ...comment, content: editedContent };
      await axiosInstance.put(`/comments/${comment.id}`, updatedComment);
      onCommentUpdated(updatedComment);
      setIsEditing(false);
    } catch (err) {
      console.error('Error updating comment:', err);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this comment?')) return;
    try {
      await axiosInstance.delete(`/comments/${comment.id}`);
      onCommentDeleted(comment.id);
    } catch (err) {
      console.error('Error deleting comment:', err);
    }
  };

  return (
    <Box
      style={{
        marginBottom: '1rem',
        padding: '1rem',
        border: '1px solid #ddd',
        borderRadius: '8px',
      }}
    >
      {isEditing ? (
        <Box>
          <TextField
            fullWidth
            multiline
            rows={3}
            value={editedContent}
            onChange={(e) => setEditedContent(e.target.value)}
            style={{ marginBottom: '1rem' }}
          />
          <Button variant="contained" color="primary" onClick={handleEdit}>
            Save
          </Button>
          <Button variant="text" color="secondary" onClick={() => setIsEditing(false)}>
            Cancel
          </Button>
        </Box>
      ) : (
        <Link to={`/thread/${comment.thread_id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
          <Box
            style={{
              cursor: 'pointer',
              transition: 'box-shadow 0.3s',
            }}
            sx={{
              '&:hover': {
                boxShadow: '0px 4px 10px rgba(0, 0, 0, 0.2)',
              },
            }}
          >
            <Typography variant="body1">{comment.content}</Typography>
            <Typography variant="body2" color="textSecondary">
              By {username || 'Unknown User'} at {new Date(comment.created_at).toLocaleString()}
            </Typography>
          </Box>
        </Link>
      )}

      {userId === comment.user_id.toString() && !isEditing && (
        <Box style={{ marginTop: '0.5rem' }}>
          <Button variant="text" onClick={() => setIsEditing(true)}>
            Edit
          </Button>
          <Button variant="text" color="error" onClick={handleDelete}>
            Delete
          </Button>
        </Box>
      )}
    </Box>
  );
};

export default CommentCard;

