# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.1.2 AS xx
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS base

COPY --from=goreleaser-xx / /
COPY --from=xx / /
RUN apk add --no-cache clang git lld llvm

WORKDIR /src 