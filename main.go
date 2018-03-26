package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr string
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var sourceConn *websocket.Conn
var sinkConn *websocket.Conn

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "link address")
}

type Pipe struct {
	sourceConn *websocket.Conn
	sinkConn   *websocket.Conn
}

var pipe = new(Pipe)

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/source", source)
	http.HandleFunc("/sink", sink)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Printf("serving static page.")
	homeTemplate.Execute(w, "ws://"+r.Host+"/sink")
}

func source(w http.ResponseWriter, r *http.Request) {
	s, err := upgrader.Upgrade(w, r, nil)
	pipe.sourceConn = s
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer s.Close()
	select {}

}

func sink(w http.ResponseWriter, r *http.Request) {
	s, err := upgrader.Upgrade(w, r, nil)
	pipe.sinkConn = s
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer sinkConn.Close()
	go func(p *Pipe) {
		for {
			mt, message, err := p.sourceConn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			err = p.sinkConn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}

		}
	}(pipe)
	defer s.Close()
	select {}
}

var homeTemplate = template.Must(template.New("").Parse(`
	<!DOCTYPE html>
  <head>
    <script type="text/javascript">
      var ws
      function fetchImage() {
       
        ws = new WebSocket("{{.}}");  
        ws.onopen = function(evt) {
            console.log("OPEN");
        }
        ws.onclose = function(evt) {
          console.log("CLOSE");
          ws = null;
        }
        ws.onmessage = function(evt) {
          drawImage(evt.data);
        }
        ws.onerror = function(evt) {
          console.log("ERROR: " + evt.data);
        }

      }

      function drawImage(data){
        document.images["screen"].src = URL.createObjectURL(data)
      }

      document.addEventListener("DOMContentLoaded", fetchImage)
    </script>
  </head>
  <body>
    <img width="100%" id="screen">
  </body>
</html>

	`))

/*
img, err := png.Decode(bytes.NewReader(message))
			fileName := fmt.Sprintf("%s.png", time.Now().String())
			log.Printf("Writing %d bytes to %s", len(message), fileName)
			file, err := os.Create(fileName)
			if err != nil {
				log.Fatalf("%v", err)
			}
			png.Encode(file, img)
			file.Close()
*/
