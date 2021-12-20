#!/bin/bash -eux

pushd dp-frontend-filter-flex-dataset
  make build
  cp build/dp-frontend-filter-flex-dataset Dockerfile.concourse ../build
popd
