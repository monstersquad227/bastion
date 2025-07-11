package main

import (
	"bastion/repository"
	"bastion/service"
	"bastion/utils"
	gliderssh "github.com/gliderlabs/ssh"
	"log"
)

func main() {
	repository.InitMysql()
	service.InitOpenLdap()
	bastionSvc := &service.BastionService{
		BastionRepository: &repository.BastionRepository{},
	}
	userSvc := &service.UserService{
		UserRepository: &repository.UserRepository{},
	}

	server := &gliderssh.Server{
		Addr:    ":42678",
		Handler: utils.GliderSSHHandler,
		PasswordHandler: func(ctx gliderssh.Context, pass string) bool {
			return utils.ValidateUser(ctx.User(), pass, bastionSvc, userSvc)
		},
	}
	log.Printf("starting gliderssh server at %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
