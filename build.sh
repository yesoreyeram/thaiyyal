#!/bin/bash
# Build script for Thaiyyal workflow engine
# This script builds the frontend and backend

set -e

echo "🏗️  Building Thaiyyal Workflow Engine"
echo "=================================="
echo

# Step 1: Build frontend and copy to backend static directory
echo "📦 Building frontend and copying to backend..."
npm run build
echo "✅ Frontend build complete and files copied to backend/pkg/server/static"
echo

# Step 2: Build backend
echo "🔨 Building backend..."
cd backend/cmd/server
go build -o ../../../thaiyyal-server .
cd ../../..
echo "✅ Backend build complete"
echo

echo "🎉 Build complete!"
echo "   Run './thaiyyal-server' to start the application."
echo "   Server will be available at http://localhost:8080"
