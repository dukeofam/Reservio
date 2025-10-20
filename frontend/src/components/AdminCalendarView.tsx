import { Fragment, useState, useEffect } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { XMarkIcon, UserGroupIcon, CalendarIcon } from '@heroicons/react/24/solid';
import api from '../api/axios';

interface Props {
  open: boolean;
  date: string;
  onClose: () => void;
}

export default function AdminCalendarView({ open, date, onClose }: Props) {
  const [reservations, setReservations] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (open && date) {
      fetchReservationsForDate();
    }
  }, [open, date]);

  const fetchReservationsForDate = async () => {
    setLoading(true);
    try {
      const response = await api.get(`/admin/reservations?date=${date}`);
      setReservations(response.data.data || []);
    } catch (error) {
      console.error('Failed to fetch reservations:', error);
      setReservations([]);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'approved':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'rejected':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <Transition appear show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onClose}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-200"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-150"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black/50" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-200"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-150"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="w-full max-w-2xl rounded-xl bg-white p-6 shadow-xl">
                <div className="flex items-center justify-between mb-6">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                      <CalendarIcon className="w-6 h-6 text-blue-600" />
                    </div>
                    <div>
                      <Dialog.Title className="text-xl font-bold text-gray-900">
                        Reservations for {formatDate(date)}
                      </Dialog.Title>
                      <p className="text-sm text-gray-600">
                        {reservations.length} total reservations
                      </p>
                    </div>
                  </div>
                  <button
                    onClick={onClose}
                    className="text-gray-400 hover:text-gray-600 transition-colors"
                  >
                    <XMarkIcon className="w-6 h-6" />
                  </button>
                </div>

                {loading ? (
                  <div className="flex items-center justify-center py-12">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                  </div>
                ) : reservations.length > 0 ? (
                  <div className="space-y-3">
                    {reservations.map((reservation) => (
                      <div
                        key={reservation.id}
                        className="flex items-center justify-between p-4 bg-gray-50 rounded-lg border border-gray-200"
                      >
                        <div className="flex items-center gap-4">
                          <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                            <UserGroupIcon className="w-5 h-5 text-blue-600" />
                          </div>
                          <div>
                            <div className="font-semibold text-gray-900">
                              {reservation.child?.name || 'Unknown Child'}
                            </div>
                            <div className="text-sm text-gray-600">
                              Parent: {reservation.child?.parent?.firstName} {reservation.child?.parent?.lastName}
                            </div>
                            <div className="text-sm text-gray-600">
                              Email: {reservation.child?.parent?.email}
                            </div>
                          </div>
                        </div>
                        <div className="text-right">
                          <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(reservation.status)}`}>
                            {reservation.status}
                          </span>
                          <div className="text-sm text-gray-600 mt-1">
                            Slot: {reservation.slot?.capacity} spots
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <UserGroupIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">No reservations</h3>
                    <p className="text-gray-600">No children are registered for this day.</p>
                  </div>
                )}
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
