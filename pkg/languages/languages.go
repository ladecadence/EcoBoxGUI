package languages

var eu = map[string]string{
	"welcome.welcome": "Ongi etorri",
	"welcome.info":    "Hasteko, eskaneatu zure QR kodea",
}

var es = map[string]string{
	"welcome.welcome": "Bienvenido/a",
	"welcome.info":    "Escanea tu código QR para empezar",
}

var en = map[string]string{
	"welcome.welcome": "Welcome",
	"welcome.info":    "To start scan your QR to code",
}

func GetString(id string, lang string) string {
	var out map[string]string

	switch lang {
	case "es":
		out = es
	case "eu":
		out = eu
	case "en":
		out = en
	default:
		out = eu
	}
	s, ok := out[id]
	if !ok {
		return "---"
	}
	return s
}
