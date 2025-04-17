package staticfiles

type Endpoint struct {
	cdn string
}

var endpoint = Endpoint{}

func Init(cdnEndpoint string) {
	endpoint.cdn = cdnEndpoint
}
