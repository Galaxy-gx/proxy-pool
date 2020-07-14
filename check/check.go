package check

import (
	"fmt"
	"proxy-pool/config"
	"proxy-pool/databases"
	"proxy-pool/model"
	"sync"
	"time"

	"github.com/google/martian/log"
)

// Checker 检查IP可用性
type Checker struct {
	DB   *databases.DB
	Conf *config.Config
}

// NewChecker 检查IP可用性
func NewChecker(db *databases.DB, conf *config.Config) *Checker {
	return &Checker{
		DB:   db,
		Conf: conf,
	}
}

// CheckAll 检查所有IP的可用性
func (c *Checker) CheckAll() {
	log.Infof("check all ip avaliable start...")
	// TODO: to config file
	var wg sync.WaitGroup
	ch := make(chan struct{}, 10) // 并发数量10控制

	proxys := make([]*model.Proxy, 64)
	if err := c.DB.Mysql.Where("is_deleted=?", 0).Find(&proxys).Error; err != nil {
		log.Errorf("get proxys from db %#v", err.Error())
		return
	}
	for _, proxy := range proxys {
		ch <- struct{}{}
		wg.Add(1)
		go func(proxy *model.Proxy) {
			ok, err := c.CheckProxyAvailable(proxy)
			// 代理失效 标记删除
			if err != nil || !ok {
				fmt.Println(proxy.ID, proxy.IP)
				if err := c.DB.Mysql.Table("proxy").
					Where("id=?", proxy.ID).
					Updates(map[string]interface{}{"is_deleted": true}).Error; err != nil {
					log.Errorf("update error %#v", err.Error())
				}
				log.Infof("proxy check faild, IP:%s, Port:%d, ok:%t\n", proxy.IP, proxy.Port, ok)
			} else {
				// 可用更新check时间
				err := c.DB.Mysql.Table("proxy").Where("id=?", proxy.ID).Updates(map[string]interface{}{"check_time": time.Now()}).Error
				if err != nil {
					log.Infof("Updates check_time faild %#v", err.Error())
				}
				log.Infof("proxy check ok,IP:%s, Port:%d\n", proxy.IP, proxy.Port)
			}
			<-ch
			wg.Done()
		}(proxy)
	}
	wg.Wait()
	log.Infof("check all ip avaliable end.")
}
