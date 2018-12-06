# derek

[![Build Status](https://travis-ci.org/alexellis/derek.svg?branch=master)](https://travis-ci.org/alexellis/derek)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![](https://godoc.org/github.com/alexellis/derek?status.svg)](https://godoc.org/github.com/alexellis/derek)


It's derek ![](https://pbs.twimg.com/media/DPo4OyrWsAAOk_i.png). Nice to meet you. I'd like to help you with Pull Requests and Issues on your GitHub project.

> Please show support for the project and **Star** the repo.

From the team that brought you [OpenFaaS](https://www.openfaas.com) - Serverless Functions Made Simple.

### Looking for the User Guide?

Existing users of Derek can read the [user-guide here](./USER_GUIDE.md).

## Our users

Some of our users include:

* Docker / Moby:

https://github.com/moby/moby/issues/35736

* OpenFaaS

https://github.com/openfaas/faas-cli/issues/85

* [Subsurface Diving app](https://subsurface-divelog.org)

Example: https://github.com/Subsurface-divelog/subsurface/pull/1748

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

* URL redirection of configuration files is supported via the "redirect" field:

```
redirect: https://github.com/<some-user>/<some-repo>/.DEREK.yaml
```

If this optional field is non-empty, Derek will read it's configuration from another location. This allows multiple projects to use the same configuration.
Please note that redirection is only supported for GitHub repository URLs.

* Command triggers

By default, Derek commands can be called with `Derek <some-command>`. The prefix `Derek ` is the default trigger, but the bot also supports the `/` trigger which can be enabled by setting the `use_slash_trigger` environment variable to `true`.

### Examples:

* Update the title of a PR or issue

Let's say a user raised an issue with the title `I can't get it to work on my computer`

```
Derek set title: Question - does this work on Windows 10?
```
or
```
Derek edit title: Question - does this work on Windows 10?
```

* Triage and organise work through labels

Labels can be used to triage work or help sort it.

```
Derek add label: proposal
Derek add label: help wanted
Derek remove label: bug
```

* Set milestones for issues

You can organize your issues in groups through existing milestones

```
Derek set milestone: example
Derek remove milestone: example
```

* Assign work

You can assign work to people too

```
Derek assign: alexellis
Derek unassign: me
```

* Add a reviewer to a PR

You can assign people for a PR review as well

```
Derek set reviewer: alexellis
Derek clear reviewer: me
```

* Open and close issues and PRs

Sometimes you may want to close or re-open issues or Pull Requests:

```
Derek close
Derek reopen
```

* Lock/un-lock conversation/threads

This is useful for when conversations are going off topic or an old thread receives a lot of comments that are better placed in a new issue.

```
Derek lock
Derek unlock
```

### Backlog:

* [x] Derek as a managed GitHub App
* [x] Lock thread
* [x] Edit title
* [x] Toggle the DCO-feature

Future work:

* [ ] Caching of customers / .DEREK.yml file
* [ ] Observability of GitHub API Token rate limit
* [ ] Add roles & actions
* [ ] Branch Checking

[Live demo here](https://twitter.com/alexellisuk/status/905694832445804544)

## How do I work?

I'm designed to be installed as a [GitHub App](https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/), but don't worry - I don't need a lot of permissions. Just access to issues and Pull Requests will do.

I'm not a long-running daemon.. I'd get bored that way. I work with webhooks - so stick me in a serverless framework like [OpenFaaS](https://github.com/alexellis/faas) and forget about me. Just apply oil from time to time.

This is me in action! Normally contributors edit and re-push within a few minutes after re-reading the contribution guide.

![](https://user-images.githubusercontent.com/6358735/29704343-542a36da-8971-11e7-871e-da30c8e86cae.png)

*Inspiration for Derek*

The idea for a bot that could comment on issues or respond to activity is from the docker/docker or Moby project's bot called [Poule](https://github.com/icecrime/poule). It's a much more complex long-running daemon which uses Personal Access Tokens (so needs to run as a full GitHub login). Derek is much simpler (so hackable) and can be installed with granular permissions.

### Where is Derek working now?

Derek is active and operating 24/7 helping the award-winning OpenFaaS project!

* https://github.com/moby/moby
* http://github.com/openfaas/faas
* http://github.com/openfaas/faas-netes
* http://github.com/openfaas/faas-cli

### Maintainers / contributors

* Alex Ellis - author
* Richard Gee (@rgee0) - co-maintainer
* John Mccabe (@johnmccabe) - contributor

Alex Ellis created Derek to automate project maintainer duties around licensing and to help bring granular permissions back to GitHub. Derek has empowered contributors in the OpenFaaS community to run and maintain the project without needing full write access. OpenFaaS contributors continue to improve Derek so they can get the job done without fuss.

### Contributions

Please follow the [OpenFaaS contribution guide](https://github.com/openfaas/faas/blob/master/CONTRIBUTING.md).

## Get your own Derek robot

To use our managed service (recommended) get in touch with [Alex Ellis](mailto:alex@openfaas.com) for more info. Once you have installed the GitHub App you will need to send a PR to the [customers file](https://github.com/alexellis/derek/blob/master/.CUSTOMERS) with your username or organisation. The final step is to add your .DEREK.yml - you can use the file from this repository as an example.

You can host and manage your own Derek robot using [these instuctions](GET.md), or use our managed service.

