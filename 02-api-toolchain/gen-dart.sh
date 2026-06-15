#!/usr/bin/env sh
# Dartクライアントを生成。openapi-generator が必要。
#   brew install openapi-generator     (mac)
# もしくは Docker 版を使う。
set -e
openapi-generator generate \
  -i openapi.yaml \
  -g dart-dio \
  -o ./client
echo "Dart client generated in ./client"
