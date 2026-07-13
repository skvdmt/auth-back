# Подготовка.
FROM golang:alpine AS preper
ARG NAME
WORKDIR /usr/src/${NAME}
COPY . .
COPY ./config /etc
RUN go mod download

# Тестирование.
FROM preper AS testing
ARG POSTGRES_PASSWORD
RUN go test --tags=unit -v ./...
RUN go test --tags=integration -v ./...
RUN go test --tags=e2e -v ./...

# Сборка.
FROM preper AS building
RUN go build -v -o /usr/local/bin/${NAME} ./cmd/main.go

# Релиз.
FROM alpine AS release
ARG NAME
# Настройки.
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
# Копирование файлов.
COPY ./config /etc
COPY --from=building /usr/local/bin/${NAME} /usr/local/bin/${NAME}
# Создание точки входа.
COPY ./docker-entrypoint.sh /usr/local/bin
RUN echo "exec ${NAME}" >> /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
EXPOSE 8000
ENTRYPOINT [ "docker-entrypoint.sh" ]
