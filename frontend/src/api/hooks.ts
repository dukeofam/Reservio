import { useEffect, useState } from 'react';
import api from './axios';
import { useAuth } from '../store/useAuth';

export function useDashboardStats() {
  const [stats, setStats] = useState<any>(null);
  useEffect(() => {
    api.get('/dashboard/stats').then(r => setStats(r.data)).catch(() => {});
  }, []);
  return stats;
}

export function useAnnouncements() {
  const [announcements, setAnnouncements] = useState<any[]>([]);
  useEffect(() => {
    api.get('/announcements').then(r => setAnnouncements(r.data.data)).catch(() => {});
  }, []);
  return announcements;
}

export function useSlots() {
  const [slots, setSlots] = useState<any[]>([]);
  useEffect(() => {
    api.get('/slots').then(r => setSlots(r.data.data)).catch(() => {});
  }, []);
  return slots;
}

export function useReservations() {
  const [reservations, setReservations] = useState<any[]>([]);
  const { user } = useAuth();
  useEffect(() => {
    const path = user?.role === 'admin' ? '/admin/reservations' : '/parent/reservations';
    api.get(path).then(r => setReservations(r.data.data || [])).catch(() => {});
  }, []);
  return reservations;
}

export function useChildren() {
  const [children, setChildren] = useState<any[]>([]);
  const { user } = useAuth();
  const refresh = () => {
    const path = user?.role === 'admin' ? '/admin/children' : '/parent/children';
    api.get(path).then(r => setChildren(r.data.data || r.data.children || [])).catch(() => {});
  };
  useEffect(() => { refresh(); }, []);
  return { children, refresh };
}

export function useAdminChildren() {
  const [children, setChildren] = useState<any[]>([]);
  useEffect(() => {
    api.get('/admin/children').then(r => setChildren(r.data.data)).catch(() => {});
  }, []);
  return children;
}

export function useCalendarSlots() {
  const [calendar, setCalendar] = useState<any>({});
  useEffect(() => {
    api.get('/slots/calendar').then(r => setCalendar(r.data.calendar)).catch(() => {});
  }, []);
  return calendar;
}

export function useAdminSlots() {
  const [slots, setSlots] = useState<any[]>([]);
  const refresh = () => api.get('/admin/slots').then(r => setSlots(r.data.data)).catch(() => {});
  useEffect(() => { refresh(); }, []);
  return { slots, refresh };
}

export function useAdminUsers() {
  const [users, setUsers] = useState<any[]>([]);
  const refresh = () => api.get('/admin/users').then(r => setUsers(r.data.data || [])).catch(() => {});
  useEffect(() => { refresh(); }, []);
  return { users, refresh };
} 