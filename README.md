# Contact Form API

A production-ready backend service in Go following **Clean Architecture** principles for handling portfolio website contact form submissions via email.

## Features

- ✅ **Clean Architecture** - Domain, Use Case, Delivery, Infrastructure layers
- ✅ **RESTful JSON API** - Fiber framework
- ✅ **SMTP Email Delivery** - Gmail App Password support
- ✅ **Input Validation** - Email format, message length
- ✅ **Rate Limiting** - IP-based protection
- ✅ **CORS Support** - Configurable origins
- ✅ **Request Logging** - Detailed logging with timing
- ✅ **Graceful Shutdown** - Signal handling
- ✅ **Docker Ready** - Multi-stage build, health checks
- ✅ **No Hardcoded Secrets** - Environment-based configuration

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                         Delivery Layer                           │
│  (HTTP Handlers, Middleware, Router, DTOs)                       │
├──────────────────────────────────────────────────────────────────┤
│                         Use Case Layer                           │
│  (Application Business Rules, Input/Output DTOs)                 │
├──────────────────────────────────────────────────────────────────┤
│                         Domain Layer                             │
│  (Entities, Repository Interfaces, Business Rules)               │
├──────────────────────────────────────────────────────────────────┤
│                      Infrastructure Layer                         │
│  (SMTP Repository, Config, External Services)                    │
└──────────────────────────────────────────────────────────────────┘
```

### Dependency Flow

```
Delivery → Use Case → Domain ← Infrastructure
```

- **Domain** contains entities and repository interfaces (no external dependencies)
- **Use Case** depends only on Domain
- **Delivery** depends on Use Case (calls use cases, transforms HTTP ↔ DTOs)
- **Infrastructure** implements Domain interfaces (Dependency Inversion)

## Project Structure

```
go-mail-server/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Domain Layer (innermost)
│   │   ├── entity/
│   │   │   └── contact.go          # Contact entity + validation
│   │   └── repository/
│   │       └── email_repository.go # Repository interface
│   ├── usecase/                    # Use Case Layer
│   │   └── contact/
│   │       ├── interface.go        # Use case interface + DTOs
│   │       └── contact_usecase.go  # Use case implementation
│   ├── delivery/                   # Delivery Layer (HTTP)
│   │   └── http/
│   │       ├── dto/
│   │       │   ├── request.go      # Request DTOs
│   │       │   └── response.go     # Response DTOs
│   │       ├── handler/
│   │       │   ├── contact_handler.go
│   │       │   └── health_handler.go
│   │       ├── middleware/
│   │       │   ├── rate_limiter.go
│   │       │   ├── logger.go
│   │       │   └── recover.go
│   │       └── router/
│   │           └── router.go       # Route configuration
│   └── infrastructure/             # Infrastructure Layer
│       ├── config/
│       │   └── config.go           # Environment configuration
│       └── email/
│           └── smtp_repository.go  # SMTP implementation
├── Dockerfile                      # Multi-stage Docker build
├── docker-compose.yml              # Docker Compose configuration
├── .dockerignore                   # Docker ignore rules
├── Makefile                        # Development commands
├── .env.example                    # Example environment file
├── .gitignore                      # Git ignore rules
├── go.mod                          # Go module
├── go.sum                          # Dependency checksums
└── README.md                       # This file
```

## Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose (for containerized deployment)
- Gmail account with App Password

## Gmail App Password Setup

1. Go to your [Google Account settings](https://myaccount.google.com/)
2. Navigate to **Security → 2-Step Verification** (enable if not already)
3. Go to **Security → App passwords**
4. Generate a new app password for "Mail"
5. Use this 16-character password in `SMTP_PASSWORD`

## Quick Start

### Local Development

1. **Clone and setup:**
```bash
git clone <repository-url>
cd go-mail-server
cp .env.example .env
```

2. **Configure `.env`:**
```env
APP_PORT=3000
APP_ENV=development
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
RECEIVER_EMAIL=where-to-receive@example.com
ALLOWED_ORIGINS=http://localhost:3000
RATE_LIMIT=10
```

3. **Run:**
```bash
go mod tidy
go run ./cmd/api
```

### Docker Deployment

1. **Configure environment:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

2. **Build and run:**
```bash
# Using Docker Compose
docker-compose up -d

