package main

import (
	"com.enesuysal/go-chat/api"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/rsms/gotalk"
	"log"
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
	Text string `json:"texts"`
	Tkn string `json:"token"`
	Receiver string `json:"rcv"`
}

type MsgSender struct {
	Text string `json:"texts"`
	Sender string `json:"sender"`
}

type OnLines struct {
	Users []User `json:"users"`
}

func server() {
	db := api.OpenDb()
	defer db.Close()

	gotalk.Handle("sendMsg", func(in Msg) error{
		u, _ := api.QueryUserbyToken(context.Background(), in.Tkn, db)
		if u != nil {
			api.CreateMessage(context.Background(), db, u.Username, in.Receiver, in.Text)
		}

		return nil
	})

	gotalk.Handle("getMsgs", func(in Token) (Messages, error){
		u, _ := api.QueryUserbyToken(context.Background(), in.Tkn, db)

		msgs, _ := api.QueryMessagesUsers(context.Background(), u)

		if msgs == nil {
			msgs = []string{"Mesaj yok."}
		}

		for _, x := range msgs {
			log.Println(x)
		}

		return Messages{Texts: msgs}, nil
	})

	gotalk.Handle("login", func(in User) (Token, error) {
		user, _ := api.QueryUser(context.Background(), in.Username, db)
		if user == nil {
			log.Println("Kullanıcı bulunamadı, oluşturuluyor...")
			user, _ = api.CreateUser(context.Background(), in.Username, in.Name, in.Surname, db)
		}

		t := md5.Sum([]byte(in.Username))
		token := Token{Tkn: hex.EncodeToString(t[:])}
		user.Update().SetToken(token.Tkn).SetIsOnline(1).Save(context.Background())

		ol, _ := api.QueryOnlineUsers(context.Background(), db)

		for _, receiver := range ol {
			if user.Username != receiver.Username {
				api.CreateMessage(context.Background(), db, "Broadcast", receiver.Username, user.Username+" is online.\n")
			}
		}

		return token, nil
	})

	gotalk.Handle("logout", func(token Token) error{
		u, _ := api.QueryUserbyToken(context.Background(), token.Tkn, db)
		u.Update().SetIsOnline(0).Save(context.Background())
		ol, _ := api.QueryOnlineUsers(context.Background(), db)

		for _, receiver := range ol {
			api.CreateMessage(context.Background(), db, "Broadcast", receiver.Username, u.Username+" has left.\n")
		}

		return nil
	})

	gotalk.Handle("online", func() (OnLines, error) {
		o, _ := api.QueryOnlineUsers(context.Background(), db)

		users := make([]User, 0)

		for _, x := range(o) {
			users = append(users, User{
				Name:     x.Name,
				Surname:  x.Surname,
				Username: x.Username,
			})
		}

		return OnLines{Users: users}, nil
	})

	gotalk.Handle("listen", func(token Token)([]MsgSender, error) {
		u, _ := api.QueryUserbyToken(context.Background(), token.Tkn, db)
		msgs, _ := api.QueryLastMessages(context.Background(), u)
		result := make([]MsgSender, 0)
		var tmp MsgSender
		for _, msg := range(msgs) {
			tmp = MsgSender{
				Text:   msg.Message,
				Sender: msg.SenderUsername,
			}
			result = append(result, tmp)
		}

		return result, nil
	})

	fmt.Printf("%s\n", "Listening on 0.0.0.0:1234")
	if err := gotalk.Serve("tcp", "0.0.0.0:1234", nil); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	server()
}