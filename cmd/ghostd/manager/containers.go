package manager

import (
	"log"
	"os/exec"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/maincontainer"
)

type Containers struct {
	List  map[uint32]*Container
	Id    uint32
	Count uint32
}

func NewContainers() *Containers {
	return &Containers{Id: 1, Count: 0, List: make(map[uint32]*Container)}
}

func (containers *Containers) GetContainer(id uint32) *Container {
	return containers.List[id]
}

func (containers *Containers) ReleaseContainer(id uint32) bool {
	if container, exist := containers.List[id]; exist {
		if container.Cmd != nil {
			container.Cmd.Process.Kill()
		}
		containers.Count--
		delete(containers.List, id)
	} else {
		return false
	}
	return true
}

func (containers *Containers) CheckExistAddr(host, port string) bool {
	for _, container := range containers.List {
		if container.Port == port && container.Ip == host {
			return true
		}
	}
	return false
}

func (containers *Containers) LoginContainer(password []byte, username, host, port string) *Container {
	if containers.CheckExistAddr(host, port) {
		log.Printf("already allocated port number = %s", port)
		return nil
	}

	id := containers.Id
	container := NewContainer(id, host, port)
	if container == nil {
		return nil
	}
	container.Username = username

	container.Client.ConnectServer()
	defer container.Client.CloseServer()
	if !container.Client.LoginContainer(password, username, host, port) {
		log.Println("Login Fail ", username)
		return nil
	}

	containers.List[id] = container
	containers.Id++
	containers.Count++
	return container
}

func (containers *Containers) ForkContainer(password []byte, username, host, port string) *Container {
	if containers.CheckExistAddr(host, port) {
		log.Printf("already allocated port number = %s", port)
		return nil
	}

	wg := sync.WaitGroup{}
	id := containers.Id
	containers.Id++
	containers.Count++

	container := NewContainer(id, host, port)
	if container == nil {
		return nil
	}
	container.Username = username
	containers.List[id] = container
	wg.Add(1)

	go func(cid uint32) {
		log.Printf("execute node [%s] addr = %s:%s\n", username, host, port)
		args := []string{
			"--port=" + port,
			"-u=" + username,
		}
		if host != "" {
			args = append(args, host)
		}
		cmd := exec.Command("ghostnet", args...)
		out, err := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		log.Println(cmd.Process.Pid)
		container.PID = cmd.Process.Pid
		container.Cmd = cmd

		outBuf := make([]byte, 128)
		wg.Done()
		for {
			_, err := out.Read(outBuf)
			log.Printf("[%d] %s", cid, string(outBuf))
			if err != nil {
				log.Fatal(err)
			}
		}
	}(id)

	wg.Wait()
	container.Client.ConnectServer()
	defer container.Client.CloseServer()
	if !container.Client.LoginContainer(password, username, host, port) {
		log.Println("Login Fail ", username)
	}

	return container
}

func (containers *Containers) CreateContainer(password []byte, username, host, port string) *Container {
	if containers.CheckExistAddr(host, port) {
		log.Printf("already allocated port number = %s", port)
		return nil
	}

	wg := sync.WaitGroup{}
	id := containers.Id
	containers.Id++
	containers.Count++

	container := NewContainer(id, host, port)
	if container == nil {
		return nil
	}
	container.Username = username
	containers.List[id] = container
	wg.Add(1)

	go func(cid uint32) {
		log.Printf("execute node [%s] addr = %s:%s\n", username, host, port)
		cfg := gconfig.NewDefaultConfig()
		cfg.Username = username
		cfg.Password = password
		cfg.Ip = host
		cfg.Port = port
		main := maincontainer.NewMainContainer(cfg)
		main.StartContainer()
		wg.Done()
	}(id)

	wg.Wait()
	container.Client.ConnectServer()
	defer container.Client.CloseServer()
	if !container.Client.LoginContainer(password, username, host, port) {
		log.Println("Login Fail ", username)
	}

	return container
}
