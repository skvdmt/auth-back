# Подготовка.
FROM golang:alpine AS preper
ARG NAME
WORKDIR /usr/src/${NAME}
# Создание директории конфигурации.
RUN mkdir /etc/${NAME}
# Копирование файлов.
COPY . .
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
RUN go mod download

# Тестирование.
FROM preper AS testing
ARG POSTGRES_PASSWORD
ARG ACCESS_KEY
ARG REFRESH_KEY
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
# Создание директории журналов.
RUN mkdir /var/log/${NAME}
# Создание директории конфигурации.
RUN mkdir /etc/${NAME}
# Копирование файлов.
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
COPY --from=building /usr/local/bin/${NAME} /usr/local/bin/${NAME}
# Создание точки входа.
COPY ./docker-entrypoint.sh /usr/local/bin
RUN echo "exec ${NAME}" >> /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
EXPOSE 8000
ENTRYPOINT [ "docker-entrypoint.sh" ]
