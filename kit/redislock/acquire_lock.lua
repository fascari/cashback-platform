local acquired = redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
if acquired then
    return redis.call('INCR', KEYS[2])
end
return nil
