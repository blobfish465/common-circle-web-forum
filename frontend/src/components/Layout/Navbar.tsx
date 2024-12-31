import React from 'react';
import { AppBar, Toolbar, Typography, Button } from '@mui/material';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext'


const Navbar: React.FC = () => {
  const { userId, logout } = useAuth();
  return (
    <AppBar position="sticky">
      <Toolbar>
        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          CommonCircle
        </Typography>
        <Button color="inherit" component={Link} to="/">Home</Button>
        {userId ? (
          <>
            <Button color="inherit" component={Link} to={`/create-thread`}>
              Create Thread
            </Button>
            <Button color="inherit" component={Link} to={`/user/${userId}`}>
              Profile
            </Button>
            <Button color="inherit" onClick={logout}>
              Logout
            </Button>
          </>
        ) : (
          <>
            <Button color="inherit" component={Link} to="/login">
              Login
            </Button>
            <Button color="inherit" component={Link} to="/create-account">
              Create Account
            </Button>
          </>
        )}
      </Toolbar>
    </AppBar>
  );
};

export default Navbar;
