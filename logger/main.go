package main
import(
	"fmt"
	"os"
	"time"
)

type LogMessage struct{
	Level string
	Message string
	Time time.Time
}

func main(){
	logChannel := make(chan LogMessage, 100)
	go logger(logChannel)

	logChannel <- LogMessage{"INFO","Server started",time.Now()}
	logChannel <- LogMessage{"ERROR","Database failed",time.Now()}
	time.Sleep(1*time.Second)
}

func logger(ch <-chan LogMessage){
	file,err := os.OpenFile("app.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
	
	if err!=nil{
		panic(err)
	}
	defer file.Close()

	for msg := range ch{
		line := fmt.Sprintf("[%s] %s: %s\n",msg.Time.Format(time.RFC3339),msg.Level,msg.Message)
		file.WriteString(line)
	}
}