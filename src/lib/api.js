import axios from 'axios';
import { auth } from './stores/auth';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add auth token to requests
api.interceptors.request.use(
  (config) => {
    const token = auth.getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      auth.logout();
      if (browser) {
        goto('/login');
      }
    }
    return Promise.reject(error);
  }
);

export default api;

