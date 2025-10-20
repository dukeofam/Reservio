import { Link } from 'react-router-dom';
import { useAuth } from '../store/useAuth';
import LanguageSwitcher from './LanguageSwitcher';
import { useTranslation } from 'react-i18next';
import { UserCircleIcon, Bars3Icon, BellIcon, Cog6ToothIcon, XMarkIcon } from '@heroicons/react/24/solid';
import { Menu, Transition } from '@headlessui/react';
import { Fragment, useState, useEffect, useRef } from 'react';

interface NavbarProps { onBurger: () => void; }

export default function Navbar({ onBurger }: NavbarProps) {
  const { user, logout } = useAuth();
  const { t } = useTranslation();
  const [notificationsOpen, setNotificationsOpen] = useState(false);
  const notificationsRef = useRef<HTMLDivElement>(null);
  
  // Close notifications when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (notificationsRef.current && !notificationsRef.current.contains(event.target as Node)) {
        setNotificationsOpen(false);
      }
    };

    if (notificationsOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [notificationsOpen]);
  
  // Sample notifications data - in a real app, this would come from an API
  const notifications = [
    {
      id: 1,
      title: "New Reservation Request",
      message: "Emma Johnson has requested a slot for tomorrow",
      time: "2 minutes ago",
      unread: true
    },
    {
      id: 2,
      title: "Slot Capacity Updated",
      message: "Monday slots now have 5 available spots",
      time: "1 hour ago",
      unread: true
    },
    {
      id: 3,
      title: "System Maintenance",
      message: "Scheduled maintenance completed successfully",
      time: "3 hours ago",
      unread: false
    }
  ];
  
  return (
    <nav className="sticky top-0 z-30 bg-white/80 backdrop-blur-md border-b border-gray-200 shadow-sm">
      <div className="container-modern">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/dashboard" className="flex items-center gap-3 group">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg group-hover:shadow-xl transition-all duration-200">
              <span className="text-white font-bold text-lg">R</span>
            </div>
            <div>
              <span className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
                Reservio
              </span>
              <p className="text-xs text-gray-500 -mt-1">Kindergarten Management</p>
            </div>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-6">
            <LanguageSwitcher />
            
            {/* Notifications */}
            <div className="relative" ref={notificationsRef}>
              <button 
                onClick={() => setNotificationsOpen(!notificationsOpen)}
                className="relative p-2 text-gray-400 hover:text-gray-600 transition-colors"
              >
                <BellIcon className="w-5 h-5" />
                {notifications.some(n => n.unread) && (
                  <span className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full"></span>
                )}
              </button>
              
              {/* Notifications Dropdown */}
              <Transition
                show={notificationsOpen}
                as={Fragment}
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <div className="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-modern-lg border border-gray-200 z-50">
                  <div className="p-4 border-b border-gray-100">
                    <div className="flex items-center justify-between">
                      <h3 className="text-lg font-semibold text-gray-900">Notifications</h3>
                      <button
                        onClick={() => setNotificationsOpen(false)}
                        className="p-1 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
                      >
                        <XMarkIcon className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                  
                  <div className="max-h-96 overflow-y-auto">
                    {notifications.length > 0 ? (
                      notifications.map((notification) => (
                        <div
                          key={notification.id}
                          className={`p-4 border-b border-gray-100 hover:bg-gray-50 transition-colors ${
                            notification.unread ? 'bg-blue-50' : ''
                          }`}
                        >
                          <div className="flex items-start gap-3">
                            <div className={`w-2 h-2 rounded-full mt-2 ${
                              notification.unread ? 'bg-blue-500' : 'bg-gray-300'
                            }`}></div>
                            <div className="flex-1 min-w-0">
                              <h4 className="text-sm font-medium text-gray-900 mb-1">
                                {notification.title}
                              </h4>
                              <p className="text-sm text-gray-600 mb-2">
                                {notification.message}
                              </p>
                              <p className="text-xs text-gray-500">
                                {notification.time}
                              </p>
                            </div>
                          </div>
                        </div>
                      ))
                    ) : (
                      <div className="p-8 text-center text-gray-500">
                        <BellIcon className="w-8 h-8 mx-auto mb-2 text-gray-300" />
                        <p>No notifications</p>
                      </div>
                    )}
                  </div>
                  
                  {notifications.length > 0 && (
                    <div className="p-3 border-t border-gray-100">
                      <button className="w-full text-sm text-blue-600 hover:text-blue-800 font-medium">
                        Mark all as read
                      </button>
                    </div>
                  )}
                </div>
              </Transition>
            </div>

            {/* User Menu */}
            {user && (
              <Menu as="div" className="relative">
                <Menu.Button className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition-colors">
                  {user.profilePicture ? (
                    <img 
                      src={user.profilePicture} 
                      className="h-8 w-8 rounded-full object-cover border-2 border-gray-200" 
                      alt="Profile"
                    />
                  ) : (
                    <div className="h-8 w-8 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-full flex items-center justify-center">
                      <UserCircleIcon className="w-5 h-5 text-white" />
                    </div>
                  )}
                  <div className="text-left hidden lg:block">
                    <p className="text-sm font-medium text-gray-900">
                      {user.firstName && user.lastName 
                        ? `${user.firstName} ${user.lastName}` 
                        : user.email
                      }
                    </p>
                    <p className="text-xs text-gray-500 capitalize">{user.role}</p>
                  </div>
                </Menu.Button>
                
                <Transition
                  as={Fragment}
                  enter="transition ease-out duration-100"
                  enterFrom="transform opacity-0 scale-95"
                  enterTo="transform opacity-100 scale-100"
                  leave="transition ease-in duration-75"
                  leaveFrom="transform opacity-100 scale-100"
                  leaveTo="transform opacity-0 scale-95"
                >
                  <Menu.Items className="absolute right-0 mt-2 w-56 bg-white rounded-xl shadow-modern-lg border border-gray-200 py-2 z-50">
                    <div className="px-4 py-3 border-b border-gray-100">
                      <p className="text-sm font-medium text-gray-900">
                        {user.firstName && user.lastName 
                          ? `${user.firstName} ${user.lastName}` 
                          : 'User'
                        }
                      </p>
                      <p className="text-sm text-gray-500">{user.email}</p>
                    </div>
                    
                    <Menu.Item>
                      {({ active }) => (
                        <Link 
                          to="/profile" 
                          className={`flex items-center gap-3 px-4 py-2 text-sm transition-colors ${
                            active ? 'bg-gray-50 text-gray-900' : 'text-gray-700'
                          }`}
                        >
                          <UserCircleIcon className="w-4 h-4" />
                          {t('profile')}
                        </Link>
                      )}
                    </Menu.Item>
                    
                    <Menu.Item>
                      {({ active }) => (
                        <button 
                          onClick={logout} 
                          className={`w-full flex items-center gap-3 px-4 py-2 text-sm text-left transition-colors ${
                            active ? 'bg-gray-50 text-gray-900' : 'text-gray-700'
                          }`}
                        >
                          <Cog6ToothIcon className="w-4 h-4" />
                          {t('logout')}
                        </button>
                      )}
                    </Menu.Item>
                  </Menu.Items>
                </Transition>
              </Menu>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="flex items-center gap-3 md:hidden">
            <LanguageSwitcher />
            <button 
              onClick={onBurger}
              className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <Bars3Icon className="h-6 w-6" />
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
} 