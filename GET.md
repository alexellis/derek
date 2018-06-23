## Get your own Derek robot

* The easy way

To use our hosted/managed Derek robot service for free then get in touch with [Alex Ellis](mailto:alex@openfaas.com) for more information. (setup time 5 minutes)

* The harder way

Read on if you want to setup your own cluster, OpenFaaS and a private GitHub App. (est. setup time 30-60mins depending on your experience-level)

### Ready.. Set.. Derek!

* Setup [OpenFaaS](https://github.com/openfaas/faas) and the [faas-cli](https://github.com/openfaas/faas-cli)

* Now get your publicly-available URL for OpenFaaS (or one punched out with an ngrok.io tunnel)

* Install Derek as a GitHub app and get your private key, save it as "derek.pem".

* It is also recommended that you set a webhook secret within the GitHub application settings. If applying secrets from files, store this in a file called "derek-secret-key".

### Add your secrets :
  
Using the method appropriate to the orchestrator chosen during the OpenFaaS setup add `derek.pem` and the GitHub webhook secret as `derek-private-key` and `derek-secret-key` respectively.

Using Docker with files:
```
$ docker secret create derek-private-key derek.pem && \
    docker secret create derek-secret-key derek-secret-key
```

Using Kubernetes:
```
$ kubectl create secret generic derek-private-key --from-file=path/to/derek.pem && \
    kubectl create secret generic derek-secret-key --from-file=path/to/derek-secret-key
```

### Configure Docker image:

Most of the configuration is outside of the image, the exception being the version of the OpenFaaS watchdog which is pulled from [GitHub](https://github.com/openfaas/faas/releases); check the desired version is being pulled into the image and go ahead and build:  

```
$ docker build -t derek .
```

### Configure stack.yml:

This is where Derek finds the details he needs to do the work he does.  The main areas that will need to be updated are the `application` and `customer_url` variables.  The gateway value may also need amending if the gateway is remote.

``` yml
provider:
  name: faas
  gateway: http://localhost:8080  \# can be a remote server
  
functions:
  derek:
    handler: ./derek
    image: derek
    lang: Dockerfile
    environment:
      application: <your_GH_applicationID>
      customer_url: <your_customers_file_url>
      validate_hmac: true
      debug: true
    secrets:
      - derek-secret-key
      - derek-private-key
```
Fill out the `application` variable with the ID of the registered Derek GitHub App.

Provide an `https` location of your `.CUSTOMERS` file.  If hosted on GitHub then this should be location obtained from raw.

Validating via a symmetric key is also known as HMAC. If the webhook secret wasn't set earlier and you want to turn this off (to edit and debug) then set `validate_hmac="false"`

Now deploy Derek:
```
$ faas-cli deploy
```

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