# Or build manually
docker build -t contact-form-api:1.0.0 .
docker run -d -p 3000:3000 --env-file .env contact-form-api:1.0.0
```

3. **View logs:**
```bash
docker-compose logs -f
```

4. **Stop:**
```bash
docker-compose down
```

### Using Makefile

```bash
make help          # Show available commands
make build         # Build binary
make run           # Run locally
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
make docker-stop   # Stop containers
make docker-logs   # View logs
```

## API Endpoints

### Health Check
```http
GET /health
```
Response:
```json
{
  "status": "healthy",
  "service": "contact-form-api",
  "version": "1.0.0"
}
```

### Readiness Check (Kubernetes)
```http
GET /ready
```

### Submit Contact Form
```http
POST /api/contact
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "subject": "Hello",
  "message": "This is a test message from my portfolio."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

**Error Response (400/429/500):**
```json
{
  "success": false,
  "message": "error description"
}
```

## Frontend Integration

### JavaScript Fetch
```javascript
async function submitContactForm(formData) {
  try {
    const response = await fetch('https://your-api.com/api/contact', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: formData.name,
        email: formData.email,
        subject: formData.subject,
        message: formData.message,
      }),
    });

    const data = await response.json();

    if (data.success) {
      console.log('Email sent successfully!');
    } else {
      console.error(`Error: ${data.message}`);
    }

    return data;
  } catch (error) {
    console.error('Network error:', error);
    throw error;
  }
}
```

### React Hook
```jsx
import { useState } from 'react';

const useContactForm = (apiUrl) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  const submit = async (formData) => {
    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      const response = await fetch(`${apiUrl}/api/contact`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (data.success) {
        setSuccess(true);
      } else {
        setError(data.message);
      }

      return data;
    } catch (err) {
      setError('Network error. Please try again.');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { submit, loading, error, success };
};
```

## Security Best Practices

### Implemented
- ✅ **Rate Limiting** - IP-based, configurable (default: 10 req/min)
- ✅ **Input Validation** - Email format, length limits
- ✅ **HTML Escaping** - XSS prevention in emails
- ✅ **CORS** - Configurable allowed origins
- ✅ **Non-root Docker** - Container runs as unprivileged user
- ✅ **Graceful Shutdown** - Proper signal handling

### Recommended Additions
- **reCAPTCHA/hCaptcha** - Add bot protection
- **Honeypot Fields** - Hidden fields to catch bots
- **Request Size Limiting** - Already default 4MB in Fiber
- **HTTPS** - Use reverse proxy (Nginx/Traefik)
- **IP Blacklisting** - Block known bad actors

## Extending Email Providers

The architecture uses interfaces for easy provider switching:

```go
// Domain interface (internal/domain/repository/email_repository.go)
type EmailRepository interface {
    Send(contact *entity.Contact) error
}

// To add SendGrid, create:
// internal/infrastructure/email/sendgrid_repository.go
type sendGridRepository struct {
    apiKey string
    // ...
}

func (r *sendGridRepository) Send(contact *entity.Contact) error {
    // SendGrid API implementation
}

// Then in main.go, switch the repository:
emailRepo := email.NewSendGridRepository(cfg)
```

## Production Deployment

### Deploy to Server (Docker Run)

1. **Build and push image:**
```bash
docker login
docker build -t andrianprasetya/go-mail-server:latest .
docker push andrianprasetya/go-mail-server:latest
```

2. **On the server, pull and run:**
```bash
docker pull andrianprasetya/go-mail-server:latest

docker run -d \
  --name contact-form-api \
  --restart unless-stopped \
  -p 8090:3000 \
  -e APP_PORT=3000 \
  -e APP_ENV=production \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=your-email@gmail.com \
  -e SMTP_PASSWORD=your-app-password \
  -e RECEIVER_EMAIL=receiver@example.com \
  -e ALLOWED_ORIGINS=https://yourdomain.com \
  -e RATE_LIMIT=2 \
  -e RATE_LIMIT_EXPIRATION_HOURS=24 \
  andrianprasetya/go-mail-server:latest
```

Or using an `.env` file on the server:
```bash
docker run -d \
  --name contact-form-api \
  --restart unless-stopped \
  -p 8090:3000 \
  --env-file .env \
  andrianprasetya/go-mail-server:latest
```

3. **Useful commands:**
```bash
docker logs -f contact-form-api    # View logs
docker restart contact-form-api    # Restart
docker stop contact-form-api       # Stop
docker rm contact-form-api         # Remove container
```

### Docker Compose with Reverse Proxy

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  api:
    image: contact-form-api:1.0.0
    environment:
      - APP_ENV=production
      # ... other env vars
    networks:
      - web
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`api.yourdomain.com`)"
      - "traefik.http.routers.api.tls.certresolver=letsencrypt"
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: contact-form-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: contact-form-api
  template:
    spec:
      containers:
      - name: api
        image: contact-form-api:1.0.0
        ports:
        - containerPort: 3000
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
        envFrom:
        - secretRef:
            name: contact-form-secrets
```

## License

MIT License
