


FROM alpine:3.5  AS builder
RUN apk add --no-cache ca-certificates bash
WORKDIR /go/src/github.com/alllomancer/k8s-slackbot
ADD . /go/src/github.com/alllomancer/k8s-slackbot
RUN cd /go/src/github.com/alllomancer/k8s-slackbot

RUN  /go/src/github.com/alllomancer/k8s-slackbot/build/build-go.sh 
RUN  /go/src/github.com/alllomancer/k8s-slackbot/build/build.sh 
RUN /go/src/github.com/alllomancer/k8s-slackbot/build/finalize.sh 

RUN chmod +x /go/src/github.com/alllomancer/webserver/ldd-cp.sh
RUN /go/src/github.com/alllomancer/k8s-slackbot/ldd-cp.sh ldd-cp  /k8s-slackbot /temp


# Create a small image
FROM busybox AS default-image

COPY --from=builder /temp/ /
ENTRYPOINT ["/k8s-slackbot"]