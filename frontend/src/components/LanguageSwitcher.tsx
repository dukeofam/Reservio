import { useTranslation } from 'react-i18next';

export default function LanguageSwitcher() {
  const { i18n, t } = useTranslation();

  const change = (e: React.ChangeEvent<HTMLSelectElement>) => {
    i18n.changeLanguage(e.target.value);
  };

  return (
    <select
      onChange={change}
      value={i18n.language}
      className="bg-transparent text-white border border-white rounded px-2 py-1 text-sm focus:outline-none"
    >
      <option value="en">EN</option>
      <option value="cs">CZ</option>
    </select>
  );
} 