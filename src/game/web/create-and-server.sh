#!/bin/sh

set -e

make clean all
main -c /tmp/test.pile $1 8000 5000
main /tmp/test.pile

