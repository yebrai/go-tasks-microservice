<template>
  <div class="events">
    <div class="events-header">
      <h2>📡 Event Stream</h2>
      <div class="connection-status">
        <span :class="['status', connected ? 'status-completed' : 'status-pending']">
          {{ connected ? 'Connected to RabbitMQ' : 'Disconnected' }}
        </span>
      </div>
    </div>

    <div class="grid grid-2">
      <!-- Real-time Events -->
      <div class="card">
        <h3>🔔 Real-time Events</h3>
        <div class="events-controls">
          <button @click="toggleConnection" class="btn btn-primary">
            {{ connected ? 'Disconnect' : 'Connect' }}
          </button>
          <button @click="clearEvents" class="btn btn-secondary">
            Clear Events
          </button>
        </div>
        
        <div class="events-list">
          <div v-if="events.length === 0" class="no-events">
            No events received yet...
          </div>
          <div
            v-for="event in events.slice().reverse()"
            :key="event.id"
            class="event-item"
          >
            <div class="event-header">
              <span class="event-type">{{ event.type }}</span>
              <span class="event-time">{{ formatTime(event.timestamp) }}</span>
            </div>
            <div class="event-details">
              <div class="event-aggregate">
                Aggregate ID: {{ event.aggregateId }}
              </div>
              <div class="event-payload">
                <pre>{{ JSON.stringify(event.payload, null, 2) }}</pre>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Event Statistics -->
      <div class="card">
        <h3>📊 Event Statistics</h3>
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-label">Total Events</div>
            <div class="stat-value">{{ eventStats.total }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">Task Created</div>
            <div class="stat-value">{{ eventStats.taskCreated }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">Task Updated</div>
            <div class="stat-value">{{ eventStats.taskUpdated }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">Task Completed</div>
            <div class="stat-value">{{ eventStats.taskCompleted }}</div>
          </div>
        </div>

        <div class="event-types">
          <h4>Event Types Distribution</h4>
          <div
            v-for="(count, type) in eventTypeDistribution"
            :key="type"
            class="event-type-bar"
          >
            <span class="event-type-name">{{ type }}</span>
            <div class="event-type-progress">
              <div
                class="event-type-fill"
                :style="{ width: getProgressWidth(count) + '%' }"
              ></div>
            </div>
            <span class="event-type-count">{{ count }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Event Details Modal (if needed) -->
    <div class="card">
      <h3>🔍 Event Stream Monitor</h3>
      <div class="monitor-info">
        <p>
          <strong>⚠️ Development Mode:</strong> This panel shows events from your task interactions.
          In production, events would be streamed from RabbitMQ through WebSocket/SSE endpoints.
        </p>
        <div class="monitor-details">
          <div><strong>Current Mode:</strong> localStorage simulation</div>
          <div><strong>Backend Events:</strong> Published to RabbitMQ exchange</div>
          <div><strong>Event Types:</strong> task.created, task.updated, task.completed</div>
          <div><strong>Note:</strong> Backend WebSocket/SSE endpoints needed for real-time streaming</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'

export default {
  name: 'Events',
  setup() {
    const events = ref([])
    const connected = ref(false)
    const ws = ref(null)
    let eventIdCounter = 1

    const eventStats = computed(() => {
      return {
        total: events.value.length,
        taskCreated: events.value.filter(e => e.type === 'task.created').length,
        taskUpdated: events.value.filter(e => e.type === 'task.updated').length,
        taskCompleted: events.value.filter(e => e.type === 'task.completed').length
      }
    })

    const eventTypeDistribution = computed(() => {
      const distribution = {}
      events.value.forEach(event => {
        distribution[event.type] = (distribution[event.type] || 0) + 1
      })
      return distribution
    })

    const getProgressWidth = (count) => {
      const max = Math.max(...Object.values(eventTypeDistribution.value))
      return max > 0 ? (count / max) * 100 : 0
    }

    // Load real events from localStorage (temporary solution until backend provides event endpoints)
    const loadEventsFromStorage = () => {
      const storedEvents = localStorage.getItem('taskEvents')
      if (storedEvents) {
        try {
          const parsed = JSON.parse(storedEvents)
          events.value = parsed || []
        } catch (error) {
          console.error('Error loading events from storage:', error)
          events.value = []
        }
      }
    }

    const connectToEventStream = () => {
      console.log('⚠️  Backend WebSocket/SSE endpoints not implemented yet. Showing events from localStorage.')
      connected.value = true
      loadEventsFromStorage()

      // Set up storage listener to detect new events
      const handleStorageChange = (e) => {
        if (e.key === 'taskEvents') {
          loadEventsFromStorage()
        }
      }

      window.addEventListener('storage', handleStorageChange)

      // Check for updates every few seconds
      const checkForUpdates = () => {
        if (!connected.value) return
        loadEventsFromStorage()
        setTimeout(checkForUpdates, 3000)
      }
      
      setTimeout(checkForUpdates, 1000)
    }

    const disconnectFromEventStream = () => {
      connected.value = false
      if (ws.value) {
        ws.value.close()
        ws.value = null
      }
    }

    const toggleConnection = () => {
      if (connected.value) {
        disconnectFromEventStream()
      } else {
        connectToEventStream()
      }
    }

    const clearEvents = () => {
      events.value = []
      localStorage.removeItem('taskEvents')
      console.log('🗑️ Events cleared from localStorage')
    }

    const formatTime = (date) => {
      return date.toLocaleTimeString()
    }

    onMounted(() => {
      // Auto-connect on mount
      connectToEventStream()
    })

    onUnmounted(() => {
      disconnectFromEventStream()
    })

    return {
      events,
      connected,
      eventStats,
      eventTypeDistribution,
      toggleConnection,
      clearEvents,
      formatTime,
      getProgressWidth
    }
  }
}
</script>

<style scoped>
.events-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.events-header h2 {
  margin: 0;
  color: #333;
}

.events-controls {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #545b62;
}

.events-list {
  max-height: 500px;
  overflow-y: auto;
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 1rem;
}

.no-events {
  text-align: center;
  color: #666;
  font-style: italic;
  padding: 2rem;
}

.event-item {
  border-bottom: 1px solid #eee;
  padding: 1rem 0;
  margin-bottom: 1rem;
}

.event-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.event-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.event-type {
  font-weight: bold;
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

.event-time {
  font-size: 0.875rem;
  color: #666;
}

.event-details {
  font-size: 0.875rem;
}

.event-aggregate {
  color: #666;
  margin-bottom: 0.5rem;
}

.event-payload pre {
  background: #f8f9fa;
  padding: 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  overflow-x: auto;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-item {
  text-align: center;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.stat-label {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #667eea;
}

.event-types h4 {
  margin-bottom: 1rem;
  color: #333;
}

.event-type-bar {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.event-type-name {
  min-width: 120px;
  font-size: 0.875rem;
}

.event-type-progress {
  flex: 1;
  height: 20px;
  background: #e9ecef;
  border-radius: 10px;
  overflow: hidden;
}

.event-type-fill {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s ease;
}

.event-type-count {
  min-width: 30px;
  text-align: right;
  font-weight: bold;
}

.monitor-info {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 4px;
  border-left: 4px solid #667eea;
}

.monitor-details {
  margin-top: 1rem;
  font-size: 0.875rem;
}

.monitor-details div {
  margin-bottom: 0.25rem;
}
</style>