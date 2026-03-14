# coredns-pg-backend

A [Postgres](https://postgresql.org) backend plugin for [CoreDNS](https://coredns.io)
with an HTTP API, WebUI, and CLI for managing entries.

It is a monorepo with five components:
- Common go code (libcoredns-pg)
- The CoreDNS plugin (plugin)
- The REST API (api)
- The CLI (cli)
- The Web UI (web)

All components are available under the MIT license.
