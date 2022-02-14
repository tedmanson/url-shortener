############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/woolworthslimited/app/
COPY . .

# Fetch dependencies.

# Using go mod.
RUN go mod download
RUN go mod verify

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app

############################
# STEP 2 build a small image
############################
FROM gcr.io/distroless/base-debian10

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/app /go/bin/app

# Use an unprivileged user.
USER appuser:appuser

# Port on which the service will be exposed.
EXPOSE 9292

# Run the binary.
ENTRYPOINT ["/go/bin/app"]