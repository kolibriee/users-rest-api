import React, { useState, useEffect } from 'react'
import {
  Container,
  Paper,
  Typography,
  Box,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Link
} from '@mui/material'
import { useNavigate, Link as RouterLink } from 'react-router-dom'
import { toast } from 'react-toastify'
import { getCurrentUser, updateCurrentUser, deleteCurrentUser } from '../api'
import useAuthStore from '../store/authStore'

const Profile = () => {
  const navigate = useNavigate()
  const { logout, user, isAdmin } = useAuthStore()
  const [userData, setUserData] = useState(null)
  const [isEditing, setIsEditing] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
  const [editData, setEditData] = useState({
    name: '',
    username: '',
    email: '',
    password: '',
    city: ''
  })

  useEffect(() => {
    fetchUserData()
  }, [])

  const fetchUserData = async () => {
    try {
      const response = await getCurrentUser()
      console.log('Raw response:', response)
      console.log('User data response:', response.data)
      const userData = {
        id: response.data.ID,
        name: response.data.Name,
        username: response.data.Username,
        email: response.data.Email || '', // Добавляем fallback на пустую строку
        city: response.data.City,
        role: response.data.Role
      }
      console.log('Processed userData:', userData)
      setUserData(userData)
      setEditData({
        name: userData.name,
        username: userData.username,
        email: userData.email || '', // Добавляем fallback на пустую строку
        password: '',
        city: userData.city || ''
      })
    } catch (error) {
      console.error('Error fetching user data:', error.response || error)
      const message = error.response?.data?.message || 'Failed to load user data'
      toast.error(message)
    }
  }

  const handleEditChange = (e) => {
    setEditData({
      ...editData,
      [e.target.name]: e.target.value
    })
  }

  const handleEditSubmit = async () => {
    try {
      await updateCurrentUser(editData)
      setIsEditing(false)
      await fetchUserData()  // Ждем завершения обновления данных
      toast.success('Profile updated successfully')
    } catch (error) {
      console.error('Error updating profile:', error)
      const message = error.response?.data?.message || 'Failed to update profile'
      toast.error(message)
    }
  }

  const handleDelete = async () => {
    try {
      await deleteCurrentUser()
      logout()
      navigate('/login')
      toast.success('Account deleted successfully')
    } catch (error) {
      console.error('Error deleting account:', error)
      const message = error.response?.data?.message || 'Failed to delete account'
      toast.error(message)
    }
  }

  return (
    <Container maxWidth="sm">
      <Paper elevation={3} sx={{ p: 4, mt: 4 }}>
        <Typography variant="h4" gutterBottom>
          Profile
        </Typography>
        {isAdmin && (
          <Box sx={{ mb: 2 }}>
            <Button
              component={RouterLink}
              to="/admin"
              variant="contained"
              color="primary"
            >
              Go to Admin Panel
            </Button>
          </Box>
        )}
        {userData ? (
          <>
            {!isEditing ? (
              <Box>
                <Typography variant="h6">Name: {userData.name}</Typography>
                <Typography variant="h6">Username: {userData.username}</Typography>
                <Typography variant="h6">Email: <a href={`mailto:${userData.email}`}>{userData.email}</a></Typography>
                {userData.city && (
                  <Typography variant="h6">City: {userData.city}</Typography>
                )}
                <Typography variant="h6">Role: {userData.role}</Typography>
                <Box sx={{ mt: 3, display: 'flex', gap: 2, flexDirection: 'column' }}>
                  <Box sx={{ display: 'flex', gap: 2 }}>
                    <Button
                      variant="contained"
                      onClick={() => setIsEditing(true)}
                      fullWidth
                    >
                      Edit Profile
                    </Button>
                    <Button
                      variant="contained"
                      color="error"
                      onClick={() => setIsDeleteDialogOpen(true)}
                      fullWidth
                    >
                      Delete Account
                    </Button>
                  </Box>
                  <Button
                    variant="outlined"
                    color="secondary"
                    onClick={() => {
                      logout();
                      navigate('/login');
                    }}
                    fullWidth
                  >
                    Logout
                  </Button>
                </Box>
              </Box>
            ) : (
              <Box component="form">
                <TextField
                  fullWidth
                  label="Name"
                  name="name"
                  value={editData.name}
                  onChange={handleEditChange}
                  margin="normal"
                />
                <TextField
                  fullWidth
                  label="Username"
                  name="username"
                  value={editData.username}
                  onChange={handleEditChange}
                  margin="normal"
                />
                <TextField
                  fullWidth
                  label="Email"
                  name="email"
                  value={editData.email}
                  onChange={handleEditChange}
                  margin="normal"
                />
                <TextField
                  fullWidth
                  label="Password"
                  name="password"
                  type="password"
                  value={editData.password}
                  onChange={handleEditChange}
                  margin="normal"
                  helperText="Leave empty to keep current password"
                />
                <TextField
                  fullWidth
                  label="City"
                  name="city"
                  value={editData.city}
                  onChange={handleEditChange}
                  margin="normal"
                />
                <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
                  <Button
                    variant="contained"
                    onClick={handleEditSubmit}
                    fullWidth
                  >
                    Save Changes
                  </Button>
                  <Button
                    variant="outlined"
                    onClick={() => setIsEditing(false)}
                    fullWidth
                  >
                    Cancel
                  </Button>
                </Box>
              </Box>
            )}
          </>
        ) : (
          <Typography>Loading...</Typography>
        )}
      </Paper>

      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => setIsDeleteDialogOpen(false)}
      >
        <DialogTitle>Delete Account</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete your account? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setIsDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleDelete} color="error">
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  )
}

export default Profile
