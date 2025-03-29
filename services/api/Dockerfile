# Stage 1: Build the Go binary
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./go.mod ./go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api .

# Stage 2: Run the Go binary in a minimal image
FROM alpine

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api ./

RUN mkdir -p /app/templates

COPY --from=builder /app/server/routes/admin/templates/ /app/templates/
COPY --from=builder /app/core/mailer/email_templates/ /app/email_templates/

RUN chmod +x api

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/api"]
