# davni

Hey duder. 

Let me know when you fork/clone this so I can delete this repo :)

Anywho ... lets get started. 

A few quick things to note. 

1. Something like this I'd probably create a microservice to serve up the pdfs. Basically, a website you can curl that will generate the contents and then let you download it(even if it is programatically). I can show you how I'd do that if you're interested, but basically I'd use [Genie](https://github.com/kcmerrill/genie) or some similiar [FAAS](https://github.com/search?utf8=%E2%9C%93&q=faas&type=)


2. I'm going to show you the "wrong" way first, only so you can see something I found really cool with docker. Basically, you mount the docker binary from your host machine inside a container(that doesn't have docker installed), along with it's socket. This lets you play with docker inside a container, without actually installing anything inside the container. It also is neat because whenever you exec `docker` it's running on the `host` machine, and not in the container. I'm on a mac, but this should work for all `*nix` env's. 

3. I'll also show you the right way. What you're looking for is a thing called data volumes. Read more about them [here](https://docs.docker.com/engine/admin/volumes/volumes/#create-and-manage-volumes). The TL;DR is that docker creates a "container" and it's soul purpose in life is to store data that multiple containers can read from. So, you can create a volume mounted folder, lets for our example use `/pdfs` in both of our services, so anytime data gets written into `/pdfs` all the containers that have that volume can see it as if it were a regular directory. 


# First Example

Probably the easiest. I think you already know how to shell out in go lang, if not, it looks like this:

```golang
package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
    // local /tmp/pdfs ->container /pdfs
    // notice, /tmp/pdfs .... this binary is running INSIDE a container, and htis container has /tmp/pdfs mounted
	cmd := exec.Command("sh", "-c", "docker run --rm -v /tmp/pdfs:/pdfs -w /pdfs madnight/docker-alpine-wkhtmltopdf", "google.com", "google.pdf")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}
```

But here is the secret sauce. When you run a container(alpine images are wonky, some quick googling you should be able to figure this out for alpine, but for time sake i'm going to use an ubuntu based container).

```-v /usr/bin/docker:/usr/bin/docker -v /var/run/docker.sock:/var/run/docker.sock```

So, that will volume mount the docker binary ... and it will also volume mount the docker socket so when you're _inside_ the container, and you run `docker run` whatever, it's as if you're running docker on your host machine, even though you're inside the container. 

Give this command a go. 
```docker run -ti --rm -v /usr/bin/docker:/usr/bin/docker -v /var/run/docker.sock:/var/run/docker.sock ubuntu```

Now, hop into the container, and play around with docker. Keep in mind, it's as if you're running docker on the _HOST_ machine. 

Here is what mine looks like:

```sh
02:52 PM ✔ kcmerrill  davni ] docker run -ti --rm -v /usr/bin/docker:/usr/bin/docker -v /var/run/docker.sock:/var/run/docker.sock ubuntu
root@51dad591094b:/# docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
51dad591094b        ubuntu              "/bin/bash"         3 seconds ago       Up 3 seconds                            hopeful_lichterman
root@51dad591094b:/#
```

One quick note, mine just worked, so yours should too, but if you're getting a weird error, about not able to talk to the docker socket, or something .... try running this: `apt-get update && apt-get install -y libltdl7` and then try again.

Alright, so lets add all this up .... and this is without any gocode(again, it's a simple shell command at this point)

```sh
root@050d0cc0a6dc:/pdfs# ls
root@050d0cc0a6dc:/pdfs# # notice, I'm INSIDE a container .... 
# ok, this looks a little weird, but you're inside a container, but this command gets run ON THE HOST, hence the /tmp/pdf 
root@050d0cc0a6dc:/pdfs# docker run -ti -v /tmp/pdfs:/pdfs -w /pdfs madnight/docker-alpine-wkhtmltopdf google.com test.pdf
Loading pages (1/6)
Warning: SSL error ignored
Counting pages (2/6)
Resolving links (4/6)
Loading headers and footers (5/6)
Printing pages (6/6)
Done
root@050d0cc0a6dc:/pdfs# ls
test.pdf
root@050d0cc0a6dc:/pdfs#
```

2. Ok, so here is the pure docker implimentation. If you do a `docker-compose up -d` you'll see that the golang service sits and waits, and the htmltopdf will create tester.pdf in the `/pdfs` directory. After a few seconds have passed, the golang service will spit out the contents of it's `/pdfs` folder, and you can see that `tester.pdf` is there. 

Hopefully that shows how you would go about "linking" the two together. 

3. Tying it all together ... 

For this demo, we'll use `docker-compose2.yml` so run the following command: `docker-compose -f docker-compose2.yml up`.

```sh
09:12 PM ✔ kcmerrill  davni ] docker-compose -f docker-compose2.yml up
Recreating davni_golangservice_1 ...
Recreating davni_golangservice_1 ... done
Attaching to davni_golangservice_1
golangservice_1  | test.pdf
golangservice_1  |
davni_golangservice_1 exited with code 0
09:12 PM ✔ kcmerrill  davni ]
```

Ok, so hopefully that makes sense, but if not, let me walk through it once more, hopefully a bit more clearer.

First, write your golang application inside a container. When you run the container, make sure you can play with docker by mounting the docker socket and the docker binary. The next part, make sure to volume mount a shared folder on your host machine, in this case, the "shared" folder between the two is `/tmp/pdfs`. Now, inside your golang service container, run a docker container(which is really running on the host machine) making sure to volume mount the "shared" folder between the two. 


Anywho, EZPZ. Hope this helps :)


*** NOTE ***

When I say the "right"/"wrong" way, I don't mean the volume mounting of the docker binary or the docker socket, but how you go about mounting the volume(you can either use data volumes, the "right" way, or you can just volume mount a folder on run, the "wrong" way).
