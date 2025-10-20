import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../store/useAuth';
import { useToast } from '../components/Toast';
import { useTranslation } from 'react-i18next';
import { EyeIcon, EyeSlashIcon, UserCircleIcon, LockClosedIcon } from '@heroicons/react/24/solid';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
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
    <div className="min-h-screen bg-modern flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg">
            <span className="text-white font-bold text-2xl">R</span>
          </div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Welcome to Reservio</h1>
          <p className="text-gray-600">Kindergarten Management System</p>
        </div>

        {/* Login Form */}
        <div className="card p-8 animate-fade-in">
          <div className="text-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-2">{t('welcome')}</h2>
            <p className="text-gray-600">{t('signInSubtitle')}</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('email')}
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <UserCircleIcon className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Enter your email"
                  className="input pl-10"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('password')}
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <LockClosedIcon className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter your password"
                  className="input pl-10 pr-10"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  {showPassword ? (
                    <EyeSlashIcon className="h-5 w-5 text-gray-400 hover:text-gray-600" />
                  ) : (
                    <EyeIcon className="h-5 w-5 text-gray-400 hover:text-gray-600" />
                  )}
                </button>
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary btn-lg w-full"
            >
              {loading ? (
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
              ) : (
                t('login')
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              {t('dontHaveAccount')}{' '}
              <button
                onClick={() => navigate('/register')}
                className="text-blue-600 hover:text-blue-800 font-medium transition-colors"
              >
                {t('register')}
              </button>
            </p>
          </div>
        </div>

        {/* Demo Accounts */}
        <div className="mt-6 card p-4 bg-blue-50 border-blue-200">
          <h3 className="text-sm font-medium text-blue-900 mb-2">Demo Accounts</h3>
          <div className="text-xs text-blue-700 space-y-1">
            <p><strong>Admin:</strong> admin@reservio.com / admin123</p>
            <p><strong>Parent:</strong> parent@test.com / parent123</p>
          </div>
        </div>
      </div>
    </div>
  );
} 