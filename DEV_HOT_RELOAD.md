# Hot Module Reloading for Thaiyyal Development

This document describes how to use hot module reloading (HMR) for frontend development in Thaiyyal.

## Overview

Thaiyyal now supports hot module reloading for frontend development, allowing you to see changes instantly without rebuilding or restarting the server. The implementation supports multiple development workflows:

1. **Local Development**: Running both Next.js dev server and Go backend locally
2. **Docker Development**: Running in Docker with volume mounts for hot reloading
3. **Docker Compose**: Full stack development environment with hot reloading

## How It Works

### Architecture

- **Production Mode** (default): Frontend is built as static files and embedded into the Go binary
- **Development Mode** (`-tags dev`): Go backend proxies frontend requests to Next.js dev server running on port 3000

### Build Tags

The implementation uses Go build tags to separate dev and production modes:

- `embed.go`: Production mode (embeds static files)
- `embed_dev.go`: Development mode (proxies to Next.js dev server)

When built with `-tags dev`, the Go backend:
1. Skips embedding static files
2. Proxies all frontend requests to `http://localhost:3000`
3. Serves API routes normally on port 8080

## Quick Start

### Method 1: Using the Dev Script (Recommended)

The easiest way to start development with hot reloading:

```bash
# Make the script executable (first time only)
chmod +x dev.sh

# Start both frontend and backend in dev mode
./dev.sh
```

Or using npm:

```bash
npm run dev:full
```

This script:
- ✅ Installs dependencies if needed
- ✅ Starts Next.js dev server on port 3000
- ✅ Builds and starts Go backend in dev mode on port 8080
- ✅ Manages both processes with proper cleanup
- ✅ Shows helpful URLs and instructions

**Access the application:**
- Frontend: http://localhost:3000 (Next.js dev server with HMR)
- Backend: http://localhost:8080 (proxies to Next.js + serves API)
- API: http://localhost:8080/api/v1/workflow/execute

**Hot Reloading:**
- Frontend changes: Instant reload via Next.js HMR
- Backend changes: Restart the Go server manually

### Method 2: Manual Setup

If you prefer more control, run the servers separately:

**Terminal 1 - Frontend:**
```bash
npm run dev
# Next.js dev server starts on http://localhost:3000
```

**Terminal 2 - Backend:**
```bash
cd backend/cmd/server
go run -tags dev . -addr :8080
# Backend starts on http://localhost:8080 and proxies to Next.js
```

### Method 3: Docker Compose (Full Environment)

For a complete development environment in Docker:

```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up

# Or in detached mode
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Stop environment
docker-compose -f docker-compose.dev.yml down
```

The Docker setup includes:
- ✅ Volume mounts for hot reloading
- ✅ Both frontend and backend
- ✅ Automatic dependency installation
- ✅ Proper port mapping (3000, 8080)

### Method 4: Docker Only

Build and run the development image:

```bash
# Build development image
docker build -f Dockerfile.dev -t thaiyyal-dev .

# Run container
docker run -it --rm \
  -p 3000:3000 \
  -p 8080:8080 \
  -v $(pwd)/src:/app/src \
  -v $(pwd)/backend:/app/backend \
  -v $(pwd)/public:/app/public \
  thaiyyal-dev
```

## Development Workflow

### Making Frontend Changes

1. Edit files in `src/` directory
2. Save the file
3. Next.js automatically reloads the page
4. Changes appear instantly in the browser

Example:
```bash
# Edit a component
vim src/components/WorkflowCanvas.tsx

# Save and see changes immediately in browser
```

### Making Backend Changes

1. Edit files in `backend/` directory
2. Stop the Go server (Ctrl+C if running via dev.sh)
3. The script will handle restart, or manually restart:

```bash
# If using dev.sh, just Ctrl+C and run again
./dev.sh

# If running manually
cd backend/cmd/server
go run -tags dev . -addr :8080
```

### Making Both Frontend and Backend Changes

The dev.sh script manages both, so you can:
1. Make frontend changes → Instant HMR
2. Make backend changes → Restart dev.sh
3. All changes preserved, no manual coordination needed

## Testing Hot Reload

### Frontend Hot Reload Test

1. Start the dev environment:
   ```bash
   ./dev.sh
   ```

2. Open http://localhost:8080 in your browser

3. Edit a frontend file, e.g., `src/app/page.tsx`:
   ```tsx
   // Add a test message
   <h1>Hot Reload Test - Updated!</h1>
   ```

4. Save the file and observe the browser automatically refreshing

