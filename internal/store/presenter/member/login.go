package member

type (
    LoginRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    LoginResponse struct {
        IdUser int    `json:"id_user"`
        Token  string `json:"token"`
    }
)
