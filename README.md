# chef2consul: copies chef data struct to consul format (key => val)


## To run:

```bash
KNIFERB_FILE=<path to knife.rb file> CONSUL_PREFIX=<destination prefix> CONSUL_HOST=<consul server> CONSUL_TOKEN=<consul token> chef2consul <chef node name> <chef attribute>
```

Un example:

```sh
KNIFERB_FILE=~/.chef/knife.rb CONSUL_PREFIX=prefix/path CONSUL_HOST=consul.example.com:8500 CONSUL_TOKEN=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee chef2consul.go chef-node-01.example.com attribute
```
