import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import { useCalendarSlots, useReservations } from '../api/hooks';
import { useAuth } from '../store/useAuth';
import { useTranslation } from 'react-i18next';
import tippy from 'tippy.js';
import 'tippy.js/dist/tippy.css';
import { useState } from 'react';
import api from '../api/axios';

import ReservationDialog from '../components/ReservationDialog';
import AdminCalendarView from '../components/AdminCalendarView';

// Workaround: FullCalendar typings cause issues with strict JSX checks
const CalendarComponent: any = FullCalendar;

export default function ReservationsPage() {
  const calendar = useCalendarSlots();
  const reservations = useReservations();
  const { user } = useAuth();
  const { t } = useTranslation();

  const [dialogOpen, setDialogOpen] = useState(false);
  const [adminViewOpen, setAdminViewOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<string>('');
  const [selectedSlots, setSelectedSlots] = useState<any[]>([]);

  const refreshReservations = () => {
    window.location.reload();
  };

  const cancelReservation = async (reservationId: number) => {
    if (!confirm('Cancel this reservation?')) return;
    try {
      await api.delete(`/parent/reservations/${reservationId}`);
      refreshReservations();
    } catch (err: any) {
      console.error('Failed to cancel reservation:', err);
    }
  };

  const events = [] as any[];
  // Reservations
  reservations?.forEach((r: any) => {
    events.push({
      id: `res-${r.id}`,
      title: `${r.status}`,
      start: r.date,
      end: r.date,
      allDay: true,
      backgroundColor: '#60a5fa',
      extendedProps: { type: 'reservation', status: r.status }
    });
  });
  // Availability markers (show to everyone so parents know which days have capacity)
  Object.keys(calendar).forEach(dateStr => {
    const remaining = calendar[dateStr]?.reduce((acc: number, s: any) => acc + s.remaining, 0) || 0;
    // Skip days with zero capacity to reduce clutter
    if (remaining > 0) {
      events.push({
        id: `avail-${dateStr}`,
        title: `${remaining} free`,
        start: dateStr,
        allDay: true,
        backgroundColor: '#34d399',
        extendedProps: { type: 'availability', remaining }
      });
    }
  });

  return (
    <div className="min-h-screen bg-modern">
      <div className="container-modern py-8">
        <div className="page-header">
          <h1 className="page-title">{t('reservations')}</h1>
          <p className="page-subtitle">
            {user?.role === 'admin' 
              ? 'View and manage all reservations' 
              : 'Book and manage your reservations'
            }
          </p>
        </div>
        
        <div className="card p-6">
          <CalendarComponent
            plugins={[dayGridPlugin, interactionPlugin]}
            initialView="dayGridMonth"
            events={events}
            height={600}
            dayMaxEvents={false}
            moreLinkClick="popover"
            dateClick={(info: any) => {
              const dateStr = info.dateStr;
              const slotsForDate = calendar[dateStr] || [];
              const remaining = slotsForDate.reduce((acc: number, s: any) => acc + s.remaining, 0);
              
              if (user?.role === 'admin') {
                // Admin view: show children with reservations for this day
                setSelectedDate(dateStr);
                setAdminViewOpen(true);
              } else {
                // Parent view: show reservation form
                if (slotsForDate.length === 0 || remaining === 0) {
                  alert('No available slots on this day');
                  return;
                }
                setSelectedDate(dateStr);
                setSelectedSlots(slotsForDate);
                setDialogOpen(true);
              }
            }}
            eventDidMount={(arg: any) => {
              const { type, status, remaining } = arg.event.extendedProps as any;
              let content = '';
              if (type === 'reservation') {
                content = `Reservation: ${status}`;
              } else if (type === 'availability') {
                content = `${remaining} free slots`;
              }
              if (content) {
                tippy(arg.el, { content });
              }
            }}
            dayCellClassNames={(arg: any) => {
              const dateStr = arg.date.toISOString().split('T')[0];
              const slotsForDate = calendar[dateStr] || [];
              const remaining = slotsForDate.reduce((acc: number, s: any) => acc + s.remaining, 0);
              
              if (slotsForDate.length > 0 && remaining > 0) {
                return 'cursor-pointer hover:bg-blue-50 hover:border-blue-300 transition-all duration-200';
              }
              return '';
            }}
          />
        </div>

        {/* Reservations List */}
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">My Reservations</h2>
          <div className="bg-white/80 rounded-xl shadow p-4">
            {reservations && reservations.length > 0 ? (
              <div className="space-y-2">
                {reservations.map((reservation: any) => (
                  <div key={reservation.id} className="flex items-center justify-between p-4 bg-white rounded-lg shadow-sm border border-gray-200">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <div className="font-semibold text-gray-800">
                          {new Date(reservation.date).toLocaleDateString('en-US', { 
                            weekday: 'long', 
                            month: 'short', 
                            day: 'numeric' 
                          })}
                        </div>
                        <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                          reservation.status === 'approved' ? 'bg-green-100 text-green-800' :
                          reservation.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                          'bg-red-100 text-red-800'
                        }`}>
                          {reservation.status}
                        </span>
                      </div>
                      <div className="text-sm text-gray-600">
                        Child: {reservation.child?.name || 'Unknown'} • 
                        Capacity: {reservation.slot?.capacity || 'N/A'} spots
                      </div>
                    </div>
                    {reservation.status === 'pending' && (
                      <button
                        onClick={() => cancelReservation(reservation.id)}
                        className="text-red-600 hover:text-red-800 text-sm font-medium px-3 py-1 rounded hover:bg-red-50 transition-colors"
                      >
                        Cancel
                      </button>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-500 text-center py-4">No reservations yet</p>
            )}
          </div>
        </div>

        {/* Optionally show per-date remaining capacity */}
        {user?.role === 'admin' && (
          <div className="mt-8">
            <h2 className="text-xl font-semibold mb-2">Daily Availability</h2>
            <ul className="flex flex-wrap gap-4">
              {Object.keys(calendar).map(date => {
                const totalRemaining = calendar[date].reduce((acc: number, s: any) => acc + s.remaining, 0);
                return (
                  <li key={date} className="bg-white/80 rounded shadow px-4 py-2">
                    <div className="font-bold">{date}</div>
                    <div className="text-sm text-gray-600">Remaining: {totalRemaining}</div>
                  </li>
                );
              })}
            </ul>
          </div>
        )}

        <ReservationDialog
          open={dialogOpen}
          date={selectedDate}
          slots={selectedSlots}
          onClose={() => setDialogOpen(false)}
          onReserved={() => { refreshReservations(); }}
        />
        
        <AdminCalendarView
          open={adminViewOpen}
          date={selectedDate}
          onClose={() => setAdminViewOpen(false)}
        />
      </div>
    </div>
  );
} 