package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	typ "github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
	"golang.org/x/net/context"
	"io"
	"os"
	"strings"
)

//A Package wrapper for handling docker containers
type Dcontainer struct {
	ID          string
	Ctx         context.Context
	LogFd       *os.File
	HijackedRes types.HijackedResponse
	Cli         *client.Client
}

//Run will PUll an image if its not available and start a container and attach to it
//It will also set the max memory limit of the image.
func (d *Dcontainer) Run(name, image string, cmd []string, mem int64, logFileName string) error {

	var err error
	d.LogFd, err = os.Create(logFileName)
	if err != nil {
		logs.Printf("Unable to open the logfile")
		return err
	}

	//Create a Dummy contexts and a Client handler for this connection
	d.Ctx = context.Background()
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.18", nil, defaultHeaders)
	if err != nil {
		d.Close(false)
		return err
	}

	//Try to PUll the image and check the return response it ca be either of the two only
	resp, err := cli.ImagePull(d.Ctx, image, types.ImagePullOptions{All: false})
	if err != nil {
		logs.Printf("Error PUlling %v\n", err)
		d.Close(false)
		return err
	}
	var eoferr error
	var line, prev_line []byte
	rImage := bufio.NewReader(resp)
	for eoferr == nil {
		line, eoferr = rImage.ReadBytes('\n')
		if eoferr == nil {
			prev_line = line
		}
	}
	line_str := string(prev_line)
	if !strings.Contains(line_str, "Image is up to date for") && !strings.Contains(line_str, "Downloaded newer image for") {
		d.Close(false)
		return fmt.Errorf("%s pull failed\n", image)
	}

	//The steps are as follows
	// 1) Create a Container Entry with the desired attributes
	// 2) Attach to the container id and get the Reader handle to get the console output
	// 3) Actually start the container
	// 4) Now start reading the stream 'HijackedResponse' from the ContainerAttach Handle

	//CREATE
	mem = mem * 1024 * 1024 //Argument supplied in MB
	cconfig := container.Config{Image: image, Cmd: cmd}
	hconfig := container.HostConfig{NetworkMode: "host", Resources: container.Resources{Memory: mem}}
	r, err := cli.ContainerCreate(d.Ctx, &cconfig, &hconfig, nil, name)
	if err != nil {
		logs.Printf("Error creating a container %v\n", err)
		d.Close(false)
		return err
	}

	//ATTACH
	d.HijackedRes, err = cli.ContainerAttach(d.Ctx, r.ID, types.ContainerAttachOptions{Stdout: true, Stderr: true, Stream: true})
	if err != nil {
		logs.Printf("Unable to attach the container\n")
		d.Close(true)
		return err
	}

	//START
	err = cli.ContainerStart(d.Ctx, r.ID, types.ContainerStartOptions{})
	if err != nil {
		logs.Printf("Unable to start a docker container\n")
		d.Close(true)
		return err
	}

	d.ID = r.ID
	d.Cli = cli
	return nil
}

//Getstats get the docker container stats and stores into StatsInfo
func (d *Dcontainer) GetStats() (typ.StatsInfo, error) {

	var data typ.StatsInfo

	//start getting the docker container stats
	resp, err := d.Cli.ContainerStats(d.Ctx, d.ID, true)
	if err != nil {
		logs.Println("Container stats error", err)
		return typ.StatsInfo{}, err
	}

	defer resp.Body.Close()
	body := io.Reader(resp.Body)

	if err := json.NewDecoder(body).Decode(&data); err != nil {
		logs.Printf("Json Unmarshall error = %v", err)
		return typ.StatsInfo{}, err
	}
	return data, nil
}

func (d *Dcontainer) Wait() int {
	if d.ID == "" {
		return -1
	}
	//Start READING the stream untl EOF
	go func() {
		defer d.Close(true)
		fdw := bufio.NewWriter(d.LogFd)
		_, err := io.Copy(fdw, d.HijackedRes.Reader)
		if err != nil {
			logs.Printf("Error copying to STDOUT %v", err)
		}
	}()
	retVal, _ := d.Cli.ContainerWait(d.Ctx, d.ID)
	return retVal
}

func (d *Dcontainer) Close(HijackRes bool) {
	if HijackRes {
		d.HijackedRes.Close()
	}
	d.LogFd.Close()
}

//Kill will kill the docker container using container id
func (d *Dcontainer) Kill() error {
	if d.ID == "" {
		return fmt.Errorf("Invalid Container")
	}
	return d.Cli.ContainerKill(d.Ctx, d.ID, "KILL")
}
