import React from 'react';
import { Comment } from '../../types/comment';
import CommentCard from './CommentCard';
import { List } from '@mui/material';

interface CommentListProps {
  comments: Comment[];
  onCommentUpdated: (updatedComment: Comment) => void;
  onCommentDeleted: (commentId: number) => void;
}

const CommentList: React.FC<CommentListProps> = ({ comments, onCommentUpdated, onCommentDeleted }) => {
  return (
    <List>
      {comments.map((comment) => (
        <CommentCard 
            key={comment.id} 
            comment={comment} 
            onCommentUpdated={onCommentUpdated}
            onCommentDeleted={onCommentDeleted}
        />
      ))}
    </List>
  );
};

export default CommentList;
