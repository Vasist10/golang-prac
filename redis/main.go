package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main(){
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong , err := rdb.Ping(ctx).Result()
    if err != nil {
        panic(err)
    }

    fmt.Println("redis connected",pong)
    err = rdb.Set(ctx, "name", "Vasist", 0).Err()
    if err != nil {
        panic(err)
    }

    val,err := rdb.Get(ctx, "name").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("name",val)
    rdb.LPush(ctx,"fruits","apple","banana","mango")
    list,err := rdb.LRange(ctx,"fruits",0,-1).Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("fruits",list)
}