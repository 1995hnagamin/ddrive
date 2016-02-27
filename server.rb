require 'sinatra'
require 'sinatra/reloader'

def filepath(filename)
  return "./objects/#{filename}"
end

get '/' do
  'Hello World'
end

get '/get/:filename' do |filename|
  path = filepath(filename)
  if File.exist? path
    send_file(path)
  else
    status 404
    "Not Found"
  end
end
