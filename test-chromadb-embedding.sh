#!/bin/bash

echo "=== Testing ChromaDB Embedding Functions ==="

echo "1. Testing collection creation with default embedding..."
curl -s -w "\nStatus: %{http_code}\n" \
  -X POST http://localhost:8000/api/v2/tenants/default_tenant/databases/default_database/collections \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test_embedding_collection",
    "metadata": {"description": "Test collection with embeddings"},
    "get_or_create": true
  }'

echo -e "\n2. Testing document add without embeddings..."
curl -s -w "\nStatus: %{http_code}\n" \
  -X POST http://localhost:8000/api/v2/tenants/default_tenant/databases/default_database/collections/test_embedding_collection/add \
  -H "Content-Type: application/json" \
  -d '{
    "ids": ["test1"],
    "documents": ["This is a test document"],
    "metadatas": [{"source": "test"}]
  }'

echo -e "\n3. Testing document add with empty embeddings..."
curl -s -w "\nStatus: %{http_code}\n" \
  -X POST http://localhost:8000/api/v2/tenants/default_tenant/databases/default_database/collections/test_embedding_collection/add \
  -H "Content-Type: application/json" \
  -d '{
    "ids": ["test2"],
    "documents": ["This is another test document"],
    "metadatas": [{"source": "test"}],
    "embeddings": null
  }'

echo -e "\n4. Testing collection creation with explicit embedding function..."
curl -s -w "\nStatus: %{http_code}\n" \
  -X POST http://localhost:8000/api/v2/tenants/default_tenant/databases/default_database/collections \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test_embedding_collection2",
    "metadata": {"description": "Test collection with explicit embedding function"},
    "embedding_function": "default",
    "get_or_create": true
  }'