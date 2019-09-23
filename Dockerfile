# Stage 1
FROM golang:1.13.0-alpine3.10 as builder

# Add git
RUN apk update && \
    apk add git && \
    apk add openssl-dev && \
    apk add gcc && \
    apk add libc-dev

RUN mkdir $GOPATH/src/gitlab.com
RUN mkdir $GOPATH/src/gitlab.com/go-pher
RUN mkdir $GOPATH/src/gitlab.com/go-pher/go-auth

ADD . $GOPATH/src/gitlab.com/go-pher/go-auth/
#RUN git clone https://oauth2:$BUILD_TOKEN@gitlab.com/go-pher/go-auth.git $GOPATH/src/gitlab.com/go-pher/go-auth

WORKDIR $GOPATH/src/gitlab.com/go-pher/go-auth

#RUN echo $GOPATH

RUN go get ./

RUN go build

# Stage 2

FROM alpine:3.10

RUN apk update && \
    apk add openssl-dev && \
    apk add ca-certificates

COPY --from=builder /go/bin/go-auth /

EXPOSE 60061

CMD ["./go-auth"]
