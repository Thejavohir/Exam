FROM golang:latest
RUN mkdir customer
COPY . /customer
WORKDIR /customer
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 3000