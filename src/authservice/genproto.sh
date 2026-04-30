#!/bin/bash -e

# python -m grpc_tools.protoc -I../../pb \
#   --python_out=./authservice \
#   --grpc_python_out=./authservice \
#   ../../pb/chirp.proto
python -m grpc_tools.protoc -I../../pb \
  --python_out=. \
  --grpc_python_out=. \
  ../../pb/chirp.proto
