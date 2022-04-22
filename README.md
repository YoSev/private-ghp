# private-ghp

Serves static sites from private repositories to members with read access (or higher), secured using GitHub OAuth2.\
The server is written in [Go](https://go.dev/).

## Why

Github forces users to pay for an enterprise license in order to share github pages only with members of an organization.\
It is not even included in the paid Team plan.

This repository is a workaround for that - it also works for free tier plans as it doesn't rely on github pages nor public repositories.
## Features

- Supports multiple pages, one per subdomain
- Define your own branch
- No need to publish a repository nor use of github pages
- Works with github free tier
- Selfhosted

## Config

You need to create a [Github OAuth2 App](https://docs.github.com/en/developers/apps/building-oauth-apps/creating-an-oauth-app).

Check out the sample.config.yaml for more informations.\
The callback URL must point to **http(s)://domain:publicPort/login/github/callback** (not subdomain).

## Secure using HTTPs

To protect your site using SSL, we advice to use a reverse proxy like [Traefik](https://traefik.io/).

## Usage
- First, build using:
  - make prepare
  - make build+linux or 
  - make build+docker (optional)
- Second, set up a configuration 
  - checkout sample.config.yaml for more informations
- Third, execute the binary or docker image
  - ./prviate-ghp --config=\<path_to_config\>

## Architecture

This is a high level explanation of how this project works.

For more information of how GitHub OAuth works, see [the official documentation](https://developer.github.com/apps/building-github-apps/identifying-and-authorizing-users-for-github-apps/).

- The client requests a resource
  - If the session cookie is present and valid, skip the next two steps
  - Otherwise, redirects to the provider's OAuth page
- Provider's (e.g. GitHub) OAuth page
  - If successful, redirects to the callback URL (this service)
- The callback request is received from the OAuth provider
  - Get an OAuth token, then store it client-side in a cookie
- A call is performed to the Github API using the token the client sends with each request (as cookie) to get the resource, which is then served to the client

## Serving from a documentation directory

GitHub Pages allows to serve from a **/docs** directory, which is supported by private-ghp too, if the hosted page uses /docs as basePath.

For any question, create an issue here on Github.

