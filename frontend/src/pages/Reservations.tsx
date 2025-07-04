import { Calendar, dateFnsLocalizer } from 'react-big-calendar';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import { useSlots, useReservations } from '../api/hooks';
import { useAuth } from '../store/useAuth';
import { useTranslation } from 'react-i18next';
import { format, parse, startOfWeek, getDay } from 'date-fns';
import enUS from 'date-fns/locale/en-US';

const locales = { 'en': enUS };
const localizer = dateFnsLocalizer({
  format,
  parse,
  startOfWeek: () => startOfWeek(new Date(), { weekStartsOn: 1 }),
  getDay,
  locales,
});

export default function ReservationsPage() {
  const slots = useSlots();
  const reservations = useReservations();
  const { user } = useAuth();
  const { t } = useTranslation();

  // Map reservations to calendar events
  const events = reservations?.map((r: any) => ({
    id: r.id,
    title: r.status,
    start: new Date(r.date),
    end: new Date(r.date),
    allDay: true
  })) || [];

  return (
    <div className="p-6 bg-amber-50 min-h-screen bg-[url('/wave.svg')] bg-cover bg-center">
      <h1 className="text-3xl font-extrabold text-primary mb-4">{t('reservations')}</h1>
      <div className="bg-white/80 rounded-xl shadow p-4">
        <Calendar
          localizer={localizer}
          events={events}
          startAccessor="start"
          endAccessor="end"
          style={{ height: 500 }}
        />
      </div>
      <div className="mt-8">
        <h2 className="text-xl font-semibold mb-2">Slot Availability</h2>
        <ul className="flex flex-wrap gap-4">
          {slots?.map((slot: any) => (
            <li key={slot.id} className="bg-white/80 rounded shadow px-4 py-2">
              <div className="font-bold">{slot.date}</div>
              <div className="text-sm text-gray-600">Capacity: {slot.capacity}</div>
            </li>
          ))}
        </ul>
      </div>
      {/* Add/Edit/Delete reservation modals would go here */}
    </div>
  );
} 