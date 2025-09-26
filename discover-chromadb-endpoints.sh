#!/bin/bash

echo "=== Discovering ChromaDB v2 API Endpoints ==="

echo "1. Testing basic endpoints..."
curl -s -w "\nStatus: %{http_code}\n" http://localhost:8000/api/v2/heartbeat
echo ""

echo "2. Testing different collection endpoint formats..."
for endpoint in "collections" "collection" "databases" "database"; do
    echo "Testing /api/v2/$endpoint..."
    curl -s -w "Status: %{http_code}\n" http://localhost:8000/api/v2/$endpoint
    echo ""
done

echo "3. Testing root v2 API..."
curl -s -w "Status: %{http_code}\n" http://localhost:8000/api/v2/
echo ""

echo "4. Testing OpenAPI spec..."
curl -s -w "Status: %{http_code}\n" http://localhost:8000/docs
echo ""

echo "5. Testing Swagger UI..."
curl -s -w "Status: %{http_code}\n" http://localhost:8000/openapi.json
echo ""