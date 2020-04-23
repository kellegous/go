# A "go" short-link service

## Background
The first time I encountered "go" links was at Google. Anyone on the corporate
network could register a URL shortcut and it would redirect the user to the
appropriate page. So for instance, if you wanted to find out more about BigTable,
you simply directed your browser at http://go/bigtable and you would be redirected to
something about the BigTable data storage system. I was later told that the
first go service at Google was written by [Benjamin Staffin](https://www.linkedin.com/in/benjaminstaffin)
to end the never-ending stream of requests for internal CNAME entries. He
described it as AOL keywords for the corporate network. These days if you go to
any reasonably sized company, you are likely to find a similar system. Etsy made
one after seeing that Twitter had one ... it's a contagious and useful little
tool. So contagious, in fact, that many former Googlers that I know have built
or contributed to a similar system post-Google. I am no different, this is my
"go" link service.

One slight difference between this go service and Google's is that this one is also
capable of generating short links for you.

## Installation
This tool is written in Go (ironically) and can be easily installed and started
with the following commands.

```
GOPATH=`pwd` go install github.com/kellegous/go
bin/go
```

By default, the service will put all of its data in the directory `data` and will
listen to requests on the port `8067`. Both of these, however, are easily configured
using the `--data=/path/to/data` and `--addr=:80` command line flags.

## DNS Setup
To get the most benefit from the service, you should setup a DNS entry on your
local network, `go.corp.mycompany.com`. Make sure that corp.mycompany.com is in
the search domains for each user on the network. This is usually easily accomplished
by configuring your DHCP server. Now, simply typing "go" into your browser should
take you to the service, where you can register shortcuts. Obviously, those
shortcuts will also be available by typing "go/shortcut".

## Using the Service
Once you have it all setup, using it is pretty straight-forward.

#### Create a new shortcut
Type `go/edit/my-shortcut` and enter the URL.

#### Visit a shortcut
Type `go/my-shortcut` and you'll be redirected to the URL.

#### Shorten a URL
Type `go` and enter the URL.
