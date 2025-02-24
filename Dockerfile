FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .


ARG APP_PORT
ENV APP_PORT=${APP_PORT}

# Expose the dynamic port
EXPOSE ${APP_PORT}

CMD ["./main"]
