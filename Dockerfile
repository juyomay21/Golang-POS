FROM golang:1.17-alpine AS base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o build/main main.go
ENV RUN_ENV=docker_dev

EXPOSE 3030

CMD [ "build/main" ]




FROM alpine:latest as prod

COPY --from=base /app/build/main /usr/local/bin/pos
EXPOSE 3030

ENTRYPOINT ["/usr/local/bin/pos"]


# FROM golang:1.16-alpine AS base
# WORKDIR /app

# ENV GO111MODULE="on"
# ENV GOOS="linux"
# ENV CGO_ENABLED=0

# RUN apk update \
#     && apk add --no-cache \
#     ca-certificates \
#     curl \
#     tzdata \
#     git \
#     && update-ca-certificates

# FROM base AS dev
# WORKDIR /app

# RUN go get -u github.com/cosmtrek/air && go install github.com/go-delve/delve/cmd/dlv@latest
# #RUN go get -u github.com/cosmtrek/air
# EXPOSE 5000
# EXPOSE 2345

# ENTRYPOINT ["air"]

# FROM base AS builder
# WORKDIR /app

# COPY . /app
# RUN go mod download \
#     && go mod verify

# RUN go build -o pos -a .

# FROM alpine:latest as prod

# COPY --from=builder /app/pos /usr/local/bin/pos
# EXPOSE 5000

# ENTRYPOINT ["/usr/local/bin/pos"]

