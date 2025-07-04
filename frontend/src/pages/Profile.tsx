import { useState } from 'react';
import { useAuth } from '../store/useAuth';
import { useToast } from '../components/Toast';
import api from '../api/axios';
import { useTranslation } from 'react-i18next';

export default function ProfilePage() {
  const { user, fetchProfile } = useAuth();
  const [email, setEmail] = useState(user?.email || '');
  const [password, setPassword] = useState('');
  const toast = useToast();
  const { t } = useTranslation();

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.put('/user/profile', { email: email.trim(), password: password.trim() || undefined });
      await fetchProfile();
      toast('Profile saved', 'success');
      setPassword('');
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  return (
    <div className="p-6 max-w-xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold mb-4">{t('profile') ?? 'Profile'}</h1>
      <form onSubmit={handleSave} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Email</label>
          <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" className="w-full border px-3 py-2 rounded" required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">New password</label>
          <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" className="w-full border px-3 py-2 rounded" />
        </div>
        <button type="submit" className="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded">Save</button>
      </form>
    </div>
  );
} 