FROM golang:1.19.3-alpine
RUN mkdir post
COPY . /post
WORKDIR /post
RUN go mod tidy && go mod vendor
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 4000