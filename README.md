# auth

## client

A simple library for authentication via:

* SASL Plain authentication via Unix Socket (non-standardized)
* SMTPS
* STARTTLS

## server

A simple service which processes authentication requests received via:

* SASL Plain authentication via Unix Socket (non-standardized)

### Security Considerations

* Even though the credentials are usually the same, authentication via SMTP is preferred over IMAP. We don't want access to stored emails.
