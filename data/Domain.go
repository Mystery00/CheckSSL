package data

import (
	"context"
	"crypto/tls"
	json2 "encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"

	"CheckSSL/utils"
)

type DomainSave struct {
	Name          string
	LastCheckTime string
	Subject       string
	From          time.Time
	Until         time.Time
	Issuer        string
}

type Domain struct {
	Name          string
	LastCheckTime string
	Subject       string
	From          string
	Until         string
	Remain        string
	Issuer        string
	Valid         bool
}

func (d *Domain) DoCheck(rdb *redis.Client, ctx context.Context) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	seedUrl := "https://" + d.Name
	resp, err := client.Get(seedUrl)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	certInfo := resp.TLS.PeerCertificates[0]
	save := DomainSave{
		Name:          d.Name,
		LastCheckTime: time.Now().Format("2006-02-01 15:04:05.000"),
		Subject:       certInfo.Subject.String(),
		From:          certInfo.NotBefore,
		Until:         certInfo.NotAfter,
		Issuer:        certInfo.Issuer.String(),
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
	d.From = save.From.Format("2006-02-01 15:04:05.000")
	d.Until = save.Until.Format("2006-02-01 15:04:05.000")
	remain := int(time.Until(save.Until).Hours() / 24)
	d.Remain = fmt.Sprintf("%d days", remain)
	d.LastCheckTime = save.LastCheckTime
	d.Subject = save.Subject
	d.Issuer = save.Issuer
	d.Valid = remain > 0
}
