import React, { useState, useEffect } from 'react'
import {
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Typography,
  Box,
  Select,
  MenuItem,
  FormControl,
  InputLabel
} from '@mui/material'
import { toast } from 'react-toastify'
import { getAllUsers, createUser, updateUser, deleteUser } from '../api'
import { useNavigate } from 'react-router-dom'

const Admin = () => {
  const navigate = useNavigate()
  const [users, setUsers] = useState([])
  const [openDialog, setOpenDialog] = useState(false)
  const [editingUser, setEditingUser] = useState(null)
  const [formData, setFormData] = useState({
    name: '',
    username: '',
    email: '',
    password: '',
    city: '',
    role: 'user'
  })

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      const response = await getAllUsers()
      console.log('All users response:', response.data)
      // Преобразуем поля к нижнему регистру
      const formattedUsers = response.data.map(user => ({
        id: user.ID,
        name: user.Name,
        username: user.Username,
        email: user.Email,
        city: user.City,
        role: user.Role
      }))
      setUsers(formattedUsers)
    } catch (error) {
      console.error('Error fetching users:', error.response || error)
      const message = error.response?.data?.message || 'Failed to fetch users'
      toast.error(message)
    }
  }

  const handleInputChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
  }

  const handleSubmit = async () => {
    try {
      if (editingUser) {
        await updateUser(editingUser.id, formData)
        toast.success('User updated successfully')
      } else {
        await createUser(formData)
        toast.success('User created successfully')
      }
      setOpenDialog(false)
      setEditingUser(null)
      setFormData({
        name: '',
        username: '',
        email: '',
        password: '',
        city: '',
        role: 'user'
      })
      await fetchUsers()  // Ждем завершения обновления списка
    } catch (error) {
      console.error('Error:', error)
      const message = error.response?.data?.message || (editingUser ? 'Failed to update user' : 'Failed to create user')
      toast.error(message)
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this user?')) return
    try {
      await deleteUser(id)
      toast.success('User deleted successfully')
      await fetchUsers()  // Ждем завершения обновления списка
    } catch (error) {
      console.error('Error deleting user:', error)
      const message = error.response?.data?.message || 'Failed to delete user'
      toast.error(message)
    }
  }

  const handleEdit = (user) => {
    setEditingUser(user)
    setFormData({
      name: user.name,
      username: user.username,
      email: user.email,
      password: '',
      city: user.city || '',
      role: user.role
    })
    setOpenDialog(true)
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ mb: 4, mt: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4">Admin Panel</Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => navigate('/profile')}
        >
          Back to Profile
        </Button>
      </Box>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h4">User Management</Typography>
          <Button
            variant="contained"
            onClick={() => {
              setEditingUser(null)
              setFormData({
                name: '',
                username: '',
                email: '',
                password: '',
                city: '',
                role: 'user'
              })
              setOpenDialog(true)
            }}
          >
            Add User
          </Button>
        </Box>

        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>Username</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>City</TableCell>
                <TableCell>Role</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.id}</TableCell>
                  <TableCell>{user.name}</TableCell>
                  <TableCell>{user.username}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{user.city}</TableCell>
                  <TableCell>{user.role}</TableCell>
                  <TableCell>
                    <Button onClick={() => handleEdit(user)}>Edit</Button>
                    <Button
                      color="error"
                      onClick={() => handleDelete(user.id)}
                    >
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
          <DialogTitle>
            {editingUser ? 'Edit User' : 'Create User'}
          </DialogTitle>
          <DialogContent>
            <TextField
              fullWidth
              label="Name"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              margin="normal"
            />
            <TextField
              fullWidth
              label="Username"
              name="username"
              value={formData.username}
              onChange={handleInputChange}
              margin="normal"
            />
            <TextField
              fullWidth
              label="Email"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleInputChange}
              margin="normal"
            />
            <TextField
              fullWidth
              label="Password"
              name="password"
              type="password"
              value={formData.password}
              onChange={handleInputChange}
              margin="normal"
            />
            <TextField
              fullWidth
              label="City"
              name="city"
              value={formData.city}
              onChange={handleInputChange}
              margin="normal"
            />
            <FormControl fullWidth margin="normal">
              <InputLabel>Role</InputLabel>
              <Select
                name="role"
                value={formData.role}
                onChange={handleInputChange}
              >
                <MenuItem value="user">User</MenuItem>
                <MenuItem value="admin">Admin</MenuItem>
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
            <Button onClick={handleSubmit} variant="contained">
              {editingUser ? 'Update' : 'Create'}
            </Button>
          </DialogActions>
        </Dialog>
      </Paper>
    </Container>
  )
}

export default Admin
