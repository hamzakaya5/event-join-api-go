#!/bin/bash

sleep 5

BACKEND_UPSTREAM=""
for container in $(getent hosts backend | awk '{print $1}'); do
    BACKEND_UPSTREAM+="server ${container}:8080;\n"
done

sed "s|{{BACKEND_SERVERS}}|$BACKEND_UPSTREAM|" /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/nginx.conf

nginx -g 'daemon off;'
