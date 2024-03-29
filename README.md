# backend

Main backend for inventory hub application

## DISCLAIMER

The code in this repository is obsolete. The documentation is still somewhat accurate and usable.

## Table of Contents

- [Setup](#setup)
  - [Docker](#docker)
  - [Local Development](#local-development)
- [Diagrams](#diagrams)
  - [Architecture](#architecture)
  - [White Space Analysis](#white-space-analysis)
  - [Registration Sequence Diagram](#registration-sequence-diagram)
  - [Order State Machine](#order-state-machine)
- [Api Documentation](#api-documentation)
  - [/api/auth](#apiauth)
    - /api/auth/login [[POST](#apiauthlogin-post)]
    - /api/auth/invite [[POST](#apiauthinvite-post)]
    - /api/auth/refresh [[POST](#apiauthrefresh-post)]
    - /api/auth/register [[POST](#apiauthregister-post)]
  - [/api/users](#apiusers)
    - /api/users [[GET](#apiusers-get)]
    - /api/users/:id [[GET](#apiusersid-get), [PUT](#apiusersid-put), [DELETE](#apiusersid-delete)]
  - [/api/categories](#apicategories)
    - /api/categories [[GET](#apicategories-get), [POST](#apicategories-post)]
    - /api/categories/:name [[DELETE](#apicategoriesname-delete)]
  - [/api/products](#apiproducts)
    - /api/products [[GET](#apiproducts-get), [POST](#apiproducts-post)]
    - /api/products/:id [[GET](#apiproductsid-get), [PUT](#apiproductsid-put), [DELETE](#apiproductsid-delete)]
  - [/api/products/orders](#apiproductsorders)
    - [order states](#order-states)
    - /api/products/orders [[GET](#apiproductsorders-get), [POST](#apiproductsorders-post)]
    - /api/products/orders/:id [[GET](#apiproductsordersid-get) [PUT](#apiproductsordersid-put), [DELETE](#apiproductsordersid-delete)]
    - /api/products/orders/:id/state [[PATCH](#apiproductsordersidstate-patch)]

## Setup

After cloning the repo, run `cp .env.example .env` in the root directory and fill in or replace the environment variables.

### Docker

Install docker and docker-compose for your OS.

To start everything in production mode, run `docker-compose up -d --build` in the root directory.

### Local Development

Start the services with `docker-compose -f=compose-services.yml up -d` in the root directory. (or have them running locally / on cloud + change the env variables)

Install the go sdk.

...WIP

## Diagrams

### Architecture

![Architecture Diagram](./.github/diagrams/architecture.drawio.png)

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

### Order State Machine

```mermaid
stateDiagram-v2
    state Choice <<choice>>
    [*] --> Pending
    note right of Pending
      Order is waiting approval
    end note
    Pending --> Ready : decrease item quantity [manger/admin approval]
    Pending --> Cancelled : [manager/admin/user rejection]
    Ready --> Choice
    Choice --> Completed : complete order [enough quantity]
    Choice --> Cancelled : cancel order, increase item quantity
    Cancelled --> [*]
    Completed --> [*]
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
│   ├── [GET]
│   └── /:id [GET, PUT, DELETE]
├── /products
|   ├── /categories
│   │   ├── [GET, POST]
│   │   └── /:name [DELETE]
|   ├── /orders
│   │   ├── [GET, POST]
│   │   └── /state/:id [PATCH]
│   ├── [GET, POST]
│   ├── /:id [GET, PUT, DELETE]
│   └── /quantity/:id [PATCH]

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

### /api/users

#### /api/users [GET]

Get the list of users with pagination, searching, maybe sorting (change the contract and add defaults) and minimal information required.

Authorization: Authorized (results show inferiors and peers)

Example query parameters:

```yaml
page: 1
pageSize: 10
search: "John"
```

Example success response:

```json
{
  "users": [
    {
      "id": "<id>",
      "username": "john_admin88",
      "firstName": "John",
      "lastName": "Admin",
      "role": "Admin",
      "createdAt": "2021-01-01T00:00:00.000Z"
    },
    {
      "id": "<id>",
      "username": "john_manager88",
      "firstName": "John",
      "lastName": "Doe",
      "role": "ReadonlyUser",
      "createdAt": "2021-01-01T00:00:00.000Z"
    }
  ],
  "totalPages": 1
}
```

Example error response:

```json
{
  "errors": {
    "page": ["The page must be a positive integer"]
  }
}
```

#### /api/users/:id [GET]

Get the user by id.

Authorization: Authorized (fails if the user is an inferior role)

Example success response:

```json
{
  "id": "<id>",
  "username": "john_admin88",
  "firstName": "John",
  "lastName": "Admin",
  "role": "Admin",
  "email": "john.admin@example.com",
  "createdAt": "2021-01-01T00:00:00.000Z"
}
```

Example error response:

```json
{
  "errors": {
    "id": ["The user with id '<id>' does not exist"]
  }
}
```

#### /api/users/:id [PUT]

Update the user by id.

Authorization: Authorized (fails if the user is an inferior or equal role)

Example payload:

```json
{
  "firstName": "John",
  "lastName": "Admin",
  "role": "Admin",
  "email": "johnny.admin@example.com"
}
```

Example success response: 204 No Content

Example error response:

```json
{
  "errors": {
    "id": ["The user with id '<id>' does not exist"],
    "role": ["The role 'BigBoss' is not valid"]
  }
}
```

#### /api/users/:id [DELETE]

Delete the user by id.

Authorization: Authorized (fails if the user is an inferior or equal role)

Example success response: 204 No Content

Example error response:

```json
{
  "errors": {
    "id": ["The user with id '<id>' does not exist"]
  }
}
```

### /api/categories

#### /api/categories [GET]

Get the list of categories available.

Authorization: Authorized

Example success response:

```json
{
  "categories": [
    {
      "id": "<id>",
      "name": "Electronics",
      "itemquantity": 10
    },
    {
      "id": "<id>",
      "name": "Furniture",
      "itemquantity": 5
    }
  ]
}
```

#### /api/products/categories [POST]

Create a new category.

Authorization: [Admin, Manager]

Example payload:

```json
{
  "name": "Electronics"
}
```

Example success response: 201 Created

```json
{
  "id": "<id>",
  "name": "Electronics"
}
```

Example error response:

```json
{
  "errors": {
    "name": ["The category with name 'Electronics' already exists"]
  }
}
```

#### /api/products/categories/:name [DELETE]

Delete the category by name and all products in it.

Authorization: [Admin, Manager]

Example success response: 204 No Content

Example error response:

```json
{
  "errors": {
    "name": ["The category with name 'Electronics' does not exist"]
  }
}
```

### /api/products

#### /api/products [GET]

Get the list of products with pagination, searching, filtering and maybe sorting (change the contract and add defaults).

Authorization: Authorized

Example query parameters:

```yaml
page: 1
pageSize: 10
search: "iPhone"
category: "Electronics" # optional, defaults to all categories
```

Example success response:

```json
{
  "products": [
    {
      "id": "<id>",
      "name": "iPhone 12",
      "category": "Electronics",
      "quantity": 10,
      "description": "Better than iPhone 11 (maybe)",
      "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
      "createdAt": "2021-01-01T00:00:00.000Z",
      "updatedAt": "2021-01-01T00:00:00.000Z"
    },
    {
      "id": "<id>",
      "name": "iPhone 11",
      "category": "Electronics",
      "quantity": 5,
      "imageUrl": null,
      "description": "Better than iPhone 10",
      "createdAt": "2021-01-01T00:00:00.000Z",
      "updatedAt": "2021-01-01T00:00:00.000Z"
    }
  ]
}
```

Example error response:

```json
{
  "errors": {
    "page": ["The page must be a positive integer"],
    "category": ["The category with name 'Electronics' does not exist"]
  }
}
```

#### /api/products [POST]

Create a new item.

Authorization: [Admin, Manager, User]

Example payload (form data):

```yml
name: iPhone 12
category: Electronics
quantity: 10
description: Better than iPhone 11 (maybe)
image: <binary data> | null
```

> Note: if implementing form data is too difficult, use JSON instead and the image will be encoded in base64

Example success response: 201 Created

```json
{
  "id": "<id>",
  "name": "iPhone 12",
  "category": "Electronics",
  "quantity": 10,
  "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
  "description": "Better than iPhone 11 (maybe)",
  "createdAt": "2021-01-01T00:00:00.000Z",
  "updatedAt": "2021-01-01T00:00:00.000Z"
}
```

Example error response:

```json
{
  "errors": {
    "category": ["The category with name 'Electronics' does not exist"]
  }
}
```

#### /api/products/:id [GET]

Get the item by id.

Authorization: Authorized

Example success response:

```json
{
  "id": "<id>",
  "name": "iPhone 12",
  "category": "Electronics",
  "quantity": 10,
  "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
  "description": "Better than iPhone 11 (maybe)",
  "createdAt": "2021-01-01T00:00:00.000Z",
  "updatedAt": "2021-01-01T00:00:00.000Z"
}
```

Example error response:

```json
{
  "errors": {
    "id": ["The item with id '<id>' does not exist"]
  }
}
```

#### /api/products/:id [PUT]

Update the item by id.

Authorization: [Admin, Manager, User]

Example payload (form data):

```yml
name: iPhone 12
category: Electronics
quantity: 10
description: Better than iPhone 11 (maybe)
image: <binary data> | null
```

> Note: if implementing form data is too difficult, use JSON instead and the image will be encoded in base64

Example success response: 200 OK

```json
{
  "id": "<id>",
  "name": "iPhone 12",
  "category": "Electronics",
  "quantity": 10,
  "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
  "description": "Better than iPhone 11 (maybe)",
  "createdAt": "2021-01-01T00:00:00.000Z",
  "updatedAt": "2023-01-01T00:00:00.000Z"
}
```

Example error response:

```json
{
  "errors": {
    "id": ["The item with id '<id>' does not exist"],
    "category": ["The category with name 'Electronics' does not exist"]
  }
}
```

#### /api/products/:id [DELETE]

Delete the item by id.

Authorization: [Admin, Manager]

Example success response: 204 No Content

Example error response:

```json
{
  "errors": {
    "id": ["The item with id '<id>' does not exist"]
  }
}
```

### /api/products/orders

#### /api/products/orders [GET]

Get the list of orders with limited information with pagination, searching, filtering and maybe sorting (change the contract and add defaults).

Authorization: Authorized

Example query parameters:

```yaml
page: 1
pageSize: 10
search: "iPhone"
category: "Electronics" # optional, defaults to all categories
```

Example success response:

```json
{
  "orders": [
    {
      "id": "<transactionId>",
      "client": "John Doe Enterprises",
      "description": "Will buy 5 IPhones in lease",
      "quantity": 5,
      "state": "Draft",
      "createdAt": "2021-01-01T00:00:00.000Z",
      "updatedAt": "2021-01-01T00:00:00.000Z",
      "item": {
        "name": "iPhone 12",
        "category": "Electronics",
        "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png"
      }
    },
    {
      "id": "<transactionId>",
      "client": "Main Office Floor 4",
      "description": "",
      "quantity": 12,
      "state": "Pending",
      "createdAt": "2021-01-01T00:00:00.000Z",
      "updatedAt": "2021-02-01T00:00:00.000Z",
      "item": {
        "name": "Toilet Paper",
        "category": "Hygiene",
        "quantity": 100
      }
    }
  ],
  "totalPages": 1
}
```

Example error response:

```json
{
  "errors": {
    "page": ["The page must be a positive integer"]
  }
}
```

#### /api/products/orders [POST]

Create a new transaction.

Authorization: [Admin, Manager, User]

Example payload:

```json
{
  "itemId": "<itemId>",
  "client": "John Doe Enterprises",
  "description": "Will buy 5 IPhones in lease",
  "quantity": 5,
  "initialState": "Pending"
}
```

Example success response:

```json
{
  "id": "<transactionId>",
  "client": "John Doe Enterprises",
  "description": "Will buy 5 IPhones in lease",
  "quantity": 5,
  "state": "Pending",
  "createdAt": "2021-01-01T00:00:00.000Z",
  "updatedAt": "2021-01-01T00:00:00.000Z",
  "item": {
    "id": "<itemId>",
    "category": "Electronics",
    "quantity": 10,
    "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
    "description": "Better than iPhone 11 (maybe)",
    "createdAt": "2021-01-01T00:00:00.000Z",
    "updatedAt": "2021-01-01T00:00:00.000Z"
  }
}
```

Example error response:

```json
{
  "errors": {
    "initialState": ["The state 'Completed' cannot be an initial state"]
  }
}
```

#### /api/products/orders/:id [GET]

Get the transaction by id with all the information.

Authorization: [Admin, Manager, User]

Example success response:

```json
{
  "id": "<transactionId>",
  "client": "John Doe Enterprises Incorporated",
  "description": "Will buy 5 IPhones in lease",
  "quantity": 6,
  "state": "Pending",
  "createdAt": "2021-01-01T00:00:00.000Z",
  "updatedAt": "2021-01-01T00:00:00.000Z",
  "item": {
    "id": "<itemId>",
    "name": "iPhone 12",
    "category": "Electronics",
    "quantity": 10,
    "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
    "description": "Better than iPhone 11 (maybe)",
    "createdAt": "2021-01-01T00:00:00.000Z",
    "updatedAt": "2021-01-01T00:00:00.000Z"
  }
}
```

Example error response:

```json
{
  "errors": {
    "id": ["The transaction with id '<transactionId>' does not exist"]
  }
}
```

#### /api/products/orders/:id [PUT]

Update the transaction by id.

Authorization: [Admin, Manager]

Example payload:

```json
{
  "client": "John Doe Enterprises Incorporated",
  "description": "Will buy 5 IPhones in lease",
  "quantity": 6
}
```

Example response:

```json
{
  "id": "<transactionId>",
  "client": "John Doe Enterprises Incorporated",
  "description": "Will buy 5 IPhones in lease",
  "quantity": 6,
  "state": "Pending",
  "createdAt": "2021-01-01T00:00:00.000Z",
  "updatedAt": "2021-01-01T00:00:00.000Z",
  "item": {
    "id": "<itemId>",
    "name": "iPhone 12",
    "category": "Electronics",
    "quantity": 10,
    "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
    "description": "Better than iPhone 11 (maybe)",
    "createdAt": "2021-01-01T00:00:00.000Z",
    "updatedAt": "2021-01-01T00:00:00.000Z"
  }
}
```

Example error response:

```json
{
  "errors": {
    "quantity": [
      "You specified the quantity 11, but the item 'iPhone 12' has only 10 unit(s) available"
    ]
  }
}
```

#### /api/products/orders/:id [DELETE]

Delete the transaction by id. orders should not be deleted, to be kept for historic record, but admins can still do that.

Authorization: [Admin]

Example success response: 204 No Content

Example error response: 401 Unauthorized

#### /api/products/orders/:id/state [PATCH]

Update the transaction state by id.

Authorization: [Admin, Manager]

Example payload:

```json
{
  "transition": "Completed"
}
```

Example success response:

```json
{
  "id": "<transactionId>",
  "client": "John Doe Enterprises Incorporated",
  "description": "Will buy 5 IPhones in lease",
  "quantity": 6,
  "state": "Completed",
  "createdAt": "2021-01-01T00:00:00.000Z",
  "updatedAt": "2021-01-01T00:00:00.000Z",
  "item": {
    "id": "<itemId>",
    "name": "iPhone 12",
    "category": "Electronics",
    "quantity": 4,
    "imageUrl": "cdn.inventory-hub.space/uploads/products/iPhone12-<id>.png",
    "description": "Better than iPhone 11 (maybe)",
    "createdAt": "2021-01-01T00:00:00.000Z",
    "updatedAt": "2021-01-01T00:00:00.000Z"
  }
}
```

Example error response:

```json
{
  "errors": {
    "transition": [
      "The transaction cannot go from the state 'Draft' to the state 'Completed'"
    ]
  }
}
```
