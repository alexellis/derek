# derek
It's derek. Nice to meet you.

## Early instructions

Get a JWT:

* In auth folder, insert your pem/private key

You'll get this when you create your Github App

* Get a JWT

This needs the Ruby version on brew, install required gems too.

Run export JWT=(ruby app.rb)

* Now get your bearer token:

```
$ curl -i -X POST -H "Authorization: Bearer $JWT" -H "Accept: application/vnd.github.machine-man-preview+json" https://api.github.com/installations/<id>/access_tokens
```

The <id> is where your app was installed, you can find this via the webhooks sent to your endpoint or via GitHub profile page.

Save the resulting token into your access_token.txt file.

* Save a test event

You can save a test event from a webhook (re-delivery page or the live endpoint) or edit sample/cli.json.

Run `derek`:

```
$ export access_token=$(cat access_token.txt) ; ./derek < sample/cli.json 
```

