#!/usr/bin/bash

arch=(amd64 arm64)
os=(linux windows darwin)

for a in ${arch[@]}
do
    for o in ${os[@]}
    do
        env GOOS=${o} GOARCH=${a} go build -o build/mdn_${a}-${o}
    done
done
