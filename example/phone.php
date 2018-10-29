<?php

try {
    $redis_handle = new Redis();
    $redis_handle->connect('127.0.0.1', 8799, 30);
    $result = $redis_handle->hgetall('15618228951');
    print_r($result);
} catch(Throwable $e) {
    echo $e->getMessage() . PHP_EOL;
}