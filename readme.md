# Simple reverse proxy for local development

if you have several services running on localhost or in docker containers and you want to route different urls to specific addresses.

## Usage
run this application and path needed routes
`'{Desired URL}->{Target URL}'`

Desired URL should be specify in the hosts file to point to the localhost.

## Example:
`local-routing 'foo.loc->http://localhost:8080/foo/htdocs/'`
`local-routing 'foo.loc->http://localhost:8080/foo/htdocs/' 'bar.loc->http://localhost:8080/bar/link/'`
### Hosts:
`127.0.0.1	foo.loc`
`127.0.0.1	bar.loc`
