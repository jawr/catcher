FROM node:alpine AS frontend-builder
WORKDIR /build
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install --frozen-lockfile
COPY frontend/ ./
RUN NODE_ENV=production npm run build

FROM golang:alpine AS service-builder
WORKDIR /build
COPY service/go.mod service/go.sum ./
RUN go mod download
COPY service/ ./
RUN CGO_ENABLED=0 go test ./... 
RUN CGO_ENABLED=0 go build -o catcher cmd/inmem/main.go

FROM alpine
WORKDIR /app
RUN addgroup -g 1001 -S catcher
RUN adduser -S catcher -u 1001
COPY --from=frontend-builder /build/out ./frontend
COPY --from=service-builder /build/catcher ./
COPY config.yaml ./
USER catcher
EXPOSE 8080
EXPOSE 2525
CMD ["./catcher"]

