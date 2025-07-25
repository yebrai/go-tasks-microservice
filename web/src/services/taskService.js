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
      return {
        status: 'error',
        database: 'disconnected',
        rabbitmq: 'disconnected'
      }
    }
  },

  /**
   * Get all tasks from the backend
   */
  async getTasks() {
    try {
      const response = await api.get('/v1/tasks')
      return response.data.data || []
    } catch (error) {
      throw new Error(`Failed to fetch tasks: ${error.message}`)
    }
  },

  /**
   * Get a specific task by ID
   */
  async getTask(id) {
    try {
      const response = await api.get(`/v1/tasks/${id}`)
      return response.data.data
    } catch (error) {
      throw new Error(`Failed to fetch task ${id}: ${error.message}`)
    }
  },

  /**
   * Create a new task
   */
  async createTask(taskData) {
    try {
      const response = await api.post('/v1/tasks', taskData)
      return response.data
    } catch (error) {
      throw new Error(`Failed to create task: ${error.message}`)
    }
  },

  /**
   * Update an existing task
   */
  async updateTask(id, updates) {
    try {
      const response = await api.put(`/v1/tasks/${id}`, updates)
      return response.data
    } catch (error) {
      throw new Error(`Failed to update task ${id}: ${error.message}`)
    }
  },

  /**
   * Delete a task
   */
  async deleteTask(id) {
    try {
      const response = await api.delete(`/v1/tasks/${id}`)
      return response.data
    } catch (error) {
      throw new Error(`Failed to delete task ${id}: ${error.message}`)
    }
  },

  /**
   * Get task statistics
   */
  async getTaskStats() {
    try {
      // Since there's no stats endpoint, calculate from tasks
      const tasks = await this.getTasks()
      return {
        total: tasks.length,
        completed: tasks.filter(t => t.status === 'completed').length,
        pending: tasks.filter(t => t.status === 'pending').length,
        in_progress: tasks.filter(t => t.status === 'in-progress').length
      }
    } catch (error) {
      throw new Error(`Failed to fetch task stats: ${error.message}`)
    }
  }
}

export default taskService