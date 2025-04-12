# Notification System

## Description

This is a personal project of a notification system having the following requirements

- The system needs to be able to send notifications via several different channels (email,
  sms, slack) and be easily extensible to support more channels in the future.
- The system needs to be horizontally scalable.
- The system must guarantee an "at least once" SLA for sending the message.
- The interface for accepting notifications to be sent is an HTTP API.

## Architecture

## Prerequisites

- Go
- Docker

## How to configure .env

### RabbitMQ

### Twilio

### Slack

### Email

## How to run locally?

- Run docker-compose
- Run the service
- Run the worker

## How to run in production?

## Future improvements

- Protect the API with API keys per user
- Protect the API with rate limiting - both overall and per user
- Write more tests for better coverage
- Support rich formatting of the Slack messages
- Add Slack channel validation in both format, existence and permissions to send messages to that channel
