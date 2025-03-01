FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o recipe-manager cmd/main.go

FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/recipe-manager .

RUN mkdir -p configs
COPY --from=builder /app/configs/props.* ./configs/
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENV PORT=8080

RUN adduser -D -g '' appuser
USER appuser

CMD ["./recipe-manager"]
