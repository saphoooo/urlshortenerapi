# URL shortener API

A basic implementation of an URL shortener in Go using a Redis storage

The schema of the API is quite simple:

| URL          | Verbs | Success | Failure  |
|--------------|:-----:|--------:|---------:|
| /api/v1/new  | POST  | 201     | 500, 404 |
| /api/v1/:url | GET   | 200     | 404      |
