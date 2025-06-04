package mq

import "log"
var done chan bool

func StartConsumer(qName, cName string, callback func(msg []byte)bool){
	msgs,err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil{
		log.Fatal(err)
		return
	}

	done = make(chan bool)
	go func(){
		for d := range msgs{
			processErr := callback(d.Body)
			if !processErr {
				// TODO: 将任务写入错误队列中
				log.Panicln("callback err!")
			}
		}
	}()

	<-done

	channel.Close()
}