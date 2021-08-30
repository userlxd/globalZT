package office

var office Office

type Office struct {
	UUID string
}

func init() {
	office = Office{
		UUID: "123",
	}
}

func Run() {

}
