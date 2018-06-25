# hmac

Validate HMAC in Golang.

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
