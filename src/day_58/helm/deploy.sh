#!/bin/bash
helm upgrade --install nats nats/nats \
  --set nats.jetstream.enabled=true \
  --set cluster.enabled=true \
  --set replicas=3