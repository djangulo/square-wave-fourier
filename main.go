package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/zserge/webview"
)

const (
	windowWidth  = 800
	windowHeight = 600
	changeCount  = "changeCount:"
	changeFreq   = "changeFreq:"
	changePhase  = "changePhase:"
)

var indexHTML string

func init() {

	indexHTML = fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
		</head>
		<body>
			<label for="harmonics">
				Harmonics:
				<input id="harmonics" value="10" min="1" max="100" type="number"/>
			</label>
			<label for="frequency">
				Frequency:
				<input id="frequency" value="5" min="0" max="100" type="number"/>
			</label>
			<label for="phase">
				Phase:
				<input id="phase" value="0" min="0" max="100" type="number"/>
			</label>
			<button onclick="external.invoke('render')">Render</button>
	
			<img id="harmonics-image" src="" />
			<img id="square-wave-image" src="" />
	
		<script type="text/javascript">
			function RefreshImages(){
				var harmonics = document.getElementById("harmonics-image");
				harmonics.setAttribute("src", "data:image/png;base64," + window.HarmonicsImageBase64);
				var square = document.getElementById("square-wave-image");
				square.setAttribute("src", "data:image/png;base64," + window.SquareWaveImageBase64);
			}
	
			function ChangeHarmonics() {
				var count = document.getElementById("harmonics");
				external.invoke('%s'+count.value);
				RefreshImages();
			}
	
			function ChangeBaseFreq() {
				var freq = document.getElementById("frequency");
				external.invoke('%s'+freq.value);
				RefreshImages();
			}
	
			function ChangePhase() {
				var phase = document.getElementById("phase");
				external.invoke('%s'+phase.value);
				RefreshImages();
			}
	
			document.getElementById("harmonics").onchange = ChangeHarmonics
			document.getElementById("frequency").onchange = ChangeBaseFreq;
			document.getElementById("phase").onchange = ChangePhase;
		</script>
		</body>
	</html>
	`, changeCount, changeFreq, changePhase)
}

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func handleRPC(w webview.WebView, data string) {
	switch {
	// case data == "close":
	// 	w.Terminate()
	// case data == "save":
	// 	log.Println("save", w.Dialog(webview.DialogTypeSave, 0, "Save file", ""))
	case strings.HasPrefix(data, changeFreq):
		n, err := strconv.Atoi(strings.TrimPrefix(data, changeFreq))
		if err != nil {
			log.Println(err)
			return
		}
		if n > 50 {
			baseFreq = 50
			setInput(w, "frequency", 50)
		} else if n < 0 {
			baseFreq = 0
			setInput(w, "frequency", 0)
		} else {
			baseFreq = n
		}
		buildHarmonics()
		sumHarmonics()

		refreshValues(w)
	case strings.HasPrefix(data, changeCount):
		n, err := strconv.Atoi(strings.TrimPrefix(data, changeCount))
		if err != nil {
			log.Println(err)
			return
		}
		if n > 50 {
			harmonicCount = 50
			setInput(w, "harmonics", 50)
		} else if n < 1 {
			harmonicCount = 1
			setInput(w, "harmonics", 1)
		} else {
			harmonicCount = n
		}
		buildHarmonics()
		sumHarmonics()
		refreshValues(w)
	case strings.HasPrefix(data, changePhase):
		n, err := strconv.Atoi(strings.TrimPrefix(data, changePhase))
		if err != nil {
			log.Println(err)
			return
		}
		if n > 50 {
			basePhase = 50
			setInput(w, "phase", 50)
		} else if n < 0 {
			basePhase = 0
			setInput(w, "phase", 0)
		} else {
			basePhase = n
		}
		buildHarmonics()
		sumHarmonics()
		refreshValues(w)
	case data == "render":
		refreshValues(w)
	}
}

func refreshValues(w webview.WebView) {
	w.Eval(
		fmt.Sprintf(`(function(harmonics, squarewave){
			window.HarmonicsImageBase64 = harmonics;
			window.SquareWaveImageBase64 = squarewave;
			RefreshImages();
		})("%s", "%s")`, getHarmonics(), getSquareWave()),
	)
}

func setInput(w webview.WebView, which string, value int) {
	w.Eval(
		fmt.Sprintf(`(function(id, value){
			var el = document.getElementById(id);
			el.value = value;
		})("%s", "%d")`, which, value),
	)
}

func main() {
	url := startServer()
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "Square wave renderer",
		Resizable:              false,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(200, 200, 200, 255)
	defer w.Exit()
	w.Run()
}
