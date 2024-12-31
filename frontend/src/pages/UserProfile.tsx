import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { User } from '../types/user';
import { Thread } from '../types/thread';
import { Comment } from '../types/comment';
import { getUserById } from '../api/usersAPI';
import { getThreadsByUserId } from '../api/threadsAPI';
import { getCommentsByUserId } from '../api/commentsAPI'; // Optional if needed
import { Typography, Card, CardContent, Divider } from '@mui/material';
import ThreadCard from '../components/Thread/ThreadCard';
import CommentList from '../components/Comment/CommentList';

const UserProfile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<User | null>(null);
  const [threads, setThreads] = useState<Thread[]>([]);
  const [comments, setComments] = useState<Comment[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
         // Fetch user details
        const userData = await getUserById(Number(id));
        setUser(userData);
        
        // Fetch threads created by the user
        const userThreads = await getThreadsByUserId(Number(id));
        setThreads(userThreads);

        // Fetch all comments made by the user
        const userComments = await getCommentsByUserId(Number(id));
        setComments(userComments);
      } catch (error) {
        console.error('Error fetching user data:', error);
        setError('Failed to fetch user data');
      }
    };

    if (id) fetchUserData();
  }, [id]);

  // Handlers for updating and deleting comments
  const handleCommentUpdated = (updatedComment: Comment) => {
    setComments((prev) =>
      prev.map((comment) => (comment.id === updatedComment.id ? updatedComment : comment))
    );
  };

  const handleCommentDeleted = (commentId: number) => {
    setComments((prev) => prev.filter((comment) => comment.id !== commentId));
  };

  const handleThreadDeleted = (threadId: number) => {
    // Remove the deleted thread
    setThreads((prevThreads) => prevThreads.filter((thread) => thread.id !== threadId));
  
    // Remove comments associated with the deleted thread
    setComments((prevComments) =>
      prevComments.filter((comment) => comment.thread_id !== threadId)
    );
  };

  if (error) {
    return <Typography color="error">{error}</Typography>;
  }

  if (!user) {
    return <Typography>Loading...</Typography>;
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h4" gutterBottom>
          {user.username}'s Profile
        </Typography>
        <Typography variant="body1" gutterBottom>
          Email: {user.email}
        </Typography>
        <Divider style={{ margin: '1rem 0' }} />

        <Typography variant="h5">Your Threads</Typography>
        {threads.length === 0 ? (
          <Typography>No threads created.</Typography>
        ) : (
          threads.map((thread) => 
            <ThreadCard 
              key={thread.id} 
              thread={thread} 
              onThreadUpdated={(updatedThread) => {
                setThreads((prev) =>
                  prev.map((t) => (t.id === updatedThread.id ? updatedThread : t))
                );
              }}
              onThreadDeleted={handleThreadDeleted}
            />)
        )}

        <Divider style={{ margin: '1rem 0' }} />
        <Typography variant="h5">Comments</Typography>
        {comments.length === 0 ? (
          <Typography>No comments made.</Typography>
        ) : (
          <CommentList 
            comments={comments}
            onCommentUpdated={handleCommentUpdated}
            onCommentDeleted={handleCommentDeleted}
          />
        )}
      </CardContent>
    </Card>
  );
};

export default UserProfile;



