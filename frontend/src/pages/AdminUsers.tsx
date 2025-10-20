import { useState } from 'react';
import { useAdminUsers } from '../api/hooks';
import Modal from '../components/Modal';
import AddUserModal from '../components/AddUserModal';
import api from '../api/axios';
import { useToast } from '../components/Toast';
import { useTranslation } from 'react-i18next';
import { PlusIcon, UserCircleIcon, TrashIcon, PencilIcon } from '@heroicons/react/24/solid';

export default function AdminUsersPage() {
  const { users, refresh } = useAdminUsers();
  const [open, setOpen] = useState(false);
  const [addUserOpen, setAddUserOpen] = useState(false);
  const [editUser, setEditUser] = useState<any>(null);
  const [role, setRole] = useState('parent');
  const toast = useToast();
  const { t } = useTranslation();

  const openEdit = (user: any) => {
    setEditUser(user);
    setRole(user.role);
    setOpen(true);
  };

  const saveUser = async () => {
    if (!editUser) return;
    try {
      await api.put(`/admin/users/${editUser.id}/role`, { role });
      toast('User role updated', 'success');
      setOpen(false);
      refresh();
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  const deleteUser = async (id: number) => {
    if (!confirm('Delete this user? This will also delete all their children and reservations.')) return;
    try {
      await api.delete(`/admin/users/${id}`);
      toast('User deleted', 'success');
      refresh();
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  return (
    <div className="min-h-screen bg-modern">
      <div className="container-modern py-8">
        <div className="page-header">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="page-title">Manage Users</h1>
              <p className="page-subtitle">Add, edit, and manage user accounts</p>
            </div>
            <button
              onClick={() => setAddUserOpen(true)}
              className="btn btn-primary btn-lg group"
            >
              <PlusIcon className="w-5 h-5 mr-2" />
              Add User
            </button>
          </div>
        </div>

        <div className="card">
          {users && users.length > 0 ? (
            <div className="divide-y divide-gray-200">
              {users.map((user: any, index: number) => (
                <div 
                  key={user.id} 
                  className="p-6 hover:bg-gray-50 transition-colors animate-fade-in"
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-full flex items-center justify-center">
                        <UserCircleIcon className="w-6 h-6 text-white" />
                      </div>
                      <div>
                        <h3 className="font-semibold text-gray-900 text-lg">{user.email}</h3>
                        <p className="text-gray-600">
                          {user.firstName && user.lastName 
                            ? `${user.firstName} ${user.lastName}` 
                            : 'No name provided'
                          }
                        </p>
                        <p className="text-sm text-gray-500">ID: {user.id}</p>
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-4">
                      <span className={`badge ${
                        user.role === 'admin' ? 'badge-danger' : 'badge-info'
                      }`}>
                        {user.role}
                      </span>
                      
                      <div className="flex gap-2">
                        <button
                          onClick={() => openEdit(user)}
                          className="btn btn-secondary btn-sm"
                        >
                          <PencilIcon className="w-4 h-4 mr-1" />
                          Edit Role
                        </button>
                        <button
                          onClick={() => deleteUser(user.id)}
                          className="btn btn-danger btn-sm"
                        >
                          <TrashIcon className="w-4 h-4 mr-1" />
                          Delete
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <UserCircleIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No users found</h3>
              <p className="text-gray-600 mb-4">Get started by adding your first user.</p>
              <button
                onClick={() => setAddUserOpen(true)}
                className="btn btn-primary"
              >
                <PlusIcon className="w-4 h-4 mr-2" />
                Add User
              </button>
            </div>
          )}
        </div>

        <Modal open={open} onClose={() => setOpen(false)}>
          <div className="text-center mb-6">
            <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-xl flex items-center justify-center mx-auto mb-4">
              <PencilIcon className="w-6 h-6 text-white" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Edit User Role</h2>
            <p className="text-gray-600">Update the role for this user</p>
          </div>
          
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
              <input
                type="text"
                value={editUser?.email || ''}
                disabled
                className="input bg-gray-100"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Role</label>
              <select
                value={role}
                onChange={e => setRole(e.target.value)}
                className="input"
              >
                <option value="parent">Parent</option>
                <option value="admin">Administrator</option>
              </select>
            </div>
            
            <div className="flex gap-3 pt-4">
              <button
                type="button"
                onClick={() => setOpen(false)}
                className="btn btn-secondary flex-1"
              >
                Cancel
              </button>
              <button
                onClick={saveUser}
                className="btn btn-primary flex-1"
              >
                Update Role
              </button>
            </div>
          </div>
        </Modal>

        <AddUserModal
          open={addUserOpen}
          onClose={() => setAddUserOpen(false)}
          onUserAdded={refresh}
        />
      </div>
    </div>
  );
}