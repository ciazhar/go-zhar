package model

type HealthData struct {
	Year                int32
	Gender              string
	Age                 int32
	Location            string
	RaceAfricanAmerican bool
	RaceAsian           bool
	RaceCaucasian       bool
	RaceHispanic        bool
	RaceOther           bool
	Hypertension        bool
	HeartDisease        bool
	SmokingHistory      string
	BMI                 float32
	Hba1cLevel          float32
	BloodGlucoseLevel   float32
	Diabetes            bool
}

type JsonBody struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Base64      string `json:"base64"`
}
