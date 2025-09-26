#!/bin/bash

echo "=== Testing ChromaDB v2 API ==="
echo "Checking v2 heartbeat..."
curl -s -w "\nHTTP Status: %{http_code}\n" http://localhost:8000/api/v2/heartbeat

echo -e "\n=== Testing v2 Version ==="
curl -s -w "\nHTTP Status: %{http_code}\n" http://localhost:8000/api/v2/version

echo -e "\n=== Testing v2 Collections Endpoint ==="
curl -s -w "\nHTTP Status: %{http_code}\n" -X GET http://localhost:8000/api/v2/collections

echo -e "\n=== Testing Collection Creation ==="
curl -s -w "\nHTTP Status: %{http_code}\n" \
  -X POST http://localhost:8000/api/v2/collections \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test_collection",
    "metadata": {"description": "Test collection"},
    "get_or_create": true
  }'

echo -e "\n=== Listing Collections Again ==="
curl -s -w "\nHTTP Status: %{http_code}\n" -X GET http://localhost:8000/api/v2/collections