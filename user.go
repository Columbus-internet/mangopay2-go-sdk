// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"fmt"
)

type UserList []*User

// User is used by the user activity API and describe common fields to
// both natural and legal users.
type User struct {
	PersonType   string
	Email        string
	Id           string
	Tag          string
	CreationDate int
}

func (u *User) String() string {
	return fmt.Sprintf(`
Person type             : %s
Email                   : %s
Id                      : %s
Tag                     : %s
CreationDate            : %d`, u.PersonType, u.Email, u.Id, u.Tag, u.CreationDate)
}

// Users returns a list of all registered users, either natural
// or legal.
func (m *MangoPay) Users() (UserList, error) {
	resp, err := m.request(actionAllUsers, nil)
	if err != nil {
		return nil, err
	}
	ul := UserList{}
	if err := m.unMarshalJSONResponse(resp, &ul); err != nil {
		return nil, err
	}
	return ul, nil
}

// User fetch a user (natural or legal) using the Id attribute.
func (m *MangoPay) User(id string) (*User, error) {
	u, err := m.userRequest(actionFetchUser, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *MangoPay) userRequest(action mangoAction, data JsonObject) (*User, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(User)
	if err := m.unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}