package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func NewUserDB(dbPath string) (db *UserDB, err error) {
	db = &UserDB{
		RWMutex:  &sync.RWMutex{},
		users:    []User{},
		usersMap: map[interface{}]User{},
		dbPath:   dbPath,
	}

	loadErr := db.Load()
	if !os.IsNotExist(loadErr) {
		err = loadErr
	}

	return
}

type UserDB struct {
	*sync.RWMutex
	users    []User
	usersMap map[interface{}]User
	dbPath   string
}

func (u *UserDB) Create(user User) error {

	if u.Exists(user.ID) {
		return fmt.Errorf("users with ID %d already exists", user.ID)
	}
	// new user, make sure they can't spoof the access rights
	user.Access = UANone
	u.Lock()
	u.users = append(u.users, user)
	u.usersMap[user.ID] = user
	u.Unlock()
	return u.Save()
}

func (u *UserDB) Update(user User) error {
	if !u.Exists(user.ID) {
		return fmt.Errorf("users with ID %d doesn't exist", user.ID)
	}
	u.Lock()
	for i := 0; i < len(u.users); i++ {
		if u.users[i].ID == user.ID {
			u.users[i] = user
			break
		}
	}
	u.Unlock()
	u.usersMap[user.ID] = user
	return u.Save()
}

func (u *UserDB) Delete(user User) error {
	if !u.Exists(user.ID) {
		return fmt.Errorf("users with ID %d doesn't exist", user.ID)
	}
	u.Lock()
	for i, usr := range u.users {
		if user.ID == usr.ID {
			u.users = append(u.users[:i], u.users[i+1:]...)
			break
		}
	}
	u.Unlock()
	delete(u.usersMap, user.ID)

	return u.Save()
}

func (u *UserDB) User(id interface{}) (User, bool) { // @todo: refactor to a nil check?
	u.RLock()
	defer u.RUnlock()
	user, exists := u.usersMap[id]
	return user, exists
}

func (u *UserDB) Exists(id interface{}) bool {
	u.RLock()
	defer u.RUnlock()
	_, ok := u.usersMap[id]
	return ok
}

func (u *UserDB) Users() []User {
	return u.users
}

func (u *UserDB) Admins() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAAdmin {
			result = append(result, user)
		}
	}
	return result
}

func (u *UserDB) Members() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAMember {
			result = append(result, user)
		}
	}
	return result
}

func (u *UserDB) Revoked() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UARevoked {
			result = append(result, user)
		}
	}
	return result
}

func (u *UserDB) IsAdmin(id int) bool {
	u.RLock()
	defer u.RUnlock()
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAAdmin
}

func (u *UserDB) IsMember(id int) bool {
	u.RLock()
	defer u.RUnlock()
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAMember
}

func (u *UserDB) IsRevoked(id int) bool {
	u.RLock()
	defer u.RUnlock()
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UARevoked
}

func (u *UserDB) Save() error {
	u.Lock()
	defer u.Unlock()

	// Open a temporary file to hold the new database
	var tempDb *os.File
	tempDb, err := ioutil.TempFile(filepath.Dir(u.dbPath), filepath.Base(u.dbPath))
	if err != nil {
		return err
	}

	var db = struct {
		Users []User `json:"users"`
	}{
		u.users,
	}

	// Write the data to the new file
	if err = json.NewEncoder(tempDb).Encode(db); err != nil {
		return err
	}

	// Close the file if we succeeded in opening one
	if err = tempDb.Close(); err != nil {
		return err
	}

	return os.Rename(tempDb.Name(), u.dbPath)
}

func (u *UserDB) Load() error {
	u.Lock()
	defer u.Unlock()

	raw, err := ioutil.ReadFile(u.dbPath)
	if err != nil {
		return err
	}

	var db struct {
		Users []User
	}

	if err = json.Unmarshal(raw, &db); err != nil {
		return err
	}
	u.users = db.Users
	u.usersMap = map[interface{}]User{}
	for _, user := range db.Users {
		u.usersMap[user.ID] = user
	}

	return nil
}
