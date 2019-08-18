
#build stage
FROM golang:alpine AS builder
ENV GO111MODULE=on
WORKDIR /go/src/jobregistration-app
COPY ./src/jobregistration-app .
RUN apk --no-cache add git
RUN go get -d -v ./...
RUN go install -v ./...
RUN ls /go/bin

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/jobregistration-app /app/jobregistration-api
RUN ls -l /app
ENTRYPOINT ["/app/jobregistration-api"]
LABEL Name=jobregistration-api Version=0.0.1
EXPOSE 8001

