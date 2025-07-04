import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

const resources = {
  en: {
    translation: {
      welcome: 'Welcome back',
      signInSubtitle: 'Sign in to manage your reservations',
      email: 'Email',
      password: 'Password',
      login: 'Login',
      dontHaveAccount: "Don't have an account?",
      register: 'Register',
      createAccount: 'Create account',
      firstName: 'First name',
      lastName: 'Last name',
      childName: 'Child name',
      phone: 'Phone number',
      passwordAgain: 'Password again',
      alreadyAccount: 'Already have an account?',
      dashboard: 'Dashboard',
      children: 'Children',
      reservations: 'Reservations',
      adminUsers: 'Admin Users',
      logout: 'Logout',
      language: 'Language'
    }
  },
  cs: {
    translation: {
      welcome: 'Vítejte zpět',
      signInSubtitle: 'Přihlaste se pro správu rezervací',
      email: 'E-mail',
      password: 'Heslo',
      login: 'Přihlásit se',
      dontHaveAccount: 'Nemáte účet?',
      register: 'Registrace',
      createAccount: 'Vytvořit účet',
      firstName: 'Jméno',
      lastName: 'Příjmení',
      childName: 'Jméno dítěte',
      phone: 'Telefon',
      passwordAgain: 'Heslo znovu',
      alreadyAccount: 'Už máte účet?',
      dashboard: 'Přehled',
      children: 'Děti',
      reservations: 'Rezervace',
      adminUsers: 'Správa uživatelů',
      logout: 'Odhlásit',
      language: 'Jazyk'
    }
  }
};

i18n.use(initReactI18next).init({
  resources,
  lng: 'en',
  fallbackLng: 'en',
  interpolation: {
    escapeValue: false
  }
});

export default i18n; 