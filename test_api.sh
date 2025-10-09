#!/bin/bash

# Test script for the Prodyo Backend API
echo "Testing Prodyo Backend API..."

BASE_URL="http://localhost:8081/api/v1"

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/../health" | jq .

echo -e "\n2. Testing User endpoints..."

# Create a user
echo "Creating user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }')
echo $USER_RESPONSE | jq .
USER_ID=$(echo $USER_RESPONSE | jq -r '.id')

# Get all users
echo -e "\nGetting all users..."
curl -s "$BASE_URL/users" | jq .

# Get user by ID
echo -e "\nGetting user by ID..."
curl -s "$BASE_URL/users/$USER_ID" | jq .

echo -e "\n3. Testing Project endpoints..."

# Create a project
echo "Creating project..."
PROJECT_RESPONSE=$(curl -s -X POST "$BASE_URL/projects" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Test Project\",
    \"description\": \"A test project\",
    \"color\": \"#FF5733\",
    \"prod_range\": {
      \"ok\": 80,
      \"alert\": 60,
      \"critical\": 40
    },
    \"member_ids\": [\"$USER_ID\"]
  }")
echo $PROJECT_RESPONSE | jq .
PROJECT_ID=$(echo $PROJECT_RESPONSE | jq -r '.id')

# Get all projects (now with pagination)
echo -e "\nGetting all projects (paginated)..."
curl -s "$BASE_URL/projects" | jq .

# Get projects with pagination
echo -e "\nGetting projects with pagination..."
curl -s "$BASE_URL/projects?page=1&page_size=5" | jq .

# Get all users (now with pagination)
echo -e "\nGetting all users (paginated)..."
curl -s "$BASE_URL/users" | jq .

echo -e "\nAPI testing completed!"
