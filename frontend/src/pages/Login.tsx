import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../store/useAuth';
import { useToast } from '../components/Toast';
import { useTranslation } from 'react-i18next';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const { login, loading } = useAuth();
  const toast = useToast();
  const navigate = useNavigate();
  const { t } = useTranslation();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await login(email, password);
      toast(`Welcome back, ${email}!`, 'success');
      navigate('/dashboard');
    } catch (err: any) {
      toast(err.response?.data?.error || 'Login failed', 'error');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-amber-50 px-4 bg-[url('/wave.svg')] bg-cover bg-center">
      <form onSubmit={handleSubmit} className="bg-white/80 backdrop-blur-sm p-8 rounded-xl shadow w-full max-w-md space-y-4">
        <h1 className="text-3xl font-extrabold text-center text-blue-800">{t('welcome')}</h1>
        <p className="text-center text-gray-700 text-sm">{t('signInSubtitle')}</p>
        <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" placeholder={t('email')} className="w-full border px-3 py-2 rounded" required />
        <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder={t('password')} className="w-full border px-3 py-2 rounded" required />
        <button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 transition text-white py-2 rounded" disabled={loading}>{t('login')}</button>
        <p className="text-center text-sm">{t('dontHaveAccount')} <span className="text-blue-600 cursor-pointer" onClick={() => navigate('/register')}>{t('register')}</span></p>
      </form>
    </div>
  );
} 