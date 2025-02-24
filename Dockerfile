FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .


ARG PORT
ENV APORT=${PORT}

# Expose the dynamic port
EXPOSE ${PORT}

CMD ["./main"]
