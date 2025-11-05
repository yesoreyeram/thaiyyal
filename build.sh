#!/bin/bash
# Build script for Thaiyyal workflow engine
# This script builds the frontend and copies files to backend static directory

set -e

echo "🏗️  Building Thaiyyal Workflow Engine"
echo "=================================="
echo

# Step 1: Build frontend
echo "📦 Building frontend..."
npm run build
echo "✅ Frontend build complete"
echo

# Step 2: Prepare backend static directory
echo "📂 Preparing backend static directory..."
mkdir -p backend/pkg/server/static

# Step 3: Copy frontend build output
echo "📋 Copying frontend files to backend..."
cp -r .next/standalone/.next/server/app/*.html backend/pkg/server/static/ 2>/dev/null || true
cp -r .next/static backend/pkg/server/static/_next 2>/dev/null || true
cp -r public/* backend/pkg/server/static/ 2>/dev/null || true
echo "✅ Frontend files copied"
echo

# Step 4: Build backend
echo "🔨 Building backend..."
cd backend/cmd/server
go build -o ../../../server .
cd ../../..
echo "✅ Backend build complete"
echo

echo "🎉 Build complete! Run './server' to start the application."
