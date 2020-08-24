/*
@Time : 2020/8/24 17:20
@Author : ZhouHui2
@File : consul_register.go
@Software: GoLand
*/
package register

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hxzhouh/my-kit/kit/discovery"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type ConsulRegister struct {
	Target string
	Ttl    int
}

func NewConsulRegister(target string, ttl int) *ConsulRegister {
	return &ConsulRegister{Target: target, Ttl: ttl}
}

func (cr *ConsulRegister) Register(info discovery.RegisterInfo) error {

	// initial consul client config
	config := api.DefaultConfig()
	config.Address = cr.Target
	client, err := api.NewClient(config)
	if err != nil {
		log.Println("create consul client error:", err.Error())
	}

	serviceId := generateServiceId(info.ServiceName, info.Host, info.Port)

	reg := &api.AgentServiceRegistration{
		ID:      serviceId,
		Name:    info.ServiceName,
		Tags:    []string{info.ServiceName},
		Port:    info.Port,
		Address: info.Host,
	}

	if err = client.Agent().ServiceRegister(reg); err != nil {
		panic(err)
	}

	// initial register service check
	check := api.AgentServiceCheck{TTL: fmt.Sprintf("%ds", cr.Ttl), Status: api.HealthPassing}
	err = client.Agent().CheckRegister(
		&api.AgentCheckRegistration{
			ID:                serviceId,
			Name:              info.ServiceName,
			ServiceID:         serviceId,
			AgentServiceCheck: check})
	if err != nil {
		return fmt.Errorf("LearnGrpc: initial register service check to consul error: %s", err.Error())
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		x := <-ch
		log.Println("LearnGrpc: receive signal: ", x)
		// un-register service
		cr.DeRegister(info)

		s, _ := strconv.Atoi(fmt.Sprintf("%d", x))
		os.Exit(s)
	}()

	go func() {
		ticker := time.NewTicker(info.UpdateInterval)
		for {
			<-ticker.C
			err = client.Agent().UpdateTTL(serviceId, "", check.Status)
			if err != nil {
				log.Println("LearnGrpc: update ttl of service error: ", err.Error())
			}
		}
	}()

	return nil
}

func (cr *ConsulRegister) DeRegister(info discovery.RegisterInfo) error {

	serviceId := generateServiceId(info.ServiceName, info.Host, info.Port)

	config := api.DefaultConfig()
	config.Address = cr.Target
	client, err := api.NewClient(config)
	if err != nil {
		log.Println("create consul client error:", err.Error())
	}

	err = client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		log.Println("LearnGrpc: deregister service error: ", err.Error())
	} else {
		log.Println("LearnGrpc: deregistered service from consul server.")
	}

	err = client.Agent().CheckDeregister(serviceId)
	if err != nil {
		log.Println("LearnGrpc: deregister check error: ", err.Error())
	}

	return nil
}

func generateServiceId(name, host string, port int) string {
	return fmt.Sprintf("%s-%s-%d", name, host, port)
}
