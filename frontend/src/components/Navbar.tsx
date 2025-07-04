import { Link } from 'react-router-dom';
import { useAuth } from '../store/useAuth';
import LanguageSwitcher from './LanguageSwitcher';
import { useTranslation } from 'react-i18next';
import { UserCircleIcon } from '@heroicons/react/24/solid';

export default function Navbar() {
  const { user, logout } = useAuth();
  const { t } = useTranslation();
  return (
    <nav className="bg-primary text-white px-6 py-4 flex items-center justify-between shadow-lg">
      <Link to="/dashboard" className="flex items-center gap-2">
        <span className="text-3xl font-extrabold tracking-wide">Reservio</span>
      </Link>
      <div className="flex gap-6 items-center">
        <LanguageSwitcher />
        {user && (
          <>
            <span>{user.email}</span>
            <Link to="/profile"><UserCircleIcon className="h-8 w-8 text-white hover:text-amber-200" /></Link>
            <button onClick={logout} className="bg-white/20 hover:bg-white/30 transition px-3 py-1 rounded">{t('logout')}</button>
          </>
        )}
      </div>
    </nav>
  );
} 