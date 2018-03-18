## Get your own Derek robot

To use our managed Derek robot service get in touch with [Alex Ellis](mailto:alex@openfaas.com) for more info. (setup time 5 minutes)

Read on if you want to setup your own cluster, OpenFaaS and a private GitHub App. (est. setup time several hours)

### The hard way

* Setup [OpenFaaS](https://github.com/openfaas/faas) and the [faas-cli](https://github.com/openfaas/faas-cli)

* Now get your publically-available URL for OpenFaaS (or one punched out with an ngrok.io tunnel)

* Install Derek as a Github app and get your private key, save it as "derek.pem" and put it into the auth folder.

### Configure Docker image:

We have to build a Docker image with your .pem file included

We'll also set the symmetric key or secret that you got from GitHub as the `secret_key` environmental variable. Validating via a symmetric key is also known as HMAC. If you want to turn this off (to edit and debug) then set `validate_hmac="false"`

Fill out the `application` variable with the ID of the registered Derek Github App, and the `installation` variable with the installation ID you got when adding Derek to your account.

Set the following in your Dockerfile

```
ENV secret_key="docker"
ENV application=4385
ENV installation=45362
ENV private_key="derek.pem"

ENV validate_hmac="true"
```

Now, build and deploy Derek:

```
$ docker build -t derek .
$ faas-cli deploy --name derek --image derek --fprocess=./derek
```

Finally configure the features that you want to enable within your GitHub repo by creating a `.DEREK.yml` file.
The file should detail which features you wish to enable and the maintainer names; for example this repo would look as follows:

```yml
maintainers:
 - alexellis
 - rgee0
 - johnmccabe
features:
 - dco_check
 - comments
 - commit_linting
```

Features:

* `dco_check` - checks that each commit finishes with a "Signed-off-by:" statement
* `comments` - allows `maintainers` to issue commands to Derek to add labels etc
* `commit_linting` - applies linting rules to commit messages

Commit linting ensures:

* subjects start with a capital letter
* subject lines are less than 50 characters

**Testing**

Create a label of "no-dco" within every project you want Derek to help you with.

Head over to your GitHub repository and raise a Pull Request from the web-UI for your README file. This will not sign-off the commit, so you'll have Derek on your case.
