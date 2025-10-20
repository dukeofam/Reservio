import { Fragment, useState } from 'react';
import { Dialog, Transition, Listbox } from '@headlessui/react';
import { CheckIcon, ChevronUpDownIcon } from '@heroicons/react/24/solid';
import api from '../api/axios';
import { useToast } from './Toast';
import { useChildren } from '../api/hooks';

interface SlotSummary { id: number; remaining: number; capacity: number; }

interface Props {
  open: boolean;
  date: string; // YYYY-MM-DD
  slots: SlotSummary[];
  onClose: () => void;
  onReserved: () => void;
}

export default function ReservationDialog({ open, date, slots, onClose, onReserved }: Props) {
  const childList = useChildren() as any[];
  const [selectedChild, setSelectedChild] = useState<any>(null);
  const [selectedSlot, setSelectedSlot] = useState<SlotSummary | null>(null);
  const toast = useToast();

  const reserve = async () => {
    if (!selectedChild || !selectedSlot) return;
    try {
      await api.post('/parent/reserve', { slot_id: selectedSlot.id, child_id: selectedChild.id });
      toast('Reservation requested', 'success');
      onReserved();
      onClose();
    } catch (err: any) {
      toast(err.response?.data?.error || 'Error', 'error');
    }
  };

  return (
    <Transition appear show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onClose}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-200" enterFrom="opacity-0" enterTo="opacity-100"
          leave="ease-in duration-150" leaveFrom="opacity-100" leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black/50" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-200" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100"
              leave="ease-in duration-150" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="w-full max-w-md rounded bg-white p-6 shadow-lg">
                <Dialog.Title className="text-lg font-semibold mb-4">Reserve {date}</Dialog.Title>
                {/* Child select */}
                <div className="mb-4">
                  <label className="block text-sm font-medium mb-1">Child</label>
                  <Listbox value={selectedChild} onChange={setSelectedChild}>
                    <div className="relative">
                      <Listbox.Button className="w-full bg-white border px-3 py-2 rounded flex justify-between items-center">
                        <span>{selectedChild ? selectedChild.name : 'Select child'}</span>
                        <ChevronUpDownIcon className="h-5 w-5" />
                      </Listbox.Button>
                      <Transition
                        as={Fragment}
                        enter="transition duration-100" enterFrom="opacity-0" enterTo="opacity-100"
                        leave="transition duration-75" leaveFrom="opacity-100" leaveTo="opacity-0"
                      >
                        <Listbox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded bg-white shadow-lg border text-sm">
                          {childList.map((c: any) => (
                            <Listbox.Option
                              key={c.id}
                              value={c}
                              className={({ active }: { active: boolean }) => `cursor-pointer px-3 py-2 ${active ? 'bg-primary/10' : ''}`}
                            >
                              {({ selected }: { selected: boolean }) => (
                                <span className={selected ? 'font-semibold flex items-center gap-1' : ''}>
                                  {selected && <CheckIcon className="h-4 w-4" />} {c.name}
                                </span>
                              )}
                            </Listbox.Option>
                          ))}
                        </Listbox.Options>
                      </Transition>
                    </div>
                  </Listbox>
                </div>
                {/* Slot select */}
                <div className="mb-4">
                  <label className="block text-sm font-medium mb-1">Available Time Slots</label>
                  <Listbox value={selectedSlot} onChange={setSelectedSlot}>
                    <div className="relative">
                      <Listbox.Button className="w-full bg-white border px-3 py-2 rounded flex justify-between items-center">
                        <span>{selectedSlot ? `${selectedSlot.capacity} spots available (${selectedSlot.remaining} remaining)` : 'Select time slot'}</span>
                        <ChevronUpDownIcon className="h-5 w-5" />
                      </Listbox.Button>
                      <Transition
                        as={Fragment}
                        enter="transition duration-100" enterFrom="opacity-0" enterTo="opacity-100"
                        leave="transition duration-75" leaveFrom="opacity-100" leaveTo="opacity-0"
                      >
                        <Listbox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded bg-white shadow-lg border text-sm">
                          {slots.map((s, index) => (
                            <Listbox.Option
                              key={s.id}
                              value={s}
                              className={({ active }: { active: boolean }) => `cursor-pointer px-3 py-2 ${active ? 'bg-primary/10' : ''}`}
                            >
                              {({ selected }: { selected: boolean }) => (
                                <span className={selected ? 'font-semibold flex items-center gap-1' : ''}>
                                  {selected && <CheckIcon className="h-4 w-4" />} 
                                  {s.capacity} spots available ({s.remaining} remaining)
                                </span>
                              )}
                            </Listbox.Option>
                          ))}
                        </Listbox.Options>
                      </Transition>
                    </div>
                  </Listbox>
                </div>
                <button onClick={reserve} className="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded shadow w-full disabled:opacity-50" disabled={!selectedChild || !selectedSlot}>
                  Reserve
                </button>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
} 