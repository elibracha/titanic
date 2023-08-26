FROM golang:1.21 as builder
LABEL stage=builder
WORKDIR /usr/src/app

COPY . .
RUN make test && make build


FROM golang:1.21
WORKDIR /app/

ARG port
ARG csv_path
ARG sqlite_path

ENV API_PORT=${port}
ENV CSV_STORE_PATH=${csv_path}
ENV SQLITE_STORE_PATH=${sqlite_path}

COPY --from=builder /usr/src/app/app .
COPY --from=builder /usr/src/app/config.yaml config.yaml
COPY --from=builder /usr/src/app/docs ./docs
COPY --from=builder /usr/src/app/data ./data

EXPOSE $port
ENTRYPOINT ["./app"]