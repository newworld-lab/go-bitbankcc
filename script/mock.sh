#!/bin/sh

for path in `go list ./... | egrep -v 'vendor|_test|_mock'`; do
  for gofile in $( ls ${GOPATH}/src/${path} | grep .go$ ); do
    package=${path##*/}
    if egrep -q "type.*interface" ${GOPATH}/src/${path}/${gofile}; then
      filename=${gofile%.*}
      from=${GOPATH}/src/${path}/${filename}.go
      to=${GOPATH}/src/${path}/${filename}_mock.go
      mockgen -source ${from} -destination ${to} -package ${package}
      sed -i "" "2d" ${to}
      echo "${to}"
    fi
  done
done