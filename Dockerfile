FROM golang:1.10 AS builder

# Install dep
RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    wget -O dep.zip https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64.zip && \
    echo '96c191251164b1404332793fb7d1e5d8de2641706b128bf8d65772363758f364  dep.zip' | sha256sum -c - && \
    unzip -d /usr/bin dep.zip && rm dep.zip

# Create classic GOPATH structure & set that as our working dir
RUN mkdir -p /go/src/github.com/***
WORKDIR /go/src/github.com/***

# Copy our dep files, specifying our dependencies
COPY Gopkg.toml Gopkg.lock ./

# Install dependencies
RUN dep ensure -vendor-only
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o ***

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Switch to /root/
WORKDIR /root/

COPY --from=builder /go/src/github.com/***  .

ENTRYPOINT ["./***"]