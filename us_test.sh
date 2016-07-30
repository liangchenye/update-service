#!/bin/sh
set -e

echo "start test, killing the exist 'upserver' server"
pidof upserver | xargs kill -9

TMPDIR=$(mktemp -d)
TMPSTORAGE=${TMPDIR}/storage
TMPKM=${TMPDIR}/km
echo "creating tmp storage dir: " $TMPSTORAGE
mkdir -p $TMPSTORAGE
echo "creating tmp keymanager dir: " $TMPKM
mkdir -p $TMPKM


echo "start to compile server"
cd server
make

echo "start the update server"
./upserver web --storage "local:/""$TMPSTORAGE" --keymanager "local:/""$TMPKM"  &

echo "start to compile client"
cd ../client
make
cd ..

echo "\nset enviornment and start to run tests"
export US_TEST_SERVER="appV1://localhost:1234"
echo "---------------------------------------------"
go test -v $(go list ./... | grep -v /vendor/)

echo "\nstart to run client command line"
echo "---------------------------------------------"
cd client
./upclient push README.md "appV1://localhost:1234/citest/official"
./upclient pull README.md "appV1://localhost:1234/citest/official"

echo "\n---------------------------------------------"
echo "killing the testing 'upserver' server"
killall upserver

echo "clean all the generated data"
rm -fr $TMPDIR
rm -fr ~/.dockyard/cache/citest

echo "end of test"
