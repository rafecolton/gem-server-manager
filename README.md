gem-server-with-a-twist
=======================

Wraps your gem server and helps you pull in new gems

0. It wraps your `gem server` in a pretty golang executable
0. It consumes instructions from an
   [AMQP](http://en.wikipedia.org/wiki/Advanced_Message_Queuing_Protocol)
queue regarding which of your applications have been updated, so it can
pull down the requisite gems into your GEMDIR.

## Server Wrapper

The first question you may ask is, "Is this necessary?"  The answer is
"no." The reason I chose two wrap `gem server` in `gswat` is for
convenience.  Since the gem library updating component is meant to aid
your `gem server` running on the same machine, why not package them
together for ease of use.

For more information about `gem server` and your `GEMDIR`, look
[here](http://guides.rubygems.org/run-your-own-gem-server/) or run `gem
server --help`

## Gem Library Updater

If desired, this part can be used without the server wrapper.  The Gem
Library Updater (GLU) consumes from an AMQP queue.  It expects the
message body to be a JSON payload that looks something like this:

```json
//TODO: fill this in
```

It then users a GitHub API key to pull down the raw `Gemfile` and
`Gemfile.lock` from your repo.  It then runs `bundle install`, using
some args (TBD) to place the install gems into your GEMDIR.
