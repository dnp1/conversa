package client

type ResponseBody struct {
    Success bool `json:"success"`
    Message string `json:"message"`
}

type EmptyResponse struct {
    ResponseBody
}

type RoomItem struct {
    Username string `json:"username"`
    Name     string `json:"name"`
}

type RoomData struct {
    ResponseBody
    Items []RoomItem `json:"data"`
}

type RoomBody struct {
    Name string `json:"name"`
}

type MessageBody struct {
    Content string `json:"content"`
}

type LoginBody struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type SignUp struct {
    Username             string `json:"username"`
    Password             string `json:"password"`
    PasswordConfirmation string `json:"passwordConfirmation"`
}