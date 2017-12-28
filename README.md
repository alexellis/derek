# derek

[![Build Status](https://travis-ci.org/alexellis/derek.svg?branch=master)](https://travis-ci.org/alexellis/derek)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

It's derek ![](https://pbs.twimg.com/media/DPo4OyrWsAAOk_i.png). Nice to meet you. I'd like to help you with Pull Requests and Issues on your GitHub project.

> Please show support for the project and **Star** the repo.

From the team that bought you [OpenFaaS](https://www.openfaas.com) - Serverless Functions Made Simple.

## What can I do?

* Check that commits are signed-off

When someone sends a PR without a sign-off, I'll apply a label `no-dco` and also send them a comment pointing them to the contributor guide. Most of the time when I've been helping the OpenFaaS project - people read my message and fix things up without you having to get involved.

* Allow users in a specified .DEREK.yml file to manage issues and pull-requests

You don't have to give people full write access anymore to help you manage issues and pull-requests. I'll do that for you, just put them in a .DEREK.yml file in the root and when they comment on an issue then I'll use my granular permissions instead.

* Wait.. doesn't the term "maintainer" mean write access in GitHub?

No this is what Derek sets out to resolve. The users in your maintainers list have granular permissions which you'll see in detail when you add the app to your repo org.

```
maintainers:
- alexellis
- rgee0
```

You can use the alias "curators" instead for the exact same behaviour:

```
curators:
- alexellis
- rgee0
```

* What about roles?

We are planning to add roles in the ROADMAP which will mean you can get even more granular and have folks who can only add labels but not close issues for instance. If you feel you need to make that distinction. It will also let you call the roles whatever you think makes sense. 

> Note that the assign/unassign commands provides the shortcut `me` to assign to the commenter

Example:

* Titles

```
Derek set title: This is a more meaningful title
```

* Labels

```
Derek add label: awesome
Derek remove label: awesome
```

* Assign work

```
Derek assign: alexellis
Derek unassign: me
```

* Manage issues and PRs

```
Derek close
Derek reopen
Derek lock
Derek unlock
```

Coming soon
* [x] Derek as a managed GitHub App
* [x] Lock thread
* [x] Edit title
* [x] Toggle the DCO-feature
* [ ] Add roles & actions
* [ ] Branch Checking

[Live demo here](https://twitter.com/alexellisuk/status/905694832445804544)

Example from a live project:

https://github.com/openfaas/faas-cli/issues/85

## How do I work?

I'm designed to be installed as a [GitHub App](https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/), but don't worry - I don't need a lot of permissions. Just access to issues and Pull Requests will do.

I'm not a long-running daemon.. I'd get bored that way. I work with webhooks - so stick me in a serverless framework like [OpenFaaS](https://github.com/alexellis/faas) and forget about me. Just apply oil from time to time.

This is me in action! Normally contributors edit and re-push within a few minutes after re-reading the contribution guide.

![](https://user-images.githubusercontent.com/6358735/29704343-542a36da-8971-11e7-871e-da30c8e86cae.png)

*Inspiration for Derek*

The idea for a bot that could comment on issues or respond to activity is from the docker/docker or Moby project's bot called [Poule](https://github.com/icecrime/poule). It's a much more complex long-running daemon which uses Personal Access Tokens (so needs to run as a full GitHub login). Derek is much simpler (so hackable) and can be installed with granular permissions.

### Where is Derek working now?

Derek is active and operating 24/7 helping the award-winning OpenFaaS project!

* http://github.com/openfaas/faas
* http://github.com/openfaas/faas-netes
* http://github.com/openfaas/faas-cli

### Maintainers

* Alex Ellis - author
* Richard Gee (@rgee0) - co-maintainer
* John Mccabe (@johnmccabe) - co-maintainer

Alex Ellis created Derek to automate project maintainer duties around licensing and to help bring granular permissions back to GitHub. Derek has empowered contributors in the OpenFaaS community to run and maintain the project without needing full write access. OpenFaaS contributors continue to improve Derek so they can get the job done without fuss.

### Contributions

Please follow the [OpenFaaS contribution guide](https://github.com/openfaas/faas/blob/master/CONTRIBUTING.md).

## Get your own Derek robot

You can host and manage your own Derek robot using [these instuctions](GET.md), or use our managed service. To use our managed service get in touch with [Alex Ellis](mailto:alex@openfaas.com) for more info.
