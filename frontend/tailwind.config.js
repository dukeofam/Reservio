module.exports = {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}'
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#2563eb',
          dark: '#1e40af'
        },
        accent: {
          DEFAULT: '#f59e0b'
        }
      }
    }
  },
  plugins: []
}; 