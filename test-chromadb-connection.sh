#!/bin/bash

echo "=== Testing ChromaDB Connection ==="
echo "Checking if ChromaDB is running..."
curl -s -w "HTTP Status: %{http_code}\n" http://localhost:8000/api/v1/heartbeat

echo -e "\n=== Testing API Version ==="
curl -s -w "HTTP Status: %{http_code}\n" http://localhost:8000/api/v1/version

echo -e "\n=== Testing Collections Endpoint ==="
curl -s -w "HTTP Status: %{http_code}\n" -X GET http://localhost:8000/api/v1/collections

echo -e "\n=== Testing Root API ==="
curl -s -w "HTTP Status: %{http_code}\n" http://localhost:8000/

echo -e "\n=== Checking Docker Logs ==="
docker logs $(docker ps -q --filter "name=chroma") --tail 10