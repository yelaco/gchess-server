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
	"time"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
	"github.com/yelaco/go-chess-server/internal/corenet"
	"github.com/yelaco/go-chess-server/internal/session"
	"github.com/yelaco/go-chess-server/pkg/config"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type matchResponse struct {
	Type        string              `json:"type"`
	SessionID   string              `json:"session_id"`
	GameState   session.GameState   `json:"game_state"`
	PlayerState session.PlayerState `json:"player_state"`
}

var currentUser *User
var playerID string
var loginForm *tview.Form
var registerForm *tview.Form

func main() {
	app := tview.NewApplication()

	loginForm = tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			username := loginForm.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := loginForm.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			user := User{Username: username, Password: password}
			token, err := login(user)
			if err != nil {
				log.Printf("Login failed: %v", err)
				return
			}
			currentUser = &user
			playerID = token
			app.SetRoot(postLoginMenu(app), true).Run()
		}).
		AddButton("Back", func() {
			app.SetRoot(mainMenu(app), true).Run()
		})

	registerForm = tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Register", func() {
			username := registerForm.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := registerForm.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			user := User{Username: username, Password: password}
			err := register(user)
			if err != nil {
				log.Printf("Register failed: %v", err)
				return
			}
			fmt.Println("Register successful!")
			app.SetRoot(mainMenu(app), true).Run()
		}).
		AddButton("Back", func() {
			app.SetRoot(mainMenu(app), true).Run()
		})

	app.SetRoot(mainMenu(app), true).Run()
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
				fmt.Println(state.Status)
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
			return
		}
	}
}

func mainMenu(app *tview.Application) *tview.Flex {
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

func postLoginMenu(app *tview.Application) *tview.Flex {
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
		AddItem("Join Match", "Join a new match", '1', func() {
			app.Stop()
			joinMatch()
		}).
		AddItem("Other Option", "Perform another action", '2', func() {
			fmt.Println("Other Option selected!")
		}).
		AddItem("Logout", "Logout from your account", '3', func() {
			currentUser = nil
			app.SetRoot(mainMenu(app), true).Run()
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerBox, 3, 1, false).
		AddItem(headerText, 1, 1, false).
		AddItem(menu, 0, 1, true)

	return flex
}

func register(user User) error {
	url := "http://localhost:7202/register"
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		return fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("register failed: %s", body)
	}

	fmt.Println("User registered successfully")
	return nil
}

func login(user User) (string, error) {
	// url := "http://localhost:7202/login"
	// userJSON, err := json.Marshal(user)
	// if err != nil {
	// 	return "", fmt.Errorf("error marshalling user: %w", err)
	// }

	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(userJSON))
	// if err != nil {
	// 	return "", fmt.Errorf("error making POST request: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	body, _ := io.ReadAll(resp.Body)
	// 	return "", fmt.Errorf("login failed: %s", body)
	// }

	// var result map[string]interface{}
	// err = json.NewDecoder(resp.Body).Decode(&result)
	// if err != nil {
	// 	return "", fmt.Errorf("error decoding response: %w", err)
	// }

	// // Assuming the token is part of the response
	// token, ok := result["token"].(string)
	// if !ok {
	// 	return "", fmt.Errorf("token not found in response")
	// }

	// fmt.Println("User logged in successfully")
	// return token, nil

	return "11234", nil
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
	cmd := exec.Command("clear") // for Windows use "cls"
	cmd.Stdout = os.Stdout
	cmd.Run()
}
