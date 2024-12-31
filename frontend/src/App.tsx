import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { CssBaseline, Container } from '@mui/material';

import Navbar from './components/Layout/Navbar';
import Footer from './components/Layout/Footer';

import Home from './pages/Home';
import CreateThread from './pages/CreateThread';
import UserProfile from './pages/UserProfile';
import ThreadDetails from './pages/ThreadDetails';
import Login from './pages/Login';
import CreateAccount from './pages/CreateAccount';
import ProtectedRoute from './components/ProtectedRoute/ProtectedRoute';
import './App.css';

const App: React.FC = () => {
  return (
    <Router>
      <CssBaseline />
      <Navbar />
      <Container>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} /> 
          <Route path="/create-account" element={<CreateAccount />} />
          <Route path="/thread/:id" element={<ThreadDetails />} />
          <Route 
            path="/create-thread" 
            element={
              <ProtectedRoute>
                <CreateThread />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/user/:id" 
            element={
              <ProtectedRoute>
                <UserProfile />
              </ProtectedRoute>
            } 
          />
        </Routes>
      </Container>
      <Footer />
    </Router>
  );
};

export default App;
