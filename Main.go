package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Kitchen = "кухня"
	Street  = "улица"
	Bedroom = "комната"
	Hall    = "коридор"
)

var (
	character       *MainCharacter
	itemsAndObjects map[string]string
	mapAndAction    map[string]string
	questItems      []*Item
)

/*
код писать в этом файле
наверняка у вас будут какие-то структуры с методами, глобальные переменные ( тут можно ), функции
*/

type Item struct {
	name string
}

type Door struct {
	from string
	to   string
	lock bool
}

type Furniture struct {
	name                string
	nameInPrepositional string
	items               []Item
	storages            []Storage
}

type Room struct {
	name           string
	rooms          []*Room
	furniture      []Furniture
	massageLooking func() string
	doors          *Door
	infoAboutRoom  string
}

type MainCharacter struct {
	name        string
	currentRoom *Room
	storage     Storage
	readyToGo   bool
}

type Storage struct {
	name  string
	items map[string]Item
}

func NewDoor(from string, to string) *Door {
	return &Door{
		from: from,
		to:   to,
		lock: false,
	}
}

func NewStorage(name string, items map[string]Item) *Storage {
	return &Storage{
		name:  name,
		items: items,
	}
}

func NewItem(name string) *Item {
	return &Item{name: name}
}

func NewFurniture(name string, items []Item, storages []Storage, nameInPrepositional string) *Furniture {
	return &Furniture{
		name:                name,
		items:               items,
		storages:            storages,
		nameInPrepositional: nameInPrepositional,
	}
}

func NewRooms(name string, rooms []*Room, furniture []Furniture, massageString string, doors *Door, infoAboutRoom string) *Room {
	return &Room{
		name:          name,
		rooms:         rooms,
		furniture:     furniture,
		doors:         doors,
		infoAboutRoom: infoAboutRoom,
		massageLooking: func() string {
			var result = massageString
			var sizeFurnitures = 0
			for _, fur := range furniture {
				if len(fur.items) == 0 && len(fur.storages) == 0 {
					sizeFurnitures++
					continue
				}

				result += fur.nameInPrepositional

				for _, item := range fur.items {
					result += " " + item.name + ","
				}

				for _, stor := range fur.storages {
					result += " " + stor.name + ","
				}
			}

			if sizeFurnitures == len(furniture) {
				result += "пустая " + name + "."
			}

			if character.currentRoom.name == Bedroom {
				runes := []rune(result)
				runes[len(runes)-1] = '.'
				result = string(runes)
			}

			if name == Kitchen {

				if !character.readyToGo {
					result += " надо собрать рюкзак и идти в универ."
				} else {
					result += " надо идти в универ."
				}
			}

			result += character.MoveStringMakerExtra(character.currentRoom.name)

			return result
		},
	}
}

func NewMainCharacter(name string, currentRoom *Room, storage Storage) *MainCharacter {
	return &MainCharacter{
		storage:     storage,
		currentRoom: currentRoom,
		name:        name,
		readyToGo:   true,
	}
}

func (character *MainCharacter) CheckGoOut(room string) bool {
	return room == Street && character.currentRoom.name == Hall &&
		character.currentRoom.doors != nil && !character.currentRoom.doors.lock
}

func (character *MainCharacter) TakeStorage(storageName string) string {
	for itFur, fur := range character.currentRoom.furniture {

		for it, stor := range fur.storages {
			if stor.name == storageName {
				character.currentRoom.furniture[itFur].storages = append(character.currentRoom.furniture[itFur].storages[:it], character.currentRoom.furniture[itFur].storages[it+1:]...)
				character.storage = stor
				return "вы надели: " + storageName
			}
		}
	}
	return "нет такого"
}

func (character *MainCharacter) UseItem(item string, object string) string {
	var flag = false
	for _, it := range character.storage.items {
		if it.name == item {
			flag = true
		}
	}

	if !flag {
		return "нет предмета в инвентаре - " + item
	}

	if itemsAndObjects[item] == object {
		if character.currentRoom.doors != nil {
			character.currentRoom.doors.lock = true
		}
		return mapAndAction[item+" "+object]
	} else {
		return "не к чему применить"
	}
}

