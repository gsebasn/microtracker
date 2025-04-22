#!/bin/bash

# Make 200 requests quickly
for i in {1..200}; do
    curl -s -w "%{http_code}\n" -o /dev/null http://localhost:8080/api/v1/packages &
done

# Wait for all requests to complete
wait 