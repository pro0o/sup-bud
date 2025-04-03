FROM golang:1.23.5 AS builder
WORKDIR /app

COPY go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server server.go

RUN GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/sup-bud.wasm cmd/sup-bud/main.go && \
    mkdir -p web/js && \
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/js/ && \
    if [ -f "app.js" ]; then cp app.js web/js/; fi && \
    if [ -f "index.html" ]; then cp index.html web/; fi

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/web/ /app/web/

RUN chmod +x /app/server

EXPOSE 8080

CMD ["/app/server"]