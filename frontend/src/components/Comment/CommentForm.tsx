import React, { useState } from 'react';
import { TextField, Button, Box } from '@mui/material';
import { createComment } from '../../api/commentsAPI';
import { Comment } from '../../types/comment';

interface CommentFormProps {
  threadId: number;
  userId: number;
  onCommentAdded: (comment: Comment) => void;
}

const CommentForm: React.FC<CommentFormProps> = ({ threadId, userId, onCommentAdded }) => {
  const [content, setContent] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!content.trim()) {
      alert('Comment cannot be empty.');
      return;
    }

    try {
      const newComment = {
        thread_id: threadId,
        user_id: userId,
        content,
        created_at: new Date().toISOString(),
      };

      const createdComment = await createComment(newComment);
      onCommentAdded(createdComment);
      setContent(''); // Clear the input field
    } catch (err) {
      console.error('Error creating comment:', err);
      alert('Failed to create comment.');
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit} style={{ marginTop: '1rem' }}>
      <TextField
        label="Write a comment"
        fullWidth
        multiline
        rows={3}
        value={content}
        onChange={(e) => setContent(e.target.value)}
        style={{ marginBottom: '1rem' }}
      />
      <Button type="submit" variant="contained" color="primary">
        Add Comment
      </Button>
    </Box>
  );
};

export default CommentForm;
