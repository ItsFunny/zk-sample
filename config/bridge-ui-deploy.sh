#!/bin/sh

# Build app
cd /app && npm run build

# Copy nginx config
cp /app/deployment/nginx.conf /etc/nginx/conf.d/default.conf

# Copy app dist
cp -r /app/dist/. /usr/share/nginx/html

# Delete source code
# rm -rf /app

# Run nginx
nginx -g 'daemon off;'
