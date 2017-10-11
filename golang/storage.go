package d7024e

import (
	"crypto/sha1"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"time"
)

/*
* A list of stored files, change accordingly later.
 */
type Storage struct {
	Files  []file
	mutex  sync.Mutex
	timers []timerHolder
}

func NewStorage() Storage {
	return Storage{Files: []file{}, mutex: sync.Mutex{}, timers: []timerHolder{}}
}

type file struct {
	Name []uint8
	Text []byte
}

type timerHolder struct {
	timer  *time.Timer
	fileID string
}

func (storage *Storage) deleteTimer(filename string) {
	for i, item := range storage.timers {
		if item.fileID == filename {
			storage.timers = append(storage.timers[:i], storage.timers[i+1:]...)
			return
		}
	}
}

func (storage *Storage) updateTimer(time time.Duration, filename string) {
	for _, item := range storage.timers {
		if item.fileID == filename {
			if !item.timer.Stop() {
				<-item.timer.C
			}
			item.timer.Reset(time)
			return
		}
	}
}

func (storage *Storage) addTimer(timer *time.Timer, filename string) {
	for i, item := range storage.timers {
		if item.fileID == filename {
			if !item.timer.Stop() {
				<-item.timer.C
			}
			storage.timers = append(storage.timers[:i], storage.timers[i+1:]...)
		}
	}
	storage.timers = append(storage.timers, timerHolder{timer, filename})
}

/*
* When a file misses the update, it is no longer prioritized and
* moved to Memory from RAM.
 */
func (storage *Storage) MoveToMemory(name []byte) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	file := file{}
	//compareName := storage.HashFile(name)
	for i, v := range storage.Files {
		if reflect.DeepEqual(v.Name, name) {
			file.Name = v.Name
			file.Text = v.Text

			//Delete file in Files
			storage.Files = append(storage.Files[:i], storage.Files[i+1:]...)

			//Insert into Memory
			storage.Memory(name, file.Text)
			//break out of for loop.
			break
		}
	}
}

/*
* Deletes a file from ram. returns true if it existed, false otherwise.
 */
func (storage *Storage) deleteFromRam(name string) bool {
	for i, v := range storage.Files {
		if reflect.DeepEqual(v.Name, []byte(name)) {
			//Delete file in Files
			storage.Files = append(storage.Files[:i], storage.Files[i+1:]...)
			return true
		}
	}
	return false
}

/*
* Check if file is in RAM
 */
func (storage *Storage) ReadRAM(name []byte) *file {
	file := file{}
	//compareName := storage.HashFile(name)
	for _, v := range storage.Files {
		if reflect.DeepEqual(v.Name, name) {
			file.Name = v.Name
			file.Text = v.Text
			return &file
		}
	}
	return nil
}

/*
* Look if the RAM storage include a certain file, if so return file, else
check Memory if it's there and return.
*/
func (storage *Storage) Search(name []byte) *file {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	returnedFile := storage.ReadRAM(name)
	if returnedFile == nil {
		// Check if memory has file.
		returnedFile = storage.ReadMemory(name)
	}
	return returnedFile
}

/*
* Read Memory and see if a file is there, if so, add file to RAM
and return it.
*/
func (storage *Storage) ReadMemory(name []byte) *file {
	filename := "../newfiles/" + string(name)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	//fmt.Printf("File contents: %s", content,"\n")
	returnedFile := &file{name, []byte(content)}

	/*
	 * When a file is retrieved from Memory, add it from Memory to RAM for "caching"
	 */
	storage.RAM(name, []byte(content))
	return returnedFile
}

func (storage *Storage) Store(name []byte, text []byte) bool {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	okRam := storage.RAM(name, text)
	okMem := storage.Memory(name, text)
	return (okRam && okMem)
}

/*
* Store a file into RAM, does not return anything.

 */
func (storage *Storage) RAM(name []byte, text []byte) bool {
	for _, v := range storage.Files {
		if reflect.DeepEqual(v.Name, name) {
			return false
		}
	}
	newFile := file{name, text}
	storage.Files = append(storage.Files, newFile)
	return true
}

/*
* Store a file into Memory, does not return anything.
 */
func (storage *Storage) Memory(name []byte, text []byte) bool {
	filename := "../newfiles/" + string(name)
	ok := false
	file := storage.ReadMemory(name)
	if file == nil {
		ok = true
	}
	err2 := ioutil.WriteFile(filename, text, 0644)
	if err2 != nil {
		panic(err2)
	}
	return ok
}

/*
* Convert a file name of type []byte into []uint8 (SHA-1)!
 */
func (storage *Storage) HashFile(name []byte) []uint8 {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	hashing := sha1.New()
	return hashing.Sum(name)
}

/*
* TODO: Dont always return true, check if file on memory exists.
 */
func (storage *Storage) DeleteFile(name string) bool {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	ramOK := storage.deleteFromRam(name)
	path := "../newfiles/" + name
	os.Remove(path) // clean up
	return (ramOK || true)
}
