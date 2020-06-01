# ecrlogin
This GO project will help login to an AWS ECR by passing the AWS profile credentials and a region.

The usage is:
`ecrlogin [PROFILE] [REGION]`

## Required
* [Docker](https://docs.docker.com/get-docker/)

* [GO](https://golang.org/doc/install)

* Add the following to `.bashrc`, `.bash_profile`, or run via CLI

```
export GOPATH="$HOME/go"

export PATH="$PATH:$GOPATH/bin"
```

## Installing
* `go get github.com/FabulousGinger/ecrlogin`

## Config
* Have profiles listed in `~/.aws/credentials`
```
[dev]
aws_access_key_id = VALUE
aws_secret_access_key = VALUE

[staging]
aws_access_key_id = VALUE
aws_secret_access_key = VALUE

[production]
aws_access_key_id = VALUE
aws_secret_access_key = VALUE
```

* The AWS region is optional, if there is a defualt region, it can be set in `~/.aws/config`, 
and this GO project will use that region.
```
[default]
region = us-east-1
```

## Usage
Provide the profile to login to ECR.
The region is optional and if not set the `~/.aws/config` will be loaded.

`ecrlogin [PROFILE] [REGION]`
