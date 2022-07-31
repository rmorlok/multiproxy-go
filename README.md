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
level environment, but stores files in S3 or another object store. Because your database likely has references to objects
that aren't in your development S3 bucket, you can use this proxy to make them appear seamless in the app, while new
uploads could go to your development S3 bucket. Of course, this only works if the files are publicly accessible.

Multiproxy supports http and https, and path prefixes on the hosts so that the hosts do not need to be aligned at the
root level.

## Real World Example

To demonstrate a real example of multiplexing, use the following example:

```bash
# Start the server
go run multiproxy.go https://go.dev/ https://hmpioneers.net/hockey/

# Run a request that will find a valid response on the first website
curl -f -s "http://0.0.0.0:8888/learn/" | grep "Go"

# Run a request that will find a valid response on the second website
curl -f -s "http://0.0.0.0:8888/history.htm" | grep "Hill-Murray"
```

Note how the multiproxy falls back to the second site when it can't find the content on the first. This tool supports
an arbitrary number of fallbacks.

## Running with a Let's Encrypt Certificate

Follow the instructions on [Certbot](https://certbot.eff.org/) to receive a certificate.

## Manually generating Certificates Locally

Per [this StackOverflow question](https://stackoverflow.com/questions/8169999/how-can-i-create-a-self-signed-cert-for-localhost),
use the following commands on MacOS:

```bash
# Use 'localhost' for the 'Common name'
openssl req -x509 -sha256 -nodes -newkey rsa:2048 -days 365 -keyout localhost.key -out localhost.crt

# Trust the certificate locally
sudo security add-trusted-cert -p ssl -d -r trustRoot -k ~/Library/Keychains/login.keychain localhost.crt
```

You can then launch `multiproxy` using the generated cert with the following command:

```bash
go run multiproxy.go --pem ./localhost.crt  --key ./localhost.key https://example.com/
```