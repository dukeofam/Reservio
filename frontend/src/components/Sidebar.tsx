import { NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

export default function Sidebar() {
  const { t } = useTranslation();
  return (
    <aside className="w-48 bg-gray-100 h-full p-4">
      <ul className="space-y-2">
        <li><NavLink to="/dashboard" className="block p-2 rounded hover:bg-gray-200">{t('dashboard')}</NavLink></li>
        <li><NavLink to="/children" className="block p-2 rounded hover:bg-gray-200">{t('children')}</NavLink></li>
        <li><NavLink to="/reservations" className="block p-2 rounded hover:bg-gray-200">{t('reservations')}</NavLink></li>
        <li><NavLink to="/admin/users" className="block p-2 rounded hover:bg-gray-200">{t('adminUsers')}</NavLink></li>
      </ul>
    </aside>
  );
} 