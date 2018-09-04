# Use Go image for building of binary executable.
FROM golang:1.10 AS builder

# Build arg used to specify which cmd inside cmd/ to build for/use as entrypoint.
ARG ROLE

# Install dep so that dependencies can be installed.
RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
 	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Create classic GOPATH structure.
RUN mkdir -p /go/src/github.com/sweettea-io/rest-api

# Switch to project dir as new working dir.
WORKDIR /go/src/github.com/sweettea-io/rest-api

# Copy files needed by dep in order to install dependencies.
COPY Gopkg.toml Gopkg.lock ./

# Install dependencies.
RUN dep ensure -vendor-only

# Copy this project into working dir.
COPY . .

# Build Go binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -o main ./cmd/$ROLE

# Switch over to alpine base image (*much* lighter) for running of Go binary.
FROM alpine:latest

# Build arg used to specifiy which environment this image is for (i.e. dev, prod, etc.).
ARG BUILD_ENV

RUN apk --no-cache add ca-certificates

# Set working dir to /root inside alpine image.
WORKDIR /root

# Copy kubeconfig file for this BUILD_ENV.
COPY ./tmp/kubeconfigs/$BUILD_ENV ./.kubeconfig

# Copy Go binary built in first image over to this image.
COPY --from=builder /go/src/github.com/sweettea-io/rest-api/main ./main

# Execute Go binary.
ENTRYPOINT ["./main"]