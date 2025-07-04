import { createContext, useContext, useState } from 'react';

interface Toast {
  id: number;
  message: string;
  type: 'success' | 'error';
}

const ToastContext = createContext<(msg: string, type?: 'success' | 'error') => void>(() => {});

export function useToast() {
  return useContext(ToastContext);
}

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToasts((prev) => [...prev, { id: Date.now(), message, type }]);
    setTimeout(() => setToasts((prev) => prev.slice(1)), 3000);
  };

  return (
    <ToastContext.Provider value={addToast}>
      {children}
      <div className="fixed bottom-4 right-4 space-y-2">
        {toasts.map((t) => (
          <div key={t.id} className={`px-4 py-2 rounded shadow text-white ${t.type === 'success' ? 'bg-gradient-to-r from-green-400 to-green-600' : 'bg-red-600'}`}>{t.message}</div>
        ))}
      </div>
    </ToastContext.Provider>
  );
} 