import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import { useCalendarSlots, useReservations } from '../api/hooks';
import { useAuth } from '../store/useAuth';
import { useTranslation } from 'react-i18next';
import tippy from 'tippy.js';
import 'tippy.js/dist/tippy.css';
import { useState } from 'react';

import ReservationDialog from '../components/ReservationDialog';

// Workaround: FullCalendar typings cause issues with strict JSX checks
const CalendarComponent: any = FullCalendar;

export default function ReservationsPage() {
  const calendar = useCalendarSlots();
  const reservations = useReservations();
  const { user } = useAuth();
  const { t } = useTranslation();

  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<string>('');
  const [selectedSlots, setSelectedSlots] = useState<any[]>([]);

  const refreshReservations = () => window.location.reload();

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
    <div className="p-6 bg-amber-50 min-h-screen bg-[url('/wave.svg')] bg-cover bg-center">
      <h1 className="text-3xl font-extrabold text-primary mb-4">{t('reservations')}</h1>
      <div className="bg-white/80 rounded-xl shadow p-4">
        <CalendarComponent
          plugins={[dayGridPlugin, interactionPlugin]}
          initialView="dayGridMonth"
          events={events}
          height={600}
          dateClick={(info: any) => {
            const dateStr = info.dateStr;
            const slotsForDate = calendar[dateStr] || [];
            const remaining = slotsForDate.reduce((acc: number, s: any) => acc + s.remaining, 0);
            if (slotsForDate.length === 0 || remaining === 0) {
              alert('No available slots on this day');
              return;
            }
            setSelectedDate(dateStr);
            setSelectedSlots(slotsForDate);
            setDialogOpen(true);
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
        />
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
    </div>
  );
} 