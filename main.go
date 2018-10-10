package main

import (
	"github.com/SomniaStellarum/StellarUtilities/slog"
	"github.com/ethereum/go-ethereum/common"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("html/*.html"))
}

func index(w http.ResponseWriter, r *http.Request) {
	slog.DebugPrint("Path requested: ", r.URL.String())
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Panic("Template error on index: ", err)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	slog.DebugPrint("Path requested: ", r.URL.String())
	err := tpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		log.Panic("Template error on register: ", err)
	}
}

func doRegister(w http.ResponseWriter, r *http.Request) {
	slog.DebugPrint("doRegister called")

	ethAddress := common.HexToAddress(r.PostFormValue("ethAddress"))
	email := r.PostFormValue("email")

	putEmailReqChan <- &putEmailRequest{ethAddress, email}
}

func getUserEmail(w http.ResponseWriter, r *http.Request) {
	slog.DebugPrint("getUserEmail called")

	ethAddress := common.HexToAddress(r.PostFormValue("ethAddress"))

	getEmailReqChan <- ethAddress
	result := <-getEmailResChan
	if result.err != nil {
		// Crude, but return empty string to represent no email found
		w.Write([]byte(""))
		return
	}
	email := result.email

	w.Write([]byte(email))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test func finished"))
}

type logRequester struct {
	h http.Handler
}

func (l logRequester) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.DebugPrint("Static Path requested: ", r.URL.String())
	l.h.ServeHTTP(w, r)
}

func logRequest(f http.Handler) http.Handler {
	return logRequester{h: f}
}

func main() {
	slog.SetDebug()

	initEthStuff()
	go dbRequestsHandler()
	go ethReadLoop(4200000)

	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/doRegister", doRegister)
	http.HandleFunc("/getUserEmail", getUserEmail)
	http.HandleFunc("/test", test)
	http.Handle("/css/", slog.DebugPrintURL(http.StripPrefix("/css", http.FileServer(http.Dir("./css")))))
	http.Handle("/js/", slog.DebugPrintURL(http.StripPrefix("/js", http.FileServer(http.Dir("./js")))))
	http.ListenAndServe(":8080", nil)
}
