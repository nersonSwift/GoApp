FROM golang:1.13.6-alpine AS build-stage

WORKDIR /go/src/application

COPY servises servises
COPY go.mod go.mod
COPY go.sum go.sum

COPY second_app .

RUN go get -u 
RUN go mod tidy 
RUN CGO_ENABLED=0 GOOS=linux go build -o application -ldflags="-w -s" -v


FROM alpine:latest AS final-stage
RUN apk --no-cache add ca-certificates tini su-exec
COPY --from=build-stage /go/src/application/application /application

ENTRYPOINT ["/application", "--workers", "$WORKERS"]
