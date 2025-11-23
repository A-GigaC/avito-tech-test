FROM golang

#WORKDIR /app

# Копируем гомод и госам
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

# Копируем сурсы
COPY ./src .

# Собираем приложение
RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]