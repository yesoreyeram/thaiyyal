#!/bin/bash
# Development mode script for Thaiyyal workflow engine
# This script starts both the Next.js dev server and the Go backend in dev mode

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Thaiyyal in Development Mode${NC}"
echo "=========================================="
echo

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${YELLOW}⚠️  Node.js is not installed. Please install Node.js 20.x or later.${NC}"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}⚠️  Go is not installed. Please install Go 1.24.7 or later.${NC}"
    exit 1
fi

# Install frontend dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo -e "${GREEN}📦 Installing frontend dependencies...${NC}"
    npm install
    echo
fi

# Function to cleanup background processes
cleanup() {
    echo
    echo -e "${YELLOW}🛑 Shutting down development servers...${NC}"
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    echo -e "${GREEN}✅ Cleanup complete${NC}"
    exit 0
}

# Register cleanup function
trap cleanup SIGINT SIGTERM EXIT

# Start Next.js dev server in background
echo -e "${GREEN}🎨 Starting Next.js dev server on http://localhost:3000${NC}"
npm run dev > /tmp/thaiyyal-frontend.log 2>&1 &
FRONTEND_PID=$!
echo "Frontend PID: $FRONTEND_PID"
echo

# Wait for Next.js to be ready (check for port to be listening)
echo -e "${YELLOW}⏳ Waiting for Next.js dev server to be ready...${NC}"
for i in {1..30}; do
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Next.js dev server is ready${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${YELLOW}⚠️  Next.js dev server is taking longer than expected to start${NC}"
        echo -e "${YELLOW}   Check /tmp/thaiyyal-frontend.log for details${NC}"
    fi
    sleep 1
done
echo

# Build and start Go backend in dev mode
echo -e "${GREEN}🔨 Building Go backend in development mode...${NC}"
cd backend/cmd/server
go build -tags dev -o ../../../server-dev .
cd ../../..
echo -e "${GREEN}✅ Backend build complete${NC}"
echo

echo -e "${GREEN}🚀 Starting Go backend on http://localhost:8080${NC}"
echo -e "${BLUE}📝 API requests to http://localhost:8080 will proxy to Next.js dev server${NC}"
echo
./server-dev -addr :8080 > /tmp/thaiyyal-backend.log 2>&1 &
BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"
echo

# Print helpful information
echo -e "${GREEN}✨ Development servers are running!${NC}"
echo "=========================================="
echo
echo -e "${BLUE}Frontend (Next.js):${NC}  http://localhost:3000"
echo -e "${BLUE}Backend (Go):${NC}        http://localhost:8080"
echo -e "${BLUE}API Endpoint:${NC}        http://localhost:8080/api/v1/workflow/execute"
echo -e "${BLUE}Health Check:${NC}        http://localhost:8080/health"
echo -e "${BLUE}Metrics:${NC}             http://localhost:8080/metrics"
echo
echo -e "${GREEN}Features:${NC}"
echo "  ✅ Hot module reloading for frontend (automatic)"
echo "  ✅ Backend proxies frontend requests to Next.js dev server"
echo "  ✅ API routes available on :8080"
echo
echo -e "${YELLOW}Logs:${NC}"
echo "  Frontend: tail -f /tmp/thaiyyal-frontend.log"
echo "  Backend:  tail -f /tmp/thaiyyal-backend.log"
echo
echo -e "${YELLOW}Press Ctrl+C to stop all servers${NC}"
echo

# Wait for either process to exit
wait $FRONTEND_PID $BACKEND_PID
