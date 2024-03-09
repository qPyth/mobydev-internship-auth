# Simple User Management Application

This application is a simple user management system that provides APIs for registering a new user, authenticating a user, and updating a user's profile. It's built using [Go](https://golang.org/) and uses a Makefile for easy running.

## Features

- User registration
- User authentication
- User profile update

## Prerequisites

Before running this application, make sure you have [Go](https://golang.org/dl/) installed on your machine.

## Installation

Clone the repository to your local machine:

```bash
git clone https://github.com/qPyth/mobydev-internship-auth
```


## Configuration

This application uses environment variables for configuration. Copy the `.env.example` file to `.env` and adjust the settings according to your environment:

```bash
cp .env.example .env
```


Edit the `.env` file and set your secret key for JWTAuth, and edit the `config.yaml` and set your credentials and other necessary configurations.


## Running the Application

To run the application, use the `make run` command from the root directory of the project:

```bash
make run
```


This command compiles the Go application and starts the server on the default port.

## API Endpoints

The application exposes the following endpoints:

- `POST /user/signup`: Register a new user. Requires a JSON body with `email`, `password` and `pass_conf` fields.
- `POST /user/signin`: Authenticate a user. Requires a JSON body with `email` and `password`. Returns a JWT token upon successful authentication.
- `PUT /user/profile/update`: Update an existing user's profile. Requires a JWT token for authorization and a JSON body with fields you want to update.
