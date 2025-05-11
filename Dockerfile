######################Stage 1#################################
FROM golang:1.22.5 as baseImage

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN ls -al /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
######################Stage 2#################################

FROM scratch

WORKDIR /app

# The other files wont go in the binary , only the compiled go files would be there
COPY --from=baseImage /app/main .
COPY --from=baseImage /app/config.json .


EXPOSE 9090

ENTRYPOINT ["/app/main"]



