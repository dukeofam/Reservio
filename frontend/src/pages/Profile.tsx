import { useState } from 'react';
import { useAuth } from '../store/useAuth';
import { useToast } from '../components/Toast';
import api from '../api/axios';
import { useTranslation } from 'react-i18next';

export default function ProfilePage() {
  const { user, fetchProfile } = useAuth();
  const [email, setEmail] = useState(user?.email || '');
  const [firstName, setFirstName] = useState(user?.firstName || '');
  const [lastName, setLastName] = useState(user?.lastName || '');
  const [phone, setPhone] = useState(user?.phone || '');
  const [password, setPassword] = useState('');
  const [children, setChildren] = useState(user?.children || []);
  const [childName, setChildName] = useState('');
  const [childAge, setChildAge] = useState('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const toast = useToast();
  const { t } = useTranslation();

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      const url = URL.createObjectURL(file);
      setPreviewUrl(url);
    }
  };

  const handleUploadPicture = async () => {
    if (!selectedFile) {
      toast('Please select a file', 'error');
      return;
    }

    setUploading(true);
    try {
      const formData = new FormData();
      formData.append('file', selectedFile);

      await api.post('/user/profile-picture', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      await fetchProfile();
      toast('Profile picture uploaded successfully', 'success');
      setSelectedFile(null);
      setPreviewUrl(null);
    } catch (err: any) {
      toast(err.response?.data?.error || 'Upload failed', 'error');
    } finally {
      setUploading(false);
    }
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.put('/user/profile', {
        email: email.trim(),
        firstName: firstName.trim(),
        lastName: lastName.trim(),
        phone: phone.trim(),
        password: password.trim() || undefined
      });
      await fetchProfile();
      toast('Profile saved', 'success');
      setPassword('');
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  const handleAddChild = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!childName || !childAge) {
      toast('Please enter child name and age', 'error');
      return;
    }
    try {
      const res = await api.post('/parent/children', { name: childName, age: Number(childAge) });
      setChildren((prev: any) => [...prev, res.data.child]);
      setChildName('');
      setChildAge('');
      toast('Child added', 'success');
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  return (
    <div className="p-6 max-w-xl mx-auto space-y-8">
      <h1 className="text-2xl font-bold mb-4">{t('profile') ?? 'Profile'}</h1>
      
      {/* Profile Picture Section */}
      <div className="space-y-4">
        <h2 className="text-lg font-semibold">Profile Picture</h2>
        <div className="flex items-center space-x-4">
          {user?.profilePicture && (
            <img 
              src={user.profilePicture} 
              alt="Profile" 
              className="w-20 h-20 rounded-full object-cover border-2 border-gray-200"
            />
          )}
          <div className="flex-1">
            <input
              type="file"
              accept="image/*"
              onChange={handleFileSelect}
              className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary file:text-white hover:file:bg-primary-dark"
            />
            {previewUrl && (
              <div className="mt-2">
                <img src={previewUrl} alt="Preview" className="w-16 h-16 rounded-full object-cover" />
                <button
                  type="button"
                  onClick={handleUploadPicture}
                  disabled={uploading}
                  className="mt-2 bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded text-sm disabled:opacity-50"
                >
                  {uploading ? 'Uploading...' : 'Upload Picture'}
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      <form onSubmit={handleSave} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">First name</label>
            <input value={firstName} onChange={e => setFirstName(e.target.value)} className="w-full border px-3 py-2 rounded" required />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Last name</label>
            <input value={lastName} onChange={e => setLastName(e.target.value)} className="w-full border px-3 py-2 rounded" required />
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Phone</label>
          <input value={phone} onChange={e => setPhone(e.target.value)} className="w-full border px-3 py-2 rounded" required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Email</label>
          <input value={email} onChange={e => setEmail(e.target.value)} type="email" className="w-full border px-3 py-2 rounded" required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">New password</label>
          <input value={password} onChange={e => setPassword(e.target.value)} type="password" className="w-full border px-3 py-2 rounded" />
        </div>
        <button type="submit" className="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded">Save</button>
      </form>

      <div className="mt-8">
        <h2 className="text-xl font-semibold mb-2">Children</h2>
        <ul className="mb-4 space-y-1">
          {children && children.length > 0 ? children.map((child: any) => (
            <li key={child.id} className="bg-gray-100 rounded px-3 py-2 flex items-center gap-2">
              <span className="font-medium">{child.name}</span>
              <span className="text-gray-500 text-sm">({child.age})</span>
            </li>
          )) : <li className="text-gray-500">No children yet.</li>}
        </ul>
        <form onSubmit={handleAddChild} className="flex gap-2 items-end">
          <input value={childName} onChange={e => setChildName(e.target.value)} placeholder="Child name" className="border px-3 py-2 rounded w-1/2" required />
          <input value={childAge} onChange={e => setChildAge(e.target.value)} placeholder="Age" type="number" min="0" max="18" className="border px-3 py-2 rounded w-1/4" required />
          <button type="submit" className="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded">Add</button>
        </form>
      </div>
    </div>
  );
} 