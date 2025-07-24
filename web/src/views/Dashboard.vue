<template>
  <div class="dashboard">
    <h2>📊 Dashboard</h2>
    
    <div class="grid grid-3">
      <div class="card">
        <h3>📋 Total Tasks</h3>
        <div class="metric">{{ stats.totalTasks }}</div>
      </div>
      
      <div class="card">
        <h3>✅ Completed</h3>
        <div class="metric">{{ stats.completedTasks }}</div>
      </div>
      
      <div class="card">
        <h3>⏳ Pending</h3>
        <div class="metric">{{ stats.pendingTasks }}</div>
      </div>
    </div>

    <div class="grid grid-2">
      <div class="card">
        <h3>🔄 System Status</h3>
        <div class="status-grid">
          <div class="status-item">
            <span>Backend API:</span>
            <span :class="['status', systemStatus.backend ? 'status-completed' : 'status-pending']">
              {{ systemStatus.backend ? 'Online' : 'Offline' }}
            </span>
          </div>
          <div class="status-item">
            <span>MongoDB:</span>
            <span :class="['status', systemStatus.database ? 'status-completed' : 'status-pending']">
              {{ systemStatus.database ? 'Connected' : 'Disconnected' }}
            </span>
          </div>
          <div class="status-item">
            <span>RabbitMQ:</span>
            <span :class="['status', systemStatus.rabbitmq ? 'status-completed' : 'status-pending']">
              {{ systemStatus.rabbitmq ? 'Connected' : 'Disconnected' }}
            </span>
          </div>
        </div>
      </div>

      <div class="card">
        <h3>📈 Recent Activity</h3>
        <div class="activity-list">
          <div v-for="event in recentEvents" :key="event.id" class="activity-item">
            <span class="activity-time">{{ formatTime(event.timestamp) }}</span>
            <span class="activity-desc">{{ event.description }}</span>
          </div>
          <div v-if="recentEvents.length === 0" class="no-activity">
            No recent activity
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { taskService } from '../services/taskService'

export default {
  name: 'Dashboard',
  setup() {
    const stats = ref({
      totalTasks: 0,
      completedTasks: 0,
      pendingTasks: 0
    })
    
    const systemStatus = ref({
      backend: false,
      database: false,
      rabbitmq: false
    })
    
    const recentEvents = ref([])

    const checkSystemHealth = async () => {
      try {
        const health = await taskService.getHealth()
        systemStatus.value.backend = health.status === 'ok'
        systemStatus.value.database = health.database === 'connected'
        systemStatus.value.rabbitmq = health.rabbitmq === 'connected'
      } catch (error) {
        console.error('Health check failed:', error)
        systemStatus.value = { backend: false, database: false, rabbitmq: false }
      }
    }

    const loadStats = async () => {
      try {
        // Simulated stats - replace with real API calls
        stats.value = {
          totalTasks: 12,
          completedTasks: 8,
          pendingTasks: 4
        }
        
        // Simulated recent events
        recentEvents.value = [
          { id: 1, timestamp: new Date(), description: 'Task "Fix login bug" completed' },
          { id: 2, timestamp: new Date(Date.now() - 300000), description: 'New task created: "Update documentation"' },
          { id: 3, timestamp: new Date(Date.now() - 600000), description: 'Task "Deploy to staging" started' }
        ]
      } catch (error) {
        console.error('Failed to load stats:', error)
      }
    }

    const formatTime = (date) => {
      return date.toLocaleTimeString()
    }

    onMounted(async () => {
      await Promise.all([checkSystemHealth(), loadStats()])
    })

    return {
      stats,
      systemStatus,
      recentEvents,
      formatTime
    }
  }
}
</script>

<style scoped>
.dashboard h2 {
  margin-bottom: 2rem;
  color: #333;
}

.metric {
  font-size: 2.5rem;
  font-weight: bold;
  color: #667eea;
  margin-top: 0.5rem;
}

.status-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.activity-list {
  max-height: 300px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid #eee;
}

.activity-time {
  font-size: 0.875rem;
  color: #666;
  min-width: 80px;
}

.activity-desc {
  flex: 1;
  margin-left: 1rem;
}

.no-activity {
  text-align: center;
  color: #666;
  font-style: italic;
  padding: 2rem;
}
</style>