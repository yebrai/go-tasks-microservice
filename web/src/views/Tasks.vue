<template>
  <div class="tasks">
    <div class="tasks-header">
      <h2>📋 Tasks Management</h2>
      <button @click="showCreateForm = !showCreateForm" class="btn btn-primary">
        {{ showCreateForm ? 'Cancel' : 'Create Task' }}
      </button>
    </div>

    <!-- Create Task Form -->
    <div v-if="showCreateForm" class="card create-form">
      <h3>Create New Task</h3>
      <form @submit.prevent="createTask">
        <div class="form-group">
          <label for="title">Title:</label>
          <input
            id="title"
            v-model="newTask.title"
            type="text"
            required
            placeholder="Enter task title"
          />
        </div>
        
        <div class="form-group">
          <label for="description">Description:</label>
          <textarea
            id="description"
            v-model="newTask.description"
            placeholder="Enter task description (optional)"
          ></textarea>
        </div>
        
        <div class="form-group">
          <label for="due_date">Due Date:</label>
          <input
            id="due_date"
            v-model="newTask.due_date"
            type="date"
          />
        </div>
        
        <div class="form-actions">
          <button type="submit" class="btn btn-success" :disabled="loading">
            {{ loading ? 'Creating...' : 'Create Task' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="error">
      {{ error }}
    </div>

    <!-- Tasks List -->
    <div class="tasks-grid">
      <div v-if="loading && tasks.length === 0" class="loading">
        Loading tasks...
      </div>
      
      <div v-else-if="tasks.length === 0" class="card no-tasks">
        <h3>No tasks found</h3>
        <p>Create your first task to get started!</p>
      </div>
      
      <div v-else class="grid grid-2">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="card task-card"
        >
          <div class="task-header">
            <h4>{{ task.title }}</h4>
            <span :class="['status', getStatusClass(task.status)]">
              {{ task.status || 'pending' }}
            </span>
          </div>
          
          <p v-if="task.description" class="task-description">
            {{ task.description }}
          </p>
          
          <div class="task-meta">
            <div v-if="task.due_date" class="task-due-date">
              📅 Due: {{ formatDate(task.due_date) }}
            </div>
            <div class="task-created">
              🕒 Created: {{ formatDate(task.created_at) }}
            </div>
          </div>
          
          <div class="task-actions">
            <button
              v-if="task.status !== 'completed'"
              @click="completeTask(task.id)"
              class="btn btn-success"
              :disabled="loading"
            >
              Mark Complete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, reactive } from 'vue'
import { taskService } from '../services/taskService'

export default {
  name: 'Tasks',
  setup() {
    const tasks = ref([])
    const loading = ref(false)
    const error = ref('')
    const showCreateForm = ref(false)
    
    const newTask = reactive({
      title: '',
      description: '',
      due_date: ''
    })

    const loadTasks = async () => {
      loading.value = true
      error.value = ''
      try {
        tasks.value = await taskService.getTasks()
      } catch (err) {
        error.value = 'Failed to load tasks: ' + err.message
        console.error('Failed to load tasks:', err)
      } finally {
        loading.value = false
      }
    }

    const createTask = async () => {
      if (!newTask.title?.trim()) {
        error.value = 'Task title is required'
        return
      }

      loading.value = true
      error.value = ''
      try {
        const taskData = {
          title: newTask.title,
          description: newTask.description || undefined,
          due_date: newTask.due_date || undefined
        }
        
        await taskService.createTask(taskData)
        
        // Reset form
        Object.assign(newTask, { title: '', description: '', due_date: '' })
        showCreateForm.value = false
        
        // Reload tasks
        await loadTasks()
      } catch (err) {
        error.value = 'Failed to create task: ' + err.message
        console.error('Failed to create task:', err)
      } finally {
        loading.value = false
      }
    }

    const completeTask = async (taskId) => {
      loading.value = true
      try {
        try {
          await taskService.updateTask(taskId, { status: 'completed' })
          await loadTasks()
        } catch (err) {
          error.value = 'Failed to complete task: ' + err.message
          console.error('Failed to complete task:', err)
        }
      } catch (err) {
        error.value = 'Failed to complete task: ' + err.message
      } finally {
        loading.value = false
      }
    }

    const getStatusClass = (status) => {
      switch (status) {
        case 'completed': return 'status-completed'
        case 'in-progress': return 'status-in-progress'
        default: return 'status-pending'
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString()
    }

    onMounted(() => {
      loadTasks()
    })

    return {
      tasks,
      loading,
      error,
      showCreateForm,
      newTask,
      createTask,
      completeTask,
      getStatusClass,
      formatDate
    }
  }
}
</script>

<style scoped>
.tasks-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.tasks-header h2 {
  margin: 0;
  color: #333;
}

.create-form {
  margin-bottom: 2rem;
}

.create-form h3 {
  margin-bottom: 1.5rem;
  color: #333;
}

.form-actions {
  display: flex;
  gap: 1rem;
}

.task-card {
  position: relative;
  transition: transform 0.2s, box-shadow 0.2s;
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.task-header h4 {
  margin: 0;
  color: #333;
  flex: 1;
  margin-right: 1rem;
}

.task-description {
  color: #666;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.task-meta {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 1rem;
}

.task-due-date {
  margin-bottom: 0.25rem;
}

.task-actions {
  display: flex;
  gap: 0.5rem;
}

.no-tasks {
  text-align: center;
  padding: 3rem;
}

.no-tasks h3 {
  color: #666;
  margin-bottom: 1rem;
}

.no-tasks p {
  color: #999;
}
</style>