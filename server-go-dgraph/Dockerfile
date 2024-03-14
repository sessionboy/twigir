FROM golang:alpine

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /home/twigir/server
# install packages
# RUN go mod download
COPY . .
# RUN go build

EXPOSE 8035
ENTRYPOINT ["./server"]