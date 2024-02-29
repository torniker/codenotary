# The Problem

Please have a look at our cloud service immudb (Vault)[http://vault.immudb.io/] and build a simple application (Frontend + Backend) around it with the following requirements:

- Application is storing accounting information within immudb Vault with the following structure: account number (unique), account name, iban, address, amount, type (sending, receiving)

- Application has an API to add and retrieve accounting information

- Application has a frontend that displays accounting information and allows to create new records.

The solution should:

- Have a readme

- Have a documented API

- Have docker-compose so it is easy to run

Should be easily configurable to run on non-localhost URLs.

Resources:
immudb Vault documentation: https://vault.immudb.io/docs/
API reference: https://vault.immudb.io/docs/api/v1

# Setup

Run the following command to start the application

```bash
docker compose build
docker compose up
```

# API

`GET /api/accounting` - returns accounting information
query params:

```
{
  "page": "int",
  "perPage": "int"
}
```

`GET /api/accounting/count` - returns count of accounting information

`POST /api/accounting` - creates a new accounting information
body:

```
{
  "accountNumber": "string",
  "accountName": "string",
  "iban": "string",
  "address": "string",
  "amount": "int",
  "type": "string"
}
```

# Useful commands

Run the following command to create accounting collection

```bash
make install-wand && ./build/wand create-accounting-collection
```

Run the following command to get accounting collection

```bash
make install-wand && ./build/wand get-accounting-collection
```

Run the following command to drop accounting collection

```bash
make install-wand && ./build/wand drop-accounting-collection
```

Run the following command to add random accounting document

```bash
make install-wand && ./build/wand add-accounting-document
```

Run the following command to get the list of accounting documents

```bash
make install-wand && ./build/wand search-accounting-document | jq
```

Run the following command to replace documents. Replace id with existing id

```bash
make install-wand && ./build/wand replace-accounting-document {id}
```