func (character *MainCharacter) TakeItem(itemName string) string {

	if character.storage.name == "" {
		return "некуда класть"
	}

	for itFur, fur := range character.currentRoom.furniture {
		for it, item := range fur.items {
			if item.name == itemName {
				character.storage.items[itemName] = item
				character.currentRoom.furniture[itFur].items = append(character.currentRoom.furniture[itFur].items[:it], character.currentRoom.furniture[itFur].items[it+1:]...)
				return "предмет добавлен в инвентарь: " + item.name
			}
		}
	}
	return "нет такого"
}

func (character *MainCharacter) updateReadyToGo() {

	result := true
	for _, quIt := range questItems {
		if quIt.name != character.storage.items[quIt.name].name {
			result = false
		}
	}

	character.readyToGo = result
}

func (character *MainCharacter) LookAround() string {
	character.updateReadyToGo()
	return character.currentRoom.massageLooking()
}

func (character *MainCharacter) MoveStringMakerExtra(room string) string {

	if room == Street {
		return " можно пройти - домой"
	}

	var result string
	for it, curRoom := range character.currentRoom.rooms {
		if it == len(character.currentRoom.rooms)-1 {
			result += curRoom.name
		} else {
			result += curRoom.name + ", "
		}
	}

	return " можно пройти - " + result

}

func (character *MainCharacter) MoveStringMaker(room string) string {

	return character.currentRoom.infoAboutRoom + character.MoveStringMakerExtra(room)

}

func (character *MainCharacter) Move(room string) string {
	for _, iteratorRoom := range character.currentRoom.rooms {
		if iteratorRoom.name == room {

			if character.CheckGoOut(room) {
				return "дверь закрыта"
			}

			character.currentRoom = iteratorRoom

			return character.MoveStringMaker(room)
		}
	}

	return "нет пути в " + room
}

func main() {
	/*
		в этой функции можно ничего не писать,
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/

	initGame()
	var command string
	for command != "стоп" {
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		command = reader.Text()
		fmt.Println(handleCommand(command))
	}

}

func initGame() {
	/*
		эта функция инициализирует игровой мир - все комнаты
		если что-то было - оно корректно перезатирается
	*/
	itemsAndObjects = map[string]string{
		"ключи":   "дверь",
		"телефон": "шкаф",
	}
	mapAndAction = map[string]string{
		"ключи дверь":  "дверь открыта",
		"телефон шкаф": "телефон положили на шкаф",
	}

	var key = NewItem("ключи")
	var notes = NewItem("конспекты")
	var bag = NewStorage("рюкзак", make(map[string]Item))
	var cupOfTea = NewItem("чай")

	questItems = []*Item{
		key, notes,
	}

	var table = NewFurniture("стол", []Item{*cupOfTea}, make([]Storage, 0), "на столе:")
	var wardrobe = NewFurniture("шкаф", make([]Item, 0), make([]Storage, 0), " в шкафу:")
	var chair = NewFurniture("стул", make([]Item, 0), []Storage{*bag}, " на стуле:")
	var bedroomTable = NewFurniture("стол", []Item{*key, *notes}, make([]Storage, 0), "на столе:")

	var doorHallToStreet = NewDoor("коридор", "улица")

	var hall = NewRooms("коридор", make([]*Room, 0), []Furniture{*wardrobe}, "", doorHallToStreet, "ничего интересного.")
	var kitchen = NewRooms("кухня", []*Room{hall}, []Furniture{*table}, "ты находишься на кухне, ", nil, "кухня, ничего интересного.")
	var bedroom = NewRooms("комната", []*Room{hall}, []Furniture{*bedroomTable, *chair}, "", nil, "ты в своей комнате.")
	var street = NewRooms("улица", []*Room{hall}, make([]Furniture, 0), "", nil, "на улице весна.")

	hall.rooms = append(hall.rooms, kitchen, bedroom, street)

	character = NewMainCharacter("Ванька", kitchen, Storage{"", make(map[string]Item)})
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/

	fmt.Printf("\n")
	var parseString = strings.Split(command, " ")

	if len(parseString) == 1 && parseString[0] == "осмотреться" {
		return character.LookAround()
	}

	if len(parseString) == 2 && parseString[0] == "идти" {
		return character.Move(parseString[1])
	}

	if len(parseString) == 2 && parseString[0] == "надеть" {
		return character.TakeStorage(parseString[1])
	}

	if len(parseString) == 2 && parseString[0] == "взять" {
		return character.TakeItem(parseString[1])
	}

	if len(parseString) == 3 && parseString[0] == "применить" {
		return character.UseItem(parseString[1], parseString[2])
	}

	return "неизвестная команда"
}
