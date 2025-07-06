import { NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../store/useAuth';
import { XMarkIcon } from '@heroicons/react/24/solid';

interface SidebarProps {
  open?: boolean; // for mobile overlay
  onClose?: () => void;
  className?: string;
}

export default function Sidebar({ open = false, onClose, className = '' }: SidebarProps) {
  const { t } = useTranslation();
  const { user } = useAuth();

  // mobile overlay hidden when !open
  const base = `fixed inset-y-0 left-0 w-64 bg-white shadow-xl z-40 transform transition-transform duration-300 md:relative md:translate-x-0 ${className}`;
  const translate = open ? 'translate-x-0' : '-translate-x-full md:translate-x-0';

  return (
    <aside className={`${base} ${translate}`}>
      {/* Close button mobile */}
      {onClose && (
        <button onClick={onClose} className="md:hidden absolute top-3 right-3 text-gray-600 hover:text-gray-800">
          <XMarkIcon className="h-6 w-6" />
        </button>
      )}
      <ul className="space-y-2 mt-12 md:mt-0">
        <li><NavLink to="/dashboard" className="block p-2 rounded hover:bg-gray-200 flex items-center gap-2">{t('dashboard')}</NavLink></li>
        {user?.role === 'admin' && (
          <li>
            <NavLink to="/children" className="block p-2 rounded hover:bg-gray-200 flex items-center gap-2">
              {t('children')}
            </NavLink>
          </li>
        )}
        <li><NavLink to="/reservations" className="block p-2 rounded hover:bg-gray-200 flex items-center gap-2">{t('reservations')}</NavLink></li>
        {user?.role === 'admin' && (
          <>
            <li><NavLink to="/admin/slots" className="block p-2 rounded hover:bg-gray-200 flex items-center gap-2">Slots</NavLink></li>
            <li><NavLink to="/admin/users" className="block p-2 rounded hover:bg-gray-200 flex items-center gap-2">{t('adminUsers')}</NavLink></li>
          </>
        )}
      </ul>
    </aside>
  );
} 