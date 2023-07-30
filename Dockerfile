FROM golang:1.20 as builder
LABEL stage=builder
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies
# and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

# copy source files and build the binary
COPY . .
RUN make build


FROM scratch
WORKDIR /app/

ARG port
ARG csv_path
ARG sqlite_path

ENV API_PORT=${port}
ENV CSV_STORE_PATH=${csv_path}
ENV SQLITE_STORE_PATH=${sqlite_path}
ENV STORE_TYPE=${store_type}

COPY --from=builder /usr/src/app/app .
COPY --from=builder /usr/src/app/docs ./docs
COPY --from=builder /usr/src/app/data ./data

EXPOSE $port
ENTRYPOINT ["./app"]