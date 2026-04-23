# PaseLista - Checador en Línea

Sistema completo de registro de asistencia con geolocalización y fotografías.

## Stack

- **Backend:** Go + Gin + PostgreSQL
- **Frontend:** Vue 3 + Vite + TailwindCSS

## Funcionalidades

- ✅ Registro de perfil con captcha matemático
- ✅ Login con JWT
- ✅ Generación de QR que apunta al checador
- ✅ Registro de entrada/salida con:
  - 📍 Solicitud de ubicación en tiempo real con tracking del recorrido
  - 📸 Fotografía del sitio (cámara trasera)
  - 🤳 Selfie (cámara frontal)
  - ⚡ Soporte de flash y cambio de cámara
  - 🕐 Hora y fecha exactas del registro
- ✅ Historial de registros con fotos y puntos de recorrido

## Inicio Rápido con Docker

```bash
docker-compose up --build
```

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## Desarrollo Local

### Backend

```bash
cd backend
go mod tidy
# Requiere PostgreSQL corriendo en localhost:5432
go run main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Variables de Entorno (backend/.env)

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=paselista
JWT_SECRET=changeme_secret_key_32chars_long!
PORT=8080
FRONTEND_URL=http://localhost:5173
```

## API Endpoints

| Método | Ruta                  | Auth | Descripción              |
| ------ | --------------------- | ---- | ------------------------ |
| POST   | /api/auth/register    | No   | Crear perfil             |
| POST   | /api/auth/login       | No   | Iniciar sesión           |
| GET    | /api/auth/me          | Sí   | Perfil actual            |
| GET    | /api/qr               | No   | Imagen QR PNG            |
| GET    | /api/checkin-url      | No   | URL del checador         |
| POST   | /api/checks           | Sí   | Registrar entrada/salida |
| GET    | /api/checks           | Sí   | Mis registros            |
| GET    | /api/checks/:id/route | Sí   | Puntos del recorrido     |
| POST   | /api/location-points  | Sí   | Agregar punto GPS        |
