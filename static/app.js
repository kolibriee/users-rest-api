// API Base URL
const API_BASE_URL = 'http://localhost:8080/api/v1';

// Bootstrap components
const toast = new bootstrap.Toast(document.getElementById('toast'));
const editModal = new bootstrap.Modal(document.getElementById('editUserModal'));

// Show notification
function showNotification(message, success = true) {
    const toastElement = document.getElementById('toast');
    const toastMessage = document.getElementById('toastMessage');
    toastMessage.textContent = message;
    toastElement.classList.remove('bg-danger', 'bg-success');
    toastElement.classList.add(success ? 'bg-success' : 'bg-danger');
    toast.show();
}

// Fetch all users
async function fetchUsers() {
    try {
        const response = await fetch(`${API_BASE_URL}/users`);
        if (!response.ok) throw new Error('Failed to fetch users');
        const users = await response.json();
        displayUsers(users);
    } catch (error) {
        showNotification(error.message, false);
    }
}

// Display users in the table
function displayUsers(users) {
    const usersList = document.getElementById('usersList');
    usersList.innerHTML = '';
    
    users.forEach(user => {
        const row = document.createElement('tr');
        row.classList.add('user-row');
        row.innerHTML = `
            <td>${user.id}</td>
            <td>${user.name}</td>
            <td>${user.email}</td>
            <td>
                <button class="btn btn-sm btn-primary me-1" onclick="openEditModal(${JSON.stringify(user).replace(/"/g, '&quot;')})">
                    Edit
                </button>
                <button class="btn btn-sm btn-danger" onclick="deleteUser(${user.id})">
                    Delete
                </button>
            </td>
        `;
        usersList.appendChild(row);
    });
}

// Create new user
document.getElementById('createUserForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const userData = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value
    };

    try {
        const response = await fetch(`${API_BASE_URL}/users`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        });

        if (!response.ok) throw new Error('Failed to create user');
        
        showNotification('User created successfully');
        document.getElementById('createUserForm').reset();
        fetchUsers();
    } catch (error) {
        showNotification(error.message, false);
    }
});

// Open edit modal with user data
function openEditModal(user) {
    document.getElementById('editUserId').value = user.id;
    document.getElementById('editName').value = user.name;
    document.getElementById('editEmail').value = user.email;
    editModal.show();
}

// Update user
async function updateUser() {
    const userId = document.getElementById('editUserId').value;
    const userData = {
        name: document.getElementById('editName').value,
        email: document.getElementById('editEmail').value
    };

    try {
        const response = await fetch(`${API_BASE_URL}/users/${userId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        });

        if (!response.ok) throw new Error('Failed to update user');
        
        showNotification('User updated successfully');
        editModal.hide();
        fetchUsers();
    } catch (error) {
        showNotification(error.message, false);
    }
}

// Delete user
async function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user?')) return;

    try {
        const response = await fetch(`${API_BASE_URL}/users/${userId}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete user');
        
        showNotification('User deleted successfully');
        fetchUsers();
    } catch (error) {
        showNotification(error.message, false);
    }
}

// Refresh users list
function refreshUsers() {
    fetchUsers();
}

// Initial load
document.addEventListener('DOMContentLoaded', () => {
    fetchUsers();
});
