package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"net/http"
	"sync"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	scale float64
	x, y  float64
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	var wg sync.WaitGroup
	//ch := make(chan int, 10)

	g := &Game{scale: 0}

	wg.Add(2)
	go StartServer("16100", &wg, g)
	StartEbiten(&wg, g)

	wg.Wait()
}

func StartEbiten(wg *sync.WaitGroup, g *Game) {
	defer wg.Done()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Server")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// StartServer Запуск сервера
func StartServer(port string, wg *sync.WaitGroup, g *Game) {
	defer wg.Done()
	fmt.Println("Start server on port: ", port)

	http.HandleFunc("/health-check", HealthCheck)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocketConnection(w, r, g)
	})
	http.ListenAndServe(":"+port, nil)
}

// HandleWebSocketConnection Обработка входящих соединений и обновление их до вебсокета
func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request, g *Game) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("web socket upgrade error")
	}

	go HandleConnection(conn, g)
}

type RequestMessage struct {
	Method string `json:"method"`
	Point  Point  `json:"point"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// HandleConnection Обработка сообщений от клиента
func HandleConnection(conn *websocket.Conn, g *Game) {
	for {
		_, messageBytes, err := conn.ReadMessage()
		var request RequestMessage
		err = json.Unmarshal(messageBytes, &request)

		g.x = request.Point.X
		g.y = request.Point.Y

		fmt.Println("x: ", g.x)
		fmt.Println("y: ", g.y)

		if err != nil {
			log.Println(err)
		}

		//err = conn.WriteMessage(websocket.TextMessage, messageBytes)
		//if err != nil {
		//	log.Println(err)
		//}
	}
}

// HealthCheck Проверка состояния сервера
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("url: ", r.URL)

	responseData := map[string]bool{
		"success": true,
	}

	response, err := json.Marshal(responseData)

	if err != nil {
		log.Println("json marshal response error")
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		log.Println("response write error")
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.scale += 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.scale -= 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {

	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.x -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.x += 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	p1 := NewEdge(color.RGBA64{uint16(65535), uint16(0), uint16(0), 1}, 5, g.scale, screen)
	p1.DrawImage(10+g.x, 40+g.y)

	p2 := NewEdge(color.RGBA64{uint16(65535), uint16(65535), uint16(0), 0}, 5, g.scale, screen)
	p2.DrawImage(30+g.x, 40+g.y)

	p3 := NewEdge(color.RGBA64{uint16(65535), uint16(0), uint16(0), 1}, 5, g.scale, screen)
	p3.DrawImage(170+g.x, 170+g.y)

	vector.StrokeLine(screen, p1.x, p1.y, p2.x, p2.y, float32(1+g.scale), color.White, true)
	vector.StrokeLine(screen, p2.x, p2.y, p3.x, p3.y, float32(1+g.scale), color.White, true)
	vector.StrokeLine(screen, p3.x, p3.y, p1.x, p1.y, float32(1+g.scale), color.White, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) SetXYCords(x, y float64) {
	g.x = x
	g.y = y
}
