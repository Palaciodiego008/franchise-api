package main

type Franchise struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Location Location `json:"location"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Address string `json:"address"`
	ZipCode string `json:"zip_code"`
}

type Company struct {
	Owner      Owner       `json:"owner"`
	Info       Info        `json:"information"`
	Franchises []Franchise `json:"franchises"`
}

type Owner struct {
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type Info struct {
	Name      string   `json:"name"`
	TaxNumber string   `json:"tax_number"`
	Location  Location `json:"location"`
}

var franchises map[string]Company

func main() {

}
