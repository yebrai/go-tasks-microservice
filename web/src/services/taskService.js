import api from './api'

// Helper function to save events to localStorage
const saveEvent = (type, taskData, aggregateId = null) => {
  try {
    const event = {
      id: Date.now() + Math.random(), // Simple unique ID
      type: type,
      timestamp: new Date(),
      aggregateId: aggregateId || taskData?.id || 'unknown',
      payload: {
        ...taskData,
        timestamp: new Date().toISOString()
      }
    }

    const storedEvents = localStorage.getItem('taskEvents')
    const events = storedEvents ? JSON.parse(storedEvents) : []
    events.push(event)

    // Keep only last 100 events
    if (events.length > 100) {
      events.splice(0, events.length - 100)
    }

    localStorage.setItem('taskEvents', JSON.stringify(events))
    console.log(`📡 Event saved: ${type}`, event)
  } catch (error) {
    console.error('Error saving event:', error)
  }
}

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
   * NOTE: Backend doesn't have GET endpoint yet, so using mock data
   */
  async getTasks() {
    try {
      const response = await api.get('/v1/tasks')
      return response.data || []
    } catch (error) {
      console.error('Failed to fetch tasks:', error)
      if (USE_MOCK_DATA) {
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
        ];
      }
      return [];
    }

    return baseTasks
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
      console.log('✅ Task created successfully via API:', response.data)

      // Since backend doesn't return the created task, we'll create it locally
      const newTask = {
        id: Date.now().toString(), // Simple ID generation
        title: taskData.title,
        description: taskData.description || '',
        status: 'pending',
        created_at: new Date().toISOString(),
        due_date: taskData.due_date || null
      }

      // Save to localStorage to simulate persistence
      const storedTasks = localStorage.getItem('tasks')
      const tasks = storedTasks ? JSON.parse(storedTasks) : []
      tasks.push(newTask)
      localStorage.setItem('tasks', JSON.stringify(tasks))

      // Save event for task creation
      saveEvent('task.created', newTask)

      console.log('✅ Task saved locally:', newTask)
      return { ...response.data, task: newTask }
    } catch (error) {
      console.error('❌ Failed to create task:', error)
      throw error
    }
  },

  /**
   * Update an existing task
   * NOTE: Backend doesn't have PUT endpoint yet, so using localStorage
   */
  async updateTask(id, updates) {
    console.log('⚠️  Backend PUT /v1/tasks/:id endpoint not implemented yet. Using localStorage.')

    try {
      // Get tasks from localStorage
      const storedTasks = localStorage.getItem('tasks')
      let tasks = storedTasks ? JSON.parse(storedTasks) : []

      // Find and update the task
      const taskIndex = tasks.findIndex(task => task.id === id)
      if (taskIndex === -1) {
        throw new Error(`Task with id ${id} not found`)
      }

      // Update the task
      tasks[taskIndex] = { ...tasks[taskIndex], ...updates }

      // Save back to localStorage
      localStorage.setItem('tasks', JSON.stringify(tasks))

      // Save event for task update
      const eventType = updates.status === 'completed' ? 'task.completed' : 'task.updated'
      saveEvent(eventType, tasks[taskIndex])

      console.log('✅ Task updated locally:', tasks[taskIndex])
      return {
        message: 'Task updated successfully',
        success: true,
        task: tasks[taskIndex]
      }
    } catch (error) {
      console.error(`❌ Failed to update task ${id}:`, error)
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