# backend

Main backend for inventory hub application

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
