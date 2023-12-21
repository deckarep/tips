# tips OR tailtop
The command-line tool to wrangle your Tailscale tailnet cluster whether large or small.

### What is tips?
Any Tailscale user whether a hobbyist with a 3 node cluster or a seasoned cloud professional managing thousands of 
production nodes can benefit from this tool. `tips` is the go-to tool to quickly and effectively manage a `tailnet`
cluster of any size. It allows you to confidently slice and dice nodes, filter/nodes, remotely execute 
commands and manage your nodes collectively using an effective pattern modeled after cloud automation software.

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

### Supported/Tested OS's
- [x] MacOS (tested, actively developed)
- [ ] Linux (planned soon)
- [ ] PC (future planned: contributions welcome)

### Built with ❤️
* by deckarep

### Disclaimer: Independent Project
Please note that this project is a personal and independent initiative. It is not endorsed, sponsored, affiliated with, 
or otherwise associated with any company or commercial entity. This project is developed and maintained by individual 
contributors in their personal capacity. The views and opinions expressed here are those of the individual contributors 
and do not reflect those of any company or professional organization they may be associated with.