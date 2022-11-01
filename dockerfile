FROM golang:1.19-alpine as builder

RUN mkdir -p /prophet
COPY . /prophet
WORKDIR /prophet

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd

CMD ["./server"]