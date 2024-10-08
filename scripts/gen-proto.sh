#!/bin/sh

CURRENT_DIR=$(pwd)

rm -rf ${CURRENT_DIR}/genprotos
mkdir -p ${CURRENT_DIR}/genprotos

for x in $(find ${CURRENT_DIR}/gym_protos* -type d ! -path "*/.git*" ! -path "*/.git/*"); do
  if ls ${x}/*.proto 1> /dev/null 2>&1; then
    protoc -I=${x} -I/usr/local/include \
      --go_out=${CURRENT_DIR}/genprotos --go_opt=paths=source_relative \
      --go-grpc_out=${CURRENT_DIR}/genprotos --go-grpc_opt=paths=source_relative ${x}/*.proto
  fi
done
