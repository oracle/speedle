#!/bin/bash

docker run -d \
  -p 2379:2379 \
  -p 2380:2380 \
  --name speedle_etcd quay.io/coreos/etcd:latest \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://0.0.0.0:2380 


