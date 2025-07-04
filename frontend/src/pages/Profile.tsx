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
  const toast = useToast();
  const { t } = useTranslation();

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