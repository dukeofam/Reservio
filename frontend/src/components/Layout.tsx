import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Navbar from './Navbar';

export default function Layout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  return (
    <div className="h-screen flex overflow-hidden bg-amber-50">
      {/* Mobile overlay */}
      <Sidebar open={sidebarOpen} onClose={() => setSidebarOpen(false)} />

      {/* Desktop sidebar */}
      <Sidebar className="hidden md:block" />

      <div className="flex-1 flex flex-col overflow-hidden">
        <Navbar onBurger={() => setSidebarOpen(true)} />
        <main className="flex-1 overflow-y-auto p-4 md:p-6 bg-[url('/wave.svg')] bg-cover bg-center">
          <Outlet />
        </main>
      </div>
    </div>
  );
} 