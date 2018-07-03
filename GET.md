## Get your own Derek robot

* The easy way

To use our hosted/managed Derek robot service for free then get in touch with [Alex Ellis](mailto:alex@openfaas.com) for more information. (setup time 5 minutes)

* The harder way

Read on if you want to operate your own Derek bot, or deploy Derek for development. 

You will setup a single-node cluster with Kubernetes or Swarm, deploy OpenFaaS, create a GitHub App, install your GitHub App on a GitHub repo and then deploy Derek. Estimated setup time 30-60mins depending on your experience-level.

### Ready.. Set.. Derek!

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

Subscribe to these events:

- Commit comment
- Issue comment
- Issues
- Pull request
- pull request review comment (optional)

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

```
$ docker secret create derek-private-key derek-private-key
$ docker secret create derek-secret-key derek-secret-key
```

For Kubernetes:

```
$ kubectl create secret generic derek-private-key --from-file=./derek-private-key
$ kubectl create secret generic derek-secret-key --from-file=./derek-secret-key
```

### Configure Docker image:

Most of the configuration is outside of the image, the exception being the version of the OpenFaaS watchdog which is pulled from [GitHub](https://github.com/openfaas/faas/releases); check the desired version is being pulled into the image and go ahead and build.

Update derek.yml and replace `alexellis/` with your own username on the Docker Hub or your registry.

```
$ faas-cli build -f derek.yml
```

This will build a Docker image in your local library, now push it to the Docker Hub.

```
$ docker login
$ faas-cli push
```

### Configure stack.yml:

This is where Derek finds the details he needs to do the work he does.  The main areas that will need to be updated are the `application` and `customer_url` variables.  The gateway value may also need amending if the gateway is remote.

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
      application: <github_application_id>
      customer_url: <your_customers_file_url>
      validate_customers: true
      validate_hmac: true
      debug: true
      write_debug: false
    secrets:
      - derek-secret-key
      - derek-private-key
```
Fill out the `application` variable with the ID of the registered Derek GitHub App.

Provide an `https` location of your `.CUSTOMERS` file.  If hosted on GitHub then this should be location obtained from raw URL by clicking on the file in the UI then clicking "Raw". The default value is: `https://raw.githubusercontent.com/openfaas/faas/master/.DEREK.yml`

Validating via a symmetric key is also known as HMAC. If the webhook secret wasn't set earlier and you want to turn this off (to edit and debug) then set `validate_hmac="false"`

Now deploy Derek:
```
$ faas-cli deploy
```

* Config environmental options

* `application` - ID provided in UI for GitHub App
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

**Testing**

Create a label of "no-dco" within every project you want Derek to help you with.
 
Head over to your GitHub repository and raise a Pull Request from the web-UI for your README file. This will not sign-off the commit, so you'll have Derek on your case.
