FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .
# END BUILDER, RUNNING CONTAINER
FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /main .
RUN mkdir /root/static
RUN mkdir /root/uploads
RUN mkdir /root/templates
ADD ./static /root/static/
ADD ./templates /root/templates

EXPOSE 8080

CMD ["./main"]
