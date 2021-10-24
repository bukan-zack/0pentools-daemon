package cmd

import (
	"net/http"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"github.com/digitalocean/go-libvirt"
	"github.com/spf13/cobra"
	"github.com/0pentools/daemon/domain"
	"github.com/0pentools/daemon/router"
	"github.com/apex/log"
	log2 "log"
)

var rootCmd = &cobra.Command{
	Use: "daemon",
	Run: rootCmdRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log2.Fatalln(err)
	}
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	lv := libvirt.NewWithDialer(dialers.NewLocal())
	if err := lv.Connect(); err != nil {
		log.WithError(err).Fatal("failed to connect libvirt")
	}

	m := domain.NewManager(lv)

	srv := &http.Server{
		Addr: ":8000",
		Handler: router.Configure(m),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to configure http")
	}
}