# 📋 Go Tasks Microservice

> Microservicio de gestión de tareas implementado en Go con arquitectura hexagonal, CQRS y MongoDB.

## 🏗️ Arquitectura

- **Hexagonal Architecture** (Ports & Adapters)
- **CQRS** (Command Query Responsibility Segregation)
- **Domain-Driven Design** (DDD)
- **MongoDB** como persistencia

## 🚀 Quick Start

### Prerrequisitos
- Go 1.22+
- Docker
- Make (opcional)

### Ejecutar
```bash
# 1. Instalar dependencias
go mod tidy

# 2. Crear config.yaml
server:
  address: ":8080"
mongo:
  uri: "mongodb://mongoroot:secret@localhost:27017/?authSource=admin"
  database: "taskdb"

# 3. Levantar servicios
make infra      # MongoDB
make run        # Aplicación
```

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
├── cmd/api/                    # Entrada principal
├── internal/task/              # Dominio de tareas
│   ├── bootstrap/              # Configuración y DI
│   ├── creator/                # Command handlers
│   ├── http/                   # HTTP handlers
│   ├── mongo/                  # MongoDB repository
│   └── *.go                    # Entidades y eventos
├── pkg/                        # Paquetes compartidos
│   ├── cqrs/                   # Framework CQRS
│   ├── id/                     # Generación de IDs
│   └── runner/                 # Service runner
└── docker-compose.yaml        # Infraestructura
```

## 🛠️ Comandos Disponibles

```bash
make infra          # Levantar MongoDB
make infra-down     # Detener MongoDB
make run            # Ejecutar aplicación
make test           # Tests unitarios
make fmt            # Formatear código
```

## 🔧 Variables de Entorno

| Variable | Descripción | Defecto |
|----------|-------------|---------|
| `TASK_SERVER_ADDRESS` | Puerto del servidor | `:8080` |
| `TASK_MONGO_URI` | URI de MongoDB | `mongodb://mongoroot:secret@localhost:27017/?authSource=admin` |
| `TASK_MONGO_DATABASE` | Base de datos | `taskdb` |

## 🧪 Testing

```bash
make test                    # Tests unitarios
go test ./... -v             # Tests detallados
go test ./... -cover         # Con coverage
```

## 🔍 Debugging en GoLand

**Puntos clave para debugging:**
- Command Handlers: `internal/task/creator/command_handler.go`
- HTTP Handlers: `internal/task/http/task_handler.go`
- Repository: `internal/task/mongo/task_repository.go`
- Bootstrap: `internal/task/bootstrap/service.go`

## 🚀 Características

✅ **Arquitectura Hexagonal** - Separación clara de capas  
✅ **CQRS** - Comandos separados de consultas  
✅ **Event-Driven** - Preparado para eventos de dominio  
✅ **MongoDB** - Persistencia robusta  
✅ **Testing** - Tests unitarios incluidos  
✅ **Docker** - Containerización completa

## 🔄 Flujo de Trabajo

1. **Cliente** envía comando HTTP
2. **HTTP Handler** valida y procesa
3. **Command Bus** enruta al handler apropiado
4. **Command Handler** ejecuta lógica de negocio
5. **Repository** persiste en MongoDB
6. **Respuesta** regresa al cliente

## 🛠️ Troubleshooting

| Error | Solución |
|-------|----------|
| `docker daemon not running` | Iniciar Docker Desktop |
| `Authentication failed` | Verificar URI MongoDB |
| `port already in use` | Cambiar puerto en config.yaml |
| `no handler registered` | Verificar registro en bootstrap |
## 📄 Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.