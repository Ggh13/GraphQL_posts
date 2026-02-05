FROM alpine:latest AS builder

# Устанавливаем Go 1.24
RUN apk add --no-cache wget tar
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /main .
COPY config/ ./config/ 
EXPOSE 8080
CMD ["./main"]