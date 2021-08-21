FROM golang:1.17
RUN mkdir -p /workspace
ADD . /workspace
WORKDIR /workspace
RUN go build -o app ./main.go
ENTRYPOINT ["./app"]
