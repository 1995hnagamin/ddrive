require 'sinatra'
require 'sinatra/reloader'
require 'digest/sha2'
require 'fileutils'

class BadFileNameError < StandardError; end

def filepath(filename)
  head = filename[0..1]
  tail = filename[2..-1]
  return "./objects/#{head}/#{tail}"
end

get '/' do
  'Hello World'
end

get '/get/:filename' do |filename|
  begin
    raise BadFileNameError unless filename =~ /^[0-9a-f]{64}$/
    path = filepath(filename)
    if not File.exist? path
      status 404
      "Not Found"
    end
    send_file(path)
  rescue BadFileNameError => e
    status 404
    "Bad filename"
  end
end

put '/upload' do
  body = request.body.read
  shasum = Digest::SHA2.hexdigest(body)
  save_path = filepath(shasum)
  if File.exist? save_path
    status 208
    return "Already exists"
  end
  begin
    FileUtils.mkdir_p(File.dirname(save_path))
    File.open(save_path, 'wb') do |file|
      file.write body
    end
    status 201
    shasum
  rescue
    status 406
    "Something wrong"
  end
end
