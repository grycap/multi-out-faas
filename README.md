# multi-out-faas

An OpenFaaS function for filtering files in data-processing serverless workflows.

Using an intermediate bucket/folder, allows splitting a workflow by uploading files that comply with certain patterns to specific buckets/folders of multiple storage providers.

## Build

The function is publicly available in [Docker Hub](https://hub.docker.com/r/grycap/multi-out-faas), but if you prefer, you can build and push it using [`faas-cli`](https://github.com/openfaas/faas-cli) (remember to edit the file `multi-out-faas.yml` before):

 ```bash
 faas-cli build -f multi-out-faas.yml
 faas-cli push -f multi-out-faas.yml
 ```

## Usage

### Defining the configuration file

In order to use the function, you have to create a JSON configuration file specifying the storage providers with its authentication and the outputs where you want the files to be uploaded with your customised filters (using name prefixes and/or suffixes):

```json
{
  "storages":{
    "minio":[
      {
        "name":"minio-storage",
        "auth":{
          "access_key":"<MINIO_ACCESS>",
          "secret_key":"<MINIO_SECRET_KEY>",
          "endpoint":"https://<MINIO_ENDPOINT>"
        }
      }
    ]
  },
  "output":[
    {
      "storage_name":"minio-storage",
      "path":"my-bucket-1",
      "suffix":[
        "wav"
      ]
    },
    {
      "storage_name":"minio-storage",
      "path":"my-bucket-2",
      "suffix":[
        "avi"
      ],
      "prefix":[
        "video-"
      ]
    }
  ]
}
```

And create an [OpenFaaS secret](https://docs.openfaas.com/reference/secrets/) with the content of the file using [`faas-cli`](https://github.com/openfaas/faas-cli):

```bash
faas-cli secret create multi-out-faas-config --from-file=<CONFIG_FILE>
```

### Deploying the function

To deploy the function in OpenFaaS you can use our publicly available Docker image [`grycap/multi-out-faas`](https://hub.docker.com/r/grycap/multi-out-faas) or yours if you have previously generated it. In order to deploy, the file `multi-out-faas.yml` has to be edited to add the endpoint of the OpenFaaS gateway: 

```yaml
version: 1.0
provider:
  name: openfaas
  gateway: http://<OPENFAAS_GATEWAY_ENDPOINT>
functions:
  multi-out-faas:
    lang: go
    handler: .
    image: grycap/multi-out-faas
    secrets:
    - multi-out-faas-config
    environment:
      CONFIG_FILE: multi-out-faas-config
```

Finally execute:

```bash
faas-cli deploy -f multi-out-faas.yml
```

### Sending events to the function

> Currently, the function only supports [MinIO](https://min.io/) as storage provider, but integration with [Amazon S3](https://aws.amazon.com/s3/) and [Onedata](https://onedata.org/#/home) (through [OneTrigger](https://github.com/grycap/onetrigger)) is coming soon.

- **MinIO:** Configure a bucket for sending events to a webhook (the multi-out-faas function endpoint). You can follow [this guide](https://docs.min.io/docs/minio-bucket-notification-guide.html#webhooks).