## Get Derek

### The workflow

GitHub sends webhooks to Derek for different events and comments that happen across your repositories. Derek then looks for a .DEREK.yml file to see if the repository requires any response.

You can also use a single main repository and then redirect to that from the others, or use a different .DEREK.yml file in each. 

For example:

* [openfaas/faas](https://github.com/openfaas/faas/blob/master/.DEREK.yml) is the main file
* [openfaas/faas-cli](https://github.com/openfaas/faas-cli/blob/master/.DEREK.yml) is a redirect file.

Installation process:

* Install the managed or your self-hosted Derek GitHub App
* Send a PR to the [customers file](https://github.com/alexellis/derek/blob/master/.CUSTOMERS) with your GitHub username or GitHub organization
* Add a .DEREK.yml to any repositories you want to include, turn on or off any features you need as per the [user guide](./USER_GUIDE.md)
* Add in more repositories in the same organisation using the redirect feature

### Derek the SaaS (managed for you, by us)

To use our managed Derek bot service follow the instructions below which take around 5-10 minutes.

* Install this GitHub App on the individual repos (not the whole org):

https://github.com/apps/derek

You will be told what permissions are required.

* Create `.DEREK.yml` in your elected repos

You can use this file as a template: https://github.com/openfaas/faas/blob/master/.DEREK.yml

* Now raise a PR to the `.CUSTOMERS` file

Raise a PR to this file, and make sure you use `git commit --signoff` rather than the UI to make the PR

https://github.com/alexellis/derek/blob/master/.CUSTOMERS

* Finally test it works

Raise a new issue and type in `Derek close`, then edit your `.DEREK.yml` file to add your team and community maintainers/contributors

* Support the managed service

Show your support for Derek by [becoming a GitHub Sponsor from 5 USD / mo](https://github.com/sponsors/alexellis).

### Doing the things the hard way (self-hosting)

Read on if you want to operate your own Derek bot, or deploy Derek for development. 

You will setup a single-node cluster with Kubernetes or Swarm, deploy OpenFaaS, create a GitHub App, install your GitHub App on a GitHub repo and then deploy Derek. Estimated setup time 30-60mins depending on your experience-level.

#### Ready.. Set.. Derek!

Pre-reqs:

* Deploy Swarm or Kubernetes
* GitHub user account

Steps:

* Deploy OpenFaaS and the faas-cli - https://docs.openfaas.com/deployment/

* Now get your publicly-available URL for OpenFaaS (or if you're behind a firewall, use an ngrok.io tunnel, if ngrok.io is blocked then you should use a public cloud)

* Create a GitHub App in your GitHub account named "Derek dev"

For the homepage URL enter: https://github.com/alexellis/derek

Enter these OAuth Permissions/Scopes:

- Issues - read/write
- Pull requests - read/write
- Repository metadata - read only

If you are setting this up on a private repository you need to grant Derek permissions to download content, this is so that he can download the config file. The write permissions are so that Derek can update your release notes if you are using the `release_notes` feature.

- Repository contents - read/write

Subscribe to these events:

- Issue comment
- Pull request

Set "Where can this GitHub App be installed?" to Any account

* Set your webhook URL to the URL of your OpenFaaS cluster and the `derek` function - i.e. `http://<ip>:<port>/function/derek`

* Note down your Application ID you will need it later where we configure derek.yml with the `github_application_id` variable.

* Click "Generate a private key". This is the private key for your Derek bot, save it as `derek-private-key` and keep it private.

* Set a webhook secret in the UI and save the value as `derek-secret-key`

* Finally install the Derek app onto one of your own GitHub repos. From now on you will get events sent to Derek.

* Create .DEREK.yml in your GitHub repo referring to the docs for what to put in the file.

Next we will build and deploy Derek to your OpenFaaS Cluster and trigger a test event such as `Derek add label: testing`

### Define secrets (Swarm)

Before deploying Derek create two secrets:

* derek-private-key - used for GitHub authz
* derek-secret-key - used for verifying webhooks are from GitHub using HMAC

For Docker:

```sh
$ docker secret create derek-private-key derek-private-key
$ docker secret create derek-secret-key derek-secret-key
```

For Kubernetes:

```sh
$ kubectl create secret generic derek-private-key --from-file=./derek-private-key -n openfaas-fn
$ kubectl create secret generic derek-secret-key --from-file=./derek-secret-key -n openfaas-fn
```

Create secrets.yml:

```sh
cp secrets.example.yml secrets.yml
```

Then populate your own values in the file - set `application_id` to your GitHub application ID and set `secret_key` to the same value of the `derek-secret-key` file.

### Configure Docker image:

Most of the configuration is outside of the image, the exception being the version of the OpenFaaS watchdog which is pulled from [GitHub](https://github.com/openfaas/faas/releases); check the desired version is being pulled into the image and go ahead and build.

Update derek.yml and replace `alexellis/` with your own username on the Docker Hub or your registry.

```
$ faas-cli build -f derek.yml
```

This will build a Docker image in your local library, now push it to the Docker Hub.

```
$ docker login
$ faas-cli push -f derek.yml
```

### Configure stack.yml:

This is where Derek finds the details he needs to do the work he does.  The main areas that will need to be updated are the `application_id` and `customer_url` variables.  The gateway value may also need amending if the gateway is remote.

```yaml
provider:
  name: faas
  gateway: http://127.0.0.1:8080  \# can be a remote server
  
functions:
  derek:
    handler: ./derek
    image: derek
    lang: Dockerfile
    environment:
      application_id: <github_application_id>
      customer_url: <your_customers_file_url>
      validate_customers: true
      validate_hmac: true
      debug: true
      write_debug: false
      secret_path: /var/openfaas/secrets/
    secrets:
      - derek-secret-key
      - derek-private-key
```
Fill out the `application` variable with the ID of the registered Derek GitHub App.

Provide an `https` location of your `.CUSTOMERS` file.  If hosted on GitHub then this should be location obtained from raw URL by clicking on the file in the UI then clicking "Raw". The default value is: `https://raw.githubusercontent.com/openfaas/faas/master/.DEREK.yml`

Validating via a symmetric key is also known as HMAC. If the webhook secret wasn't set earlier and you want to turn this off (to edit and debug) then set `validate_hmac="false"`

Now deploy Derek:
```
$ faas-cli deploy -f derek.yml
```

* Config environmental options

* `application_id` - ID provided in UI for GitHub App
* `customer_url` - A text file of valid repos which can use your Derek installation, separated by new-lines normally `CUSTOMERS`
* `validate_customers` - If set to false then the `customer_url` is ignored
* `validate_hmac` - Validate all incoming webhooks are signed with the secret `derek-secret-key` that you enter in the GitHub UI
* `write_debug` - Dump the incoming request to the function logs. This is not needed since the request can be viewed in the advanced tab of the GitHub App UI

### Configure GitHub Repo:

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
```

### Testing and troubleshooting

To test: 

Create a label of "no-dco" within every project you want Derek to help you with.
 
Head over to your GitHub repository and raise a Pull Request from the web-UI for your README file. This will not sign-off the commit, so you'll have Derek on your case.

If you're not sure if things are working right then Click on the GitHub App via your account settings and then click "Advanced" and "Recent Deliveries". This will show you all the incoming messages and their responses.

You can tweak your environment and then hit "Redeliver" to send the message again.

### Appendix

#### Personal Access Tokens

> Note: you may be able to use a Personal Access Token generated through your *Developer Settings* page instead of a GitHub App. If you wish to do this then set the `personal_access_token` environmental variable. It is not recommended since the granularity is much more coarse meaning you must grant more privileges to the bot. The bot will also report as your own GitHub user account.
