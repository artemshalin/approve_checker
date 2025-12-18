FROM --platform=$BUILDPLATFORM golang:1.25.5 AS build
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /approve_checker ./cmd/approve_checker
RUN chmod +x /approve_checker

FROM alpine:3.23.2
RUN apk --no-cache add ca-certificates
COPY --from=build /approve_checker /usr/local/bin/approve_checker
RUN chmod +x /usr/local/bin/approve_checker
ENTRYPOINT ["/usr/local/bin/approve_checker"]
