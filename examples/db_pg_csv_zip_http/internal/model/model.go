package model

type HealthData struct {
	Year                int
	Gender              string
	Age                 int
	Location            string
	RaceAfricanAmerican bool
	RaceAsian           bool
	RaceCaucasian       bool
	RaceHispanic        bool
	RaceOther           bool
	Hypertension        bool
	HeartDisease        bool
	SmokingHistory      string
	BMI                 float64
	Hba1cLevel          float64
	BloodGlucoseLevel   float64
	Diabetes            bool
}

type JsonBody struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Base64      string `json:"base64"`
}
