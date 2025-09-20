# Based off https://github.com/GoogleContainerTools/distroless?tab=readme-ov-file#examples-with-docker
FROM golang:1.24-alpine AS build

WORKDIR /go/src/app
COPY . /go/src/app

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/app /

# Expose port
EXPOSE 8080
# Run the app
CMD ["/app"]
