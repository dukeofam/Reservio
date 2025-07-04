import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../store/useAuth';
import { useToast } from '../components/Toast';

export default function Register() {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [childName, setChildName] = useState('');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const { register, loading } = useAuth();
  const toast = useToast();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    // Simple client-side validation
    if (!firstName || !lastName || !childName || !phone || !email || !password || !confirm) {
      toast('Please fill all fields', 'error');
      return;
    }

    if (password !== confirm) {
      toast('Passwords do not match', 'error');
      return;
    }

    try {
      await register(email, password);
      toast(`Welcome, ${firstName}!`, 'success');
      navigate('/dashboard');
    } catch (err: any) {
      toast(err.response?.data?.error || 'Registration failed', 'error');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 px-4">
      <form onSubmit={handleSubmit} className="bg-white p-8 rounded shadow w-full max-w-md space-y-3">
        <h1 className="text-2xl font-bold text-center">Create account</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <input value={firstName} onChange={(e) => setFirstName(e.target.value)} placeholder="First name" className="border px-3 py-2 rounded" required />
          <input value={lastName} onChange={(e) => setLastName(e.target.value)} placeholder="Last name" className="border px-3 py-2 rounded" required />
        </div>
        <input value={childName} onChange={(e) => setChildName(e.target.value)} placeholder="Child name" className="w-full border px-3 py-2 rounded" required />
        <input value={phone} onChange={(e) => setPhone(e.target.value)} type="tel" placeholder="Phone number" className="w-full border px-3 py-2 rounded" required />
        <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" placeholder="Email" className="w-full border px-3 py-2 rounded" required />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="Password" className="border px-3 py-2 rounded" required />
          <input value={confirm} onChange={(e) => setConfirm(e.target.value)} type="password" placeholder="Password again" className="border px-3 py-2 rounded" required />
        </div>
        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded" disabled={loading}>Register</button>
        <p className="text-center text-sm">Already have an account? <span className="text-blue-600 cursor-pointer" onClick={() => navigate('/login')}>Login</span></p>
      </form>
    </div>
  );
} 