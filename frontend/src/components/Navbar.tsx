import { Link } from 'react-router-dom';
import { useAuth } from '../store/useAuth';
import LanguageSwitcher from './LanguageSwitcher';
import { useTranslation } from 'react-i18next';
import { UserCircleIcon, Bars3Icon } from '@heroicons/react/24/solid';
import { Menu } from '@headlessui/react';

interface NavbarProps { onBurger: () => void; }

export default function Navbar({ onBurger }: NavbarProps) {
  const { user, logout } = useAuth();
  const { t } = useTranslation();
  return (
    <nav className="sticky top-0 z-30 bg-gradient-to-r from-primary to-indigo-600 text-white px-4 md:px-6 py-3 flex items-center justify-between shadow-lg backdrop-blur">
      <Link to="/dashboard" className="flex items-center gap-2">
        <span className="text-3xl font-extrabold tracking-wide">Reservio</span>
      </Link>
      <div className="flex gap-4 items-center">
        <button onClick={onBurger} className="md:hidden">
          <Bars3Icon className="h-8 w-8" />
        </button>
        <LanguageSwitcher />
        {user && (
          <Menu as="div" className="relative">
            <Menu.Button className="flex items-center gap-2 focus:outline-none">
              {user.profilePicture ? (
                <img src={user.profilePicture} className="h-8 w-8 rounded-full object-cover" />
              ) : (
                <UserCircleIcon className="h-8 w-8" />
              )}
              <span className="hidden sm:block text-sm font-medium">{user.email}</span>
            </Menu.Button>
            <Menu.Items className="absolute right-0 mt-2 w-40 bg-white rounded shadow-lg text-gray-800 py-1">
              <Menu.Item>
                {({ active }: { active: boolean }) => (
                  <Link to="/profile" className={`block px-4 py-2 ${active && 'bg-gray-100'}`}>{t('profile')}</Link>
                )}
              </Menu.Item>
              <Menu.Item>
                {({ active }: { active: boolean }) => (
                  <button onClick={logout} className={`w-full text-left px-4 py-2 ${active && 'bg-gray-100'}`}>{t('logout')}</button>
                )}
              </Menu.Item>
            </Menu.Items>
          </Menu>
        )}
      </div>
    </nav>
  );
} 