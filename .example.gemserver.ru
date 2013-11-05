require 'rubygems'
require 'geminabox'

# optional for basic rack auth
use Rack::Auth::Basic, 'Restricted Area' do |username, password|
  [username, password] == %w(myusername mypassword)
end

Geminabox.data = '~gemserver/data'
run Geminabox
