package languages

var eus = map[string]string{
	"welcome.welcome":  "Ongi etorri",
	"welcome.info":     "Hasteko, eskaneatu zure QR kodea",
	"hello.hello":      "Kaixo",
	"hello.open":       "Ireki atea eta hartu behar duzuna",
	"nouser.error":     "Ezin dugu aurkitu sartutako erabiltzailea",
	"nouser.msg":       "Mesedez, saiatu berriro",
	"nouser.button":    "Ados",
	"open.open":        "Atea irekia",
	"open.msg":         "Ziurtatu atea itxita dagoela jarraitzeko",
	"closed.closed":    "Atea itxita",
	"closed.msg":       "Itxaron pixka bat aldaketak egiaztatzen ditugun bitartean, mesedez.",
	"confirm.confirm":  "Baieztatu zure eragiketa",
	"confirm.msg.none": "Ez duzu ontzirik erretiratu",
	"confirm.msg.sing": "Ontzia %d kendu duzu",
	"confirm.msg.plur": "Ontziak %d kendu dituzu",
	"confirm.button":   "Zuzena",
	"error.error":      "Sistemaren errorea.",
	"error.msg":        "Mesedez, saiatu berriro",
	"error.button":     "Ados",
}

var es = map[string]string{
	"welcome.welcome":  "Bienvenido/a",
	"welcome.info":     "Escanea tu código QR para empezar",
	"hello.hello":      "Hola",
	"hello.open":       "Abre la puerta y coje lo que necesites",
	"nouser.error":     "No podemos encontrar el usuario introducido",
	"nouser.msg":       "Por favor inténtalo de nuevo",
	"nouser.button":    "De acuerdo",
	"open.open":        "Puerta abierta",
	"open.msg":         "Asegúrate de que la puerta quede cerrada para continuar",
	"closed.closed":    "Puerta cerrada",
	"closed.msg":       "Espera un momento mientras verificamos los cambios, por favor.",
	"confirm.confirm":  "Confirma tu operación",
	"confirm.msg.none": "No has retirado nigún envase",
	"confirm.msg.sing": "Has retirado %d envase",
	"confirm.msg.plur": "Has retirado %d envases",
	"confirm.button":   "Correcto",
	"error.error":      "Error del sistema.",
	"error.msg":        "Por favor inténtalo de nuevo",
	"error.button":     "De acuerdo",
}

var en = map[string]string{
	"welcome.welcome":  "Welcome",
	"welcome.info":     "To start scan your QR to code",
	"hello.hello":      "Hello",
	"hello.open":       "Open the door and take what you need",
	"nouser.error":     "We can't find the user provided",
	"nouser.msg":       "Please try again",
	"nouser.button":    "Ok",
	"open.open":        "Door open",
	"open.msg":         "Be sure to close the door to continue",
	"closed.closed":    "Door closed",
	"closed.msg":       "Please wait a moment while we check the changes.",
	"confirm.confirm":  "Confirm your operation",
	"confirm.msg.none": "You haven't taken any containers",
	"confirm.msg.sing": "You have taken %d container",
	"confirm.msg.plur": "You have taken %d containers",
	"confirm.button":   "Ok",
	"error.error":      "System error.",
	"error.msg":        "Please try again",
	"error.button":     "Ok",
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
