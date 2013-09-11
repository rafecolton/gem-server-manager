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
//TODO: fill this in
```

It then users a GitHub API key to pull down the raw `Gemfile` and
`Gemfile.lock` from your repo.  It then runs `bundle install`, using
some args (TBD) to place the install gems into your GEMDIR.
