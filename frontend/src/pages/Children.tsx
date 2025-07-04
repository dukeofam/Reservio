import { useChildren } from '../api/hooks';
import { useAuth } from '../store/useAuth';
import { Navigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

export default function ChildrenPage() {
  const { user } = useAuth();
  const children = useChildren();
  const { t } = useTranslation();

  if (user?.role !== 'admin') {
    return <Navigate to="/dashboard" />;
  }

  return (
    <div className="p-6 bg-amber-50 min-h-screen bg-[url('/wave.svg')] bg-cover bg-center">
      <h1 className="text-3xl font-extrabold text-primary mb-4">{t('children')}</h1>
      <div className="bg-white/80 rounded-xl shadow p-4">
        <table className="w-full">
          <thead>
            <tr className="text-left">
              <th className="py-2">Name</th>
              <th className="py-2">Date of Birth</th>
              <th className="py-2">Parent</th>
              <th className="py-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            {children?.map((child: any) => (
              <tr key={child.id} className="border-t">
                <td>{child.name}</td>
                <td>{child.birthdate || '-'}</td>
                <td>{child.parentEmail || '-'}</td>
                <td>
                  {/* Edit/Delete buttons here */}
                  <button className="text-blue-600 hover:underline mr-2">Edit</button>
                  <button className="text-red-600 hover:underline">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {/* Add child form/modal would go here */}
      </div>
    </div>
  );
} 