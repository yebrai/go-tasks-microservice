import api from './api'

export const taskService = {
  /**
   * Get system health status
   */
  async getHealth() {
    try {
      const response = await api.get('/health')
      return response.data
    } catch (error) {
      console.error('Health check failed:', error)
      return {
        status: 'error',
        database: 'disconnected',
        rabbitmq: 'disconnected'
      }
    }
  },

  /**
   * Get all tasks
   */
  async getTasks() {
    try {
      const response = await api.get('/v1/tasks')
      return response.data || []
    } catch (error) {
      console.error('Failed to fetch tasks:', error)
      // Return mock data for demo purposes when API is not available
      return [
        {
          id: '1',
          title: 'Setup development environment',
          description: 'Configure Go, Docker, and local development tools',
          status: 'completed',
          created_at: '2024-01-15T10:00:00Z',
          due_date: '2024-01-20'
        },
        {
          id: '2',
          title: 'Implement task creation API',
          description: 'Create REST endpoint for task creation with validation',
          status: 'in-progress',
          created_at: '2024-01-16T14:30:00Z',
          due_date: '2024-01-25'
        },
        {
          id: '3',
          title: 'Add RabbitMQ integration',
          description: 'Integrate event publishing for task lifecycle events',
          status: 'pending',
          created_at: '2024-01-17T09:15:00Z',
          due_date: '2024-01-30'
        },
        {
          id: '4',
          title: 'Deploy to GCP',
          description: 'Setup Cloud Run deployment with proper configuration',
          status: 'pending',
          created_at: '2024-01-18T11:45:00Z',
          due_date: '2024-02-05'
        }
      ]
    }
  },

  /**
   * Get a specific task by ID
   */
  async getTask(id) {
    try {
      const response = await api.get(`/v1/tasks/${id}`)
      return response.data
    } catch (error) {
      console.error(`Failed to fetch task ${id}:`, error)
      throw error
    }
  },

  /**
   * Create a new task
   */
  async createTask(taskData) {
    try {
      const response = await api.post('/v1/tasks', taskData)
      console.log('✅ Task created successfully:', response.data)
      return response.data
    } catch (error) {
      console.error('Failed to create task:', error)
      throw error
    }
  },

  /**
   * Update an existing task
   */
  async updateTask(id, updates) {
    try {
      const response = await api.put(`/v1/tasks/${id}`, updates)
      console.log('✅ Task updated successfully:', response.data)
      return response.data
    } catch (error) {
      console.error(`Failed to update task ${id}:`, error)
      throw error
    }
  },

  /**
   * Delete a task
   */
  async deleteTask(id) {
    try {
      const response = await api.delete(`/v1/tasks/${id}`)
      console.log('✅ Task deleted successfully')
      return response.data
    } catch (error) {
      console.error(`Failed to delete task ${id}:`, error)
      throw error
    }
  },

  /**
   * Get task statistics
   */
  async getTaskStats() {
    try {
      const response = await api.get('/v1/tasks/stats')
      return response.data
    } catch (error) {
      console.error('Failed to fetch task stats:', error)
      // Return mock stats when API is not available
      return {
        total: 12,
        completed: 8,
        pending: 3,
        in_progress: 1
      }
    }
  }
}

export default taskService