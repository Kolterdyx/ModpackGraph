package app

type Config struct {
	Info struct {
		Version   string `json:"productVersion"`
		Copyright string `json:"copyright"`
		Comments  string `json:"comments"`
	} `json:"info"`
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
}
