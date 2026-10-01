# Подготовка.
FROM golang:alpine AS preper
ARG NAME
WORKDIR /usr/src/${NAME}
# Создание директорий.
RUN mkdir /etc/${NAME}
RUN mkdir /var/log/${NAME}
RUN mkdir /usr/local/share/doc
# Копирование файлов.
COPY . .
COPY ./openapi.yaml /usr/local/share/doc/openapi.yaml
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
RUN go mod download

# Тестирование.
FROM preper AS testing

RUN --mount=type=secret,id=postgres_password \
  --mount=type=secret,id=access_key \
  --mount=type=secret,id=refresh_key \
  export POSTGRES_PASSWORD=$(cat /run/secrets/postgres_password) && \
  export ACCESS_KEY=$(cat /run/secrets/access_key) && \
  export REFRESH_KEY=$(cat /run/secrets/refresh_key) && \
  go test --tags=unit -v ./... && \
  go test --tags=integration -v ./... && \
  go test --tags=e2e -v ./...

# Сборка.
FROM preper AS building
RUN go build -v -o /usr/local/bin/${NAME} ./cmd/main.go

# Релиз.
FROM alpine AS release
EXPOSE 8000
ARG NAME
# Настройки.
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
# Создание директорий.
RUN mkdir /var/log/${NAME}
RUN mkdir /etc/${NAME}
RUN mkdir /usr/local/share/doc
# Копирование файлов.
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
COPY --from=building /usr/local/bin/${NAME} /usr/local/bin/${NAME}
COPY ./openapi.yaml /usr/local/share/doc/openapi.yaml
COPY ./docker-entrypoint.sh /usr/local/bin
# Создание точки входа.
RUN echo "exec ${NAME}" >> /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
EXPOSE 8000
ENTRYPOINT [ "docker-entrypoint.sh" ]
