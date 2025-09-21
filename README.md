# News to Text Alert System

A comprehensive news alert system built with Go backend and React frontend that allows users to set up personalized news alerts and receive text notifications.

## Features

### Backend (Go)
- **Clean Architecture**: Proper separation of concerns with handlers, services, repositories, and models
- **Authentication System**: JWT-based authentication with password hashing using bcrypt
- **Database Layer**: MySQL with GORM ORM and proper migrations
- **Caching Layer**: Redis for session management and rate limiting
- **News Integration**: Support for NewsAPI and RSS feeds
- **Background Jobs**: Automated news checking and alert processing
- **SMS Notifications**: Ready for SMS gateway integration
- **API Documentation**: Swagger/OpenAPI documentation
- **Docker Support**: Complete Docker development environment

### Frontend (React)
- **Modern React**: Built with React 18 and functional components
- **Authentication**: JWT token management with protected routes
- **Responsive UI**: Clean and modern interface
- **Alert Management**: Full CRUD operations for news alerts
- **History Tracking**: View sent alert history
- **Real-time Features**: Test alerts and immediate feedback

## Tech Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Gin HTTP framework
- **Database**: MySQL 8.0
- **Cache**: Redis 7
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Testing**: Go testing package with coverage
- **Containerization**: Docker & Docker Compose

### Frontend
- **Framework**: React 18
- **Routing**: React Router v6
- **HTTP Client**: Axios
- **Styling**: CSS3 with responsive design
- **Build Tool**: Create React App

## Project Structure

```
NewsToText/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   ├── database/        # Database initialization
│   │   ├── cache/           # Redis cache layer
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # HTTP middleware
│   │   ├── models/          # Data models
│   │   ├── repositories/    # Data access layer
│   │   └── services/        # Business logic
│   ├── pkg/
│   │   ├── auth/            # JWT utilities
│   │   ├── utils/           # Helper functions
│   │   └── logger/          # Logging utilities
│   ├── migrations/          # Database migrations
│   ├── docs/               # API documentation
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── pages/           # Page components
│   │   ├── services/        # API services
│   │   └── utils/           # Helper functions
│   ├── public/
│   └── Dockerfile
└── docker-compose.yml
```

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)
- MySQL 8.0+ (for local development)
- Redis 7+ (for local development)

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd NewsToText
   ```

2. **Set up environment variables**
   ```bash
   cp backend/.env.example backend/.env
   # Edit backend/.env with your configuration
   # Get your free API token from https://www.thenewsapi.com/
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Documentation: http://localhost:8080/swagger/index.html
   - Database Admin (Adminer): http://localhost:8081

### Local Development

#### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

4. **Start MySQL and Redis**
   ```bash
   # Using Docker
   docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=newstotext mysql:8.0
   docker run -d --name redis -p 6379:6379 redis:7-alpine
   ```

5. **Run database migrations**
   ```bash
   make migrate-up
   ```

6. **Start the backend server**
   ```bash
   make run
   ```

#### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start the development server**
   ```bash
   npm start
   ```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/logout` - Logout user

### Alerts (Protected)
- `GET /api/v1/alerts` - Get user alerts
- `POST /api/v1/alerts` - Create new alert
- `PUT /api/v1/alerts/:id` - Update alert
- `DELETE /api/v1/alerts/:id` - Delete alert
- `GET /api/v1/alerts/history` - Get alert history
- `POST /api/v1/alerts/test` - Test alert

### System
- `GET /health` - Health check
- `GET /swagger/*` - API documentation

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | MySQL connection string | Local MySQL |
| `REDIS_URL` | Redis connection string | Local Redis |
| `JWT_SECRET` | JWT signing secret | Change in production |
| `NEWS_API_KEY` | The News API token | Optional |
| `SMS_API_KEY` | SMS provider API key | Optional |
| `LOG_LEVEL` | Logging level | `info` |

### Alert Frequencies
- **Real-time**: Checks every 5 minutes
- **Hourly**: Checks every hour
- **Daily**: Checks once daily at 9 AM

## Testing

### Backend Tests
```bash
cd backend

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./internal/services
```

### Test Coverage
The project includes comprehensive unit tests for:
- Authentication service
- Alert service
- News service
- JWT utilities
- Password utilities

## Development Commands

### Backend (using Makefile)
```bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make swagger        # Generate API documentation
make lint           # Run linter
make fmt            # Format code
make docker-build   # Build Docker image
```

### Frontend
```bash
npm start           # Start development server
npm test            # Run tests
npm run build       # Build for production
npm run eject       # Eject from Create React App
```

## External Services

### News Sources
- **The News API**: Primary news source (requires API token from thenewsapi.com)
  - 50+ languages supported
  - Real-time news from 1000+ sources
  - Advanced search and filtering
  - Better reliability than NewsAPI.org
  - Free tier with 1000 requests/month
- **RSS Feeds**: Fallback news sources when API is unavailable
  - TechCrunch
  - Reuters Technology
  - Hacker News
  - Bloomberg Markets

### SMS Integration
The system is ready for SMS provider integration. Popular options:
- Twilio
- AWS SNS
- Nexmo/Vonage
- TextMagic

## Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Tokens**: Secure token-based authentication
- **Token Blacklisting**: Logout invalidates tokens
- **CORS**: Configurable cross-origin resource sharing
- **Input Validation**: Request validation and sanitization
- **SQL Injection Protection**: GORM ORM with prepared statements

## Performance Features

- **Redis Caching**: Session management and news caching
- **Background Jobs**: Asynchronous news processing
- **Connection Pooling**: Efficient database connections
- **Graceful Shutdown**: Proper cleanup on termination

## Monitoring & Observability

- **Health Checks**: Application health endpoint
- **Structured Logging**: Configurable log levels
- **Error Handling**: Comprehensive error responses
- **API Documentation**: Swagger/OpenAPI specs
