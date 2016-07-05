#!/bin/sh

if [ x"$#" != x"1" ]; then
  echo $0 destination_path
  exit 0
fi

DIST=$1

rm -fr "$DIST"
docker build -t dcfg . && docker run -v "$DIST":/dist:rw --rm dcfg

