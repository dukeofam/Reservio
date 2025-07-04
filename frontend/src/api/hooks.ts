import { useEffect, useState } from 'react';
import api from './axios';

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
    api.get('/announcements').then(r => setAnnouncements(r.data)).catch(() => {});
  }, []);
  return announcements;
}

export function useSlots() {
  const [slots, setSlots] = useState<any[]>([]);
  useEffect(() => {
    api.get('/slots').then(r => setSlots(r.data)).catch(() => {});
  }, []);
  return slots;
}

export function useReservations() {
  const [reservations, setReservations] = useState<any[]>([]);
  useEffect(() => {
    api.get('/reservations').then(r => setReservations(r.data)).catch(() => {});
  }, []);
  return reservations;
}

export function useChildren() {
  const [children, setChildren] = useState<any[]>([]);
  useEffect(() => {
    api.get('/children').then(r => setChildren(r.data)).catch(() => {});
  }, []);
  return children;
} 