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
