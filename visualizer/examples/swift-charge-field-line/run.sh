#!/bin/bash

# Check if first argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <config-name> [additional-go-args]"
    echo "Example: $0 1-harmonic-3d-slow"
    exit 1
fi

CONFIG_NAME="$1"
DATA_DIR="./data/${CONFIG_NAME}"

# Create output directory if it doesn't exist
mkdir -p "${DATA_DIR}"

# Run the go program with optional additional arguments
echo "Running go program with config: ${CONFIG_NAME}"
if [ -n "$2" ]; then
    go run *.go --out-dir "${DATA_DIR}" --config-name "${CONFIG_NAME}" "${@:2}"
else
    go run *.go --out-dir "${DATA_DIR}" --config-name "${CONFIG_NAME}"
fi

# Check if go command was successful
if [ $? -ne 0 ]; then
    echo "Error: go command failed"
    exit 1
fi

# Convert frames to video
echo "Converting frames to video..."
ffmpeg -framerate 20.75 -i "${DATA_DIR}/frame-%04d.png" -c:v libx264 -profile:v high -crf 10 -pix_fmt yuv420p -y "${DATA_DIR}/swift.mp4"

if [ $? -eq 0 ]; then
    echo "Video created successfully at ${DATA_DIR}/swift.mp4"
else
    echo "Error: ffmpeg command failed"
    exit 1
fi
