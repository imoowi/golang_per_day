#!/bin/sh
docker build -f Dockerfile --progress=plain -t user-service --cpu-shares 2 --no-cache .
docker rm -f user-service || true
docker run -d -p 8000:8000 -v ./configs:/configs -v ./runtime:/runtime --name user-service user-service