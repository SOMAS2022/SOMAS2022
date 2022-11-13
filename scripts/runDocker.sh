#!/usr/bin/env bash

set -e

IMAGE_NAME="somas2022"

# Print working directory
echo "Current working directory is: $(pwd)"

# Check if Dockerfile exists in the working directory, abort if not
if [ ! -f "Dockerfile" ]; then
	echo "Dockerfile does not exist."
	exit 1
fi

# Build image
docker build --tag "${IMAGE_NAME}" .

# Run the built image
docker run -t \
	--rm \
	--mount type=bind,source="$(pwd)"/.env,target=/somas/.env \
	"${IMAGE_NAME}"
