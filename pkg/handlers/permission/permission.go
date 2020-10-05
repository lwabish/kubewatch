package permission

import (
	"fmt"
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

type Permission struct {
	ScName string
	chmod string
	chown string
}

func (p Permission) Init(c *config.Config) error {
	p.ScName=c.Handler.Permission.ScName
	p.chmod=c.Handler.Permission.Chmod
	p.chown=c.Handler.Permission.Chown
	//TODO 环境变量
	//TODO 空值检查
	return nil
}

func (p Permission) Handle(e event.Event)  {
	fmt.Println("收到pv：",e)
	fmt.Println(p.ScName,p.chmod,p.chown)
}