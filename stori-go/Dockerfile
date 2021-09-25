FROM golang:latest
WORKDIR /go/src
COPY . .
COPY transactions.csv /tmp/

RUN true \
    && go get -d -v ./... \
    && go install -v ./... \
    && go build . \
    && go test -v ./...
    
CMD ["stori-go"]
