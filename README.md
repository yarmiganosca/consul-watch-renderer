## Turns out, consul-template doesn't :heart: Consul Watches

`consul-template`, for all it's power, isn't the ideal handler for Consul Watches:
1. In `-once` mode (I don't know how else you'd run it as a watch handler) it doesn't always actualy render the template. I'm sure the reasons for (sometimes) not doing so are good and make sense, but (to my knowledge) there isn't a way to tell `consul-template` "no, don't think about it, just render the template".
2. To my knowledge it has no way to read it's data directly from stdin, which is where Watches send their data to handlers.
3. The added network call (because it can't render from the stdin-provided data) is both unnecessary, and a potential source of race conditions.

Note: If I'm wrong, and there is a way to make `consul-template` behave simply and predictably and read from stdin, for the love of god, please email me and tell me how to do it.

## So I wrote this very dumb tool.

You use it thusly in a consul watch config file:
```json
{
  "watches": [
    {
      "type": "service",
      "service": "my-service",
      "handler": "consul-watch-renderer path/to/template/file path/to/destination"
    }
  ]
}
```
and that's it! If you must use it yourself on the command line (for testing perhaps), here's a way to do that:
```shell
$ consul-watch-renderer template destination < <(cat file.json)
```
