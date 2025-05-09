FROM golang:1.22.5 as baseImage

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o main .

FROM gcr.io/distroless/baseImage

COPY --from=baseImage /app/main/ .

EXPOSE 9090

CMD ["./main"]



