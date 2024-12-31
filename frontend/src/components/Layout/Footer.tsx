import React from 'react';
import { Box, Typography } from '@mui/material';

const Footer: React.FC = () => {
  return (
    <Box sx={{ textAlign: 'center', padding: '16px', bgcolor: 'background.default' }}>
      <Typography variant="body2" color="textSecondary">
        © 2024 CommonCircle - All Rights Reserved
      </Typography>
    </Box>
  );
};

export default Footer;
