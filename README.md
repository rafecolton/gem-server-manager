gem-server-manager
==================

Helps manage your gem server by fetching gems

`gem-server-manager` ("`gsm`") consumes instructions from an
[AMQP](http://en.wikipedia.org/wiki/Advanced_Message_Queuing_Protocol)
queue regarding which of your applications have been updated, so it can
pull down the requisite gems into your GEMDIR.

The Gem Library Updater (GLU) consumes from an AMQP queue.  It expects
the message body to be a JSON payload that looks something like this:

```json
{
  "rev": "master",
  "repo_name": "gem-server-manager",
  "repo_owner": "rafecolton"
}
```

It then users a GitHub API key to pull down the raw `Gemfile` and
`Gemfile.lock` from your repo.  It then pulls down your gems from
Rubygems and uploads them to your
[geminabox](https://github.com/geminabox/geminabox) server.
