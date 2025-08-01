import axios from 'axios';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('authToken');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  register: (userData) => api.post('/auth/register', userData),
  login: (credentials) => api.post('/auth/login', credentials),
  getProfile: () => api.get('/profile'),
  updateProfile: (profileData) => api.put('/profile', profileData),
};

// Patient API
export const patientAPI = {
  getProfile: () => api.get('/patients/profile'),
  updateProfile: (profileData) => api.put('/patients/profile', profileData),
  getAll: () => api.get('/patients'),
  getById: (id) => api.get(`/patients/${id}`),
  getImages: (id) => api.get(`/patients/${id}/images`),
};

// Doctor API
export const doctorAPI = {
  getAll: () => api.get('/doctors'),
  getById: (id) => api.get(`/doctors/${id}`),
  getProfile: () => api.get('/doctors/profile'),
  updateProfile: (profileData) => api.put('/doctors/profile', profileData),
};

// Image API
export const imageAPI = {
  upload: (formData) => api.post('/images/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
  detect: (imageId) => api.post('/images/detect', { image_id: imageId }),
  scanWithCNN: (imageId, analysisType = 'basic') => 
    api.post('/images/scan-cnn', { 
      image_id: imageId, 
      analysis_type: analysisType 
    }),
  getAll: () => api.get('/images'),
  getById: (id) => api.get(`/images/${id}`),
  getFile: (id) => api.get(`/images/${id}/file`),
  getCNNHealth: () => api.get('/cnn/health'),
};

// Appointment API
export const appointmentAPI = {
  create: (appointmentData) => api.post('/appointments', appointmentData),
  getAll: () => api.get('/appointments'),
  getById: (id) => api.get(`/appointments/${id}`),
  update: (id, appointmentData) => api.put(`/appointments/${id}`, appointmentData),
  cancel: (id) => api.delete(`/appointments/${id}`),
};

// Analytics API
export const analyticsAPI = {
  getStats: () => api.get('/analytics/stats'),
  getPatientAnalytics: (patientId) => api.get(`/analytics/patient/${patientId}`),
  getDoctorAnalytics: (doctorId) => api.get(`/analytics/doctor/${doctorId}`),
};

// Health check
export const healthAPI = {
  check: () => axios.get('http://localhost:8080/health'),
};

export default api; 