# Stage 1
FROM golang:1.13.0-alpine3.10 as builder

# Add git
RUN apk update && \
    apk add git && \
    apk add openssl-dev && \
    apk add gcc && \
    apk add libc-dev

RUN mkdir $GOPATH/src/gitlab.com
RUN mkdir $GOPATH/src/gitlab.com/YOVO-LABS
RUN mkdir $GOPATH/src/github.com/YOVO-LABS/auth

ADD . $GOPATH/src/github.com/YOVO-LABS/auth/
#RUN git clone https://oauth2:$BUILD_TOKEN@github.com/YOVO-LABS/auth.git $GOPATH/src/github.com/YOVO-LABS/auth

WORKDIR $GOPATH/src/github.com/YOVO-LABS/auth

#RUN echo $GOPATH

RUN go get ./

RUN go build

# Stage 2

FROM alpine:3.10

RUN apk update && \
    apk add openssl-dev && \
    apk add ca-certificates

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder /go/bin/go-auth /

EXPOSE 60061

CMD ["./go-auth"]
