# derek
It's derek. Nice to meet you. I'd like to help you with Pull Requests on your project.

> Please show support for the project and **Star** the repo.

## How do I work?

I'm designed to be installed as a [GitHub App](https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/), but don't worry - I don't need a lot of permissions. Just access to issues and pull requests will do.

> When someone sends a PR without a sign-off, I'll apply a label `no-dco` and also send them a comment pointing them to the contributor guide.

I'm not a long-running daemon.. I'd get bored that way. I work with webhooks - so stick me in a serverless framework like [OpenFaaS](https://github.com/alexellis/faas) and forget about me. Just apply oil from time to time.

This is me in action! Normally contributors edit and re-push within a few minutes after re-reading the contribution guide.

![](https://user-images.githubusercontent.com/6358735/29704343-542a36da-8971-11e7-871e-da30c8e86cae.png)

*Inspiration for Derek*

The idea for a bot that could comment on issues or respond to activity is from the docker/docker or Moby project's bot called [Poule](https://github.com/icecrime/poule). It's a much more complex long-running daemon which uses Personal Access Tokens (so needs to run as a full GitHub login). Derek is much simpler (so hackable) and can be installed with granular permissions.

### Where is Derek working now?

Derek is active and operating 24/7 helping the OpenFaaS project!

* http://github.com/alexellis/faas

* http://github.com/alexellis/faas-netes

* http://github.com/alexellis/faas-cli

## Get your own Derek robot

* Setup [OpenFaaS](https://github.com/alexellis/openfaas) and the `faas-cli`

* Now get your publically-available URL for OpenFaaS (or one punched out with an ngrok.io tunnel)

* Install Derek as a Github app and get your private key, save it as "derek.pem" and put it into the auth folder.

### Configure Docker image:

We have to build a Docker image with your .pem file included

We'll also set the symmetric key or secret that you got from GitHub as the `secret_key` environmental variable. Validating via a symmetric key is also known as HMAC. If you want to turn this off (to edit and debug) then set `validate_hmac="false"`

Fill out the `installation` variable with the installation ID you got from GitHub for Derek.

Set the following in your Dockerfile

```
ENV secret_key="docker"
ENV installation=45362
ENV private_key="derek.pem"

ENV validate_hmac="true"
```

Now, build and deploy Derek:

```
$ docker build -t derek .
$ faas-cli -action build -name derek -image derek -fprocess=./derek
```

**Testing**

Create a label of "no-dco" within every project you want Derek to help you with.

Head over to your GitHub repository and raise a Pull Request from the web-UI for your README file. This will not sign-off the commit, so you'll have Derek on your case.
