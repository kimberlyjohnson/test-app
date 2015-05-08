package main

import "github.com/cloudfoundry-incubator/garden"	
import "github.com/cloudfoundry-incubator/garden/client"	
import "github.com/cloudfoundry-incubator/garden/client/connection"	
import "fmt"
import "os"
import "time"

func main() {
	fmt.Printf("Hello, world.\n")
	connection := connection.New("tcp", "localhost:9241")
	the_client := client.New(connection)
	
	err := the_client.Ping()
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}
	
	capacity, _ := the_client.Capacity()
	fmt.Printf("\n\nCapacity: %#v", capacity)
	
	var containerSpec garden.ContainerSpec
	container, err := the_client.Create(containerSpec)
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}
	fmt.Printf("\n\nContainer: %#v", container)

	file, err := os.Open("webserver/webserver.tar")
	if err != nil {
		fmt.Printf("%#v", err)
		return 
	}
	
	err = container.StreamIn("/", file)
	if err != nil {
		fmt.Printf("%#v", err)
		return 
	}

	properties, _ := container.Properties()
	fmt.Printf("\n\nContainer Properties: %#v", properties)

	container.SetProperty("hi","bye")
	container.SetProperty("hello","good-bye")

	properties, _ = container.Properties()
	fmt.Printf("\n\nContainer Properties: %#v", properties)

	_, err = container.Property("does_not_exist")
	if err != nil {
		fmt.Printf("\n%#v", err)
		//return 
	}
	
	var processSpec garden.ProcessSpec
	var processIO garden.ProcessIO
	processSpec.Path = "webserver.exe"

	proc, err := container.Run(processSpec, processIO)
	if err != nil {
		fmt.Printf("%#v", err)
		return 
	}
	fmt.Printf("\n\nProcess ID: %d", proc.ID())

	a, b, err := container.NetIn(1283,1283) //outside world, on this computer
	if err != nil {
		fmt.Printf("%#v", err)
		return 
	}
	fmt.Printf("\nA:%d", a)
	fmt.Printf("\nB:%d", b)

	a, b, err = container.NetIn(1284,1285)
	if err != nil {
		fmt.Printf("%#v", err)
		return 
	}
	fmt.Printf("\nA:%d", a)
	fmt.Printf("\nB:%d", b)

	fmt.Printf("\nKilling process %d in...", proc.ID())
	for i := 25; i > 0; i-- {
		fmt.Printf(" %#d...", i)
		time.Sleep(time.Second)
	}

	proc.Signal(garden.SignalKill)

	proc, err = container.Run(processSpec, processIO)
	if err != nil {
		fmt.Printf("%#v", err)
		return 
	}
	fmt.Printf("\n\nProcess ID: %d", proc.ID())
	
	fmt.Printf("\nTerminating process %d in...", proc.ID())
	for i := 5; i > 0; i-- {
		fmt.Printf(" %#d...", i)
		time.Sleep(time.Second)
	}

	proc.Signal(garden.SignalTerminate)

}

