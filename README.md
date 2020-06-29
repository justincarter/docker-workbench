# Docker Workbench

[![Build Status](https://travis-ci.org/justincarter/docker-workbench.svg?branch=master)](https://travis-ci.org/justincarter/docker-workbench)
[![codebeat badge](https://codebeat.co/badges/271f8ad5-385f-4edb-89da-12d8ce8fa654)](https://codebeat.co/projects/github-com-justincarter-docker-workbench)
[![Go Report Card](https://goreportcard.com/badge/github.com/justincarter/docker-workbench)](https://goreportcard.com/report/github.com/justincarter/docker-workbench)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/justincarter/docker-workbench/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/justincarter/docker-workbench.svg?maxAge=86400)](https://github.com/justincarter/docker-workbench/releases/latest)

`docker-workbench` is a utility for simplifying the creation of Docker-based development environments in VirtualBox with `docker-machine`.

The primary goals of `docker-workbench` are;

1. To make it easy to create Docker machines in VirtualBox with sensible defaults (CPUs, disk size, RAM, etc)
2. To make it easy to run multiple containerised web applications without managing DNS, hosts files or ports
3. To provide a standard `/workbench` shared folder to allow `docker-compose` volumes work the same for multiple users, cross-platform
4. To allow mobile, tablet and other network devices to easily access the containerised applications


## Installation

### 1. Install Go

`docker-workbench` is written in Go. To install and set up Go;

1. Download and install Go (https://golang.org/)
2. Create a "workspace" folder to store your Go source and binaries (e.g. `c:\workspace` or `~/workspace`)
3. Set the `GOPATH` environment variable to the path of the "workspace" folder you created
4. Append to your `PATH` environment variable `%GOPATH%\bin` (for Windows) or `$GOPATH/bin` (for Linux/Mac)

### 2. Install Docker Workbench

To install `docker-workbench` using Go;

    go get -u github.com/justincarter/docker-workbench

### 3. Install Additional Requirements

To use `docker-workbench` you will also need to install the following;

#### For Windows

1. Oracle VirtualBox 5.x (https://www.virtualbox.org/)
2. Git Bash (https://git-for-windows.github.io/)
3. Docker CLI Tools (docker, docker-machine, docker-compose) 

You can install the Docker CLI Tools by first installing Chocolatey;  
https://chocolatey.org/

Then from an Administrative Command Prompt;

    C:\> chocolatey install docker
    C:\> chocolatey install docker-machine
    C:\> chocolatey install docker-compose

Alternatively you can install the legacy Docker Toolbox for Windows (https://docs.docker.com/toolbox/toolbox_install_windows/). Note: Do not enable Hyper-V when prompted, otherwise Virtualbox will not work.

#### For Mac

1. Oracle VirtualBox 5.x (https://www.virtualbox.org/)
2. Docker CLI Tools (docker, docker-machine, docker-compose) 

You can install the Docker CLI Tools by first installing Homebrew;  
https://brew.sh/

    $ brew install docker
    $ brew install docker-machine
    $ brew install docker-compose

Alternatively you can install the legacy Docker Toolbox for Mac (https://docs.docker.com/toolbox/toolbox_install_mac/).

#### For Linux

1. Docker Engine (https://docs.docker.com/install/)
2. Docker Machine (https://docs.docker.com/machine/install-machine/)
3. Docker Compose (https://docs.docker.com/compose/install/)
4. Oracle VirtualBox 5.x (https://www.virtualbox.org/)


## Usage

    docker-workbench v1.5
    Provision a Docker Workbench for use with docker-machine and docker-compose

    Usage:
    docker-workbench [options] COMMAND

    Options:
    --help, -h    show help
    --version, -v print the version

    Commands:
    create        Create a new workbench machine in the current directory
    up            Start the workbench machine and show details
    proxy         Start a reverse proxy to the app in the current directory
    help          Shows a list of commands or help for one command

    Run 'docker-workbench help COMMAND' for more information on a command.


## Creating a Docker Workbench

A Docker Workbench is created in the context of the working directory from which the `docker-workbench create` command is run. The working directory is automatically configured as a shared folder inside the VM that is always named `/workbench`.

The following example will create a Docker Workbench called `workbench` (named automatically from `d:\workbench`);

    $ mkdir /d/workbench
    $ cd /d/workbench
    $ docker-workbench create

When `docker-workbench create` is run, it does a few simple things;

- The default CPU cores for the VM is set to 2
- The default RAM for the VM is set to 2GB
- The default disk size for the VM is set to 60GB
- The `docker-machine` command is run to create the VM
- The Docker Workbench reverse proxy container is installed and set to always run
- The `/workbench` shared folder is set to the working directory

You may use any standard Docker Machine environment variables used by the Oracle VirtualBox driver to customise the machine creation (e.g. default to use more cores, more RAM, etc) (https://docs.docker.com/machine/drivers/virtualbox/)


## Running an application

The Docker Workbench directory can contain multiple applications, each in their own sub-directory (typically cloned from a Git respoistory).

Here is an example of a trivial web application called "myapp" using Lucee 4.5 and Nginx, which has been placed inside the `workbench` directory;

    myapp
    ├── docker-compose.yml
    └── www
        └── index.cfm

The reverse proxy included with Docker Workbench will automatically route traffic to any containers that listen on port 80 and are configured with a `VIRTUAL_HOST` environment variable that specifies the host headers (wild cards allowed) that the application should respond to.

With Docker Workbench it is a requirement to use highly consistent naming for folders and host headers because this makes configuration obvious.

The `docker-compose.yml` file for "myapp" looks like this:

    myapp:
      image: lucee/lucee:nginx
      environment:
        - "VIRTUAL_HOST=myapp.*"
      volumes:
        - "/workbench/myapp/www:/var/www"

Note the consistent naming; "myapp" is the directory name of the application, which matches the service name at the top of the .yml file, the environment variable `VIRTUAL_HOST` wildcard prefix, and also the parth used in the volume which maps the "www" folder into the container.

Before running the application, we can ensure that the Docker Workbench is running and get some useful info about it by running the `docker-workbench up` command from with the "myapp" directory.

    $ cd myapp
    $ docker-workbench up
    Starting "workbench"...
    Machine "workbench" is already running.

    Run the following command to set this machine as your default:
    eval "$(docker-machine env workbench)"

    Start the application:
    docker-compose up

    Browse the workbench using:
    http://myapp.192.168.99.100.nip.io/

The next step is to set the machine as the default, which will set environment variables that allow us to work with docker, docker-machine and docker-compose. The output above tells us the command to run;

    $ eval "$(docker-machine env workbench)"

The output above also tells us the URL that the application will be available on when it is running. 

    http://myapp.192.168.99.100.nip.io/

This URL is made up of the value supplied in the VIRTUAL_HOST which must match the directory of the application ("myapp"), the IP address of the Docker Workbench VM (assigned by VirtualBox using DHCP from the Docker Machine network adapter), and ".nip.io" which is a wildcard DNS service that resolves names to their matching IP addresses. This means we do not have to manage our own DNS or hosts files or bother with unique, difficult to remember port numbers for each application.

The final step is to start the application using Docker Compose, as mentioned in the output above;

    $ docker-compose up

When the application finishes starting up you will be able to browse to the app using the URL above, and output similar to below will appear in the console;

    myapp_1 | lucee-server-root:/opt/lucee/server/lucee-server
    myapp_1 | ===================================================================
    myapp_1 | SERVER CONTEXT
    myapp_1 | -------------------------------------------------------------------
    myapp_1 | - config:/opt/lucee/server/lucee-server/context
    myapp_1 | - loader-version:6.1
    myapp_1 | ===================================================================
    myapp_1 |
    myapp_1 | 2010-01-25 08:47:17,423 INFO success: nginx entered RUNNING state, process has stayed up for > than 1 seconds (startsecs) 
    myapp_1 | 2010-01-25 08:47:17,423 INFO success: lucee entered RUNNING state, process has stayed up for > than 1 seconds (startsecs) 
    myapp_1 | Fri Jan 25 08:47:17 UTC 2019-464 using JRE Date Library
    myapp_1 | Fri Jan 25 08:47:17 UTC 2019-756 Start CFML Controller
    myapp_1 | Fri Jan 25 08:47:17 UTC 2019 Loaded Lucee Version 5.3.1.92
    myapp_1 | ===================================================================
    myapp_1 | WEB CONTEXT (cbe856ff790c9ba5208811309bdf168b)
    myapp_1 | -------------------------------------------------------------------
    myapp_1 | - config:/opt/lucee/web (custom setting)
    myapp_1 | - webroot:/var/www/
    myapp_1 | - hash:cbe856ff790c9ba5208811309bdf168b
    myapp_1 | - label:cbe856ff790c9ba5208811309bdf168b
    myapp_1 | ===================================================================
    myapp_1 |
    myapp_1 | 25-Jan-2019 08:47:18.029 INFO [main] org.apache.coyote.AbstractProtocol.start Starting ProtocolHandler ["http-apr-8888"]
    myapp_1 | 25-Jan-2019 08:47:18.034 INFO [main] org.apache.coyote.AbstractProtocol.start Starting ProtocolHandler ["ajp-apr-8009"]
    myapp_1 | 25-Jan-2019 08:47:18.035 INFO [main] org.apache.catalina.startup.Catalina.start Server startup in 1003 ms

Any containerised web application that listens on port 80 should be able to work with Docker Workbench. 


## Run a simple reverse proxy

Docker Workbench has a simple reverse proxy built-in which can be useful for allowing other network devices on your LAN (other PCs, tablets, phones, etc) to access the applications running inside your Docker Machine VMs.

The `proxy` command works just like the `up` command and will detect the application and workbench details automatically, giving you a list of addresses that the proxy is listening on;

    $ docker-workbench proxy
    Starting reverse proxy on port 8080...
    Listening on:

    http://myapp.192.168.0.10.nip.io:8080/

    Press Ctrl-C to terminate proxy

In this example another network device would be able to browse to `http://myapp.192.168.0.10.nip.io:8080/` to see the application running.

Note that the default Docker Machine network interface (usually `192.168.99.1`) and other loopback (`127.0.0.1`) or link local (`169.x.x.x`) addresses will not be shown here.

You can also start the proxy on a port number other than the default `8080` by using the `--port` or `-p` flag:

    $ docker-workbench proxy -p 9001


## Advanced Usage

Docker Workbench is basically a utility for easily creating VMs using `docker-machine` and a helper for getting the commands and URLs necessary for running web applications with minimal configuration.

You can use any `docker-machine` and `docker-compose` command directly, where the "machine name" is always the name of parent directory of your applications. Some example commands are;

    docker-machine ls
    docker-machine inspect workbench
    docker-machine ssh workbench 

    docker-compose build
    docker-compose config

For further info view the Docker Machine and Docker Compose reference documentation;

- https://docs.docker.com/machine/reference/
- https://docs.docker.com/compose/reference/overview/

### Multiple Docker Workbenches

For situations where you have many applications and you want to run them in separate VMs (e.g. a VM per client, or a VM per group of related applications) you can use `docker-workbench create` to create a workbench from any directory. A simple way of managing your workbenches might be to have a `workbench` folder with several folders inside named by client or application group, and inside each of those a folder for each application. For example;

    workbench
    ├── clientA
    │   ├── myapp1
    │   └── myapp2
    └── clientB
        └── anotherapp

In this scenario you would create two VMs by running `docker-workbench create` from inside the `clientA` and `clientB` folders. Other than this, there is no difference to creating and using just a single VM -- all configuration is done exactly the same as described above.

It's worth noting that even though your machines in this scenario would be called `clientA` and `clientB`, the shared folder inside the VM which is referred to in your `docker-compose.yml` file will always be named `/workbench` (the shared folder name inside the VM is not named after the VM). A `docker-compose.yml` file for clientB's "anotherapp" might look like this;

    anotherapp:
        image: lucee/lucee:nginx
        environment:
            - "VIRTUAL_HOST=anotherapp.*"
        volumes:
            - "/workbench/anotherapp/www:/var/www"

This also means that you can always run an app inside a Docker Workbench, regardless of what its name is, without modifying the `docker-compose.yml` file.


## Troubleshooting

### Destroy and recreate your Docker Workbench

There are a number of reasons that a Docker VM or your Docker Workbench may get into a bad state, such as invalid networking configurations, a full virtual disk, a missing or accidentally deleted docker-workbench-proxy container, etc. The Docker Workbench can and should be recreated often to update to newer versions of Docker or to resolve issues that can't be easily debugged by the end user.

To destroy your Docker Workbench and recreate it fresh with the latest version of boot2docker, use `docker-machine` to remove the machine by name (in this case our machine is relative to `/d/workbench` and is called `workbench`), and then `docker-workbench` to create it again;

    $ cd /d/workbench
    $ docker-machine rm workbench
    $ docker-workbench create

Within a few minutes you should be back up and running as normal.

### Keep your CLI Tools up to date

Remember to keep your Docker CLI tools up to date as newer versions of Docker and boot2docker are released. Upgrading to the latest versions on Windows using Chocolatey can be done from an Administrative Command Prompt;

    C:\> chocolatey upgrade docker
    C:\> chocolatey upgrade docker-machine
    C:\> chocolatey upgrade docker-compose

Or similarly on MacOS using brew;

    brew upgrade docker
    brew upgrade docker-machine
    brew upgrade docker-compose

This should avoid errors with older Docker client tools trying to connect to newer Docker servers, which may throw an error such as;

> Error checking TLS connection: Error checking and/or regenerating the certs: There was an error validating certificates for host "192.168.99.100:2376": x509: certificate has expired or is not yet valid

> ERROR: SSL error: [SSL: TLSV1_ALERT_PROTOCOL_VERSION] tlsv1 alert protocol version (_ssl.c:661)
