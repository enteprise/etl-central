---
description: >-
  This Dockerfile is a multi-stage build file that creates an image with a Go
  backend and a Node.js frontend, and it sets up the environment to run the
  application in an Alpine Linux container. Here is
---

# main app dockerfile



**First Stage: Building the Go Application**

1. `FROM golang:1.21-alpine as builder`
   * This line specifies the base image (`golang:1.21-alpine`) for the first build stage, and it names the stage `builder`. This image includes the Go programming language environment based on Alpine Linux.
2. `RUN mkdir -p /go/src/build`
   * This command creates a directory (`/go/src/build`) in the container filesystem where the Go application will be built.
3. `WORKDIR /go/src/build`
   * Sets the working directory for subsequent commands to `/go/src/build`.
4. `COPY go.mod /go/src/build/go.mod` `COPY go.sum /go/src/build/go.sum`
   * These commands copy the Go module files (`go.mod` and `go.sum`) from the host to the container's working directory.
5. `RUN go mod download`
   * This command downloads the Go module dependencies specified in `go.mod` and `go.sum`.
6. `ADD app /go/src/build/app`
   * Adds the application source code from the `app` directory on the host to the container's working directory.
7. `ARG DATAPLANE_VERSION=latest`
   * This line declares a build argument `DATAPLANE_VERSION` with a default value of `latest`.
8. `RUN CGO_ENABLED=0 go build -ldflags "-X github.com/dataplane-app/dataplane/app/mainapp/config.Version=$DATAPLANE_VERSION" -o dataplane app/mainapp/server.go`
   * This command compiles the Go application. It disables CGO, sets the `Version` variable in the `config` package to the value of `DATAPLANE_VERSION`, and outputs the binary as `dataplane`.

**Second Stage: Building the Frontend**

9. `FROM node:18.17 as builder2`
   * Specifies the base image for the second build stage (`node:18.17`) and names the stage `builder2`. This image includes the Node.js environment.
10. `RUN mkdir -p /node` `RUN mkdir -p /node/app/mainapp/frontbuild`
    * These commands create directories in the container filesystem for the Node.js application.
11. `WORKDIR /node`
    * Sets the working directory for subsequent commands to `/node`.
12. `ADD frontend/public /node/public` `ADD frontend/src /node/src` `ADD frontend/vite.config.js /node/vite.config.js` `ADD frontend/package.json /node/package.json` `ADD frontend/index.html /node/index.html`
    * These commands add the frontend source code and configuration files from the host to the container's working directory.
13. `RUN yarn add global env-cmd`
    * Installs the `env-cmd` package globally using Yarn.
14. `RUN yarn`
    * Installs the Node.js dependencies specified in `package.json`.
15. `RUN yarn builddocker`
    * Runs the `builddocker` script defined in `package.json`, which typically builds the frontend for production.

**Final Stage: Assembling the Final Image**

16. `FROM alpine:3.18`
    * Specifies the base image for the final stage, which is a minimal Alpine Linux (`alpine:3.18`).
17. `ENV TZ=UTC`
    * Sets the timezone environment variable to UTC.
18. `RUN apk update && apk add --no-cache tzdata htop`
    * Updates the package index and installs the `tzdata` (timezone data) and `htop` (an interactive process viewer) packages without caching them to minimize the image size.
19. `RUN rm -rf /var/cache/apk/*`
    * Removes the package manager cache to reduce the image size.
20. `ENV USER=appuser` `ENV UID=10001`
    * Sets environment variables for the user and user ID.
21. `RUN adduser ...`
    * Adds a new user (`appuser`) with the specified user ID (`10001`) and configures the user with no login shell and no home directory.
22. `COPY --from=builder2 /node/app/mainapp/frontbuild /dataplane/frontbuild` `COPY --from=builder go/src/build/dataplane /dataplane/dataplane`
    * Copies the frontend build artifacts and the Go binary from their respective build stages into the final image.
23. `RUN mkdir /dataplane/code-files/ && chown -R appuser:appuser /dataplane` `RUN chmod +w /dataplane/code-files/` `RUN mkdir /dataplane/dfs-code-files/ && chown -R appuser:appuser /dataplane` `RUN chmod +w /dataplane/dfs-code-files/`
    * Creates directories for code files, changes their ownership to `appuser`, and grants write permissions.
24. `WORKDIR /dataplane`
    * Sets the working directory for the container to `/dataplane`.
25. `USER appuser:appuser`
    * Sets the user and group for the container to `appuser`.
26. `CMD ["./dataplane"]`
    * Specifies the command to run when the container starts, which is to execute the `dataplane` binary.

When this Dockerfile is used to build an image and a container is started from that image, it will run the `dataplane` application with the frontend assets served from the `/dataplane/frontbuild` directory.
