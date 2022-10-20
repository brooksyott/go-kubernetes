# Dockerfile References: https://docs.docker.com/engine/reference/builder/
######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY ./bin/go-hellokube /app/go-hellokube

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
# CMD ["/app/go-kubernetes"]
CMD ["/app/go-hellokube"]

