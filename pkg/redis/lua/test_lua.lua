--参数（用户，商品）
--消费成功的用户一秒内不能再消费
--返回值 0-库存不足 1-成功 2-用户冷却中
if  redis.call("EXISTS",KEYS[1]) == 1 then
        return 2
else
        if  redis.call("EXISTS",KEYS[2])==1 and tonumber(redis.call("GET",KEYS[2]))>0  then
                --商品数量-1
                redis.call("DECR",KEYS[2])
                --用户拥有商品数量+1
                redis.call("HINCRBY",KEYS[1].."buy",KEYS[2],1)
                --set一个1s的key表示冷却时间
                redis.call("SET",KEYS[1],1)
                redis.call("EXPIRE",KEYS[1],1)
                return 1
        else
                return 0
        end
end
