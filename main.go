package main

func main(){
	// top level inits, this avoids race conditions
	Distance_chan = make(chan []byte, 10)

	go main_mqtt()
	go main_site()

	// wait forever
	for{
	}
}
