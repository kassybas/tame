FROM golang:1.13-alpine as build-env
RUN mkdir /workspace
WORKDIR /workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./tame

FROM scratch
COPY --from=build-env /workspace/tame /tame
ENTRYPOINT ["/tame"]