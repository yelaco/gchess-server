package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
	"github.com/yelaco/go-chess-server/internal/corenet"
	"github.com/yelaco/go-chess-server/internal/session"
	"github.com/yelaco/go-chess-server/pkg/config"
)

type matchResponse struct {
	Type        string              `json:"type"`
	SessionID   string              `json:"session_id"`
	GameState   session.GameState   `json:"game_state"`
	PlayerState session.PlayerState `json:"player_state"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	app          *tview.Application
	loginForm    *tview.Form
	registerForm *tview.Form
	currentUser  *User
	playerID     string
	gameResult   string
)

func main() {
	app = tview.NewApplication()

	setupForms(app)

	app.SetRoot(mainMenu(), true).Run()

}

func setupForms(app *tview.Application) {
	loginForm = tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			username := loginForm.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := loginForm.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			user := User{Username: username, Password: password}
			login(user)
		}).
		AddButton("Back", func() {
			app.SetRoot(mainMenu(), true).Run()
		})

	registerForm = tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Register", func() {
			username := registerForm.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := registerForm.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			user := User{Username: username, Password: password}
			register(user)
		}).
		AddButton("Back", func() {
			app.SetRoot(mainMenu(), true).Run()
		})
}

func mainMenu() *tview.Flex {
	headerBox := tview.NewBox().
		SetBorder(true).
		SetTitle("Go Chess Server").
		SetTitleAlign(tview.AlignLeft)

	headerText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	updateHeader := func() {
		if currentUser == nil {
			headerText.SetText("Please login to the server")
		} else {
			headerText.SetText(fmt.Sprintf("User: %s", currentUser.Username))
		}
	}

	updateHeader()

	menu := tview.NewList().
		AddItem("Login", "Login to your account", '1', func() {
			app.SetRoot(loginForm, true).Run()
		}).
		AddItem("Register", "Register a new account", '2', func() {
			app.SetRoot(registerForm, true).Run()
		}).
		AddItem("Quit", "Exit the application", '3', func() {
			app.Stop()
			os.Exit(0)
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerBox, 3, 1, false).
		AddItem(headerText, 1, 1, false).
		AddItem(menu, 0, 1, true)

	return flex
}

func postLoginMenu() *tview.Flex {
	headerBox := tview.NewBox().
		SetBorder(true).
		SetTitle("Go Chess Server").
		SetTitleAlign(tview.AlignLeft)

	headerText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	updateHeader := func() {
		if currentUser == nil {
			headerText.SetText("Please login to the server")
		} else {
			headerText.SetText(fmt.Sprintf("User: %s", currentUser.Username))
		}
	}

	updateHeader()

	menu := tview.NewList().
		AddItem("Join match", "Join a new match", '0', func() {
			app.Stop()
			app = tview.NewApplication()
			gameResult = "UNDEFINED"
			joinMatch()
			showLoginSuccessDialog("Game ended with " + gameResult)
		}).
		AddItem("View previous matches", "View previous game with moves", '2', func() {
			showLoginSuccessDialog("This feature is still in developing stage")
		}).
		AddItem("Logout", "Logout from your account", '3', func() {
			currentUser = nil
			app.SetRoot(mainMenu(), true).Run()
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerBox, 3, 1, false).
		AddItem(headerText, 1, 1, false).
		AddItem(menu, 0, 1, true)

	return flex
}

func showLoginErrorDialog(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(loginForm, true).Run()
		})
	app.SetRoot(modal, true).Run()
}

func showRegisterErrorDialog(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(registerForm, true).Run()
		})
	app.SetRoot(modal, true).Run()
}

func showLoginSuccessDialog(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(postLoginMenu(), true).Run()
		})
	app.SetRoot(modal, true).Run()
}

func register(user User) {
	url := "http://localhost:7202/api/users"
	userJSON, err := json.Marshal(user)
	if err != nil {
		showRegisterErrorDialog(fmt.Sprintf("Error marshalling user: %v", err))
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		showRegisterErrorDialog(fmt.Sprintf("Error making POST request: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		showRegisterErrorDialog(fmt.Sprintf("Register failed: %s", body))
		return
	}

	login(user)
}

func login(user User) {
	url := "http://localhost:7202/api/login"
	userJSON, err := json.Marshal(user)
	if err != nil {
		showLoginErrorDialog(fmt.Sprintf("Error marshalling user: %v", err))
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		showLoginErrorDialog(fmt.Sprintf("Error making POST request: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		showLoginErrorDialog(fmt.Sprintf("Login failed: %s", body))
		return
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		showLoginErrorDialog(fmt.Sprintf("Error decoding response: %v", err))
		return
	}

	// Assuming the token is part of the response
	pid, ok := result["player_id"].(string)
	if !ok {
		showLoginErrorDialog("Player ID not found in response")
		return
	}
	currentUser = &user
	playerID = pid

	showLoginSuccessDialog("Login successful!")
}

func joinMatch() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:" + config.Port, Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Println("Connected to game server")

	done := make(chan struct{})

	go func() {
		defer close(done)
		if err := c.WriteJSON(corenet.Message{
			Action: "matching",
			Data: map[string]interface{}{
				"player_id": playerID,
			},
		}); err != nil {
			log.Fatal("ws write", err)
		}

		clearScreen()
		log.Println("Attemp matchmaking...")

		var resp matchResponse
		if err := c.ReadJSON(&resp); err != nil {
			log.Fatal("ws match:", err)
			return
		}

		if resp.Type != "matched" {
			log.Fatal("not matched")
			return
		}

		sessionResp := session.SessionResponse{
			Type:      "session",
			GameState: resp.GameState,
		}
		var state session.GameState
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if sessionResp.Type == "session" {
				state = sessionResp.GameState
			}
			clearScreen()
			printBoard(state.Board, resp.PlayerState.IsWhiteSide)
			if state.Status != "ACTIVE" {
				gameResult = state.Status
				return
			}
			if resp.PlayerState.IsWhiteSide == state.IsWhiteTurn {
				if sessionResp.Type == "session" {
					fmt.Print("Enter your move (e.g., e2-e4): ")
				} else {
					fmt.Print("[Invalid] Enter new move (e.g., e2-e4):")
				}
				scanner.Scan()
				move := scanner.Text()

				c.WriteJSON(corenet.Message{
					Action: "move",
					Data: map[string]interface{}{
						"session_id": resp.SessionID,
						"player_id":  playerID,
						"move":       move,
					},
				})

				if err := c.ReadJSON(&sessionResp); err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Print("Wait for your opponent...")
				if err := c.ReadJSON(&sessionResp); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			os.Exit(1)
		}
	}
}

func printBoard(board [8][8]string, isWhiteSide bool) {
	fmt.Println("  +-----------------+")

	if isWhiteSide {
		for i := 7; i >= 0; i-- {
			fmt.Printf("%d | ", i+1)
			for j := 0; j < 8; j++ {
				box := board[j][i]
				if box != "" {
					fmt.Print(box + " ")
				} else {
					fmt.Print(". ")
				}
			}
			fmt.Println("|")
		}

		fmt.Println("  +-----------------+")
		fmt.Println("    a b c d e f g h")
	} else {
		for i := 0; i < 8; i++ {
			fmt.Printf("%d | ", i+1)
			for j := 7; j >= 0; j-- {
				box := board[j][i]
				if box != "" {
					fmt.Print(box + " ")
				} else {
					fmt.Print(". ")
				}
			}
			fmt.Println("|")
		}

		fmt.Println("  +-----------------+")
		fmt.Println("    h g f e d c b a ")
	}

	fmt.Println()
}

func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported OS")
		return
	}

	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error clearing screen: %v\n", err)
	}
}
