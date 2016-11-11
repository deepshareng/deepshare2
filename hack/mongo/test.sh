#!/bin/bash -e
#
# This script assumes that mongod is running in background!

if ! which mongo >/dev/null; then
    echo "missing mongo binary" && exit 1
fi

is_mongo_up() {
	mongo mydb_test --eval 'db.stats()' 2>&1 >/dev/null
	[ "x$?" == "x0" ] && return 0 || return 1
}

MONGO_UP=false
for var1 in 1 2 3
do
	is_mongo_up && MONGO_UP=true && break
	sleep 5
done
$MONGO_UP && echo "mongodb is up and running..." \
  || (echo "MongoDB isn't running!" && exit 1)

# Integration test
godep go test -v ./test/integration/mongo -timeout 90s

echo "Succuss!"
