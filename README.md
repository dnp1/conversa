# Conversa
A chat API using SSE.

#### This repository three directories with teh respective content:
- server/ -> the http server
- client/ -> the client library
- cli/ -> the command line interface

#### Building and running server
- Go to **server/** directory
- The server uses a postgres database to persist data.
- There's a ddl.sql file in the **server/** directory containing the data definition code
- After create the data definition, you **must** set **CONVERSA_DB_CONN_STR** with [database connection string](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters)
- run `go get ./...`
- run go build
- ./server


#### Building and running server
- Go to cli directory
- run `go get ./...`
- run go build
- You **must** set **CONVERSA_TARGET_URL** with a baseUrl pointing to the server
- You **may** set **CONVERSA_CLIENT_SESSION** to define the filepath where client will store session cookies. By default, the filepath is ~/.conversa
- ./cli
