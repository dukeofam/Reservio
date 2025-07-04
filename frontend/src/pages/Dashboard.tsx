import { useAuth } from '../store/useAuth';
import { useDashboardStats, useAnnouncements, useReservations, useChildren } from '../api/hooks';
import { useTranslation } from 'react-i18next';

export default function Dashboard() {
  const { user } = useAuth();
  const stats = useDashboardStats();
  const announcements = useAnnouncements();
  const reservations = useReservations();
  const children = useChildren();
  const { t } = useTranslation();

  return (
    <div className="p-6 space-y-8 bg-amber-50 min-h-screen bg-[url('/wave.svg')] bg-cover bg-center">
      <h1 className="text-3xl font-extrabold text-primary mb-4">{t('dashboard')}</h1>
      {user?.role === 'admin' ? (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white/80 rounded-xl shadow p-6 text-center">
            <div className="text-2xl font-bold text-primary">{stats?.totalChildren ?? '--'}</div>
            <div className="text-gray-700">Children</div>
          </div>
          <div className="bg-white/80 rounded-xl shadow p-6 text-center">
            <div className="text-2xl font-bold text-primary">{stats?.totalReservations ?? '--'}</div>
            <div className="text-gray-700">Reservations</div>
          </div>
          <div className="bg-white/80 rounded-xl shadow p-6 text-center">
            <div className="text-2xl font-bold text-primary">{stats?.openSlots ?? '--'}</div>
            <div className="text-gray-700">Open slots</div>
          </div>
        </div>
      ) : (
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-2">{t('children')}</h2>
          <ul className="flex flex-wrap gap-4">
            {children?.map((child: any) => (
              <li key={child.id} className="bg-white/80 rounded shadow px-4 py-2">
                <div className="font-bold">{child.name}</div>
                <div className="text-sm text-gray-600">Age: {child.age}</div>
              </li>
            ))}
          </ul>
          <h2 className="text-xl font-semibold mt-6 mb-2">{t('reservations')}</h2>
          <ul className="flex flex-wrap gap-4">
            {reservations?.map((r: any) => (
              <li key={r.id} className="bg-white/80 rounded shadow px-4 py-2">
                <div className="font-bold">{r.date}</div>
                <div className="text-sm text-gray-600">Status: {r.status}</div>
              </li>
            ))}
          </ul>
        </div>
      )}
      <div>
        <h2 className="text-xl font-semibold mb-2">Announcements</h2>
        <ul className="space-y-2">
          {announcements?.map((a: any) => (
            <li key={a.id} className="bg-white/80 rounded shadow px-4 py-2">
              <div className="font-bold">{a.title}</div>
              <div className="text-xs text-gray-500 mb-1">{a.date} • {a.author}</div>
              <div>{a.content}</div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
} 