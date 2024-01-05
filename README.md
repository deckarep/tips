[![Go](https://github.com/deckarep/tips/actions/workflows/go.yml/badge.svg)](https://github.com/deckarep/tips/actions/workflows/go.yml)

<p align="center">
  <img src="assets/tips_logo.png" width="512"/>
</p>

#                                                                                                                                                                                                                         

<p align="center">
The command-line tool to wrangle your Tailscale tailnet cluster whether large or small.
</p>

<p align="center">
  <img src="assets/main.png" width="700"/>
</p>

### What is tips?
Any Tailscale user whether a hobbyist with a 3 node cluster or a seasoned cloud professional managing thousands of 
production nodes can benefit from this tool. `tips` is the go-to tool to quickly and effectively manage a `tailnet`
cluster of any size. It allows you to confidently slice and dice nodes, filter/nodes, remotely execute 
commands and manage your nodes collectively using an effective pattern modeled after cloud automation software.

### Features
* 📝 - Nearly zero-config managed via [Viper](https://github.com/spf13/viper) and [Cobra](https://github.com/spf13/cobra), with sane defaults
* ⚡️ - Caching and indexed via [BBolt](https://github.com/etcd-io/bbolt) key-value store for fast queries
* 🔎 - Filtering, sorting, slicing and dicing built-in
* 👩‍💻‍ - Developer focused tooling
* ✅ - Remote command execution and ssh tooling built-in
* 🫦 - Beautifully rendered terminal output thanks to [Charmbracelet's lipgloss](https://github.com/charmbracelet/lipgloss)
* 👍🏼 - No deployable dependencies, single binary thanks to the [Go programming language](https://go.dev/)

### Installation

First, grab the `tips` binary by doing a typical install from source:
```sh
$ go install github.com/deckarep/tips@latest
```
Next, log into your Tailscale account and create an `API Access` token.
1. [Link to create API Access token](https://login.tailscale.com/admin/settings/keys)
2. Under `API Access Tokens` click: `Generate access token...`
3. Provide a `Description` like `tips-cli` and choose an expiry time or keep the `90` days default.
4. After clicking the Generate button copy the key that looks like:
`tskey-api-xxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` and save it in the `tips` config file.
5. Follow security best-practices by not sharing this key or copying it into a public repo!
6. `tips` uses this key to remotely access the Tailscale Web API.

At an **absolute minimum** `tips` requires the `tips_api_key` and a `tailnet` to be specified in order to use this app. 
Here is how you can set this up.

As an environment variable:
```sh
export TIPS_API_KEY=tskey-api-xxxxxxxxxxxx-xxx...
export TAILNET=user@domain.com
```

Or, in the `~/.tips.cfg` `JSON` file placed in your homedir:
```json
{
    "tips_api_key": "tskey-api-xxxxxxxxxxxx-xxx...",
    "tailnet": "user@domain.com"
}
```

Lastly, on-demand by using the relevant flags:
```sh
# This will get tedious and not recommended long term
$ ./tips --tailnet user@domain.com --tips_api_key tskey-api-xxxxxxxxxxxx-xxx...
```

### Definitions
* **[Tailscale](https://tailscale.com/)**: makes creating software-defined networks easy: securely connecting users, services, and devices
* **[tailnet](https://tailscale.com/kb/1136/tailnet)**: a single private network built from one or more nodes using Tailscale
* **[tips](https://github.com/deckarep/tips)** (this tool): a command-line tool to easily manage a tailnet cluster for use on Mac, PC, or Linux

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

How do I get a list of *all* `devices` or `nodes` in a `tailnet`?
```sh
./tips
./tips @ # This is equivalent as @ means all/everything.
```

However, it's better to query by a `full name` or `prefix` especially if you have a large infrastructure
```sh
# Simply provide a partial or full string name.
./tips [prefix-filter]

./tips blade # Find all nodes with a machine name starting with 'blade'

./tips bla # Find all nodes with a machine name starting with 'bla'
```

How do I get more details?
```sh
# Not yet supported
# ./tips --details
```

How do I sort the output?
```sh
# Partially working (some fields not supported)
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
# Not yet supported
# ./tips --columns 'name,address'
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
# Not yet supported
# ./tips --reindex
```

### Supported/Tested OS's
- [x] MacOS (actively developed, tested)
- [ ] Linux (planned soon, untested, may work to some degree)
- [ ] PC (future planned, untested, contributions welcome)

### Built with ❤️
* by deckarep

### F.A.Q.
**Q: I'm having trouble executing remote commands on some nodes in my `tailnet`.**

**A:** This is typically not a problem with this tool, but rather your `tailnet` configuration with respect to 
`permissions` or `tagging` or `ssh auth keys`. Please check that your nodes have the appropriate `ports open`, 
`permissions` and/or `public keys` to match your user `logon` credentials.

**Q: I'm having trouble `sshing` into a node with this tool.**

**A:** This tool does absolutely nothing special to manage ssh-based logins and `simply forwards` such requests to 
either the `Tailscale ssh subcommand` or the `native ssh executable`. Please see the previous question as for what 
could be wrong.

**Q: I have a massive infrastructure. Will this tool help me manage a large infrastructure with 10's of thousands of 
nodes?.**

**A:** Yes, with *robust caching and indexing built-in*, this tool ensures fast queries amongst even the largest 
clusters. One of the primary goals of this project is that it can help you manage a `tailnet` from *5 to 50,000 
nodes*. It all boils down to the indexing/caching strategies, and it largely depends on query use cases.

**Q: How can I execute remote commands with a greater degree of `parallelism` or `concurrency`?**

**A:** We have you covered: see the `-c` or `--concurrency` flag.

**Q: How can I disable colored output?**

**A:** See the `--nocolor` flag.

**Q: Isn't caching with BBolt overkill?**

**A:** Probably, but integration is easy and some users are expected to be managing thousands of nodes via Tailscale.

**Q: But why u no have unit-tests?**

**A:** Because this is a prototype! I am rapidly designing what I think this software should look like and in this
phase, unit-tests will slow me down. I've already gone through several major refactors and will likely have more coming
up but this will change. Eventually the design will be nailed down and unit-tests **will be a must**!

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