#Etapa de build do projeto
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY . .

RUN go mod download
RUN go build -buildvcs=false -o main.exe .

#Etapa de execução
FROM alpine:3.23 AS runtime

COPY --from=builder src/main.exe .
COPY /.docker/config/load_env.sh .

ENTRYPOINT [ "./main.exe" ]
