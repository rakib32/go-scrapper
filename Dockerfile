ARG GO_VERSION=1.16

# Start from base golang image
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group


# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# Unit tests
RUN CGO_ENABLED=0 go test ./... -v

# Build the Go app
RUN CGO_ENABLED=0 go build -mod=vendor -o /app .

######## Start a new stage from scratch #######
# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
# Import the compiled executable from the first stage.
COPY --from=builder /app /app
# COPY --from=builder /src/test.json /test.json

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
CMD ["/app"]
