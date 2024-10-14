# Go Style Guide

> Style Guide ini ada sebagai pedoman dalam ngoding Go berdasarkan aspek best practice, memory safety dan performance.

## Content

- [Number](#number)
- [Mutex](#mutex)
- [Slice & Map](#slice--map)
- [Test](#test)

## Integer

### Human Readable

Kita bisa menggunakan underscore untuk mempermudah membaca number
<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
number := 10000000
```

</td><td>

```go
better := 10_000_000
```

</td></tr>
<tr><td>

</td></tr>
</tbody></table>

## Mutex

### Deklarasi Mutex

Deklarasi mutex tidak perlu menggunakan pointer karena dia zero value, pointer hanya akan meningkatkan kompleksitas.
<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
mu := new(sync.Mutex)
mu.Lock()
```

</td><td>

```go
var mu sync.Mutex
mu.Lock()
```

</td></tr>
</tbody></table>

### Deklarasi Mutex di struct

Tidak diperbolehkan mengembed mutex di struct karena Lock dan Unlock belong to mutex dan bukan struct nya.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type SMap struct {
sync.Mutex

data map[string]string
}

func NewSMap() *SMap {
return &SMap{
data: make(map[string]string),
}
}

func (m *SMap) Get(k string) string {
m.Lock()
defer m.Unlock()

return m.data[k]
}
```

</td><td>

```go
type SMap struct {
mu sync.Mutex

data map[string]string
}

func NewSMap() *SMap {
return &SMap{
data: make(map[string]string),
}
}

func (m *SMap) Get(k string) string {
m.mu.Lock()
defer m.mu.Unlock()

return m.data[k]
}
```

</td></tr>
</tbody></table>

## Slice & Map

### Updating slice & map across function

Ketika kita pass slice & map ke function, sebenarnya dia pass by reference. Sehingga jika ada update value di
dalam function, variable awal juga akan berubah

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func (d *Driver) SetTrips(trips []Trip) {
d.trips = trips
}

trips := ...
d1.SetTrips(trips)

// Did you mean to modify d1.trips?
trips[0] = ...
```

</td><td>

```go
func (d *Driver) SetTrips(trips []Trip) {
d.trips = make([]Trip, len(trips))
copy(d.trips, trips)
}

trips := ...
d1.SetTrips(trips)

// We can now modify trips[0] without affecting d1.trips.
trips[0] = ...
```

</td></tr>
</tbody></table>

### Getting map & slice from function

Begitu pula data yang diambil dari slice dia juga pass by reference
<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Stats struct {
mu sync.Mutex
counters map[string]int
}

// Snapshot returns the current stats.
func (s *Stats) Snapshot() map[string]int {
s.mu.Lock()
defer s.mu.Unlock()

return s.counters
}

// snapshot is no longer protected by the mutex, so any
// access to the snapshot is subject to data races.
snapshot := stats.Snapshot()
```

</td><td>

```go
type Stats struct {
mu sync.Mutex
counters map[string]int
}

func (s *Stats) Snapshot() map[string]int {
s.mu.Lock()
defer s.mu.Unlock()

result := make(map[string]int, len(s.counters))
for k, v := range s.counters {
result[k] = v
}
return result
}

// Snapshot is now a copy.
snapshot := stats.Snapshot()
```

</td></tr>
</tbody></table>

## Test

### Bedakan nama package untuk test file

Dalam 1 folder kita bisa memiliki nama package yang berbeda misal `user` untuk base user dan `user_test` untuk file unit
test user. Hal ini berfungsi agar unexported function tidak muncul.

### Sorting variable berdasarkan size tipe data
