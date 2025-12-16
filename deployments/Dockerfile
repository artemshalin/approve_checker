# Builder stage
FROM --platform=$BUILDPLATFORM golang:1.25.5 AS build
LABEL stage=intermediate
WORKDIR /app
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/approve_checker ./cmd/approve_checker

# Runtime stage
FROM scratch
COPY --from=build /app/bin/approve_checker /bin/approve_checker
CMD ["approve_checker"]
