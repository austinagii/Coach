FROM golang:1.21.6 as base

WORKDIR /api

COPY go.mod /api/go.mod
COPY go.sum /api/go.sum
RUN go mod download -json

CMD ["/usr/bin/env", "bash"]


FROM base as devcontainer

ENTRYPOINT ["/usr/bin/env", "bash"]


FROM base as dev

COPY . /api

ENTRYPOINT ["go", "run", "/api/cmd/server/main.go"]
