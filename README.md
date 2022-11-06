# Reolink-Manager

[Actual Repo](https://git.windmaker.net/a-castellano/Reolink-Manager)

 [![pipeline status](https://git.windmaker.net/a-castellano/Reolink-Manager/badges/master/pipeline.svg)](https://git.windmaker.net/a-castellano/Reolink-Manager/-/commits/master) [![coverage report](https://git.windmaker.net/a-castellano/Reolink-Manager/badges/master/coverage.svg)](https://git.windmaker.net/a-castellano/Reolink-Manager/-/commits/master) [![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=reolink-manager&metric=alert_status)](https://sonarqube.windmaker.net/dashboard?id=reolink-manager)

Library that manages Reolink Web Cam devices

## Features

* Reboot device
* Read motion detection sensor

## Examples

### Reboot webcam
```golang
    package main
    
    import (
    	"net/http"
    	"time"
    
    	reolink "github.com/a-castellano/reolink-manager/webcam"
    )
    
    func main() {
    
    	client := http.Client{
    		Timeout: time.Second * 5, // Maximum of 5 secs
    	}
    
    	webcam := reolink.Webcam{IP: "webcamIP", User: "admin", Password: "password"}
    	connectErr := webcam.Connect(client)
    	if connectErr != nil {
    		panic(connectErr)
    	} else {
    		rebootErr := webcam.Reboot(client)
    		if rebootErr != nil {
    			panic(rebootErr)
    		}
    	}
    }
```
### Read motion sensor status
```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	reolink "github.com/a-castellano/reolink-manager/webcam"
)

func main() {

	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	webcam := reolink.Webcam{IP: "webcamIP", User: "admin", Password: "password"}
	for {
		motion, connectErr := webcam.MotionDetected(client)
		if connectErr == nil {
			fmt.Println(motion)
		} else {
			fmt.Println(connectErr)
		}
		time.Sleep(time.Second)
	}
}
```
