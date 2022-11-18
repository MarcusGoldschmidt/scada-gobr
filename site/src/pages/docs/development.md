---
title: Setup development
description: Setup development
section: Dev
layout: ../../layouts/Docs.astro
---

## Requirements for dev mode

- [Node.js](https://nodejs.org/en/) (v14.17.0 or higher)
- [Yarn](https://yarnpkg.com/) (v1.22.10 or higher)
- [Go](https://golang.org/) (v1.16.5 or higher)

### Setup

Clone the repository

```bash
git clone https://github.com/MarcusGoldschmidt/scada-gobr.git
cd scada-gobr

# Install dependencies
make install-dev
```

### Run the server

```bash
go run cmd/api/api.go
```

### Run the web client

```bash
cd scadagobr-client
npm run dev
```