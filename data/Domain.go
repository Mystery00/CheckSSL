package data

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"

	"CheckSSL/config"
	"CheckSSL/utils"
)

type DomainSave struct {
	Name          string
	LastCheckTime string
	Subject       string
	From          int64
	Until         int64
	Issuer        string
	Message       string
}

type Domain struct {
	Name          string
	LastCheckTime string
	Subject       string
	From          string
	Until         string
	Remain        string
	Issuer        string
	Message       string
}

func (d *Domain) DoCheck(rdb *redis.Client, ctx context.Context) {
	json := utils.Cmd("/bin/bash", "-c", config.EnvConfig.ScriptFile+" "+d.Name)
	save := DomainSave{
		Name:          d.Name,
		LastCheckTime: time.Now().Format("2006-02-01 15:04:05.000"),
		Subject:       gjson.Get(json, "subject").String(),
		From:          gjson.Get(json, "start").Int(),
		Until:         gjson.Get(json, "expire").Int(),
		Issuer:        gjson.Get(json, "issuer").String(),
		Message:       gjson.Get(json, "message").String(),
	}

	jsonBytes, err := json2.Marshal(save)
	if err != nil {
		panic(err)
	}
	err = rdb.HSet(ctx, utils.DomainKey, d.Name, jsonBytes).Err()
	if err != nil {
		panic(err)
	}
	utils.LogInfof(ctx, "set ssl status to redis, domainName: %s", d.Name)
	parse(save, d)
}

func (d *Domain) DoCheckLocal(ctx context.Context) {
	var rdb = utils.RedisClient()
	//检查hash是否存在
	duration, err := rdb.TTL(ctx, utils.DomainKey).Result()
	if duration < 0 {
		//key 不过期，设置一个过期时间
		utils.LogInfof(ctx, "set ssl status expire time")
		rdb.Expire(ctx, utils.DomainKey, utils.ExpireTime)
	}
	json, err := rdb.HGet(ctx, utils.DomainKey, d.Name).Result()
	if err == redis.Nil {
		d.DoCheck(rdb, ctx)
		return
	}
	var save DomainSave
	if err != nil {
		panic(err)
	}
	err = json2.Unmarshal([]byte(json), &save)
	if err != nil {
		panic(err)
	}
	utils.LogInfof(ctx, "get ssl status from redis, domainName: %s", d.Name)
	parse(save, d)
}

func parse(save DomainSave, d *Domain) {
	from := time.Unix(0, save.From*int64(time.Millisecond))
	until := time.Unix(0, save.Until*int64(time.Millisecond))
	d.From = from.Format("2006-02-01 15:04:05.000")
	d.Until = until.Format("2006-02-01 15:04:05.000")
	d.Remain = fmt.Sprintf("%d days", int(time.Until(until).Hours()/24))
	d.LastCheckTime = save.LastCheckTime
	d.Subject = save.Subject
	d.Issuer = save.Issuer
	d.Message = save.Message
}
