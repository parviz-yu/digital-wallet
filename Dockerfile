FROM golang:1.19.2-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o wallet cmd/wallet/main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/wallet .
ENTRYPOINT [ "./wallet"]