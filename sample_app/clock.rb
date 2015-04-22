require 'clockwork'
include Clockwork

every(10.seconds, 'Just say something') do
  puts "Hey Yo"
end
