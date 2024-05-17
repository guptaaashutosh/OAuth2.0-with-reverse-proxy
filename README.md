# OAuth2.0 Workflow: Secure Access Management with Ory Oathkeeper, Hydra, and Resource Server

This Project showcases how your Go server acts as a single point for both Hydra and Oathkeeper functionalities. It manages user authentication, consent, issuing access tokens, and validating them before reaching the Resource Server.

## Components
- Go Server: This server combines functionalities of Ory Hydra and Ory Oathkeeper written in Go. (Provided in this repository)
- Resource Server: A simple resource server that provides protected data (simulated in this demo).

## Workflow
The demo implements a basic authorization code flow:

1. **User Access Request**: The user tries to access a protected resource through a client application.
2. **Redirect to Login**: The client application redirects the user to your Go server's login endpoint.
3. **User Login**: The user logs in with their credentials on your Go server (acting as Hydra).
4. **Consent Screen**: Your Go server (acting as Hydra) presents the user with a consent screen detailing the permissions requested by the client application. The user grants or denies access.
5. **Authorization Code**: If consent is granted, your Go server (acting as Hydra) sends an authorization code back to the client application's redirect URI.
6. **Token Request**: The client application sends the authorization code and its client credentials to your Go server (acting as Hydra) to request an access token.
7. **Access Token Grant**: Your Go server (acting as Hydra) verifies the request and issues an access token.
8. **Access Resource**: The client application includes the access token in its request to the Resource Server.
9. **Validation**: Your Go server (acting as Oathkeeper) intercepts the request, validates the access token with its internal logic (connected to Hydra for verification), and enforces any additional authorization rules (not implemented in this basic demo).
10. **Resource Access Granted**: If valid, your Go server allows the request to reach the Resource Server, which returns the protected resource (simulated data in this demo).

<br><br>

![Workflow diagram](/oathkeeper-hydra(OAuth2.0)%20workflow%20cprt.svg)

<br><br>

## Table of Contents
- Overview of OAuth2
- Prerequisites
- Project Structure
- Hydra
- Oathkeeper
- Docker Setup
- Resource Server Setup
- OathKeeper Reverse Proxy Server Setup
- Frontend Application
- Troubleshooting



## Overview of OAuth2
OAuth2 is an authorization framework that enables applications to obtain limited access to user accounts on an HTTP service. It works by delegating user authentication to the service that hosts the user account and authorizing third-party applications to access the user account. OAuth2 is widely used to grant websites or applications limited access to user information without exposing user credentials.



## Prerequisites
Before you begin, ensure you have the following installed on your machine:
- Docker (https://docs.docker.com/get-docker/)
- Docker Compose
- Go (https://go.dev/doc/install)
- Basic understanding of OAuth2 concepts
- Git for version control
- Hydra and OathKeeper (mention below)

## Project Structure
The project structure is as follows:
```OAuth2_With_reverse_proxy_workflow/
├── cmd
│   └── server
│       └── main.go
├── config
│   └── config.go
├── controllers
│   └── example_controller.go
├── constants
│   └── constants.go
├── middleware
│   └── auth_middleware.go
├── router
│   └── index-route.go
├── oathkeeper
│   ├── oathkeeper.json
│   └── rules.json
├── docker-compose.yml
├── .env.example
├── README.md
└── go.mod
```


## Hydra
Hydra is an open-source OAuth2 and OpenID Connect server that implements the Authorization Server and OpenID Provider. It allows you to manage OAuth2 tokens and handle user consent flows securely.


### Key Features
- OAuth2 and OpenID Connect provider
- Consent app
- Token management
- High scalability

### Hydra Installation : 
1.	Installation: download from [Install Hydra](https://github.com/ory/hydra/releases)
2.	Unzip and put the file of unzip folder in GOPATH. (C:\Users\ZTI\go\bin)
3.	Run `hydra` in terminal.


## Oathkeeper
Oathkeeper is a reverse proxy and identity-aware proxy that integrates with Hydra to enforce access policies based on OAuth2 tokens.

### Key Features
- Identity and Access Proxy
- Policy enforcement
- Extensible with mutators and authenticators

### OathKeeper Installation : 
1.	Installation: download from [Install OathKeeper](https://github.com/ory/oathkeeper/releases)   or [official docs](https://www.ory.sh/docs/oathkeeper/install)
2.	Unzip and put the file of unzip folder in GOPATH. (C:\Users\ZTI\go\bin)
3.	Run `oathkeeper` in terminal.


##  Docker Setup
To set up the environment using Docker, follow these steps:

1. Clone the Repository
```
git clone https://github.com/guptaaashutosh/OAuth2.0-with-reverse-proxy.git
cd OAuth2.0-with-reverse-proxy
```

2. Configure Environment Variables

    Create a .env file in the project root directory and populate it with necessary environment variables. An example .env.example file is provided for reference.


3. Start Docker Containers

```
docker-compose up -d
```

4. Create Hydra client
```
docker-compose exec hydra hydra clients create \
--endpoint http://127.0.0.1:4445/ \
--name client_name \
--id client_id \
--secret client_secret \
--grant-types authorization_code,refresh_token \
--response-types code,id_token \
--callbacks callback_url \
--token-endpoint-auth-method client_secret_post \
--scope offline
```

##  Resource Server Setup
To run resource server.
```
cd OAuth2.0-with-reverse-proxy
go run main.go
```

##  OathKeeper Reverse Proxy Server Setup
```
cd OAuth2.0-with-reverse-proxy
oathkeeper serve proxy -c './oathkeeper/oathkeeper.json'
```

<br>

## Frontend Application
The frontend for this project is built using React. You can find the frontend application in a separate repository.

### Frontend Repository
[OAuth2 React Frontend](https://github.com/guptaaashutosh/OAuth2.0-with-reverse-proxy-client)

<br>

Note: This is a basic demo project, and the Resource Server only provides simulated protected data. You'll need to implement a real Resource Server and potentially add custom authorization rules in your Go server (acting as Oathkeeper).

<br>

## Additional Resources
- Ory Hydra Documentation: https://www.ory.sh/docs/hydra/reference/api
- Ory Oathkeeper Documentation: https://www.ory.sh/docs/oathkeeper/reference/api
- OAuth 2.0 Authorization Framework: https://auth0.com/docs/authenticate/protocols/oauth

<br><br>

By following these steps and leveraging the extensive capabilities of Ory Hydra and Oathkeeper, you can establish a secure and scalable authorization system for your Go server applications, empowering a seamless and trusted user experience.

This project provides a starting point for secure OAuth2.0 authorization in your Go server using Ory Hydra and Oathkeeper.

- Fine-grained Access Control
- Simplified User Management
- Scalable Architecture
- Rapid Development

#### Customize further by:

- Implementing a real Resource Server
- Creating custom authorization rules
- Exploring advanced Ory features

This readme provides a foundation for using your Go server with integrated Hydra and Oathkeeper functionalities. Remember to update the instructions on how to access the login endpoint and tailor the testing steps based on your specific implementation.
<br><br>

**For any questions or feedback, feel free to reach out to me:** [Aashutosh Gupta](https://www.linkedin.com/in/aashutosh-gupta-06a8b7210/)

