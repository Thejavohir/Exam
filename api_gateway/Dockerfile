FROM golang:latest
RUN mkdir api
COPY . /api
WORKDIR /api
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 9090