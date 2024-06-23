FROM golang:1.22 as build

WORKDIR /go/src/github.com/thoughtgears/demo-service-discovery/apps/store-bff
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/app .

FROM golang:1.22-alpine as artifact
WORKDIR /app
COPY --from=build /go/src/github.com/thoughtgears/demo-service-discovery/apps/store-bff/builds/app .

EXPOSE 8080
CMD ["./app"]
