#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Testing rate limits for different endpoints..."

# Test List Packages endpoint (200 requests/minute)
echo -e "\n${GREEN}Testing List Packages endpoint (200 requests/minute)${NC}"
for i in {1..210}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null http://localhost:8080/api/v1/packages)
    if [ $response -eq 429 ]; then
        echo -e "${RED}Rate limit hit at request $i${NC}"
        break
    fi
done

# Wait for rate limit to reset
echo "Waiting for rate limit to reset..."
sleep 2

# Test Search Packages endpoint (150 requests/minute)
echo -e "\n${GREEN}Testing Search Packages endpoint (150 requests/minute)${NC}"
for i in {1..160}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null "http://localhost:8080/api/v1/packages/search?query=test")
    if [ $response -eq 429 ]; then
        echo -e "${RED}Rate limit hit at request $i${NC}"
        break
    fi
done

# Wait for rate limit to reset
echo "Waiting for rate limit to reset..."
sleep 2

# Test Create Package endpoint (50 requests/minute)
echo -e "\n${GREEN}Testing Create Package endpoint (50 requests/minute)${NC}"
for i in {1..60}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null -X POST http://localhost:8080/api/v1/packages \
        -H "Content-Type: application/json" \
        -d '{"packageId":"TEST'$i'","sender":{"name":"Test Sender","address":"123 Test St"},"recipient":{"name":"Test Recipient","address":"456 Test St"},"origin":"Test Origin","destination":"Test Destination","currentStatus":"created"}')
    if [ $response -eq 429 ]; then
        echo -e "${RED}Rate limit hit at request $i${NC}"
        break
    fi
done

echo -e "\n${GREEN}Rate limit testing completed${NC}" 