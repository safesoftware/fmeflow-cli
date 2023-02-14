# argument for Go version
ARG GO_VERSION=1.19
 
# Build container stage
FROM golang:${GO_VERSION} AS build
ARG APP_VERSION=notset
WORKDIR /app
COPY . .

# Build the executable
RUN CGO_ENABLED=0 go build -o fmeserver -ldflags="-X \"github.com/safesoftare/fmeserver-cli/cmd.appVersion=${APP_VERSION}\""
 
# Use distroless for final image
FROM gcr.io/distroless/static:nonroot

# Run program as a non-root user by default
USER nonroot:nonroot
 
# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app/fmeserver /fmeserver
 
# run 
ENTRYPOINT ["/fmeserver"]