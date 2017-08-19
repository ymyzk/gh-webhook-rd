# gh-webhook-rd

[![Build Status](https://travis-ci.org/ymyzk/gh-webhook-rd.svg?branch=master)](https://travis-ci.org/ymyzk/gh-webhook-rd)

**gh-webhook-rd** is a simple application that starts a job in Rundeck when receiving a webhook from GitHub.

## Getting Started
- Go to [the releases page](https://github.com/ymyzk/gh-webhook-rd/releases) and download a binary
- [Obtain an API token](http://rundeck.org/docs/api/#token-authentication) to start jobs in Rundeck
- Write configuration.  Example:
  - Webhook URL: `http://<host>:<port>/webhook/hello`
  - If branch is `refs/heads/master`, then start a job with id `12344567-d8d5-49ff-ac95-1dfabf837561`
```toml
[[hooks]]
url = "hello"
branch = "refs/heads/master"
job_id = "12344567-d8d5-49ff-ac95-1dfabf837561"
```
- [Set up a Webhook in your repository](https://developer.github.com/webhooks/creating/)
  - Confirm that content type is `application/json`
- Start `$ ./gh-configuration -c config.toml`

## Example Configuration
```toml
[server]
# Optional (default: <empty>)
host = "127.0.0.1"
# Optional (default: 8080)
port = 8080

[rundeck]
url = "http://localhost:4440"
auth_token = "L3o8Gnn6yizg8v4vnDSSCOkpfcDFYWH5"

[[hooks]]
# Required: Webhook endpoint URL
url = "hello"
# Optional: Webhook HMAC key (https://developer.github.com/webhooks/securing/)
secret = "secret key"
# Required
branch = "refs/heads/master"
# Required: Rundeck Job ID
job_id = "12344567-d8d5-49ff-ac95-1dfabf837561"

[[hooks]]
url = "hello"
branch = "refs/heads/master"
job_id = "98765432-d8d5-49ff-ac95-1dfabf837561"
```