### Backend Proxy Test

1. Verify the backend is proxying to Next.js:
   ```bash
   curl http://localhost:8080
   # Should return the same content as http://localhost:3000
   ```

2. Test API routes still work:
   ```bash
   curl http://localhost:8080/health
   # Should return health check data
   ```

## Troubleshooting

### Port Already in Use

If port 3000 or 8080 is already in use:

```bash
# Find what's using the port
lsof -i :3000
lsof -i :8080

# Kill the process
kill -9 <PID>
```

Or change the ports in the dev script:
```bash
# Edit dev.sh and change the ports
# For Next.js: modify npm run dev to npm run dev -- -p 3001
# For backend: modify -addr :8080 to -addr :8081
```

### Next.js Dev Server Not Starting

Check the logs:
```bash
tail -f /tmp/thaiyyal-frontend.log
```

Common issues:
- Node modules not installed: Run `npm install`
- Port in use: Change the port or kill the process
- Syntax error in code: Check the error message in logs

### Backend Not Proxying

Check the logs:
```bash
tail -f /tmp/thaiyyal-backend.log
```

Common issues:
- Next.js not running: Make sure port 3000 is accessible
- Build tag missing: Make sure you're using `-tags dev`
- CORS issues: Backend has CORS enabled by default

### Docker Issues

If hot reloading doesn't work in Docker:

1. Check volume mounts:
   ```bash
   docker-compose -f docker-compose.dev.yml exec thaiyyal-dev ls -la /app/src
   ```

2. Verify ports are mapped:
   ```bash
   docker-compose -f docker-compose.dev.yml ps
   ```

3. Check logs:
   ```bash
   docker-compose -f docker-compose.dev.yml logs -f
   ```

## Production Build

When ready for production, build normally without dev tags:

```bash
# Build frontend
npm run build

# Build backend (production mode)
cd backend/cmd/server
go build -o ../../../server .
cd ../../..

# Run production server
./server
```

Or use the build script:
```bash
./build.sh
```

## Configuration

### Environment Variables

The dev script and docker-compose support these environment variables:

- `SERVER_ADDR`: Backend server address (default: `:8080`)
- `MAX_EXECUTION_TIME`: Max workflow execution time (default: `1m`)
- `MAX_NODE_EXECUTIONS`: Max nodes per workflow (default: `10000`)
- `MAX_HTTP_CALLS`: Max HTTP calls per execution (default: `100`)
- `MAX_LOOP_ITERATIONS`: Max loop iterations (default: `10000`)

### Customizing Next.js Dev Server

Edit `next.config.ts` to customize Next.js behavior:

```typescript
const nextConfig: NextConfig = {
  devIndicators: false,  // Hide dev indicators
  // Add more options...
};
```

### Customizing Backend Dev Mode

Edit `backend/pkg/server/embed_dev.go` to change proxy settings:

```go
// Change Next.js dev server URL
var nextJSDevURL = "http://localhost:3001"
```

## Tips and Best Practices

### Development

1. **Use the dev.sh script** for the best developer experience
2. **Keep both servers running** in a single terminal for easy management
3. **Check logs** in `/tmp/thaiyyal-*.log` if something goes wrong
4. **Use volume mounts** in Docker for instant file sync

### Performance

1. **Disable source maps** in production builds for faster builds
2. **Use incremental builds** in Next.js (automatic in dev mode)
3. **Don't run tests in dev mode** if performance is an issue

### Debugging

1. **Use browser DevTools** for frontend debugging
2. **Use Delve** for backend Go debugging:
   ```bash
   cd backend/cmd/server
   dlv debug -tags dev -- -addr :8080
   ```
3. **Check network tab** to see proxied requests

## Files Added/Modified

### New Files
- `backend/pkg/server/embed_dev.go` - Dev mode proxy implementation
- `dev.sh` - Development mode startup script
- `Dockerfile.dev` - Development Docker image
- `docker-compose.dev.yml` - Development compose configuration
- `DEV_HOT_RELOAD.md` - This documentation

### Modified Files
- `backend/pkg/server/embed.go` - Added build tag and helper functions
- `backend/pkg/server/routes_static.go` - Added dev mode proxy support
- `package.json` - Added `dev:full` script
- `.gitignore` - Added dev artifacts

## Further Reading

- [Next.js Fast Refresh](https://nextjs.org/docs/architecture/fast-refresh)
- [Go Build Tags](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Docker Volume Mounts](https://docs.docker.com/storage/volumes/)
- [Development Best Practices](DEV_GUIDE.md)
