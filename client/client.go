package main

import (
	"bufio"
	"fmt"
	"github.com/rsms/gotalk"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type User struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Username string `json:"username"`
}

type Messages struct {
	Texts []string `json:"texts"`
}

type Token struct {
	Tkn string `json:"token"`
}

type Msg struct {
	Text     string `json:"texts"`
	Token    string `json:"token"`
	Receiver string `json:"rcv"`
}

type MsgSender struct {
	Text string `json:"texts"`
	Sender string `json:"sender"`
}

type OnLines struct {
	Users []User `json:"users"`
}

var addr string = "127.0.0.1:1234"
var isLoaded int = 0
var isReceiverSelected int = 0

func login(uname string, name string, surname string) (Token, error){
	s, err := gotalk.Connect("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	token := Token{}
	usr := User{Name: name, Surname: surname, Username: uname}

	if err := s.Request("login", usr, &token); err != nil {
		log.Fatalln(err)
	}
	s.Close()
	return token, nil
}


func getMessages(tkn string) []string{
	s, err := gotalk.Connect("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	msgs := Messages{}
	token := Token{Tkn: tkn}
	if err := s.Request("getMsgs", token, &msgs); err != nil {
		log.Fatalln(err)
	}
	s.Close()
	return msgs.Texts
}

func sendMessage(t Token, receiver string, msg string) {
	s, err := gotalk.Connect("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	m := Msg{Text: msg, Token: t.Tkn, Receiver: receiver}
	s.Request("sendMsg", m, nil)
	s.Close()
}

func logout(t Token){
	s, err := gotalk.Connect("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	s.Request("logout", t, nil)
	s.Close()
}

func getOnlines(uname string) {
	s, err := gotalk.Connect("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	o := OnLines{}

	s.Request("online", nil, &o)
	s.Close()

	fmt.Printf("\n=====\nOnline Kullanıcılar:\n=====\n")
	for _, x := range o.Users {
		if x.Username != uname {
			fmt.Printf("%s: %s %s\n", x.Username, x.Name, x.Surname)
		}
	}
	fmt.Printf("Toplam: %d\n", len(o.Users)-1)
}

func getLastMessages(token Token) []string {
	s, err := gotalk.Connect("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	o := make([]MsgSender, 0)

	s.Request("listen", token, &o)
	s.Close()

	result := make([]string, 0)

	for _, x := range o {
		result = append(result, x.Sender+": "+x.Text)
	}

	isLoaded = 1
	return result
}

func listen(token Token) {
	for {
		if isLoaded == 1 {
			s, err := gotalk.Connect("tcp", addr)
			if err != nil {
				log.Fatalln(err)
			}

			o := make([]MsgSender, 0)

			s.Request("listen", token, &o)
			s.Close()

			if len(o) > 0 {
				for _, x := range o {
					fmt.Printf("\n%s: %s", x.Sender, x.Text)
				}
				if isReceiverSelected == 0 {
					fmt.Printf("Alıcının kullanıcı adı: ")
				} else{
					fmt.Printf("Mesaj: ")
				}
			}
			time.Sleep(time.Second)
		}
	}
}



func main() {
	var name, surname, uname string
	fmt.Printf("Ad: ")
	fmt.Scanf("%s", &name)

	fmt.Printf("Soyad: ")
	fmt.Scanf("%s", &surname)

	fmt.Printf("Kullanıcı Adı: ")
	fmt.Scanf("%s", &uname)

	token, _ := login(uname, name, surname)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigchan
		fmt.Printf("%s", s)
		logout(token)
		os.Exit(0)
	}()
	go listen(token)
	getOnlines(uname)
	var receiver, msg string

	fmt.Printf("\nMesajlarınız:\n")
	newMsgs := getLastMessages(token)
	if len(newMsgs) > 0 {
		for i, x := range newMsgs {
			fmt.Printf("%d.) %s", i+1, x)
		}
	} else {
		fmt.Printf("Yeni Mesaj Yok.\n")
	}

	fmt.Printf("Alıcının kullanıcı adı: ")
	fmt.Scanf("%s", &receiver)

	in := bufio.NewReader(os.Stdin)

	isReceiverSelected = 1
	fmt.Printf("\\q komutu ile çıkabilirsiniz.\n")

	for {
		fmt.Printf("Mesaj: ")
		msg, _ = in.ReadString('\n')
		if msg == "\\q\n" {
			break
		}
		sendMessage(token, receiver, msg)
	}
}
