import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  withCredentials: true
});

// Attach CSRF token for mutating requests
api.interceptors.request.use((config) => {
  const unsafe = ['post', 'put', 'delete'];
  if (unsafe.includes(config.method || '')) {
    const token = localStorage.getItem('csrfToken');
    if (token) {
      config.headers['X-CSRF-Token'] = token;
    }
  }
  return config;
});

// Extract new CSRF token from responses
api.interceptors.response.use(
  (res) => {
    const token = res.headers['x-csrf-token'];
    if (token) {
      localStorage.setItem('csrfToken', token);
    }
    return res;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default api; 