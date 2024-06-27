FROM golang:1.21.6 as base

WORKDIR /coach

COPY go.mod . 
COPY go.sum .
RUN go mod download -json


FROM base as dev

RUN go install github.com/codegangsta/gin

ENTRYPOINT ["bash"]


FROM base as test

COPY . /coach

ENTRYPOINT ["go", "run", "."]
