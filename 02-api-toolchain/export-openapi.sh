#!/usr/bin/env sh
# OpenAPI仕様を書き出す。go run . を別ターミナルで起動しておくこと。
set -e
# 既定は3.1。Dart生成で詰まったら下のURLを openapi-3.0.yaml に切り替える。
curl -s http://localhost:8888/openapi.yaml -o openapi.yaml
# curl -s http://localhost:8888/openapi-3.0.yaml -o openapi.yaml
echo "saved openapi.yaml"
