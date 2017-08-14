# derek
It's derek. Nice to meet you.

> Is this bot live yet?

## How do I work?

I'm designed to be installed as a GitHub App, but don't worry - I don't need a lot of permissions. Just access to issues and pull requests will do.

When someone sends a PR without a sign-off, I'll apply a label and also send them a comment pointing them to the contributor guide.

I'm not a long-running daemon.. I'd get bored that way. I work with webhooks - so stick me in a serverless framework like [OpenFaaS](https://github.com/alexellis/faas) and forget about me. Just apply oil from time to time.

*Inspiration*

The idea for a bot that could comment on issues or respond to activity is from the Moby project's bot called [Poule](https://github.com/icecrime/poule). It's a much more complex long-running daemon which uses Personal Access Tokens (so needs to run as a full GitHub login).

## Early instructions

Git clone Derek and build it with Golang.

Install Derek as a Github app and get your private key, save it as "private-key.pem" and put it into the auth folder.

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

If there's no DCO derek will add a label of no-dco and also comment on the issue.
