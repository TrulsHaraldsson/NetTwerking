package d7024e

import (
	"bytes"
	"crypto/sha1"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestStorageStoreInRAM(t *testing.T) {
	storage := Storage{}
	name := []byte("filename")
	text := []byte("File body text.")
	storage.RAM(name, text)
	file := storage.ReadRAM(name)
	bool := bytes.EqualFold(file.Text, text)
	if bool == false {
		t.Error("File content do not match!\n", string(text), "\n", string(file.Text), "\n")
	}
}

func TestStorageStoreInMemory(t *testing.T) {
	storage := Storage{}
	name := []byte("filename2")
	text := []byte("This is a test content for a temp file.")
	storage.Memory(name, text)
	file := storage.ReadMemory(name)
	if file == nil {
		t.Error("File is nil")
	} else {
		bool := bytes.EqualFold(file.Text, text)
		if bool == false {
			t.Error("File content do not match!\n", string(text), "\n", string(file.Text), "\n")
		}
	}
	path := "./../newfiles/" + string(name)
	time.Sleep(time.Millisecond * 10)
	os.Remove(path)
}
func TestStorageMoveToMemory(t *testing.T) {
	storage := Storage{}
	nameT := []byte("filenameX450")
	textT := []byte("This is the content of filenameX450!")
	storage.RAM(nameT, textT)
	storage.MoveToMemory(nameT)
	file := storage.ReadMemory(nameT)
	if file == nil {
		t.Error("File is nil")
	} else {
		bool := bytes.EqualFold(textT, file.Text)
		if bool == false {
			t.Error("File content do not match!\n", string(textT), "\n", string(file.Text), "\n")
		}
	}
	time.Sleep(time.Millisecond * 10)
	path := "../newfiles/" + string(nameT)
	os.Remove(path)
}

func TestStorageSearch(t *testing.T) {
	storage := Storage{}
	name := []byte("filename3")
	text := []byte("File content when creating StorageSearch test!")
	storage.RAM(name, text)
	file := storage.Search(name)
	if file == nil {
		t.Error("File is nil")
	} else {
		bool := bytes.EqualFold(file.Text, text)
		if bool == false {
			t.Error("File content do not match!\n", string(text), "\n", string(file.Text), "\n")
		}
	}

	name2 := []byte("filename4")
	text2 := []byte("File content when creating StorageSearch test 2!")
	storage.Memory(name2, text2)
	filed := storage.ReadMemory(name2)
	if filed == nil {
		t.Error("File is nil")
	} else {
		bool2 := bytes.EqualFold(filed.Text, text2)
		if bool2 == false {
			t.Error("File content do not match!\n", string(text2), "\n", string(filed.Text), "\n")
		}
	}
	path := "./../newfiles/" + string(name2)
	os.Remove(path)
}

func TestStorageHash(t *testing.T) {
	storage := Storage{}
	name := []byte("filenameB")
	hash := sha1.New()
	hashedName := hash.Sum(name)
	returnedHash := storage.HashFile(name)
	if reflect.DeepEqual(hashedName, returnedHash) == false {
		t.Error("Hashing is not correct! \n", hashedName, "\n", returnedHash, "\n")
	}
}
