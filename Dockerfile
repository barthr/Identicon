FROM golang:1.13 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server cmd/api/main.go


FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
