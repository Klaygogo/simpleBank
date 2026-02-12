#build stage

FROM golang:1.25.6-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./main.go
RUN apk add curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

EXPOSE 8080
CMD ["./main"]

#run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
#COPY --from=builder /app/migrate ./migrate
#COPY db/migration ./migration
COPY app.env .
#COPY start.sh .
EXPOSE 8080

#RUN apk add --no-cache ca-certificates postgresql-client
#RUN chmod +x ./main ./migrate ./start.sh
CMD ["./main"]
#ENTRYPOINT [ "./start.sh" ]

