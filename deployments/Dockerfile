# Builder stage
FROM --platform=$BUILDPLATFORM golang:1.25.5 AS build
LABEL stage=intermediate
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/approve_checker ./cmd/approve_checker

# Runtime stage
FROM alpine:3.23.2  
WORKDIR /root/  
COPY --from=build /app/bin/approve_checker ./approve_checker
RUN apk --no-cache add ca-certificates  
CMD ["./approve_checker"]