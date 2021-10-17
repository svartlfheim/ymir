# Ymir

Ymir - a terraform registry for mono-repo's of modules.

## Plans

A terraform registry essentially needs to implement 3 endpoints, according to the terraform registry protocol. See [the spec](https://www.terraform.io/docs/internals/module-registry-protocol.html). 

In short, it's these:
- `/.well-known/terraform.json`: [Service Discovery](https://www.terraform.io/docs/internals/module-registry-protocol.html#service-discovery)
- `/v1/modules/{namespace}/{name}/{provider}/versions`: [List versions](https://www.terraform.io/docs/internals/module-registry-protocol.html#list-available-versions-for-a-specific-module)
- `/v1/modules/{namespace}/{name}/{provider}/{version}/download`: [Downlaod module](https://www.terraform.io/docs/internals/module-registry-protocol.html#download-source-code-for-a-specific-module-version)

When terraform requests a module to the `Download endpoint` it expects to receieve a header in the response `X-Terraform-Get` the value of this header is a url which terraform can send a `GET` request to. This `GET` request will download an archive in `zip` or `tar.gz` format, which terraform unpacks. The contents of this archive provide the module code that terraform will use.

To achieve a registry that works with a mono-repo we should add at least these two commands:

- `ymir module add`: Defines a link between a module and directory in the mono-repo i.e.  `ord/mymodule => github.com/org/repo/mymodule`
- `ymir module publish`: This command publishes a specific version of a module which is basically a link between a commit hash in the repo and the module i.e. `org/mymodule@1.0.0 => mycommithash`.

## Storage 

To keep this registry as simple as possible we should provide 2 mechanisms for storage. The first is the local file system, the second is an S3 bucket.

When we publish a module, `ymir` will download an archive from git at a specific version. It will unpack this archive and extract the directory where the source code for the module lives. This source code will then be placed in the relevant path (within the storage service) where it can be easily located later i.e. `/org/mymodule/1.0.0/*`. 

When terraform requests to download a module we will return a URL which terraform can download the created archive from e.g. `ymir.com/module/download/org/mymodule/1.0.0`. Terraform can then unpack this module and use it as if it were any standard module.

## State 

The state of the registry can also be stored in the desired storage in a JSON/YAML configuration. This should look something like:

```json
{
    "providers": [
        {
            "name": "aws",
            "modules": [
                {
                    "namespace": "org",
                    "name": "mymodule",
                    "repository": "github.com/org/mono-repo",
                    "path": "mymodule",
                    "versions": [
                        {
                            "id": "1.0.0",
                            "commit": "somehash",
                        }
                    ]

                }
            ]
        }
    ]
}
```

> In theory we would only need this file to restore the contents of the registry in a DR scenario. We can simply iterate through the contents recreating the archives in the filesystem.

Ymir should give us the ability to create a mono-repo of terraform modules to use across paddle. 

The mono-repo would be structured something like this:

```
|-- module1/
|---- main.tf
|---- variables.tf
|-- module2/
|---- main.tf
|---- variables.tf
```

