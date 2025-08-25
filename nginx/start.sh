#!/bin/bash

# Wait for backend containers to start
sleep 5

# Generate dynamic upstream
BACKEND_UPSTREAM=""
for container in $(getent hosts backend | awk '{print $1}'); do
    BACKEND_UPSTREAM+="server ${container}:8080;\n"
done

# Replace template variable
sed "s|{{BACKEND_SERVERS}}|$BACKEND_UPSTREAM|" /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/nginx.conf

# Start Nginx
nginx -g 'daemon off;'
