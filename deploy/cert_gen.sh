#!/bin/bash

set -e
source ./deploy/.env
set +a


echo "Generating self-signed certificate for $GRAFANA_DOMAIN..."
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "./fazenda.key" \
    -out "./fazenda.crt" \
    -addext "subjectAltName=IP:$GRAFANA_DOMAIN" \
    -subj "/CN=localhost" < /dev/null

if [ -f "./fazenda.key" ]; then
    echo "Copying key to service directories..."
    cp "./fazenda.key" "./deploy/grafana/grafana.key"
    cp "./fazenda.key" "./deploy/keycloak/keycloak.key"
    cp "./fazenda.key" "./traffic_guard-front/certs/fazenda.key"
    
    chmod +r "./deploy/grafana/grafana.key"
    chmod +r "./deploy/keycloak/keycloak.key"
    chmod +r "./traffic_guard-front/certs/fazenda.key"
else
    echo "Error: Failed to generate key file "
    exit 1
fi

if [ -f "./fazenda.crt" ]; then
    #chmod 644 "$CRT_FILE"
    echo "Copying certificate to service directories..."
    cp "./fazenda.crt" "./deploy/grafana/grafana.crt"
    cp "./fazenda.crt" "./deploy/keycloak/fullchain.crt"
    cp "./fazenda.crt" "./traffic_guard-front/certs/fazenda.crt"
    

    chmod +r "./deploy/grafana/grafana.crt"
    chmod +r "./deploy/keycloak/fullchain.crt"
    chmod +r "./traffic_guard-front/certs/fazenda.crt"

else
    echo "Error: Failed to generate certificate"
    exit 1
fi

echo "Certificate generation and distribution completed successfully."
