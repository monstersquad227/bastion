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
	svc := &service.BastionService{
		BastionRepository: &repository.BastionRepository{},
	}

	server := &gliderssh.Server{
		Addr:    ":42678",
		Handler: utils.GliderSSHHandler,
		PasswordHandler: func(ctx gliderssh.Context, pass string) bool {
			return utils.ValidateUser(ctx.User(), pass, svc)
		},
	}
	log.Printf("starting gliderssh server at %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
