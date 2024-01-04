### TODO

* `--ipv6` flag for ipv6 results
* nail down the default table header/columns, provide config to enable disable for a user.
* support for themes or turning off colors all together
* ssh into node
* similar to ssh, what about curl/http requests to all selected nodes?
* support `--ssh` (single) or `--csshx` (all selected hosts)
  * `brew install parera10/csshx/csshx` (original csshx is broken on recent macos)
* filter glob syntax: `tips @ 'hostname'`, `tips blade 'hostname'`, `tips tag:peanuts 'hostname'`
  * based filter: `tag:!peanuts`
* slice syntax:
  * `tips blade[5:10]` 'hostname' - returns hosts 5-10 (can't remember if it should be inclusive or not
  * `tips blade[5:]` 'hostname' - returns host 5 on up
  * Currently support slicing via a flag: `--slice 5:10`

### Possible/Proposed Features
* ???