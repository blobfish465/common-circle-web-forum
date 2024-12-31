import React, { useEffect, useState } from 'react';
import { getAllThreads, getThreadsByCategory } from '../api/threadsAPI';
import { getCategories } from '../api/categoriesAPI';
import { Category } from '../types/category';
import { Thread } from '../types/thread';
import ThreadCard from '../components/Thread/ThreadCard';
import { Button, Box, Typography, Grid2, CircularProgress } from '@mui/material';

// Where all the threads are displayed
const Home: React.FC = () => {
  const [threads, setThreads] = useState<Thread[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [activeCategory, setActiveCategory] = useState<number | null>(null);
  const [loading, setLoading] = useState<boolean>(true);


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

  useEffect(() => {
    const fetchThreads = async () => {
      setLoading(true);
      try {
        const threadsData = activeCategory
          ? await getThreadsByCategory(activeCategory)
          : await getAllThreads();
        setThreads(threadsData);
      } catch (error) {
        console.error('Error fetching threads:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchThreads();
  }, [activeCategory]);

  const handleThreadUpdated = (updatedThread: Thread) => {
    setThreads((prev) =>
      prev.map((thread) => (thread.id === updatedThread.id ? updatedThread : thread))
    );
  };

  const handleThreadDeleted = (threadId: number) => {
    setThreads((prev) => prev.filter((thread) => thread.id !== threadId));
  };

  return (
    <div>
      {/* Categories Filter */}
      <Box mb={3} mt={2}>
        <Typography variant="h5" gutterBottom>
          Categories
        </Typography>
        <Box display="flex" flexWrap="wrap" gap={1}>
          <Button
            variant={activeCategory === null ? 'contained' : 'outlined'}
            onClick={() => setActiveCategory(null)}
          >
            All
          </Button>
          {categories.map((category) => (
            <Button
              key={category.id}
              variant={activeCategory === category.id ? 'contained' : 'outlined'}
              onClick={() => setActiveCategory(category.id)}
            >
              {category.name}
            </Button>
          ))}
        </Box>
      </Box>

      {/* Threads List */}
      <Typography variant="h5" gutterBottom>
        {activeCategory
          ? `Threads in "${categories.find((c) => c.id === activeCategory)?.name}"`
          : 'All Threads'}
      </Typography>
      {loading ? (
        <Box display="flex" justifyContent="center" mt={2}>
          <CircularProgress />
        </Box>
      ) : threads.length === 0 ? (
        <Typography>No threads available for this category.</Typography>
      ) : (
        <Grid2 container spacing={2}>
          {threads.map((thread) => (
            <Grid2 size={{xs:12, sm:6, md:4}} key={thread.id}>
              <ThreadCard
                thread={thread}
                onThreadUpdated={handleThreadUpdated}
                onThreadDeleted={handleThreadDeleted}
              />
            </Grid2>
          ))}
        </Grid2>
      )}
    </div>
  );
};

export default Home;
