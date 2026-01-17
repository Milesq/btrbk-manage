#!/usr/bin/env bash
set -euo pipefail

head -c 10 /dev/random | base64 > ./mnt/@/restoration_complete.txt
