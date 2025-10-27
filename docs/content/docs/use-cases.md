---
title: "Use Cases"
weight: 5
---

# Use Cases

Common scenarios where lanup can help streamline your development workflow.

## Mobile App Development

Test your backend API from your phone or tablet without deploying to a staging server.

### Setup

```bash
# Start your backend
npm run dev

# Expose it on your network
lanup start --watch
```

### Usage

Access your API from your mobile device:
```
http://192.168.1.100:3000/api/users
```

### Benefits

- Test on real devices
- No need for emulators
- Test network conditions
- Faster iteration

---

## Supabase Local Development

Automatically expose Supabase services for testing on mobile devices.

### Setup

```yaml
# .lanup.yaml
vars:
  SUPABASE_URL: "http://localhost:54321"
  SUPABASE_ANON_KEY: "your-anon-key"

output: ".env.local"

auto_detect:
  supabase: true
```

```bash
# Start Supabase
supabase start

# Expose with auto-detection
lanup start
```

### Result

Your Supabase services are now accessible:
- API: `http://192.168.1.100:54321`
- Studio: `http://192.168.1.100:54323`
- Database: `postgresql://...@192.168.1.100:54322/postgres`

---

## Docker Development

Auto-detect and expose Docker containers without manual configuration.

### Setup

```yaml
# .lanup.yaml
auto_detect:
  docker: true
```

```bash
# Start your Docker containers
docker-compose up

# Expose them
lanup start
```

### Result

lanup automatically detects containers and creates variables:
```
DOCKER_WEB_SERVER_PORT=http://192.168.1.100:8080
DOCKER_DATABASE_PORT=http://192.168.1.100:5432
DOCKER_REDIS_PORT=http://192.168.1.100:6379
```

---

## Team Collaboration

Share your local development environment with teammates on the same network.

### Setup

```bash
# Start with watch mode
lanup start --watch
```

### Usage

Share the generated URLs with your team:
```
API: http://192.168.1.100:8000
Frontend: http://192.168.1.100:3000
```

### Benefits

- Quick demos without deployment
- Pair programming
- Cross-device testing
- No cloud costs

---

## Multi-Device Testing

Test your responsive web app across multiple devices simultaneously.

### Setup

```yaml
# .lanup.yaml
vars:
  FRONTEND_URL: "http://localhost:3000"
  API_URL: "http://localhost:8000"

output: ".env.local"
```

```bash
lanup start --watch
```

### Usage

Access from:
- **Desktop browser:** `http://192.168.1.100:3000`
- **Tablet:** `http://192.168.1.100:3000`
- **Phone:** `http://192.168.1.100:3000`
- **Smart TV:** `http://192.168.1.100:3000`

---

## WebSocket Testing

Expose WebSocket servers for real-time testing.

### Setup

```yaml
# .lanup.yaml
vars:
  WS_URL: "ws://localhost:8080"
  API_URL: "http://localhost:8000"

output: ".env.local"
```

### Result

```
WS_URL=ws://192.168.1.100:8080
API_URL=http://192.168.1.100:8000
```

---

## CI/CD Preview

Use lanup in your CI/CD pipeline for quick previews.

### Setup

```bash
# In your CI script
lanup init --force
lanup start --no-env
```

### Benefits

- Quick environment setup
- Consistent configuration
- Easy debugging

---

## Network Switching

Automatically update environment when switching networks (Wi-Fi to Ethernet, different Wi-Fi networks).

### Setup

```bash
lanup start --watch
```

### Behavior

lanup monitors your network and automatically:
1. Detects IP address changes
2. Updates environment file
3. Notifies you of the change

### Output

```
⚠️  Network change detected!
  Old IP: 192.168.1.100
  New IP: 10.0.0.50

ℹ Regenerating environment file...
✓ Environment file updated successfully!
```

---

## Quick Service Exposure

Expose a single service without creating a configuration file.

### Usage

```bash
# Expose a service
lanup expose http://localhost:3000

# With custom name
lanup expose http://localhost:8080 --name api

# With different port
lanup expose http://localhost:5000 --port 8000
```

### Benefits

- No configuration needed
- Quick one-off exposures
- Perfect for demos
