[![Go](https://github.com/deckarep/tips/actions/workflows/go.yml/badge.svg)](https://github.com/deckarep/tips/actions/workflows/go.yml)

# tips OR tailtop OR topsntails
The command-line tool to wrangle your Tailscale tailnet cluster whether large or small.

### Definitions
* **[Tailscale](https://tailscale.com/)**: makes creating software-defined networks easy: securely connecting users, services, and devices
* **[tailnet](https://tailscale.com/kb/1136/tailnet)**: a single private network built from one or more nodes using Tailscale
* **[tips](https://github.com/deckarep/tips)** (this tool): a command-line tool to easily manage a tailnet cluster for use on Mac, PC, or Linux

### What is tips?
Any Tailscale user whether a hobbyist with a 3 node cluster or a seasoned cloud professional managing thousands of 
production nodes can benefit from this tool. `tips` is the go-to tool to quickly and effectively manage a `tailnet`
cluster of any size. It allows you to confidently slice and dice nodes, filter/nodes, remotely execute 
commands and manage your nodes collectively using an effective pattern modeled after cloud automation software.

### You'll be able to ...
* Easily view your nodes in a *beautifully rendered* and consistent table view
* View *enriched, realtime* info such as `online status` when ran from the context of a node within a tailnet
* Filter nodes based on: `tags`, `OS`, `hostname` and other fields
* *Slice or segment* nodes to work on a portion of them at a time
* Easily `ssh` into a node
* Execute *single-shot* complex commands against all matching nodes in parallel with controllable concurrency
* *Tail* the logs of long-running sessions from multiple nodes
* Broadcast commands to multiple nodes using the `csshx` power-tool if installed
* Quickly generate a `,` or `\n` delimited list of nodes for reporting or use in other apps/cli tools
* Quickly generate a `json` list of nodes

...with automatic but configurable file-system caching built-in which means fast, consistent results everytime!

### Why the name?
* The name must be short, this tool must not get in the way and will likely be often used to query infrastructure
* Simply put, this tool is about managing a (t)ailnet's distributed (ips) or nodes which shortens to: `tips`
* Lastly, what better way to show appreciation for software than to **leave a tip** especially if used in a
professional or commercial setting?

### How To
Here is a list of common commands from easy to more advanced.

How do I get a list?
```sh
./tips
```

How do I get more details?
```sh
./tips --details
```

How do I query nodes by name?
```sh
# Simply provide a partial or full string name.
./tips bla
```

How do I sort nodes?
```sh
# To sort by one column ascending (default)
./tips --sort 'name'
# To sort by multiple columns with varying order, specifically in ascending or descending order
./tips --sort 'name:dsc,email:asc'
```

How do I slice/partition nodes?
```sh
# Shows the top 5 nodes
./tips --slice 0:5
# Shows the nodes between 5-10
./tips --slice 5:10
```

How do I limit which columns to show?
```sh
./tips --columns 'name,address'
```

How do I generate a JSON-based result
```sh
./tips --json
```
How do I generate a list of ips only
```sh
# Provides a \n delimited list of ips
./tips --ips
# Provides a comma delimited list of ips
./tips --ips --delimiter ','
```
How do run a remote command on all returned nodes?
```sh
./tips [prefix-filter] [remote command here]
./tips blade "hostname" # runs the remote command 'hostname' on all nodes that start with prefix:blade
./tips bla "echo 'hello'" -c20 # same as above but does an echo with a concurrency value of 20.
```

How do I rebuild the index? Running this forces a full rebuild (fetch all remote data) and builds the index
for speedy queries. Normally you don't have to do this manually.
```sh
./tips --index
```

### Supported/Tested OS's
- [x] MacOS (tested, actively developed)
- [ ] Linux (planned soon)
- [ ] PC (future planned: contributions welcome)

### Built with ❤️
* by deckarep

### FAQ



### Alpha
This code is currently being developed in a **rapid prototyping** mode. Therefore you will not see much unit-testing as
too much code is changing and being actively refactored as the design and approach is tightened up. If you use this at
this stage please understand that things are bound to change or be broken until I promote this software to being at the
"beta" stage. Currently only MacOS builds are working.

### Disclaimer: Independent Project
Please note that this project is a personal and independent initiative. It is not endorsed, sponsored, affiliated with, 
or otherwise associated with any company or commercial entity. This project is developed and maintained by individual 
contributors in their personal capacity. The views and opinions expressed here are those of the individual contributors 
and do not reflect those of any company or professional organization they may be associated with.