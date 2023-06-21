# NexWSTCp

This program is used to convert TCP messages from the iNexBot controller into websocket data to be forwarded to the client.

## How to use

1. Follow [https://go.dev/](https://go.dev/) to install Go Environment
2. Build program
3. Run program
4. Connect to the websocket in the client, `ws://ip:9898/ws`

## Connect to controller TCP host

Send message via websocket in client

```json
{
  "command": 1,
  "data": "controller ip"
}
```

## Send message to controller

```json
{
  "command": 1234, //Command of TCP protocol
  "data": "Serialized json string"
}
```
