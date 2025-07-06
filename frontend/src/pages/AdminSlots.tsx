import { useState } from 'react';
import { useAdminSlots } from '../api/hooks';
import Modal from '../components/Modal';
import api from '../api/axios';
import { useToast } from '../components/Toast';

export default function AdminSlotsPage() {
  const { slots, refresh } = useAdminSlots();
  const [open, setOpen] = useState(false);
  const [editSlot, setEditSlot] = useState<any>(null);
  const [date, setDate] = useState('');
  const [capacity, setCapacity] = useState(1);
  const toast = useToast();

  const openNew = () => {
    setEditSlot(null);
    setDate('');
    setCapacity(1);
    setOpen(true);
  };

  const openEdit = (s: any) => {
    setEditSlot(s);
    setDate(s.date);
    setCapacity(s.capacity);
    setOpen(true);
  };

  const saveSlot = async () => {
    try {
      if (editSlot) {
        await api.put(`/admin/slots/${editSlot.id}`, { date, capacity });
        toast('Slot updated', 'success');
      } else {
        await api.post('/admin/slots', { date, capacity });
        toast('Slot created', 'success');
      }
      setOpen(false);
      refresh();
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  const deleteSlot = async (id: number) => {
    if (!confirm('Delete this slot?')) return;
    try {
      await api.delete(`/admin/slots/${id}`);
      toast('Slot deleted', 'success');
      refresh();
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  return (
    <div className="p-6 bg-amber-50 min-h-screen bg-[url('/wave.svg')] bg-cover bg-center">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-3xl font-extrabold text-primary">Manage Slots</h1>
        <button onClick={openNew} className="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded shadow-md transition">
          Add Slot
        </button>
      </div>
      <div className="bg-white/80 rounded-xl shadow p-4">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left">
              <th className="py-2">Date</th>
              <th>Capacity</th>
              <th>Remaining</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {slots.map((s: any) => (
              <tr key={s.id} className="border-t">
                <td className="py-2">{s.date}</td>
                <td>{s.capacity}</td>
                <td>{s.remaining ?? '-'}</td>
                <td>
                  <button onClick={() => openEdit(s)} className="text-blue-600 hover:underline mr-2">Edit</button>
                  <button onClick={() => deleteSlot(s.id)} className="text-red-600 hover:underline">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={open} onClose={() => setOpen(false)}>
        <h2 className="text-xl font-semibold mb-4">{editSlot ? 'Edit Slot' : 'New Slot'}</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Date</label>
            <input type="date" value={date} onChange={e => setDate(e.target.value)} className="border px-3 py-2 rounded w-full" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Capacity</label>
            <input type="number" value={capacity} min={1} max={100} onChange={e => setCapacity(Number(e.target.value))} className="border px-3 py-2 rounded w-full" />
          </div>
          <button onClick={saveSlot} className="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded shadow-md">
            Save
          </button>
        </div>
      </Modal>
    </div>
  );
} 