FROM golang:1.22

WORKDIR /app

COPY go.mod     ./
COPY go.sum     ./
COPY api        ./api/
COPY templates  ./templates/
COPY types      ./types/
COPY store      ./store/

# RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /check42

EXPOSE 2442

CMD ["/check42"]