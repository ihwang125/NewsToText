import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || '/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  register: (userData) => api.post('/auth/register', userData),
  login: (credentials) => api.post('/auth/login', credentials),
  logout: () => api.post('/auth/logout'),
};

export const alertsAPI = {
  getAlerts: () => api.get('/alerts'),
  createAlert: (alertData) => api.post('/alerts', alertData),
  updateAlert: (id, alertData) => api.put(`/alerts/${id}`, alertData),
  deleteAlert: (id) => api.delete(`/alerts/${id}`),
  getHistory: () => api.get('/alerts/history'),
  testAlert: (alertId) => api.post('/alerts/test', { alert_id: alertId }),
};

export default api;