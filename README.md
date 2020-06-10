# auth

## client

A simple library for authentication via:

* SASL Plain authentication via Unix Socket (non-standardized)
* SMTPS
* STARTTLS

## server

A simple service which processes authentication requests received via:

* SASL Plain authentication via Unix Socket (non-standardized)

## Design choices

* fail2ban integration with reverse proxies
  * Issue
    * client sends a request with a faked `X-Forwarded-For` header
    * nginx appends the client's real IP address to the header value
    * a software uses the first (faked) `X-Forwarded-For` IP address
    * fail2ban has no effect, someone else gets blacklisted
  * Decision: `client.ExtractIP` extracts the `X-Real-IP` header. Your first reverse proxy must set `X-Real-IP` to the client IP address, later proxies must not.

## Security Considerations

* Even though the credentials are usually the same, authentication via SMTP is preferred over IMAP. We don't want access to stored emails.
