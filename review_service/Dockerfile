FROM golang:1.19.3-alpine
RUN mkdir review
COPY . /review
WORKDIR /review
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 2000