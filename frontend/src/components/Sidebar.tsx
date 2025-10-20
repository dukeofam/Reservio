import { NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../store/useAuth';
import { 
  XMarkIcon, 
  HomeIcon, 
  UserGroupIcon, 
  CalendarIcon, 
  CogIcon, 
  ClipboardDocumentListIcon,
  UsersIcon,
  ChartBarIcon
} from '@heroicons/react/24/solid';

interface SidebarProps {
  open?: boolean;
  onClose?: () => void;
  className?: string;
}

export default function Sidebar({ open = false, onClose, className = '' }: SidebarProps) {
  const { t } = useTranslation();
  const { user } = useAuth();

  const base = `fixed inset-y-0 left-0 w-72 bg-white border-r border-gray-200 z-40 transform transition-transform duration-300 md:relative md:translate-x-0 ${className}`;
  const translate = open ? 'translate-x-0' : '-translate-x-full md:translate-x-0';

  const navigation = [
    {
      name: t('dashboard'),
      href: '/dashboard',
      icon: HomeIcon,
      current: false
    },
    ...(user?.role === 'admin' ? [{
      name: t('children'),
      href: '/children',
      icon: UserGroupIcon,
      current: false
    }] : []),
    {
      name: t('reservations'),
      href: '/reservations',
      icon: CalendarIcon,
      current: false
    },
    ...(user?.role === 'admin' ? [
      {
        name: 'Manage Slots',
        href: '/admin/slots',
        icon: CogIcon,
        current: false
      },
      {
        name: 'Manage Reservations',
        href: '/admin/reservations',
        icon: ClipboardDocumentListIcon,
        current: false
      },
      {
        name: t('adminUsers'),
        href: '/admin/users',
        icon: UsersIcon,
        current: false
      }
    ] : [])
  ];

  return (
    <aside className={`${base} ${translate}`}>
      {/* Header */}
      <div className="flex items-center justify-between h-16 px-6 border-b border-gray-200">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-sm">R</span>
          </div>
          <span className="text-lg font-semibold text-gray-900">Menu</span>
        </div>
        {onClose && (
          <button 
            onClick={onClose} 
            className="md:hidden p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <XMarkIcon className="h-5 w-5" />
          </button>
        )}
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-4 py-6 space-y-2">
        {navigation.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.name}
              to={item.href}
              className={({ isActive }) =>
                `group flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200 ${
                  isActive
                    ? 'bg-gradient-to-r from-blue-50 to-indigo-50 text-blue-700 border border-blue-200 shadow-sm'
                    : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                }`
              }
            >
              <Icon className="w-5 h-5 flex-shrink-0" />
              <span className="truncate">{item.name}</span>
            </NavLink>
          );
        })}
      </nav>

      {/* User Info */}
      {user && (
        <div className="p-4 border-t border-gray-200">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-full flex items-center justify-center">
              <span className="text-white font-medium text-sm">
                {user.firstName?.charAt(0) || user.email.charAt(0).toUpperCase()}
              </span>
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900 truncate">
                {user.firstName && user.lastName 
                  ? `${user.firstName} ${user.lastName}` 
                  : user.email
                }
              </p>
              <p className="text-xs text-gray-500 capitalize">{user.role}</p>
            </div>
          </div>
        </div>
      )}
    </aside>
  );
} 