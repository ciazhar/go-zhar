# Consul
> Part ini akan menjelaskan bagaimana implementasi Consul sebeagai Centralized Configuration dan Service Discovery di Go.

## Actor
Dalam Part ini ada 3 actor yang bekerja yaitu:
- **Consul Server**: akan bekerja sebagai Centralized Configuration dan Service Discovery.
- **Go Server Service**: akan berperan sebagai penyedia data bagi Client Service.
- **Go Client Service**: akan mengambil data dari Server Service tiap 10 detik.

## Configuration
Pada contoh ini, terdapat 3 konfigurasi:
- **server.json**: konfigurasi ini akan digunakan sebagai konfigurasi default sebelum mengambil konfigurasi dari Consul. Jika sudah berhasil mengambil konfigurasi, maka konfigurasi ini akan ditimpa dengan konfigurasi dari Consul.
```json
{
  "application": {
    "name": "my-service",
    "host": "localhost",
    "port": 8081
  },
  "consul": {
    "host": "localhost",
    "port": 8500,
    "schema": "http",
    "key": "example/my-service",
    "configType": "json"
  }
}
```
- **server-consul.json**: konfigurasi ini merupakan konfigurasi sebenarnya untuk server service. Konfigurasi ini harus dicopy ke Key/Value Consul agar dapat diambil nanti.
```json
{
  "application": {
    "name": "my-service",
    "host": "localhost",
    "port": 8080
  }
}
```
- **client.json**: konfigurasi ini merupakan konfigurasi client service.
```json
{
  "consul": {
    "host": "localhost",
    "port": 8500,
    "schema": "http"
  }
}
```

## Use Case
- [Run Consul via Docker Compose](#run-consul-via-docker-compose)
- [Add Server Service Config to Consul](#add-server-service-config-to-consul)
- [Connect to Consul](#connect-to-consul)
- [Retrieve configuration from Consul](#retrieve-configuration-from-consul)
- [Register service to Consul Discovery](#register-service-to-consul-discovery)
- [Deregister service from Consul Discovery](#deregister-service-from-consul-discovery)
- [Retrieve service URL by service ID](#retrieve-service-url-by-service-id)

## Run Consul via Docker Compose
Pertama, kita perlu menjalankan Consul, kita bisa menggunakan Docker Compose yang telah tersedia.
```bash
cd ../../deployments/consul && docker-compose up
```
Jika Consul berhasil berjalan, Anda dapat melihat UI-nya di http://localhost:8500/ui.

## Add Server Service Config to Consul
1. Copy server-consul.json 
2. Buka Consul UI. 
3. Buat Key Value baru. 
4. Set key ke example/my-service, lalu Paste value.

## Run Server Service
Untuk menjalankan server service, gunakan perintah berikut.
```bash
go run cmd/server.go
````

## Connect to Consul
Untuk dapat terhubung ke Consul, kita akan menggunakan fungsi Init yang akan mengembalikan instance Consul. Fungsi Init memerlukan host Consul, port Consul, dan skema protokol Consul (http, tcp, dll), yang dapat diambil dari konfigurasi lokal.

```go
func Init(host string, port int, schema string) *Consul {
	config := consul.Config{
		Address: fmt.Sprintf("%s:%d", host, port),
		Scheme:  schema,
	}
	client, err := consul.NewClient(&config)
	if err != nil {
		log.Fatalf("Error initializing Consul client: %v\n", err)
		return nil
	}

	log.Println("Consul client initialized successfully")
	return &Consul{
		client: client,
	}
}
```
## Retrieve configuration from Consul
Setelah berhasil terhubung ke Consul, kita dapat mengambil konfigurasi service kita dari Consul. Fungsi ini akan menimpa konfigurasi dari `server.json` yang sudah berjalan.
```go
func (c Consul) RetrieveConfiguration(key string, configType string) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		log.Fatalf("Error retrieving key %s: %s", key, err)
	}

	if pair != nil {
		envData := pair.Value
		viper.SetConfigType(configType)
		if err := viper.ReadConfig(bytes.NewBuffer(envData)); err != nil {
			log.Fatalf("Error reading configuration for key %s: %s", key, err)
		}

		log.Printf("Configuration registered for key: %s", key)
	}
}
```
Jika konfigurasi Consul berhasil diambil, seharusnya service akan berjalan di port 8080; jika gagal, akan berjalan di port 8081.


## Register service to Consul Discovery
Untuk menambahkan service ke Consul Discovery, gunakan fungsi ini.
```go
func (c Consul) RegisterService(id, name, host string, port int) {

	reg := &consul.AgentServiceRegistration{
		ID:   id,
		Name: name,
		Port: port,
		Check: &consul.AgentServiceCheck{
			TCP:      fmt.Sprintf("%s:%d", host, port),
			Interval: "10s",
			Timeout:  "2s",
		},
	}

	log.Printf("Registering service with ID: %s, Name: %s", id, name)

	err := c.client.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
}
```
## Deregister service from Consul Discovery
Jika service akan di shutdown, service harus dideregister dari consul agar service di delist dari consul discovery service
```
func (c Consul) DeregisterService(id string) {
	err := c.client.Agent().ServiceDeregister(id)
	if err != nil {
		log.Printf("Error deregistering service: %v", err)
	}
}
```
## Running Client Service
Untuk menjalankan client service, gunakan perintah ini.
```bash
go run cmd/client.go
```
## Retrieve service URL by service ID
Jika client service ingin menggunakan API dari server service, client service tidak perlu mengetahui di mana URL server service, tetapi hanya perlu mengambilnya dari Consul Discovery Service menggunakan fungsi berikut.
```json
func (c Consul) RetrieveServiceUrl(id string) (string, error) {
	log.Printf("Retrieving service URL for id: %s", id)

	service, _, err := c.client.Agent().Service(id, nil)
	if err != nil {
		log.Printf("Error retrieving service: %v", err)
		return "", err
	}

	url := fmt.Sprintf("http://%s:%v", service.Address, service.Port)
	log.Printf("Service URL: %s", url)

	return url, nil
}
```

Untuk full code dari project ini ada di [sini](https://github.com/ciazhar/go-zhar/tree/master/examples/consul)