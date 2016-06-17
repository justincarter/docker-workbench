# Docker Workbench

`docker-workbench` is a utility for simplifying the creation of Docker-based development environments in VirtualBox with `docker-machine`.

The primary goals of `docker-workbench` are;

1. To make it easy to create Docker machines in VirtualBox with sensible defaults (CPUs, disk size, RAM, etc)
2. To make it easy to run multiple containerised web applications without managing DNS, hosts files or ports
3. To provide a standard `/workbench` shared folder to allow `docker-compose` volumes work the same for multiple users, cross-platform 


## Installation

`docker-workbench` is written in Go. To install Go;

1. Download and install Go (https://golang.org/)
2. Create a `workspace` folder in a place of your choosing (e.g. `c:\workspace` or `~/workspace`)
3. Set the `GOPATH` environment variable to the path you created
4. Append to your `PATH` environment variable `%GOPATH%\bin` (for Windows) or `$GOPATH/bin` (for Linux/Mac)

To install `docker-workbench`;

1. Run `go get -u github.com/justincarter/docker-workbench`


## Requirements

To use `docker-workbench` you will also need to install the following;

### For Windows and Mac

1. Docker Toolbox (docker, docker-machine, docker-compose) (https://www.docker.com/products/docker-toolbox)
2. Oracle VirtualBox 5.x (https://www.virtualbox.org/)
3. Git Bash (https://git-for-windows.github.io/) (for Windows only)

### For Linux

1. Docker Engine (https://docs.docker.com/engine/installation/)
2. Docker Machine (https://docs.docker.com/machine/install-machine/)
3. Docker Compose (https://docs.docker.com/compose/install/)
4. Oracle VirtualBox 5.x (https://www.virtualbox.org/)


## Usage

    docker-workbench v1.0
    Provision a Docker Workbench for use with docker-machine and docker-compose

    Usage:
    docker-workbench [options] COMMAND

    Options:
    --help, -h    show help
    --version, -v print the version

    Commands:
    create        Create a new workbench machine in the current directory
    up            Start the workbench machine and show details
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

With Docker Workbench is a requirement to use highly consistent naming for folders and host headers because this makes configuration obvious.

The `docker-compose.yml` file for "myapp" looks like this:

    myapp:
      image: lucee/lucee4-nginx
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

When the application finishes starting up you be able to browse to the app using the URL above, and output similar to below will appear in the console;

    myapp_1 | lucee-server-root:/opt/lucee/server/lucee-server
    myapp_1 | ===================================================================
    myapp_1 | SERVER CONTEXT
    myapp_1 | -------------------------------------------------------------------
    myapp_1 | - config:/opt/lucee/server/lucee-server/context
    myapp_1 | - loader-version:4.3
    myapp_1 | ===================================================================
    myapp_1 |
    myapp_1 | 2016-06-17 08:47:17,423 INFO success: nginx entered RUNNING state, process has stayed up for > than 1 seconds (startsecs) myapp_1 | 2016-06-17 08:47:17,423 INFO success: lucee entered RUNNING state, process has stayed up for > than 1 seconds (startsecs) myapp_1 | Fri Jun 17 08:47:17 UTC 2016-464 using JRE Date Library
    myapp_1 | Fri Jun 17 08:47:17 UTC 2016-756 Start CFML Controller
    myapp_1 | Fri Jun 17 08:47:17 UTC 2016 Loaded Lucee Version 4.5.2.018
    myapp_1 | ===================================================================
    myapp_1 | WEB CONTEXT (cbe856ff790c9ba5208811309bdf168b)
    myapp_1 | -------------------------------------------------------------------
    myapp_1 | - config:/opt/lucee/web (custom setting)
    myapp_1 | - webroot:/var/www/
    myapp_1 | - hash:cbe856ff790c9ba5208811309bdf168b
    myapp_1 | - label:cbe856ff790c9ba5208811309bdf168b
    myapp_1 | ===================================================================
    myapp_1 |
    myapp_1 | 17-Jun-2016 08:47:18.029 INFO [main] org.apache.coyote.AbstractProtocol.start Starting ProtocolHandler ["http-apr-8888"]
    myapp_1 | 17-Jun-2016 08:47:18.034 INFO [main] org.apache.coyote.AbstractProtocol.start Starting ProtocolHandler ["ajp-apr-8009"]
    myapp_1 | 17-Jun-2016 08:47:18.035 INFO [main] org.apache.catalina.startup.Catalina.start Server startup in 1003 ms

Any containerised web application that listens on port 80 should be able to work with Docker Workbench. 


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

