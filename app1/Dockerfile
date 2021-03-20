FROM golang:1.13-buster as builder

WORKDIR /usr/src/app
COPY . .

# go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags '-extldflags=-static -s -w' -o main main.go

# grpcurl
# RUN go get -u github.com/fullstorydev/grpcurl
# RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl

# ### Make executable image
# FROM scratch
FROM golang:1.13-buster

COPY --from=builder /usr/src/app/main /main

# RUN apt-get update \
# &&  apt-get install -y net-tools dnsutils \
# &&  apt-get autoremove -y \
# &&  apt-get clean \
# &&  rm -rf /tmp/* /var/tmp/* \
# &&  rm -rf /var/lib/apt/lists/*
# COPY --from=builder /go/bin/grpcurl   /usr/local/bin/grpcurl
