package d7024e

import (
	"testing"
)

func TestStorageStoreInRAM(t *testing.T) {
	//fmt.Println("Storing in RAM!\n")
	storage := Storage{}
	name := []byte("filename")
	text := []byte("File body text.")
	storage.RAM(name, text)
}

func TestStorageSearchRAM(t *testing.T) {
	//fmt.Println("Search for file in RAM!\n")
	storage := Storage{}
	name := []byte("filename")
	text := []byte("File body text.")
	storage.RAM(name, text)
	//file := storage.Search(name) //TODO: uncomment
	//TODO: Fix these function, so its checks file is correct..
	//fmt.Println("File returned from RAM!\nName: ", file.Name,"\nContent: ", file.Text,"\n")
}

func TestStorageStoreInMemory(t *testing.T) {
	//fmt.Println("Storing in Memory!\n")
	storage := Storage{}
	name := []byte("filename")
	text := []byte("This is a test content for a temp file.\n")
	storage.Memory(name, text)
}

func TestStorageSearchMemory(t *testing.T) {
	//fmt.Println("Search for file in Memory!\n")
	storage := Storage{}
	name := []byte("filename4")
	text := []byte("This is a test content for a temp file.\n")
	storage.Memory(name, text)
	//file := storage.Search(name) //TODO: uncomment
	//fmt.Println("File returned from Memory!\nName: ", file.Name,"\nContent: ", file.Text,"\n")
}

func TestStorageMoveToMemory(t *testing.T) {
	//fmt.Println("Moving a file from RAM to Memory!")
	storage := Storage{}
	name := []byte("filenameX200")
	text := []byte("This is the content of filenameX200!\n")
	storage.RAM(name, text)
	storage.MoveToMemory(name)
	//fmt.Println("File has been moved from RAM to Memory!\n")
}
