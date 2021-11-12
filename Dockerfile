FROM node:alpine AS frontend-builder
WORKDIR /build
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install --frozen-lockfile
COPY frontend/ ./
RUN npm run build

FROM go:alpine AS service-builder
WORKDIR /build
COPY service/go.mod service/go.sum ./
RUN go mod install
COPY service/ ./
RUN go test ./... -race
RUN go build -o catcher cmd/inmem/main.go

FROM alpine
WORKDIR /app
RUN addgroup -g 1001 -S catcher
RUN adduser -S catcher -u 1001
COPY --from=frontend-builder /build/out ./frontend
COPY --from=service-builder /build/catcher ./
USER catcher
EXPOSE 80
CMD ["./catcher"]

