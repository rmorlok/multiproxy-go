# Multiproxy

A multiplexing reverse proxy in Go. This is useful when you want to serve requests from a list of hosts in priority
order.

For example, suppose you have `example1.com` with `dog.jpg` and `cat.jpg` and `example2.com` with `mouse.jpg` and
`horse.jpg`

Running the proxy:

```
go run multiproxy.go http://example1.com http://example2.com
```

would setup a proxy server on `localhost:8888` for which you could access all of the above files.

Concretely why would want this is if you have something like a development environment that uses some data from a higher
level enviroment, but stores files in S3 or another objectstore. Because your database likely has references to objects
that aren't in your development S3 bucket, you can use this proxy to make them appear seemless in the app, while new
uploads could go to your development S3 bucket. Of course, this only works if the files are publically accessible.

Multiproxy supports http and https, and path prefixes on the hosts so that the hosts do not need to be aligned at the
root level.