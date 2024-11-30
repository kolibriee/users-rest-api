import { Outlet, useNavigate } from 'react-router-dom'
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  Container,
} from '@mui/material'
import useAuthStore from '../store/authStore'

export default function Layout() {
  const navigate = useNavigate()
  const { isAuthenticated, isAdmin, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            User Management
          </Typography>
          {isAuthenticated() ? (
            <Box>
              <Button color="inherit" onClick={() => navigate('/profile')}>
                Profile
              </Button>
              {isAdmin && (
                <Button color="inherit" onClick={() => navigate('/admin')}>
                  Admin Panel
                </Button>
              )}
              <Button color="inherit" onClick={handleLogout}>
                Logout
              </Button>
            </Box>
          ) : (
            <Box>
              <Button color="inherit" onClick={() => navigate('/login')}>
                Login
              </Button>
              <Button color="inherit" onClick={() => navigate('/register')}>
                Register
              </Button>
            </Box>
          )}
        </Toolbar>
      </AppBar>
      <Container>
        <Outlet />
      </Container>
    </>
  )
}
