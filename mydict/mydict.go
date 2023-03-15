package mydict

import (
	"errors"
)

// type에 method를 붙일수 있음

// Dictionary type
type Dictionary map[string]string


// Error
var (
	ErrNotFound = errors.New("not Found")
	ErrWordExists = errors.New("that word already exist")
	ErrCantUpdate = errors.New("cant update non-existing word")
)
// Search for a word
func (d Dictionary) Search(word string) (string, error) {

	value, exists := d[word]

	if exists {
		return value, nil
	}
	return "", ErrNotFound
}

// Add a word to the dictionary
func(d Dictionary) Add(key, value string) error {
	_, err := d.Search(key)
	
	switch err {
		case ErrNotFound:
			d[key] = value
		case nil:
			return ErrWordExists
	}
	return nil
}

// Update a word
func (d Dictionary) Update(key, value string) error {
	_, err:= d.Search(key)

	switch err{
		case ErrNotFound:
			return ErrCantUpdate
		case nil:
			d[key] = value
	}
	return nil;
}

// Delete a word
func (d Dictionary) Delete(key string) {
	delete(d, key)
}
