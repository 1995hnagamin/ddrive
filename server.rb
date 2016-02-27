require 'sinatra'
require 'sinatra/reloader'
require 'digest/sha2'

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

put '/upload' do
  body = request.body.read
  shasum = Digest::SHA2.hexdigest(body)
  save_path = filepath(shasum)
  if File.exist? save_path
    return "Already exists"
  end
  begin
    File.open(save_path, 'wb') do |file|
      file.write body
    end
    "Successfully uploaded"
  rescue
    "Something wrong"
  end
end
