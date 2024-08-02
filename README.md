# Go Challenge for FullCycle Pós Go Expert

## Overview

This repository contains a Go application developed as a part of the "`fullcyle.rate-limiter`" (Go Challenge) from the Pós Go Expert.

## Challenge Requirements

Objective: Develop a rate limiter in Go that can be configured to limit the maximum number of requests per second based on a specific IP address or access token.

Description: The objective of this challenge is to create a rate limiter in Go that can be used to control request traffic to a web service. The rate limiter must be able to limit the number of requests based on two criteria:

- IP Address: The rate limiter must restrict the number of requests received from a single IP address within a defined time interval.
- Access Token: The rate limiter can also limit requests based on a single access token, allowing different expiration time limits for different tokens. The Token must be entered in the header in the following format:
  API_KEY: <TOKEN>
- Access token limit settings must override IP threshold settings. Ex: If the limit per IP is 10 req/s and a specific token is 100 req/s, the rate limiter must use the token information.

## Instructions

To build and run this application, follow these steps:

#### Running the Go application
This will run the application
```bash
docker-compose up -build
```

#### Request by terminal
```bash
curl -s http://localhost:8080
```
with header: API_KEY:
```bash
curl -H "API_KEY: your_token_here" http://localhost:8080 & done
```

Conferir os logs no console do terminal.