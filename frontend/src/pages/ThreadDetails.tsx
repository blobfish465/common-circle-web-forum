import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Thread } from '../types/thread';
import { Comment } from '../types/comment';
import { Category } from '../types/category';
import axiosInstance from '../api/axiosInstance';
import { Typography, Card, CardContent, Divider } from '@mui/material';
import CommentList from '../components/Comment/CommentList';
import CommentForm from '../components/Comment/CommentForm';
import { useAuth } from '../context/AuthContext';

const ThreadDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { userId } = useAuth();
  const [thread, setThread] = useState<Thread | null>(null);
  const [username, setUsername] = useState<string | null>(null);
  const [categoryName, setCategoryName] = useState<string | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchThread = async () => {
      try {
        // Fetch the thread details using thread id
        const threadRes = await axiosInstance.get(`/threads/${id}`);
        const fetchedThread = threadRes.data.payload.data;
        setThread(fetchedThread);

        // Fetch the user details based on user_id from the thread
        const userRes = await axiosInstance.get(`/users/${fetchedThread.user_id}`);
        const user = userRes.data.payload.data;
        setUsername(user.username); 

        // Fetch the category details based on category_id from the thread
        const categoryRes = await axiosInstance.get(`/categories/${fetchedThread.category_id}`);
        const category: Category = categoryRes.data.payload.data;
        setCategoryName(category.name);
        
        // Fetch the comments of the thread
        const commentsRes = await axiosInstance.get(`/threads/${id}/comments`);
        setComments(commentsRes.data.payload.data || []);

      } catch (err) {
        console.error(err);
        setError('Failed to load thread details.');
      }
    };

    if (id) fetchThread();
  }, [id]);

  // Update comment list of the thread by adding the new comment to the list
  const handleCommentAdded = (newComment: Comment) => {
    setComments((prev) => [...prev, newComment]); 
  };

  const handleCommentUpdated = (updatedComment: Comment) => {
    setComments((prev) =>
      prev.map((comment) => (comment.id === updatedComment.id ? updatedComment : comment))
    );
  };
  
  const handleCommentDeleted = (commentId: number) => {
    setComments((prev) => prev.filter((comment) => comment.id !== commentId));
  };

  // directly adjust timezone offset programmatically as 
  // using toLocaleString() does not preserve time precision when converting to a string
  // and new Date() might misinterpret the resulting string. 
  const adjustToSingaporeTime = (isoDate: string): string => {
    const date = new Date(isoDate);
    const singaporeOffset = 8 * 60; // Singapore is UTC+8
    const localOffset = date.getTimezoneOffset(); // Local timezone offset in minutes
    const offsetDifference = singaporeOffset - localOffset;
    date.setMinutes(date.getMinutes() + offsetDifference);
    return date.toLocaleString('en-SG', { hour12: false }); // 24-hour format
  };

  if (error) {
    return <Typography color="error">{error}</Typography>;
  }

  if (!thread || !username) {
    return <Typography>Loading...</Typography>;
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h4" gutterBottom>
          {thread.title}
        </Typography>

        <Typography variant="body1" gutterBottom>
          {thread.content}
        </Typography>

        <Typography variant="body2" color="textSecondary">
          Category: {categoryName}
        </Typography>

        <Typography variant="body2" color="textSecondary">
          Created by: {username}  
        </Typography>

        <Typography variant="body2" color="textSecondary">
          Created at: {adjustToSingaporeTime(thread.created_at)}
        </Typography>

        {thread.updated_at && (
          <Typography variant="body2" color="textSecondary">
            Updated at: {adjustToSingaporeTime(thread.updated_at)}
          </Typography>
        )}

        <Divider sx={{ my: 2 }} />

        <Typography variant="h5" gutterBottom style={{ marginTop: '2rem' }}>
          Comments 
        </Typography>
        
        {/* list comments if there is */}
        {comments.length === 0 ? (
          <Typography variant="body2" color="textSecondary">
            No comments
          </Typography>
        ) : (
          <CommentList 
            comments={comments}
            onCommentUpdated={handleCommentUpdated}
            onCommentDeleted={handleCommentDeleted} 
          />
        )}

        {/* Add Comment Form */}
        {userId && (
          <CommentForm
            threadId={Number(id)}
            userId={Number(userId)}
            onCommentAdded={handleCommentAdded}
          />
        )}
      </CardContent>
    </Card>
  );
};

export default ThreadDetails;
