package main

import (
	"os"
	"log"
	"time"
	// "fmt"
	// "net/mail"
	// "bytes"

	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"

	"github.com/jrootham/go-imap/imap"
)

func main() {

	app, window, widget := container()

	setup(widget)

	email(widget)

	signon(widget)

	register(widget)

	window.Show()

	app.Exec()
}


func container() (*widgets.QApplication, *widgets.QMainWindow, *widgets.QWidget) {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(250, 200)
	window.SetWindowTitle("No Password Signon")

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	return app, window, widget
}

func setup(widget *widgets.QWidget) {
	text(widget, "IMAP URL")
	url := line(widget)

	text(widget, "IMAP Port")
	port := line(widget)

	text(widget, "Email Name")
	name := line(widget)

	button := widgets.NewQPushButton2("Email Setup", nil)
	button.ConnectClicked(saveSetup(url, port, name))
	widget.Layout().AddWidget(button)

	return
}

func saveSetup(url, port, name *widgets.QLineEdit) func(bool) {
	return func(bool) {}
}

func email(widget *widgets.QWidget) {
	text(widget, "Email Password")

	input := widgets.NewQLineEdit(nil)
	widget.Layout().AddWidget(input)

	button := widgets.NewQPushButton2("Read Email", nil)
	button.ConnectClicked(makeGetEmail(""))
	widget.Layout().AddWidget(button)

	return
}

func makeGetEmail (password string) func(bool) {

	return func(bool) {
		client, dialErr := imap.DialTLS("imap.gmail.com:993", nil)
		if dialErr != nil {
			log.Fatal(dialErr)
		}

		// Don't forget to logout
		defer client.Logout(30 * time.Second)

		// Login
		_, loginErr := client.Login("jrootham@gmail.com", password)
		if loginErr != nil {
			log.Fatal(loginErr)
		}
		
		_, selectErr := client.Select("INBOX", false)
		if selectErr != nil {
			log.Fatal(selectErr)
		}


//		command, searchErr := client.Search("SUBJECT", "[#!nopassword-register=123456]")
		command, searchErr := imap.Wait(client.Search("SINCE", "1-Jan-2020"))
		if searchErr != nil {
			log.Fatal(searchErr)
		}

		log.Println(command)

		client.Recv(10000)

		results, resultErr := command.Result(0)
		if resultErr != nil {
			log.Fatal(resultErr)
		}
		log.Println(results)
		log.Println(results.String())
		log.Println(results.SearchResults())

		// // Fetch the headers of the 10 most recent messages
		// set, _ := imap.NewSeqSet("")
		// if client.Mailbox.Messages >= 10 {
		//     set.AddRange(client.Mailbox.Messages-9, client.Mailbox.Messages)
		// } else {
		//     set.Add("1:*")
		// }
		// cmd, _ := client.Fetch(set, "RFC822.HEADER")

		// // Process responses while the command is running
		// fmt.Println("\nMost recent messages:")
		// for cmd.InProgress() {
		//     // Wait for the next response (no timeout)
		//     client.Recv(-1)

		//     // Process command data
		//     for _, rsp := range cmd.Data {
		//         header := imap.AsBytes(rsp.MessageInfo().Attrs["RFC822.HEADER"])
		//         if msg, _ := mail.ReadMessage(bytes.NewReader(header)); msg != nil {
		//             fmt.Println("|--", msg.Header.Get("Subject"))
		//         }
		//     }
		//     cmd.Data = nil

		//     // Process unilateral server data
		//     for _, rsp := range client.Data {
		//         fmt.Println("Server data:", rsp)
		//     }
		//     client.Data = nil
		// }

		// // Check command completion status
		// if rsp, err := cmd.Result(imap.OK); err != nil {
		//     if err == imap.ErrAborted {
		//         fmt.Println("Fetch command aborted")
		//     } else {
		//         fmt.Println("Fetch error:", rsp.Info)
		//     }
		// }
	}
}

func signon(widget *widgets.QWidget) {
	text(widget, "Select Site")

	pick := widgets.NewQComboBox(nil)
	pick.AddItem("First", core.NewQVariant())
	pick.AddItem("Second", core.NewQVariant())
	pick.AddItem("Third", core.NewQVariant())
	widget.Layout().AddWidget(pick)

	button := widgets.NewQPushButton2("Signon", nil)
	widget.Layout().AddWidget(button)
}

func register(widget *widgets.QWidget) {
	button := widgets.NewQPushButton2("Register", nil)
	widget.Layout().AddWidget(button)
}

func connect (input *widgets.QLineEdit) func(bool) {
	box := widgets.QMessageBox__Ok
	return func (bool) {widgets.QMessageBox_Information(nil, "OK", input.Text(), box, box)}
}

func line(widget *widgets.QWidget) *widgets.QLineEdit {
	input := widgets.NewQLineEdit(nil)
	widget.Layout().AddWidget(input)
	return input
}

func text(widget *widgets.QWidget, contents string) {
	label := widgets.NewQLabel2(contents, nil, 0)
	widget.Layout().AddWidget(label)
}
