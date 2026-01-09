import axios from 'axios';
import { auth } from './stores/auth';
import { browser } from '$app/environment';
import {goto} from "$app/navigation";

// Generate or read a stable device ID (used for token tracking on backend)
let deviceId = null;
if (browser) {
  deviceId = localStorage.getItem('device_id');
  if (!deviceId && typeof crypto !== 'undefined' && crypto.randomUUID) {
    deviceId = crypto.randomUUID();
    localStorage.setItem('device_id', deviceId);
  }
}

const api = axios.create({
  baseURL: '/api/apps',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add auth and refresh tokens to every request
api.interceptors.request.use(
  (config) => {
    const token = auth.getToken();
    const refreshToken = auth.getRefreshToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    if (refreshToken) {
      config.headers['X-Refresh-Token'] = refreshToken;
    }
    if (deviceId) {
      config.headers['X-Device-Id'] = deviceId;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Handle refreshed tokens coming from backend headers (no retry needed)
api.interceptors.response.use(
  (response) => {
    const newAccessToken =
      response.headers['x-token'] || response.headers['X-Token'];
    const newRefreshToken =
      response.headers['x-refresh-token'] || response.headers['X-Refresh-Token'];

    if (newAccessToken || newRefreshToken) {
      auth.setToken(
        newAccessToken || auth.getToken(),
        newRefreshToken || auth.getRefreshToken()
      );
    }

    return response;
  },
  (error) => {
      if (error?.status===401){
          auth.logout();
      }
    return Promise.reject(error);
  }
);

export default api;

