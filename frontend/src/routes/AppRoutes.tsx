import { Route, Routes, Navigate } from 'react-router-dom';
import { useAuth } from '../store/useAuth';
import Login from '../pages/Login';
import Register from '../pages/Register';
import Dashboard from '../pages/Dashboard';
import ChildrenPage from '../pages/Children';
import ReservationsPage from '../pages/Reservations';
import AdminUsersPage from '../pages/AdminUsers';
import ProfilePage from '../pages/Profile';

export default function AppRoutes() {
  const { user } = useAuth();

  if (!user) {
    return (
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    );
  }

  return (
    <Routes>
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/children" element={<ChildrenPage />} />
      <Route path="/reservations" element={<ReservationsPage />} />
      <Route path="/profile" element={<ProfilePage />} />
      <Route path="/admin/users" element={<AdminUsersPage />} />
      <Route path="*" element={<Navigate to="/dashboard" />} />
    </Routes>
  );
} 