NexentaEdge Plugin for Docker Volumes
======================================


## Description
This plugin provides the ability to use NexentaEdge Clusters as backend
devices in a Docker environment over nbd protocol.

## Prerequisites
### Golang
[Get Go](https://golang.org/doc/install)

You will need to make sure you've added the $GOPATH/bin to your path,
AND on Ubuntu you will also need to enable the use of the GO Bin path by sudo;
either run visudo and edit, or provide an alias in your .bashrc file.

You need to pre-create a folder for the GO code.
For example in your .bashrc set the following alias after setting up PATH:
  ```
  export GOPATH=<your GO folder>
  export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin/
  alias sudo='sudo env PATH=$PATH'
  ```
Source .bashrc to update variables
```
  source .bashrc
```


### GCC
```
  apt-get install gcc
```
NOTE:
Should be run as root and command may differ depending on your OS. 

### Docker
You can find instructions and steps on the Docker website here:
[Get Docker](https://docs.docker.com/engine/)

## Driver Installation
After the above Prerequisites are met, clone repository and use the Makefile:
  ```
  git clone https://github.com/nexenta/nedge-docker-volume
  cd nedge-docker-volume
  make
  ```

In addition to providing the source, this should also build and install the
ndvol binary in your Golang bin directory.

## Configuration
Example config file can be found here:
  ```
  https://github.com/Nexenta/nedge-docker-volume/blob/master/ndvol/daemon/ndvol.json
  ```
  
Default path to config file is
  ```
  /opt/nedge/etc/ccow/ndvol.json
  ```

## Starting/stopping the daemon
After install and setting up the configuration file, start the 
nexenta-docker-driver daemon

  ```
  sudo ndvol daemon start -v
  ```

To stop the daemon

 ```
 sudo ndvol daemon stop
 ```


## Usage Examples
Now that the daemon is running, you're ready to issue calls via the Docker
Volume API and have the requests serviced by the NexentaEdge Driver.

For a list of avaialable commands run:
  ```
  docker volume --help
  ```

Here's an example of how to create a Nexenta volume using the Docker Volume
API:
  ```
  docker volume create -d ndvol --name=testvolume -o size=1GB -o repCount=2 -o ratelim=1000
  ```

Now in order to use that volume with a Container you simply specify
  ```
  docker run -v testvolume:/Data --volume-driver=ndvol -i -t ubuntu
  /bin/bash
  ```

Note that if you had NOT created the volume already, Docker will issue the
create call to the driver for you while launching the container.  The Driver
create method checks the Nexenta backend to see if the Volume already exists,
if it does it just passes back the info for the existing volume, otherwise it
runs through the create process and creates the Volume on the Nexenta
backend.
