#!/bin/bash -e

# Install NSQ, Run it up, and run integration test for it

[ "x$TEMPDIR" == "x" ] && TEMPDIR=".tmp"
# default to support mac
[ "x$NSQ_DOWNLOAD" == "x" ] && NSQ_DOWNLOAD=nsq-0.3.5.darwin-amd64.go1.4.2

cleanup() {
 	kill $(jobs -p) 2>/dev/null
}
trap cleanup INT TERM EXIT

[ -d "$TEMPDIR" ] || mkdir $TEMPDIR

pushd $TEMPDIR

[ -f "$NSQ_DOWNLOAD.tar.gz" ] || wget http://bitly-downloads.s3.amazonaws.com/nsq/$NSQ_DOWNLOAD.tar.gz -O $NSQ_DOWNLOAD.tar.gz

[ -d "$NSQ_DOWNLOAD" ] || tar zxvf $NSQ_DOWNLOAD.tar.gz

export PATH=$(pwd)/$NSQ_DOWNLOAD/bin:$PATH

# Check NSQ binaries

if ! which nsqd >/dev/null; then
    echo "missing nsqd binary" && exit 1
fi

if ! which nsqlookupd >/dev/null; then
    echo "missing nsqlookupd binary" && exit 1
fi

# run nsqlookupd
LOOKUP_LOGFILE=$(mktemp -t nsqlookupd-log.XXXXXXX)
echo "starting nsqlookupd"
nsqlookupd >$LOOKUP_LOGFILE 2>&1 &

# run nsqd
NSQD_LOGFILE=$(mktemp -t nsqd-log.XXXXXXX)
echo "starting nsqd"
rm -f *.dat
nsqd --lookupd-tcp-address=127.0.0.1:4160 >$NSQD_LOGFILE 2>&1 &

is_nsqd_up() {
	curl -I http://127.0.0.1:4151/ping 2>&1 >/dev/null
	[ "x$?" == "x0" ] && return 0 || return 1
}

sleep 0.5
is_nsqd_up || \
	(echo "NSQD Hasn't started. Consider increasing sleep time" && exit 1)

popd # TEMPDIR

# Integration test
godep go test -v ./test/integration/nsq

echo "Succuss!"
