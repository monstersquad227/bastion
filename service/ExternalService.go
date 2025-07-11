package service

import (
	"github.com/go-ldap/ldap/v3"
	"log"
)

var (
	LdapClient *ldap.Conn
)

func InitOpenLdap() {
	l, err := ldap.Dial("tcp", "192.168.1.71:389")
	if err != nil {
		log.Fatalf("初始化 OpenLdap 客户端失败: %v", err)
	}
	LdapClient = l
	log.Println("OpenLdap 客户端初始化成功")
}
