# TikTok Tech Immersion

![Tests](https://github.com/domlimm/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

## Overview

This is the completed backend assignment of 2023 TikTok Tech Immersion.

## Setup

### Clone the repository

**Using HTTPS**

```
git clone https://github.com/domlimm/assignment_demo_2023.git
```

**Using SSH**

```
git clone git@github.com:domlimm/assignment_demo_2023.git
```

### Pre-Requisites

Ensure that you have Docker installed via [download link](https://www.docker.com/products/docker-desktop/) and [Postman](https://www.postman.com/downloads/) to communicate with the RPC server via HTTP endpoints.

#### Running Locally with Development Environment

1. Ensure Go is installed via [download link](https://go.dev/doc/install).
2. Follow the rest of the steps [below](#running-locally-without-development-environment).

#### Running Locally without Development Environment

From the root directory i.e., `./assignment_demo_2023`, simply run the following command:

```
docker compose up -d
```

And to stop the server

```
docker compose down
```

#### API Specifications

The available HTTP endpoints to communicate with the RPC server.

| Endpoint  | Request Type | Response Codes |
| --------- | ------------ | -------------- |
| /api/pull | GET          | 200, 400       |
| /api/send | POST         | 200, 400       |

**/api/pull, GET**

Sample Request Body

```json
{
  "chat": "john:jane",
  "cursor": 0,
  "limit": 10,
  "reverse": false
}
```

Sample Response Body

```json
{
  "messages": [
    {
      "chat": "john:jane",
      "text": "Hello Jane",
      "sender": "john",
      "send_time": 1686403859
    },
    {
      "chat": "john:jane",
      "text": "Hello John, wick...",
      "sender": "jane",
      "send_time": 1686405006
    }
  ]
}
```

**/api/send, POST**

Sample Request Body

```json
{
  "chat": "john:jane",
  "text": "Hello John, wick...",
  "sender": "jane"
}
```

Sample Response Body

```json
empty
```
