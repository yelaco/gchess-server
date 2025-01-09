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
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
	"github.com/yelaco/go-chess-server/internal/database"
	"github.com/yelaco/go-chess-server/internal/game"
	"github.com/yelaco/go-chess-server/pkg/corenet"
	"github.com/yelaco/go-chess-server/pkg/session"
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
	prevSessions []database.Session
	playerID     string
	gameResult   string
	moveIdx      int
	boardStates  [][8][8]string
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

func viewPreviousMatches() tview.Primitive {
	list := tview.NewList()

	for _, session := range prevSessions {
		list.AddItem(fmt.Sprintf("Session ID: %s", session.SessionID), "", 0, func() {
			// When a session is selected, show a confirmation dialog
			modal := tview.NewModal().
				SetText(fmt.Sprintf("View match for session ID: %s?", session.SessionID)).
				AddButtons([]string{"Yes", "No"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						app.SetRoot(viewMatch(session.Moves), true).Run()
					} else {
						app.SetRoot(viewPreviousMatches(), true).Run()
					}
				})
			app.SetRoot(modal, true)
		})
	}

	list.AddItem("Back", "Go back to the main menu", 'b', func() {
		app.SetRoot(postLoginMenu(), true).Run()
	})

	return list
}

func viewMatch(moves []string) *tview.Flex {
	moveIdx = 0
	prevGame := game.InitGame([2]string{"-1", "-2"})
	boardStates = [][8][8]string{prevGame.GetBoard()}
	for i, move := range moves {
		pos := strings.Split(move, "-")
		if i%2 == 0 {
			if err := prevGame.MakeMove("-1", pos[0], pos[1]); err != nil {
				showViewMatchErrorDialog("Coulnd't parse move")
				os.Exit(1)
			}
		} else {
			if err := prevGame.MakeMove("-2", pos[0], pos[1]); err != nil {
				showViewMatchErrorDialog("Coulnd't parse move")
				os.Exit(1)
			}
		}
		boardStates = append(boardStates, prevGame.GetBoard())
	}
	if !prevGame.IsOver() {
		showViewMatchErrorDialog("Invalid game")
		os.Exit(1)
	}

	boardView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	updateBoardView := func() {
		board := boardStates[moveIdx]
		boardView.SetText(formatBoard(board))
	}

	updateBoardView()

	form := tview.NewForm().
		AddButton("Previous", func() {
			if moveIdx > 0 {
				moveIdx -= 1
				updateBoardView()
			}
		}).
		AddButton("Next", func() {
			if moveIdx < len(moves)-1 {
				moveIdx += 1
				updateBoardView()
			}
		}).
		AddButton("Exit", func() {
			app.SetRoot(viewPreviousMatches(), true).Run()
			boardStates = [][8][8]string{}
			moveIdx = 0
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(boardView, 0, 1, false).
		AddItem(form, 0, 1, true)

	return flex
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
			clearScreen()
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
		AddItem("Join match", "Join a new match", '1', func() {
			app.Stop()
			app = tview.NewApplication()
			gameResult = ""
			joinMatch()
			if gameResult == "timeout" {
				showMatchingErrorDialog("Matching timeout")
			} else if gameResult == "queueing" {
				showMatchingErrorDialog("You are queueing elsewhere")
			} else if gameResult == "error" {
				showMatchingErrorDialog("You are playing elsewhere")
			} else {
				showLoginSuccessDialog("Game ended with " + gameResult)
			}
		}).
		AddItem("View previous matches", "View previous game with moves", '2', func() {
			getPreviousSessions()
			if len(prevSessions) != 0 {
				app.SetRoot(viewPreviousMatches(), true).Run()
			} else {
				showLoginSuccessDialog("No match record found")
			}
		}).
		AddItem("Logout", "Logout from your account", '3', func() {
			currentUser = nil
			playerID = ""
			prevSessions = []database.Session{}
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

func showMatchingErrorDialog(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(postLoginMenu(), true).Run()
		})
	app.SetRoot(modal, true).Run()
}

func showViewMatchErrorDialog(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(viewPreviousMatches(), true).Run()
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

func getPreviousSessions() {
	if playerID == "" {
		return
	}
	url := "http://localhost:7202/api/sessions?player_id=" + playerID

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		showLoginErrorDialog("Error making the HTTP request:" + err.Error())
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		showLoginErrorDialog("Error reading the response body:" + err.Error())
		return
	}

	// Unmarshal the JSON into a slice of Session
	var sessions []database.Session
	err = json.Unmarshal(body, &sessions)
	if err != nil {
		showLoginErrorDialog("Error unmarshalling JSON:" + err.Error())
		return
	}

	prevSessions = sessions
}

func joinMatch() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:7201", Path: "/ws"}

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
			gameResult = resp.Type
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

func formatBoard(board [8][8]string) string {
	var sb strings.Builder

	sb.WriteString("  +-----------------+\n")
	for i := 7; i >= 0; i-- {
		sb.WriteString(fmt.Sprintf("%d | ", i+1))
		for j := 0; j < 8; j++ {
			box := board[j][i]
			if box != "" {
				sb.WriteString(box + " ")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteString("|\n")
	}
	sb.WriteString("  +-----------------+\n")
	sb.WriteString("    a b c d e f g h\n")

	return sb.String()
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
