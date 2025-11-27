# WNT Analytics Event Contract

*(wastingnotime.org → Lambda → SQS → Consumer → Plausible)*

## 1. Overview

This document describes the event contract and responsibilities between:

1.  **Browser tracker** on `wastingnotime.org`
2.  **AWS API Gateway** (`/event`)
3.  **Lambda**: `wnt-plausible-ingest`
4.  **SQS**: `wnt-plausible-events`
5.  **Local consumer**: `cmd/analytics-consumer`
6.  **Local Plausible** (`/api/event`)

Goal:\
\> Capture analytics events from the blog, normalize them, buffer via
SQS and finally feed them into a self-hosted Plausible instance, with
correct timestamps and replay capability.

## 2. High-level flow

**Request path**

    Browser (JS)
      → API Gateway (HTTPS)
        → Lambda (normalize + enrich)
          → SQS (event buffer)
            → Local Consumer (Go)
              → Plausible /api/event

**Time semantics**

-   **Source of truth for event time** = timestamp set by **Lambda**
    (UTC, RFC3339).
-   Plausible uses this `timestamp` for analytics, regardless of
    ingestion delays.

## 3. Browser → API Gateway → Lambda: Input payload

### 3.1. HTTP

-   Method: `POST`
-   Path: `/event` (HTTP API)
-   Content-Type: `application/json`

### 3.2. JSON body (InboundEvent)

    {
      "name": "pageview",
      "url": "https://wastingnotime.org/sagas/game-hub/the-first-breath",
      "domain": "wastingnotime.org",
      "referrer": "https://google.com",
      "screen_width": 1920,
      "props": {
        "path": "/sagas/game-hub/the-first-breath",
        "title": "The First Breath"
      }
    }

### 3.3. Derived from HTTP request

-   `User-Agent` header\
-   `X-Forwarded-For` (IP)\
-   `SourceIP` from API Gateway

## 4. Lambda → SQS: QueueEvent

    {
      "domain": "wastingnotime.org",
      "name": "pageview",
      "url": "https://wastingnotime.org/sagas/game-hub/the-first-breath",
      "referrer": "https://google.com",
      "user_agent": "Mozilla/5.0 ...",
      "screen_width": 1920,
      "ip": "203.0.113.10",
      "timestamp": "2025-11-27T13:00:00Z",
      "props": {
        "path": "/sagas/game-hub/the-first-breath",
        "title": "The First Breath"
      }
    }

## 5. Consumer → Plausible: PlausibleEvent

    {
      "name": "pageview",
      "url": "https://wastingnotime.org/sagas/game-hub/the-first-breath",
      "domain": "wastingnotime.org",
      "referrer": "https://google.com",
      "screen_width": 1920,
      "props": {
        "path": "/sagas/game-hub/the-first-breath",
        "title": "The First Breath"
      },
      "timestamp": "2025-11-27T13:00:00Z"
    }

## 6. Error handling & retries

-   Lambda → SQS: returns 500 if queueing fails
-   Consumer → Plausible: non-2xx → message is NOT deleted → retried
-   DLQ recommended after `maxReceiveCount`

## 7. Versioning

Use:

    props.schema_version

Example:

    "props": {
      "path": "/x",
      "title": "Y",
      "schema_version": "v2"
    }

## 8. Environment variables

### Lambda:

-   `EVENT_QUEUE_URL`
-   `ALLOWED_ORIGIN`

### Consumer:

-   `SQS_QUEUE_URL`
-   `PLAUSIBLE_URL`
-   `AWS_REGION`