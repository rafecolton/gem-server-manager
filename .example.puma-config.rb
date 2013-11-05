#!/usr/bin/env puma

directory '~gemserver/puma-data'

environment 'development' # default

daemonize false # default

pidfile '~gemserver/gemserver.pid'

state_path '~gemserver/gemserver.state'

threads 0, 16 # default

activate_control_app 'unix://~gemserver/pumactl.sock', { no_token: true }
