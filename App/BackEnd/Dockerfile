#
# revision 2 from 20-Oct-2019
#
# unfortunatly minishift server has only version 1.13
# and not suport the myltylayer images
# all lines below work on docker version 1.17+
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
COPY poc-cna-be.go $GOPATH/src
COPY ./configs/ /tmp/configs/
RUN ls -lias /tmp/
WORKDIR $GOPATH/src
RUN go get -d -v
RUN go build -o /tmp/poc-cna-be 

FROM alpine
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && mkdir -p /app
COPY --from=builder /tmp/poc-cna-be /app
COPY --from=builder /tmp/configs/ /app/configs/
RUN chmod a+rx /app/poc-cna-be
USER appuser
WORKDIR /app
CMD ["./poc-cna-be"]

# docker build -t poc-cna-be:0.1 .
# docker run --publish 8080:8080 --detach poc-cna-be:0.1