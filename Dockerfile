FROM golang:1.18  AS builder
WORKDIR /go/src/github.com/alllomancer/k8s-slackbot
ADD . /go/src/github.com/alllomancer/k8s-slackbot
RUN cd /go/src/github.com/alllomancer/k8s-slackbot

RUN go build -o /app/k8s-slackbot .

RUN chmod +x /go/src/github.com/alllomancer/k8s-slackbot/ldd-cp.sh
RUN /go/src/github.com/alllomancer/k8s-slackbot/ldd-cp.sh ldd-cp  /app/k8s-slackbot /temp


# Create a small image
FROM busybox AS default-image

COPY --from=builder /temp/ /
ENTRYPOINT ["/app/k8s-slackbot"]