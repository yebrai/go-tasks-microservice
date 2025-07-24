# 📋 Go Tasks Microservice

> Microservicio de gestión de tareas con frontend Vue.js, implementado en Go con arquitectura hexagonal, CQRS, MongoDB y RabbitMQ.

## 🏗️ Arquitectura

- **Backend**: Go con Hexagonal Architecture, CQRS y DDD
- **Frontend**: Vue.js 3 + Vue Router + Vite
- **Base de Datos**: MongoDB
- **Message Queue**: RabbitMQ para eventos
- **Proxy**: Nginx para servir frontend y proxy API
- **Containerización**: Docker multi-stage builds

## 🚀 Quick Start

### Prerrequisitos
- Docker & Docker Compose
- Make (opcional)
- Node.js 18+ (solo para desarrollo frontend)

### 🎯 Inicio Rápido (Todo en uno)
```bash
# Levantar todo el sistema
make dev

# En modo background
make dev-detached
```

### 🔧 Inicio por Partes
```bash
# 1. Infraestructura (MongoDB + RabbitMQ)
make infra

# 2. Backend Go
make run

# 3. Frontend Vue.js (modo desarrollo)
make frontend-dev
```

## 🌐 URLs Disponibles

| Servicio | URL | Descripción |
|----------|-----|-------------|
| **Frontend Dashboard** | http://localhost:3000 | Interfaz Vue.js completa |
| **Backend API** | http://localhost:8080 | API REST del microservicio |
| **RabbitMQ Management** | http://localhost:15672 | Interfaz de RabbitMQ (guest/guest) |
| **Health Check** | http://localhost:8080/health | Estado del sistema |

## 📋 API Endpoints

### Health Check
```http
GET /health
```

### Crear Tarea
```http
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "Nueva tarea",
  "description": "Descripción opcional",
  "due_date": "2025-12-31"
}
```

**Respuesta (201):**
```json
{
  "message": "Task created successfully",
  "success": true
}
```

## 📁 Estructura del Proyecto

```
├── cmd/api/                    # Entrada principal Go
├── internal/task/              # Dominio de tareas
│   ├── bootstrap/              # Configuración y DI
│   ├── creator/                # Command handlers
│   ├── http/                   # HTTP handlers
│   ├── mongo/                  # MongoDB repository
│   └── *.go                    # Entidades y eventos
├── pkg/                        # Paquetes compartidos
│   ├── cqrs/                   # Framework CQRS
│   ├── events/                 # Event Bus (RabbitMQ)
│   ├── id/                     # Generación de IDs
│   └── runner/                 # Service runner
├── web/                        # Frontend Vue.js
│   ├── src/
│   │   ├── views/              # Dashboard, Tasks, Events
│   │   ├── services/           # API integration
│   │   ├── router/             # Vue Router
│   │   └── components/
│   ├── package.json
│   └── vite.config.js
├── Dockerfile.backend          # Docker para Go
├── Dockerfile.frontend         # Docker para Vue + Nginx
├── docker-compose.yaml         # Orquestación completa
├── nginx.conf                  # Configuración Nginx
└── cloudbuild.yaml            # Deploy GCP
```

## 🛠️ Comandos Disponibles

### 🚀 Desarrollo
```bash
make dev                 # Todo el sistema
make dev-detached        # En background
make infra               # Solo infraestructura
make run                 # Solo backend Go
make frontend-dev        # Solo frontend (dev mode)
```

### 🏗️ Build & Deploy
```bash
make build-all           # Construir todas las imágenes
make build-backend       # Solo imagen backend
make build-frontend      # Solo imagen frontend
make gcp-build          # Deploy a Google Cloud
```

### 🔍 Monitoreo & Debug
```bash
make health              # Check salud de servicios
make logs                # Ver logs de todos
make logs-backend        # Logs solo backend
make logs-frontend       # Logs solo frontend
```

### 🧹 Limpieza
```bash
make clean               # Limpiar contenedores e imágenes
make infra-down         # Detener infraestructura
```

### 📊 Base de Datos & Cola
```bash
make db-shell           # Conectar a MongoDB
make rabbitmq-mgmt      # Abrir RabbitMQ Management UI
```

## 🎨 Frontend Features

### 📊 Dashboard
- Métricas visuales del sistema
- Estado de servicios en tiempo real
- Actividad reciente

### 📋 Gestión de Tareas
- Crear tareas con formulario
- Visualizar lista de tareas
- Estados: pending, in-progress, completed

