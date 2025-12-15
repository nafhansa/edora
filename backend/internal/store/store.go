package store

import (
    "encoding/json"
    "io/fs"
    "os"
    "path/filepath"
    "sync"
    "time"

    "edora/backend/internal/models"
)

type Store struct{
    dir string
    mu sync.RWMutex
    users []models.User
    readings []map[string]interface{}
}

func New(dir string) (*Store, error) {
    s := &Store{dir: dir}
    if err := s.load(); err != nil {
        return nil, err
    }
    return s, nil
}

func (s *Store) load() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // ensure dir exists
    if _, err := os.Stat(s.dir); os.IsNotExist(err) {
        if err := os.MkdirAll(s.dir, 0755); err != nil {
            return err
        }
    }

    // users.json
    uPath := filepath.Join(s.dir, "users.json")
    if b, err := os.ReadFile(uPath); err == nil {
        var us []models.User
        if err := json.Unmarshal(b, &us); err == nil {
            s.users = us
        }
    }

    // readings.json
    rPath := filepath.Join(s.dir, "readings.json")
    if b, err := os.ReadFile(rPath); err == nil {
        var rs []map[string]interface{}
        if err := json.Unmarshal(b, &rs); err == nil {
            s.readings = rs
        }
    }

    return nil
}

func (s *Store) SaveUsers() error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    b, err := json.MarshalIndent(s.users, "", "  ")
    if err != nil { return err }
    return os.WriteFile(filepath.Join(s.dir, "users.json"), b, fs.FileMode(0644))
}

func (s *Store) Users() []models.User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return append([]models.User(nil), s.users...)
}

func (s *Store) Readings() []map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return append([]map[string]interface{}(nil), s.readings...)
}

func (s *Store) FindUserByEmail(email string) *models.User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, u := range s.users {
        if u.Email == email {
            uu := u
            return &uu
        }
    }
    return nil
}

func (s *Store) FindUserByID(id string) *models.User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, u := range s.users {
        if u.ID == id {
            uu := u
            return &uu
        }
    }
    return nil
}

func (s *Store) AddUser(u models.User) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.users = append(s.users, u)
    return s.SaveUsers()
}

func (s *Store) AddReading(r map[string]interface{}) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    r["timestamp"] = time.Now().UTC().Format(time.RFC3339)
    s.readings = append(s.readings, r)
    b, err := json.MarshalIndent(s.readings, "", "  ")
    if err != nil { return err }
    return os.WriteFile(filepath.Join(s.dir, "readings.json"), b, fs.FileMode(0644))
}
