Go Server-Sent Events Example
=============================

# Why
Server-Sent Events allow real-time one-way communications from the server
to client over HTTP. You can send updates and notifications on-the-fly this way.

Server-sent events are implemented on top of HTTP while websockets are not.

# Installing and Running
```
$ SRCPATH=github.com/wyc/go-server-sent-events-example
$ go get $SRCPATH
$ go run $GOPATH/src/$SRCPATH/server_sent_events.go
```

Then, either visit `http://127.0.0.1:8000/messages` with a modern web browser or
```
$ curl http://127.0.0.1:8000/messages
```

# Notes
- There's a strange feature where connections take more than 10 seconds when visiting
concurrently from a browser. This is probably some caching behavior, but I'm not sure.
Concurrent connections from curl and Javascript EventSource have no delay.
- This implementation does not utilize Event IDs as per the specification:
https://html.spec.whatwg.org/multipage/comms.html#concept-event-stream-last-event-id
- Why the long project name?
[To prevent confusion](http://en.wikipedia.org/wiki/Streaming_SIMD_Extensions)
- Thanks to cronos from `irc.freenode.net/#go-nuts` for reviewing my code

Feel free to send pull requests or leave comments.
