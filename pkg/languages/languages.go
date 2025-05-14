package languages

var eus = map[string]string{
	"welcome.welcome": "Ongi etorri",
	"welcome.info":    "Hasteko, eskaneatu zure QR kodea",
	"hello.hello":     "Kaixo",
	"hello.open":      "Ireki atea eta hartu behar duzuna",
	"nouser.error":    "Ezin dugu aurkitu sartutako erabiltzailea",
	"nouser.msg":      "Mesedez, saiatu berriro",
	"nouser.button":   "Ados",
	"open.open":       "Atea ireki",
	"open.msg":        "Ziurtatu atea itxita dagoela jarraitzeko",
}

var es = map[string]string{
	"welcome.welcome": "Bienvenido/a",
	"welcome.info":    "Escanea tu código QR para empezar",
	"hello.hello":     "Hola",
	"hello.open":      "Abre la puerta y coje lo que necesites",
	"nouser.error":    "No podemos encontrar el usuario introducido",
	"nouser.msg":      "Por favor inténtalo de nuevo",
	"nouser.button":   "De acuerdo",
	"open.open":       "Puerta abierta",
	"open.msg":        "Asegúrate de que la puerta quede cerrada para continuar",
}

var en = map[string]string{
	"welcome.welcome": "Welcome",
	"welcome.info":    "To start scan your QR to code",
	"hello.hello":     "Hello",
	"hello.open":      "Open the door and take what you need",
	"nouser.error":    "We can't find the user provided",
	"nouser.msg":      "Please try again",
	"nouser.button":   "Ok",
	"open.open":       "Door open",
	"open.msg":        "Be sure to close the door to continue",
}

func GetString(id string, lang string) string {
	var out map[string]string

	switch lang {
	case "es":
		out = es
	case "eus":
		out = eus
	case "en":
		out = en
	default:
		out = eus
	}
	s, ok := out[id]
	if !ok {
		return "---"
	}
	return s
}
