# backend

Main backend for inventory hub application

## Table of Contents

- [Setup](#setup)
  - [Docker](#docker)
  - [Local Development](#local-development)
- [Diagrams](#diagrams)
  - [White Space Analysis](#white-space-analysis)
  - [Registration Sequence Diagram](#registration-sequence-diagram)
- [Api Documentation](#api-documentation)
  - [/api/auth](#apiauth)
    - [/api/auth/login [POST]](#apiauthlogin-post)
    - [/api/auth/invite [POST]](#apiauthinvite-post)
    - [/api/auth/refresh [POST]](#apiauthrefresh-post)
    - [/api/auth/register [POST]](#apiauthregister-post)

## Setup

### Docker

Install docker and docker-compose for your OS.

To start the services, run `docker-compose up -d` in the root directory.

To start the services and the backend, run `docker-compose --profile=server up -d --build` in the root directory.

### Local Development

Start the services with `docker-compose up -d` in the root directory.

Install Java XX and Maven XX.

...WIP

## Diagrams

### White Space Analysis

```mermaid
quadrantChart
    title White Space Analysis
    x-axis Lean Access --> Secure
    y-axis Expensive --> Cheap / Open Source

    quadrant-1 Secure and Cheap
    quadrant-2 Quick Setup
    quadrant-3 Wall of shame
    quadrant-4 Enterprise solutions

    Inventory Hub: [0.9, 0.9]
    Sortly: [0.9, 0.4]
    Inventree: [0.55, 0.99]
    Jira Plugin: [0.9, 0.1]
    Cin7: [0.95, 0.05]
    monday.com: [0.45, 0.8]
    StoreHub: [0.1, 0.4]
```

### Registration Sequence Diagram

```mermaid
sequenceDiagram
    title Registration Workflow
    actor Admin
    participant Backend
    participant MQ
    participant Email as Email Microservice
    participant Azure as Azure Communications
    actor User

    Admin->>+Backend: Register User Form
    Backend->>Backend: Create Draft User
    Backend-->>+MQ: Send User Created Event
    Backend-->>-Admin: User Created Response
    MQ-->>+Email: User Created Event
    deactivate MQ
    Email->>Email: Create email from template
    Email->>+Azure: Send email
    Azure-->>+User: Try to send email
    loop
        Azure->>Azure: Poll for status
    end
    deactivate User
    Azure-->>-Email: Return email sent status
    alt email sent
        Email-->>+MQ: Confirm message as processed
        deactivate MQ
    else email failed
        Email-->>+MQ: Return message to queue with N+1 retries
        deactivate MQ
    end
    deactivate Email
```

## Api Documentation

The namespace structure for the api is the following:

```
/api
├── /auth
│   ├── /login [POST]
│   ├── /invite [POST]
|   ├── /refresh [POST]
│   └── /register [POST]
├── /users
```

### /api/auth

#### /api/auth/login [POST]

Login using credentials.

Authorization: Anonymous

Example payload:

```json
{
  "email": "admin@example.com",
  "password": "password"
}
```

Example success response:

```json
{
  "accessToken": "<jwt>",
  "refreshToken": "<refreshToken>"
}
```

Example error response (you can implement it differently if you want):

```json
{
  "errors": {
    "email": ["The email is not valid"]
  }
}
```

#### /api/auth/invite [POST]

Invite a user to the application.

Authorization: [Admin, Manager]

Example payload:

```json
{
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "role": "ReadonlyUser"
}
```

Example success response: 201 Created (empty body)

Example error response:

```json
{
  "errors": {
    "role": ["The role 'Blatnoi' is not valid"]
  }
}
```

#### /api/auth/refresh [POST]

Refresh JWT.

Authorization: Anonymous

Example payload:

old JWT in the authorization header

```json
{
  "refreshToken": "<refreshToken>"
}
```

Example success response:

```json
{
  "accessToken": "<jwt>",
  "refreshToken": "<refreshToken>"
}
```

#### /api/auth/register [POST]

Register a draft user.

Authorization: Anonymous

Example payload:

```json
{
  "token": "<invitationToken>",
  "username": "tolya_perforator1996",
  "password": "Tolya123!"
}
```

Example success response: 201 Created

```json
{
  "accessToken": "<jwt>",
  "refreshToken": "<refreshToken>"
}
```

Example error response:

```json
{
  "errors": {
    "token": ["The invitation token is not valid"]
  }
}
```
