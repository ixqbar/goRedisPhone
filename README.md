
[TOC]

### purpose
* 手机信息查询

### version
```
v0.0.1
```

### usage
```
./bin/goRedisPhone --config=config.xml
```

### config
```
<?xml version="1.0" encoding="UTF-8" ?>
<config>
    <listen>0.0.0.0:8699</listen>
    <dict>/data/server/redisPhone/dict/phone.dat</dict>
</config>
```

### command
``` 
hgetall 手机号码
reload
total
```

### example
```
<?php

try {
    $redis_handle = new Redis();
    $redis_handle->connect('127.0.0.1', 8799, 30);
    $result = $redis_handle->hgetall('1367152');
    var_dump($result);
} catch(Throwable $e) {
    echo $e->getMessage() . PHP_EOL;
}
```
* type 1移动 2联通 3电信 4电信虚拟运营商 5联通虚拟运营商 6移动虚拟运营商

### deps
* https://github.com/jonnywang/go-kits/redis

### faq
* qq群 233415606