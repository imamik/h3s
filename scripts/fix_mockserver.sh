#!/bin/bash
set -e

# This script fixes all the w.Write calls in mockserver.go to use safeWrite

# Replace all w.Write calls with safeWrite
sed -i 's/w\.Write(/safeWrite(w, /g' internal/hetzner/mockhetzner/mockserver.go

echo "Fixed all w.Write calls in mockserver.go"
