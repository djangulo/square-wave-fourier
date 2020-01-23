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
)

var indexHTML = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<body>
		<label for="harmonic-count">
			Harmonic Count:
			<input id="harmonic-count" value="10" min="1" max="100" type="number"/>
		</label>
		<label>
			Base frequency:
			<input id="base-frequency" value="5" min="0" max="100" type="number"/>
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
			var count = document.getElementById("harmonic-count");
			external.invoke('changeCount:'+count.value);
			RefreshImages();
		}

		function ChangeBaseFreq() {
			var freq = document.getElementById("base-frequency");
			external.invoke('changeFreq:'+freq.value);
			RefreshImages();
		}
		document.getElementById("harmonic-count").onchange = ChangeHarmonics
		document.getElementById("base-frequency").onchange = ChangeBaseFreq;
	</script>
	</body>
</html>
`

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
	case strings.HasPrefix(data, "changeFreq:"):
		n, err := strconv.Atoi(strings.TrimPrefix(data, "changeFreq:"))
		if err != nil {
			log.Println(err)
			return
		}
		baseFreq = n
		buildHarmonics()
		sumHarmonics()

		w.Eval(
			fmt.Sprintf(`(function(harmonics, squarewave){
				window.HarmonicsImageBase64 = harmonics;
				window.SquareWaveImageBase64 = squarewave;
				RefreshImages();
			})("%s", "%s")`, getHarmonics(), getSquareWave()),
		)
	case strings.HasPrefix(data, "changeCount:"):
		n, err := strconv.Atoi(strings.TrimPrefix(data, "changeCount:"))
		if err != nil {
			log.Println(err)
			return
		}
		harmonicCount = n
		buildHarmonics()
		sumHarmonics()
		w.Eval(
			fmt.Sprintf(`(function(harmonics, squarewave){
				window.HarmonicsImageBase64 = harmonics;
				window.SquareWaveImageBase64 = squarewave;
				RefreshImages();
			})("%s", "%s")`, getHarmonics(), getSquareWave()),
		)
	case data == "render":
		w.Eval(
			fmt.Sprintf(`(function(harmonics, squarewave){
				window.HarmonicsImageBase64 = harmonics;
				window.SquareWaveImageBase64 = squarewave;
				RefreshImages();
			})("%s", "%s")`, getHarmonics(), getSquareWave()),
		)
	}
}

func main() {
	url := startServer()
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "Square wave renderer",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}
