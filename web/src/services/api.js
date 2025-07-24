import axios from 'axios'

// Create axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  config => {
    console.log(`🔄 API Request: ${config.method?.toUpperCase()} ${config.url}`)
    return config
  },
  error => {
    console.error('❌ API Request Error:', error)
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  response => {
    console.log(`✅ API Response: ${response.status} ${response.config.url}`)
    return response
  },
  error => {
    console.error('❌ API Response Error:', error)
    
    // Handle common error scenarios
    if (error.response) {
      // Server responded with error status
      const { status, data } = error.response
      console.error(`API Error ${status}:`, data)
      
      switch (status) {
        case 400:
          throw new Error(data?.message || 'Bad Request')
        case 401:
          throw new Error('Unauthorized')
        case 403:
          throw new Error('Forbidden')
        case 404:
          throw new Error('Resource not found')
        case 500:
          throw new Error('Internal Server Error')
        default:
          throw new Error(data?.message || `Server Error: ${status}`)
      }
    } else if (error.request) {
      // Network error
      throw new Error('Network Error: Unable to connect to server')
    } else {
      // Other error
      throw new Error(error.message || 'Unknown Error')
    }
  }
)

export default api