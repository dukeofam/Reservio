import Navbar from './components/Navbar';
import Sidebar from './components/Sidebar';
import AppRoutes from './routes/AppRoutes';
import { ToastProvider } from './components/Toast';
import { useAuth } from './store/useAuth';

export default function App() {
  const { user } = useAuth();

  return (
    <ToastProvider>
      {user ? (
        <div className="flex flex-col min-h-screen">
          <Navbar />
          <div className="flex flex-1">
            <Sidebar />
            <main className="flex-1 bg-white">
              <AppRoutes />
            </main>
          </div>
        </div>
      ) : (
        <AppRoutes />
      )}
    </ToastProvider>
  );
} 