import { create } from 'zustand';
import api from '../api/axios';

interface User {
  id: number;
  email: string;
  role: string;
  firstName?: string;
  lastName?: string;
  phone?: string;
  profilePicture?: string;
  children?: any[];
}

interface AuthState {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  fetchProfile: () => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
}

export const useAuth = create<AuthState>((set) => ({
  user: null,
  loading: false,
  login: async (email, password) => {
    set({ loading: true });
    try {
      const response = await api.post('/auth/login', { email, password });
      await (useAuth.getState().fetchProfile)();
    } catch (error) {
      set({ loading: false });
      throw error;
    }
    set({ loading: false });
  },
  logout: async () => {
    await api.post('/auth/logout');
    set({ user: null });
  },
  fetchProfile: async () => {
    try {
      const res = await api.get('/user/profile');
      set({ user: res.data.user });
    } catch {
      set({ user: null });
    }
  },
  register: async (email, password) => {
    set({ loading: true });
    try {
      const response = await api.post('/auth/register', { email, password });
      await (useAuth.getState().fetchProfile)();
    } catch (error) {
      set({ loading: false });
      throw error;
    }
    set({ loading: false });
  }
})); 