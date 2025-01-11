# syntax = docker/dockerfile:1.2

FROM golang:1.18-alpine
ENV CGO_ENABLED=0

# # define build args and environment variables
ARG PORT
ENV PORT $PORT

ENV VERSION 1.0.0
ENV NAME cendit

# create app directory
WORKDIR /app

# mount env file - Render cloud stores .env secrets in this location hence the need to mount unto the docker image
# RUN --mount=type=secret,id=_env,dst=/etc/secrets/.env cat /etc/secrets/.env

# install dependencies
COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/${NAME} -ldflags "-X main.Version=${VERSION}" ./cmd/${NAME}.go


EXPOSE $PORT

CMD [ "bin/cendit" ]