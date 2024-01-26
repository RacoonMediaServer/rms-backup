FROM golang as builder
WORKDIR /src/service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=`git tag --sort=-version:refname | head -n 1`" -o rms-backup -a -installsuffix cgo rms-backup.go
RUN CGO_ENABLED=0 GOOS=linux go build -o backupctl -a -installsuffix cgo ./app/backupctl/backupctl.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata && apk add --update docker openrc
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/service/rms-backup .
COPY --from=builder /src/service/backupctl .
COPY --from=builder /src/service/configs/rms-backup.json /etc/rms/
CMD ["./rms-backup"]