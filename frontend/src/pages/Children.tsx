import { useState } from 'react';
import { useChildren, useAdminChildren } from '../api/hooks';
import { useAuth } from '../store/useAuth';
import { useTranslation } from 'react-i18next';
import Modal from '../components/Modal';
import api from '../api/axios';
import { useToast } from '../components/Toast';
import { PlusIcon, UserGroupIcon } from '@heroicons/react/24/solid';

export default function ChildrenPage() {
  const { user } = useAuth();
  const childrenData = user?.role === 'admin' ? useAdminChildren() : useChildren();
  const children = user?.role === 'admin' ? childrenData : childrenData.children;
  const refresh = user?.role === 'admin' ? null : childrenData.refresh;
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [editChild, setEditChild] = useState<any>(null);
  const [name, setName] = useState('');
  const [birthdate, setBirthdate] = useState('');
  const [age, setAge] = useState('');
  const toast = useToast();

  const openNew = () => {
    setEditChild(null);
    setName('');
    setBirthdate('');
    setAge('');
    setOpen(true);
  };

  const openEdit = (child: any) => {
    setEditChild(child);
    setName(child.name);
    setBirthdate(child.birthdate || '');
    setAge(child.age?.toString() || '');
    setOpen(true);
  };

  const saveChild = async () => {
    try {
      const data: any = { name };
      if (age) data.age = parseInt(age);

      if (editChild) {
        await api.put(`/parent/children/${editChild.id}`, data);
        toast('Child updated', 'success');
      } else {
        await api.post('/parent/children', data);
        toast('Child added', 'success');
      }
      setOpen(false);
      // Reset form
      setEditChild(null);
      setName('');
      setAge('');
      setBirthdate('');
      // Refresh children data
      if (refresh) {
        refresh();
      } else {
        window.location.reload(); // Fallback for admin users
      }
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  const deleteChild = async (id: number) => {
    if (!confirm('Delete this child? This will also delete all their reservations.')) return;
    try {
      await api.delete(`/parent/children/${id}`);
      toast('Child deleted', 'success');
      // Refresh children data
      if (refresh) {
        refresh();
      } else {
        window.location.reload(); // Fallback for admin users
      }
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
              <h1 className="page-title">{t('children')}</h1>
              <p className="page-subtitle">
                {user?.role === 'admin' 
                  ? 'Manage all children in the system' 
                  : 'Manage your children\'s profiles'
                }
              </p>
            </div>
            {user?.role === 'parent' && (
              <button
                onClick={openNew}
                className="btn btn-primary btn-lg group"
              >
                <PlusIcon className="w-5 h-5 mr-2" />
                Add Child
              </button>
            )}
          </div>
        </div>
        <div className="card">
          {children && children.length > 0 ? (
            <div className="divide-y divide-gray-200">
              {children.map((child: any, index: number) => (
                <div 
                  key={child.id} 
                  className="p-6 hover:bg-gray-50 transition-colors animate-fade-in"
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-full flex items-center justify-center">
                        <span className="text-white font-bold text-lg">
                          {child.name.charAt(0).toUpperCase()}
                        </span>
                      </div>
                      <div>
                        <h3 className="font-semibold text-gray-900 text-lg">{child.name}</h3>
                        <p className="text-gray-600">Age: {child.age} years old</p>
                        {user?.role === 'admin' && (
                          <p className="text-sm text-gray-500">
                            Parent: {child.parent?.email || 'Unknown'}
                          </p>
                        )}
                      </div>
                    </div>
                    
                    {user?.role === 'parent' && (
                      <div className="flex gap-2">
                        <button 
                          onClick={() => openEdit(child)} 
                          className="btn btn-secondary btn-sm"
                        >
                          Edit
                        </button>
                        <button 
                          onClick={() => deleteChild(child.id)} 
                          className="btn btn-danger btn-sm"
                        >
                          Delete
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <UserGroupIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No children found</h3>
              <p className="text-gray-600 mb-4">
                {user?.role === 'admin' 
                  ? 'No children have been registered yet.' 
                  : 'Add your first child to get started.'
                }
              </p>
              {user?.role === 'parent' && (
                <button
                  onClick={openNew}
                  className="btn btn-primary"
                >
                  <PlusIcon className="w-4 h-4 mr-2" />
                  Add Child
                </button>
              )}
            </div>
          )}
        </div>

        {user?.role === 'parent' && (
          <Modal open={open} onClose={() => setOpen(false)}>
            <div className="text-center mb-6">
              <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-xl flex items-center justify-center mx-auto mb-4">
                <UserGroupIcon className="w-6 h-6 text-white" />
              </div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                {editChild ? 'Edit Child' : 'Add Child'}
              </h2>
              <p className="text-gray-600">
                {editChild ? 'Update your child\'s information' : 'Add a new child to your account'}
              </p>
            </div>
            
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Child's Name *</label>
                <input 
                  type="text" 
                  value={name} 
                  onChange={e => setName(e.target.value)} 
                  className="input" 
                  placeholder="Enter child's name"
                  required 
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Age *</label>
                        <input 
                          type="number" 
                          value={age} 
                          onChange={e => setAge(e.target.value)} 
                          min="2" 
                          max="5" 
                          className="input" 
                          placeholder="Enter child's age"
                          required
                        />
                        <p className="text-xs text-gray-500 mt-1">Age must be between 2 and 5 years</p>
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
                  onClick={saveChild} 
                  className="btn btn-primary flex-1"
                >
                  {editChild ? 'Update Child' : 'Add Child'}
                </button>
              </div>
            </div>
          </Modal>
        )}
      </div>
    </div>
  );
} 