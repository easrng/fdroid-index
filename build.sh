#!/bin/bash
set -euo pipefail
mkdir -p state
if [ ! -f "state/stork" ]; then
  echo "Downloading stork v1.6.0"
  wget "https://github.com/jameslittle230/stork/releases/download/v1.6.0/stork-ubuntu-20-04" -O state/stork
  if [ "$(sha256sum state/stork | awk '{print $1}')" != "cce634509194552f63eb4bc7615c65d73930eb029fbd270015facb7fcfc3ff50" ]; then
    echo "Hash mismatch, exiting."
    exit 1
  fi
  chmod +x state/stork
fi
if [ ! -f "fdroid-index" ]; then
  echo "Building fdroid-index"
  go build
fi
echo "Downloading repo index"
./fdroid-index https://f-droid.org/repo?fingerprint=43238D512C1E5EB2D6569F4A3AFBF5523418B82E0A3ED1552770ABB9A9C9CCAB >index.toml
THISHASH="$(sha256sum index.toml | awk '{print $1}')"
if [ "$THISHASH" == "$(cat state/lasthash 2>/dev/null || true)" ]; then
  echo "No changes since last build"
else
  echo "Building search index"
  state/stork build --input index.toml --output - | zstd --ultra -20 - >index.zst
  printf "%s" "$THISHASH" >state/lasthash
  echo "Built search index"
fi
rm index.toml