### 📡 Monitor de Eventos
- Stream en tiempo real de RabbitMQ
- Estadísticas de eventos
- Distribución por tipos

## 🔧 Variables de Entorno

### Backend
| Variable | Descripción | Defecto |
|----------|-------------|---------|
| `TASK_SERVER_ADDRESS` | Puerto del servidor | `:8080` |
| `TASK_MONGO_URI` | URI de MongoDB | `mongodb://mongoroot:secret@mongo:27017/taskdb?authSource=admin` |
| `TASK_RABBITMQ_ENABLED` | Habilitar RabbitMQ | `true` |
| `TASK_RABBITMQ_URL` | URL de RabbitMQ | `amqp://guest:guest@rabbitmq:5672/` |

### Frontend
| Variable | Descripción | Defecto |
|----------|-------------|---------|
| `VITE_API_URL` | URL del backend API | `/api` (producción) |
| `VITE_WS_URL` | URL WebSocket | `/ws` |

## 🧪 Testing

```bash
make test                # Tests unitarios Go
make test-coverage       # Tests con coverage
go test ./... -v         # Tests detallados
```

## 🌩️ Deploy a Google Cloud Platform

### Setup Inicial
```bash
# Configurar proyecto GCP
gcloud config set project YOUR_PROJECT_ID

# Crear repositorio Artifact Registry
gcloud artifacts repositories create taskservice-repo \
  --repository-format=docker \
  --location=us-central1

# Habilitar APIs necesarias
gcloud services enable cloudbuild.googleapis.com run.googleapis.com
```

### Deploy
```bash
# Deploy automático con Cloud Build
make gcp-build

# O manualmente
gcloud builds submit --config=cloudbuild.yaml
```

Ver guía completa en: [`gcp-setup.md`](gcp-setup.md)

## 🔍 Debugging

### Puntos Clave Backend
- Command Handlers: `internal/task/creator/command_handler.go`
- HTTP Handlers: `internal/task/http/task_handler.go`
- Repository: `internal/task/mongo/task_repository.go`
- Event Bus: `pkg/events/rabbitmq.go`

### Debugging Frontend
- API Service: `web/src/services/taskService.js`
- Router: `web/src/router/index.js`
- Components: `web/src/views/`

## 🚀 Características

### Backend
✅ **Arquitectura Hexagonal** - Separación clara de capas  
✅ **CQRS** - Comandos separados de consultas  
✅ **Event-Driven** - RabbitMQ para eventos de dominio  
✅ **MongoDB** - Persistencia robusta  
✅ **Health Checks** - Monitoreo automático  
✅ **Docker** - Containerización optimizada  

### Frontend
✅ **Vue.js 3** - Framework moderno y reactivo  
✅ **Vue Router** - Navegación SPA  
✅ **Axios** - Cliente HTTP con interceptors  
✅ **Responsive Design** - Interfaz adaptable  
✅ **Real-time Events** - Monitor de eventos live  
✅ **Error Handling** - Manejo robusto de errores  

### DevOps
✅ **Multi-stage Builds** - Imágenes optimizadas  
✅ **Nginx Proxy** - Servidor web + API proxy  
✅ **Health Checks** - Monitoreo de contenedores  
✅ **GCP Ready** - Deploy automático a Cloud Run  

## 🔄 Flujo de Trabajo

### Crear Tarea
1. **Usuario** completa formulario en Vue.js frontend
2. **Axios** envía POST a `/api/v1/tasks`
3. **Nginx** hace proxy a `taskservice:8080`
4. **HTTP Handler** valida y procesa comando
5. **Command Bus** enruta al handler apropiado
6. **Command Handler** ejecuta lógica de negocio
7. **Repository** persiste en MongoDB
8. **Event** se publica a RabbitMQ
9. **Respuesta** regresa al frontend
10. **UI** se actualiza automáticamente

## 🛠️ Troubleshooting

| Error | Solución |
|-------|----------|
| `docker daemon not running` | Iniciar Docker Desktop |
| `port already in use` | Verificar puertos con `docker ps` |
| `frontend unhealthy` | Verificar nginx.conf y build |
| `backend unhealthy` | Revisar logs con `make logs-backend` |
| `rabbitmq connection failed` | Esperar a que RabbitMQ esté healthy |

## 💡 Comandos Útiles

```bash
# Ver ayuda completa
make help

# Guía de inicio rápido
make quickstart

# Verificar salud del sistema
make health

# Ver logs en tiempo real
make logs

# Conectar a base de datos
make db-shell

# Limpiar y empezar de cero
make clean && make dev
```
## 📄 Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.