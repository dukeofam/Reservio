import { useAuth } from '../store/useAuth';
import { useDashboardStats, useAnnouncements, useReservations, useChildren, useSlots } from '../api/hooks';
import { useTranslation } from 'react-i18next';
import { 
  UserCircleIcon, 
  CalendarIcon, 
  UserGroupIcon, 
  ClipboardDocumentListIcon, 
  ExclamationTriangleIcon, 
  CheckCircleIcon, 
  ClockIcon,
  PlusIcon,
  ArrowRightIcon,
  ChartBarIcon,
  BellIcon
} from '@heroicons/react/24/solid';
import { Link } from 'react-router-dom';

export default function Dashboard() {
  const { user } = useAuth();
  const stats = useDashboardStats();
  const announcements = useAnnouncements();
  const reservations = useReservations();
  const childrenData = useChildren();
  const children = childrenData.children || [];
  const slots = useSlots();
  const { t } = useTranslation();

  const getUserDisplayName = () => {
    if (user?.firstName && user?.lastName) {
      return `${user.firstName} ${user.lastName}`;
    }
    if (user?.firstName) {
      return user.firstName;
    }
    return user?.email || 'User';
  };

  // Get upcoming slots (next 7 days)
  const upcomingSlots = slots?.slice(0, 7) || [];

  // Get recent reservations
  const recentReservations = reservations?.slice(0, 5) || [];

  // Get pending reservations count
  const pendingReservations = reservations?.filter((r: any) => r.status === 'pending').length || 0;

  return (
    <div className="min-h-screen bg-modern">
      <div className="container-modern py-8 space-y-8">
        {/* Welcome Section */}
        <div className="card card-gradient p-8 animate-fade-in">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
            <div className="flex items-center gap-6">
              {user?.profilePicture ? (
                <img
                  src={user.profilePicture}
                  alt="Profile"
                  className="w-20 h-20 rounded-2xl object-cover border-4 border-white shadow-lg"
                />
              ) : (
                <div className="w-20 h-20 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center shadow-lg">
                  <UserCircleIcon className="w-12 h-12 text-white" />
                </div>
              )}
              <div>
                <h1 className="text-4xl font-bold text-gray-900 mb-2">
                  Welcome back, {getUserDisplayName()}! 👋
                </h1>
                <p className="text-xl text-gray-600">
                  {user?.role === 'admin' ? 'Administrator Dashboard' : 'Parent Dashboard'}
                </p>
                <p className="text-sm text-gray-500 mt-1">
                  {new Date().toLocaleDateString('en-US', {
                    weekday: 'long',
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric'
                  })}
                </p>
              </div>
            </div>
            
            <div className="flex flex-wrap gap-3">
              <Link
                to="/reservations"
                className="btn btn-primary btn-lg group"
              >
                <CalendarIcon className="w-5 h-5 mr-2" />
                Make Reservation
                <ArrowRightIcon className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
              </Link>
              
              {user?.role === 'parent' && (
                <Link
                  to="/children"
                  className="btn btn-secondary btn-lg group"
                >
                  <PlusIcon className="w-5 h-5 mr-2" />
                  Add Child
                  <ArrowRightIcon className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
                </Link>
              )}
            </div>
          </div>
        </div>

        {/* Quick Actions Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Link
            to="/reservations"
            className="card card-hover p-6 group animate-slide-up"
          >
            <div className="flex items-center gap-4">
              <div className="w-14 h-14 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                <CalendarIcon className="w-7 h-7 text-white" />
              </div>
              <div>
                <h3 className="font-semibold text-gray-900 text-lg">Reservations</h3>
                <p className="text-gray-600 text-sm">Book time slots</p>
              </div>
            </div>
          </Link>

          <Link
            to="/children"
            className="card card-hover p-6 group animate-slide-up"
            style={{ animationDelay: '0.1s' }}
          >
            <div className="flex items-center gap-4">
              <div className="w-14 h-14 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                <UserGroupIcon className="w-7 h-7 text-white" />
              </div>
              <div>
                <h3 className="font-semibold text-gray-900 text-lg">Children</h3>
                <p className="text-gray-600 text-sm">Manage profiles</p>
              </div>
            </div>
          </Link>

          {user?.role === 'admin' && (
            <>
              <Link
                to="/admin/slots"
                className="card card-hover p-6 group animate-slide-up"
                style={{ animationDelay: '0.2s' }}
              >
                <div className="flex items-center gap-4">
                  <div className="w-14 h-14 bg-gradient-to-br from-purple-500 to-indigo-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                    <CogIcon className="w-7 h-7 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 text-lg">Slots</h3>
                    <p className="text-gray-600 text-sm">Manage availability</p>
                  </div>
                </div>
              </Link>

              <Link
                to="/admin/reservations"
                className="card card-hover p-6 group animate-slide-up"
                style={{ animationDelay: '0.3s' }}
              >
                <div className="flex items-center gap-4">
                  <div className="w-14 h-14 bg-gradient-to-br from-orange-500 to-red-500 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                    <ClipboardDocumentListIcon className="w-7 h-7 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 text-lg">Manage</h3>
                    <p className="text-gray-600 text-sm">Review reservations</p>
                  </div>
                </div>
              </Link>
            </>
          )}
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {user?.role === 'admin' ? (
            <>
              <div className="card p-6 animate-scale-in">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600 mb-1">Total Children</p>
                    <p className="text-3xl font-bold text-blue-600">{stats?.total_children ?? '0'}</p>
                  </div>
                  <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
                    <UserGroupIcon className="w-6 h-6 text-blue-600" />
                  </div>
                </div>
              </div>
              
              <div className="card p-6 animate-scale-in" style={{ animationDelay: '0.1s' }}>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600 mb-1">Total Reservations</p>
                    <p className="text-3xl font-bold text-green-600">{stats?.total_reservations ?? '0'}</p>
                  </div>
                  <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
                    <CalendarIcon className="w-6 h-6 text-green-600" />
                  </div>
                </div>
              </div>
              
              <div className="card p-6 animate-scale-in" style={{ animationDelay: '0.2s' }}>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600 mb-1">Available Slots</p>
                    <p className="text-3xl font-bold text-purple-600">{stats?.open_slots ?? '0'}</p>
                  </div>
                  <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center">
                    <ChartBarIcon className="w-6 h-6 text-purple-600" />
                  </div>
                </div>
              </div>
            </>
          ) : (
            <>
              <div className="card p-6 animate-scale-in">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600 mb-1">My Children</p>
                    <p className="text-3xl font-bold text-blue-600">{children?.length ?? '0'}</p>
                  </div>
                  <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
                    <UserGroupIcon className="w-6 h-6 text-blue-600" />
                  </div>
                </div>
              </div>
              
              <div className="card p-6 animate-scale-in" style={{ animationDelay: '0.1s' }}>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600 mb-1">My Reservations</p>
                    <p className="text-3xl font-bold text-green-600">{reservations?.length ?? '0'}</p>
                  </div>
                  <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
                    <CalendarIcon className="w-6 h-6 text-green-600" />
                  </div>
                </div>
              </div>
              
              <div className="card p-6 animate-scale-in" style={{ animationDelay: '0.2s' }}>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600 mb-1">Pending Approval</p>
                    <p className="text-3xl font-bold text-orange-600">{pendingReservations}</p>
                  </div>
                  <div className="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center">
                    <ClockIcon className="w-6 h-6 text-orange-600" />
                  </div>
                </div>
              </div>
            </>
          )}
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Upcoming Slots */}
          <div className="card p-6 animate-fade-in">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-gray-900">Upcoming Available Slots</h2>
              <Link 
                to="/reservations" 
                className="text-blue-600 hover:text-blue-800 text-sm font-medium flex items-center gap-1 group"
              >
                View All
                <ArrowRightIcon className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </Link>
            </div>
            
            <div className="space-y-3">
              {upcomingSlots.length > 0 ? (
                upcomingSlots.map((slot: any, index: number) => (
                  <div 
                    key={slot.id} 
                    className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                    style={{ animationDelay: `${index * 0.1}s` }}
                  >
                    <div>
                      <div className="font-semibold text-gray-900">
                        {new Date(slot.date).toLocaleDateString('en-US', {
                          weekday: 'short',
                          month: 'short',
                          day: 'numeric'
                        })}
                      </div>
                      <div className="text-sm text-gray-600">
                        {slot.available_slots} of {slot.capacity} spots available
                      </div>
                    </div>
                    <div className={`badge ${
                      slot.available_slots > 0 ? 'badge-success' : 'badge-danger'
                    }`}>
                      {slot.available_slots > 0 ? 'Available' : 'Full'}
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center py-8 text-gray-500">
                  <CalendarIcon className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                  <p className="font-medium">No upcoming slots available</p>
                  <p className="text-sm">Check back later for new availability</p>
                </div>
              )}
            </div>
          </div>

          {/* Recent Reservations */}
          <div className="card p-6 animate-fade-in">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-gray-900">Recent Reservations</h2>
              <Link 
                to="/reservations" 
                className="text-blue-600 hover:text-blue-800 text-sm font-medium flex items-center gap-1 group"
              >
                View All
                <ArrowRightIcon className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </Link>
            </div>
            
            <div className="space-y-3">
              {recentReservations.length > 0 ? (
                recentReservations.map((reservation: any, index: number) => (
                  <div 
                    key={reservation.id} 
                    className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                    style={{ animationDelay: `${index * 0.1}s` }}
                  >
                    <div>
                      <div className="font-semibold text-gray-900">
                        {new Date(reservation.date).toLocaleDateString('en-US', {
                          weekday: 'short',
                          month: 'short',
                          day: 'numeric'
                        })}
                      </div>
                      <div className="text-sm text-gray-600">
                        Child: {reservation.child?.name || 'Unknown'}
                      </div>
                    </div>
                    <div className={`badge flex items-center gap-1 ${
                      reservation.status === 'approved' ? 'badge-success' :
                      reservation.status === 'pending' ? 'badge-warning' :
                      'badge-danger'
                    }`}>
                      {reservation.status === 'approved' && <CheckCircleIcon className="w-3 h-3" />}
                      {reservation.status === 'pending' && <ClockIcon className="w-3 h-3" />}
                      {reservation.status === 'rejected' && <ExclamationTriangleIcon className="w-3 h-3" />}
                      {reservation.status}
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center py-8 text-gray-500">
                  <CalendarIcon className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                  <p className="font-medium">No reservations yet</p>
                  <Link 
                    to="/reservations" 
                    className="text-blue-600 hover:text-blue-800 text-sm font-medium inline-flex items-center gap-1 mt-2 group"
                  >
                    Make your first reservation
                    <ArrowRightIcon className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Announcements */}
        <div className="card p-6 animate-fade-in">
          <div className="flex items-center gap-3 mb-6">
            <BellIcon className="w-6 h-6 text-blue-600" />
            <h2 className="text-xl font-bold text-gray-900">Latest Announcements</h2>
          </div>
          
          <div className="space-y-4">
            {announcements && announcements.length > 0 ? (
              announcements.map((announcement: any, index: number) => (
                <div 
                  key={announcement.id} 
                  className="p-4 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg border border-blue-200 hover:shadow-md transition-shadow"
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="font-semibold text-gray-900 mb-2">{announcement.title}</h3>
                      <p className="text-gray-700 text-sm mb-3">{announcement.content}</p>
                      <div className="text-xs text-gray-500">
                        {new Date(announcement.createdAt).toLocaleDateString('en-US', {
                          year: 'numeric',
                          month: 'long',
                          day: 'numeric'
                        })} • {announcement.author?.firstName || 'Admin'}
                      </div>
                    </div>
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center py-8 text-gray-500">
                <BellIcon className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                <p className="font-medium">No announcements yet</p>
                <p className="text-sm">Check back later for updates</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}