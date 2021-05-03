ARG GO_VERSION=1.15

FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add --no-cache git
RUN apk add --no-cache make
WORKDIR /app
COPY ./ ./
RUN make all

FROM scratch
COPY --from=builder /app/dist/recipecounter /bin/recipecounter

# Use nobody user + group
USER 65534:65534

RUN ["/bin/recipecounter", "--version"]
