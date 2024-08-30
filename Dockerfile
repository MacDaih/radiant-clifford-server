FROM golang:1.22-alpine

ARG GITHUB_USERNAME
ARG GITHUB_TOKEN
ARG GITHUB_REPO_PATH

RUN apk add git
RUN git config --global \
url."https://${GITHUB_USERNAME}:${GITHUB_TOKEN}@github.com/${GITHUB_REPO_PATH}".insteadOf \
"https://github.com/${GITHUB_REPO_PATH}"
WORKDIR /go/src/ws

COPY . .
ENV GOPRIVATE=github.com/${GITHUB_REPO_PATH}

RUN go mod tidy
RUN go mod vendor
RUN go mod download
RUN go build -o main ./cmd

CMD ["/go/src/ws/main"]
