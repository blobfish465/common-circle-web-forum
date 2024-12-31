import React from 'react';
import ThreadForm from '../components/Thread/ThreadForm';
import { useAuth } from '../context/AuthContext'; 

const CreateThread: React.FC = () => {
  // retrieve currently authenticated user using useAuth hook
  // Get the current user's ID from the AuthContext
  const { userId } = useAuth(); 

  if (!userId) {
    return <p>You must be logged in to create a thread.</p>; 
  }

  return (
    <div>
      <h1>Create a New Thread</h1>
      {/* pass current userId to ThreadForm*/}
      <ThreadForm userId={Number(userId)} /> 
    </div>
  );
};

export default CreateThread;

