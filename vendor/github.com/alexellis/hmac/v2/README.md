# hmac

Validate HMAC in Golang.

## Who uses it HMAC?

[GitHub](https://developer.github.com/webhooks/securing/), Patreon and some other parties will use HMAC signing with their outgoing webhooks so that you can verify the webhook is from the expected sender.

## Who uses this project?

A few of the notable dependents on this package, but there are many more:

* [alexellis/derek](https://github.com/alexellis/derek/)
* [openfaas/faas-cli](https://github.com/openfaas/faas-cli)
* [openfaas/openfaas-cloud](https://github.com/openfaas/openfaas-cloud/)
* [crossplane/tbs](https://github.com/crossplane/tbs)
* [Qolzam/telar-cli](https://github.com/Qolzam/telar-cli)
* [s8sg/faas-flow](https://github.com/s8sg/faas-flow)

## How it works:

HMAC uses a symmetric key that both sender/receiver share ahead of time. The sender will generate a hash when wanting to transmit a message - this data is sent along with the payload. The recipient will then sign payload with the shared key and if the hash matches then the payload is assumed to be from the sender.

[Read more on Wikipedia](https://en.wikipedia.org/wiki/HMAC)

# Documentation

[![](https://godoc.org/github.com/alexellis/hmac?status.svg)](http://godoc.org/github.com/alexellis/hmac)

## Example:

```
import "github.com/alexellis/hmac"

...
var input []byte
var signature string
var secret string

valid := hmac.Validate(input, signature, secret)

fmt.Printf("Valid HMAC? %t\n")
```
