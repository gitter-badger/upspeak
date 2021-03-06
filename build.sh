#!/bin/bash

set -e

function cleanWebDist {
    rm -rf bin/web/
}

function cleanRebuild {
    rm -rf bin/
    mkdir -p bin/
}

function buildServerRelease {
    echo "Building upspeak server release"
    export CGO_ENABLED=0
    export GOOS=linux
    go build -o bin/upspeak -a -ldflags '-extldflags "-static"' .
    chmod +x bin/upspeak
}

function buildServerDev {
    echo "Building upspeak server dev"
    go build -o bin/upspeak .
    chmod +x bin/upspeak
}

function buildServerDockerDev {
    echo "Starting up with docker-compose"
    docker-compose up --build
}

function stopServerDockerDev {
    echo "Stopping with docker-compose"
    docker-compose down --remove-orphans --volumes
}

function buildWebClientDev {
    echo "Building web client dev"
    cleanWebDist
    cd web/
    npm run build-dev
    cd ..
}

function buildWebClientRelease {
    echo "Building web client release"
    cleanWebDist
    cd web/
    npm run build
    cd ..
}

function buildDev {
    cleanRebuild
    buildServerDev
}

function buildRelease {
    cleanRebuild
    buildServerRelease
}

case "$1" in
    "release")
    buildRelease
    ;;
    "dev")
    buildDev
    ;;
    "web")
    buildWebClientRelease
    ;;
    "web-dev")
    buildWebClientDev
    ;;
    "docker-start")
    buildServerDockerDev
    ;;
    "docker-stop")
    stopServerDockerDev
    ;;
    *)
    echo "Usage: ./build.sh [dev|release|web|web-dev|docker-start|docker-stop|help]"
    echo "Binaries are put in './bin/'"
    echo "Web app is packaged in './bin/web/'"
    ;;
esac
