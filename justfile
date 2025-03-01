set shell := ["bash", "-uc"]

# Show available targets/recipes
default:
    @just --choose

# Clean up old files
clean:
    rm -rf ./dist/*
    rm ./rolecule

# Build the binary for the current os/arch
build:
    go build -o bin/rolecule

# Configure your host to use this repo
setup:
    direnv allow

# Show git tags
tags:
    @git tag | sort -V

# Run unit tests
test:
    go test ./...

# Build docker images with ansible support
build-docker-ansible-images:
  docker build -t rockylinux-systemd:9.1 -f testing/ansible/rockylinux-9.1-systemd.Dockerfile .
  docker build -t ubuntu-systemd:22.04 -f testing/ansible/ubuntu-22.04-systemd.Dockerfile .

# Build podman images with ansible support
build-podman-ansible-images:
  podman build -t rockylinux-systemd:9.1 -f testing/ansible/rockylinux-9.1-systemd.Containerfile .
  podman build -t ubuntu-systemd:22.04 -f testing/ansible/ubuntu-22.04-systemd.Containerfile .

# Build all images with ansible support
build-ansible-images: build-docker-ansible-images build-podman-ansible-images

# Show help menu
help:
    @just --list --list-prefix '  ❯ '
