package main
import (
	"log"
	"strconv"
	// "os"
	// "time"


	"net/http"
	// "html/template"

	// "github.com/go-gorp/gorp"

    // "github.com/astaxie/session"
    // _ "github.com/astaxie/session/providers/memory"

)
func HandlerCount(w http.ResponseWriter, r *http.Request) {
	// クライアントからきたリクエストに埋め込まれているcookieの確認
	log.Println(" Point 1")
	for _, c := range r.Cookies() {
		log.Print("Name:[", c.Name, "] Value:", c.Value)
	}
	log.Println(" Point 2")

	sess := globalSessions.SessionStart(w, r)  // ①
	log.Println(sess.SessionID())
	log.Println(" Point 3")
	for _, c := range r.Cookies() {
		log.Print("Name:[", c.Name, "] Value:", c.Value)
	}
	log.Println(" Point 4")
	ct := sess.Get("countnum") // ②
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	// t, _ := template.ParseFiles("public/count.html")
	w.Header().Set("Content-Type", "text/html")
    // ③
	// t.Execute(w, sess.Get("countnum"))
	w.Write([]byte("Count: " + strconv.Itoa((sess.Get("countnum").(int)))))
}