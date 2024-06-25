import asyncio
import websockets
import json

async def send_move(websocket, player_id, session_id):
    while True:
        move = input("Enter your move (e.g., e2-e4): ")
        message = {
            "action": "move",
            "data": {
                "playerId": player_id,
                "sessionId": session_id,
                "move": move
            }
        }
        await websocket.send(json.dumps(message))
        print(f"Sent message: {message}")

        response = await websocket.recv()
        response_data = json.loads(response)
        print_board(response_data)
        
        response = await websocket.recv()
        response_data = json.loads(response)
        print_board(response_data)
        
def print_board(board):
    print("  +-----------------+")

    for i in range(7, -1, -1):
        print("  | ", end="")
        for j in range(8):
            box = board[j][i]
            if box != "":
                print(box + " ", end="")
            else:
                print(". ", end="")
        print("|")

    print("  +-----------------+")
    print()

async def main():
    uri = "ws://localhost:7201/ws"
    async with websockets.connect(uri) as websocket:
        # Send the matching action once
        message = {
            "action": "matching",
            "data": {
                "playerId": "11231",
            }
        }
        await websocket.send(json.dumps(message))
        print(f"Sent message: {message}")

        response = await websocket.recv()
        response_data = json.loads(response)
        print(f"Received response: {response_data}")

        session_id = response_data['sessionId']
        
        # Enter a loop for sending moves
        await send_move(websocket, "11231", session_id)

asyncio.run(main())
