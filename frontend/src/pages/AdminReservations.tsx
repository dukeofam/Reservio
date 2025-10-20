import { useState, useEffect } from 'react';
import { useAuth } from '../store/useAuth';
import { useTranslation } from 'react-i18next';
import api from '../api/axios';
import { useToast } from '../components/Toast';
import { CheckCircleIcon, XCircleIcon, ClockIcon, UserGroupIcon } from '@heroicons/react/24/solid';

export default function AdminReservationsPage() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const [reservations, setReservations] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState('all');
  const toast = useToast();

  useEffect(() => {
    fetchReservations();
  }, [statusFilter]);

  const fetchReservations = async () => {
    try {
      setLoading(true);
      const response = await api.get('/admin/reservations');
      setReservations(response.data.data || []);
    } catch (error) {
      console.error('Failed to fetch reservations:', error);
      toast('Failed to load reservations', 'error');
    } finally {
      setLoading(false);
    }
  };

  if (user?.role !== 'admin') {
    return (
      <div className="p-6 bg-gradient-to-br from-red-50 to-pink-50 min-h-screen flex items-center justify-center">
        <div className="text-center">
          <XCircleIcon className="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Access Denied</h2>
          <p className="text-gray-600">This page is only accessible to administrators.</p>
        </div>
      </div>
    );
  }

  const filteredReservations = reservations.filter((r: any) => 
    statusFilter === 'all' || r.status === statusFilter
  );

  const approveReservation = async (id: number) => {
    try {
      await api.put(`/admin/approve/${id}`);
      toast('Reservation approved', 'success');
      fetchReservations();
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  const rejectReservation = async (id: number) => {
    try {
      await api.put(`/admin/reject/${id}`);
      toast('Reservation rejected', 'success');
      fetchReservations();
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'approved': return 'bg-green-100 text-green-800 border-green-200';
      case 'pending': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'rejected': return 'bg-red-100 text-red-800 border-red-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  return (
    <div className="p-6 bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 min-h-screen">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-bold text-gray-900 mb-2">Reservation Management</h1>
            <p className="text-gray-600">Manage and approve kindergarten reservations</p>
          </div>
          <div className="flex items-center gap-4">
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-2">
              <select 
                value={statusFilter} 
                onChange={e => setStatusFilter(e.target.value)}
                className="border-0 bg-transparent text-sm font-medium text-gray-700 focus:outline-none focus:ring-0"
              >
                <option value="all">All Reservations</option>
                <option value="pending">Pending Review</option>
                <option value="approved">Approved</option>
                <option value="rejected">Rejected</option>
              </select>
            </div>
            <button
              onClick={fetchReservations}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg shadow-sm transition-colors flex items-center gap-2"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refresh
            </button>
          </div>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Reservations</p>
              <p className="text-3xl font-bold text-gray-900">{reservations.length}</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
          </div>
        </div>
        
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Pending Review</p>
              <p className="text-3xl font-bold text-yellow-600">
                {reservations.filter(r => r.status === 'pending').length}
              </p>
            </div>
            <div className="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center">
              <ClockIcon className="w-6 h-6 text-yellow-600" />
            </div>
          </div>
        </div>
        
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Approved</p>
              <p className="text-3xl font-bold text-green-600">
                {reservations.filter(r => r.status === 'approved').length}
              </p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
              <CheckCircleIcon className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>
        
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Rejected</p>
              <p className="text-3xl font-bold text-red-600">
                {reservations.filter(r => r.status === 'rejected').length}
              </p>
            </div>
            <div className="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center">
              <XCircleIcon className="w-6 h-6 text-red-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Reservations List */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        ) : filteredReservations.length > 0 ? (
          <div className="divide-y divide-gray-200">
            {filteredReservations.map((reservation: any) => (
              <div key={reservation.id} className="p-6 hover:bg-gray-50 transition-colors">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white font-bold text-lg">
                      {reservation.child?.name?.charAt(0) || '?'}
                    </div>
                    <div>
                      <h3 className="font-semibold text-gray-900 text-lg">
                        {reservation.child?.name || 'Unknown Child'}
                      </h3>
                      <p className="text-gray-600">
                        Parent: {reservation.child?.parent?.firstName} {reservation.child?.parent?.lastName}
                      </p>
                      <p className="text-sm text-gray-500">
                        Email: {reservation.child?.parent?.email}
                      </p>
                      <p className="text-sm text-gray-500">
                        Date: {new Date(reservation.date).toLocaleDateString('en-US', { 
                          weekday: 'long', 
                          year: 'numeric', 
                          month: 'long', 
                          day: 'numeric' 
                        })}
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium border ${getStatusColor(reservation.status)}`}>
                        {reservation.status === 'pending' && <ClockIcon className="w-4 h-4 mr-1" />}
                        {reservation.status === 'approved' && <CheckCircleIcon className="w-4 h-4 mr-1" />}
                        {reservation.status === 'rejected' && <XCircleIcon className="w-4 h-4 mr-1" />}
                        {reservation.status}
                      </span>
                      <p className="text-sm text-gray-500 mt-1">
                        Slot: {reservation.slot?.capacity} spots
                      </p>
                    </div>
                    
                    <div className="flex gap-2">
                      {reservation.status === 'pending' && (
                        <>
                          <button
                            onClick={() => approveReservation(reservation.id)}
                            className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
                          >
                            <CheckCircleIcon className="w-4 h-4" />
                            Approve
                          </button>
                          <button
                            onClick={() => rejectReservation(reservation.id)}
                            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
                          >
                            <XCircleIcon className="w-4 h-4" />
                            Reject
                          </button>
                        </>
                      )}
                      {reservation.status === 'approved' && (
                        <button
                          onClick={() => rejectReservation(reservation.id)}
                          className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
                        >
                          <XCircleIcon className="w-4 h-4" />
                          Reject
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <UserGroupIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No reservations found</h3>
            <p className="text-gray-600">No reservations match your current filter criteria.</p>
          </div>
        )}
      </div>
    </div>
  );
}