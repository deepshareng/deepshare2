#!/bin/bash -e

# go test
godep go test -v $(go list ./... | \
              grep -v "/vendor" | grep -v "/Godeps" | grep -v "/test/integration" | grep -v "/test/testapi")

echo "Success!"
