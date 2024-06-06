#!/bin/bash

set +e

echo "Building default binary"

CGO_ENABLED=1 go build -ldflags "-s -w" -ldflags "-X github.com/pharmaeasy/tokenizer/cmd.version=${VERSION}" -o "build/tokenizer" $PKG_SRC
./build/tokenizer start --config=tokenizer-dev.yml &
pid=$!

lastupdate=$(find . -not -path '*/\.*' | xargs stat | sed 's/"/_/g' | awk -F '_' '{print $4 $6 $8}' | md5)
while true; do
  lastafter=$(find . -not -path '*/\.*' | xargs stat | sed 's/"/_/g' | awk -F '_' '{print $4 $6 $8}' | md5)
  echo $lastupdate
  echo $lastafter
  if [ "$lastafter" != "$lastupdate" ]; then
    lastupdate=$lastafter
    if [ "$pid" != "" ]; then
      kill -kill $pid
    fi
    rm -f build/tokenizer*
    echo "Building default binary"
    CGO_ENABLED=1 go build -ldflags "-s -w" -ldflags "-X github.com/pharmaeasy/tokenizer/cmd.version=${VERSION}" -o "build/tokenizer" $PKG_SRC
    ./build/tokenizer start --config=tokenizer-dev.yml &
    pid=$!
    lastupdate=$(find . -not -path '*/\.*' | xargs stat | sed 's/"/_/g' | awk -F '_' '{print $4 $6 $8}' | md5)
  fi
  sleep 5

done
