// Package dialog is a package for managing command line interface
package dialog

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EgorKo25/GophKeeper/pkg/mycrypto"

	"github.com/EgorKo25/GophKeeper/internal/client"
	"github.com/EgorKo25/GophKeeper/internal/storage"

	"github.com/manifoldco/promptui"
)

var (
	// Command line font Style
	myStyler = promptui.Styler(promptui.FGBold, promptui.FGGreen)
)

// Manager is a struct for managing cli
type Manager struct {
	functions map[string]func(string) error
	cookie    []*http.Cookie
	user      *storage.User

	e *mycrypto.Crypto
	c *client.Client
}

// NewManager is a constructor Manager
func NewManager(c *client.Client, e *mycrypto.Crypto) *Manager {

	var dial Manager

	dial.e = e
	dial.c = c

	function := make(map[string]func(string) error)
	function["Registration"] = dial.Add
	function["Add"] = dial.Add
	function["Read"] = dial.Read
	function["Update"] = dial.Update
	function["Delete"] = dial.Delete
	function["Login"] = dial.Read

	dial.functions = function

	return &dial
}

// Run is a function for running cli
func (d *Manager) Run() error {

	d.sayHello()

	err := d.SelectAuth()
	if err != nil {
		return err
	}

	for {
		err = d.SelectFunc()
		if err != nil {
			return err
		}
	}

}

// sayHello is function for printing logotype
func (d *Manager) sayHello() {
	fmt.Println(`
╭╮╭╮╭╮╱╱╭╮╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭╮╱╱╱╱╭━━━╮╱╱╱╱╱╭╮╱╭╮╭━╮
┃┃┃┃┃┃╱╱┃┃╱╱╱╱╱╱╱╱╱╱╱╱╱╭╯╰╮╱╱╱┃╭━╮┃╱╱╱╱╱┃┃╱┃┃┃╭╯
┃┃┃┃┃┣━━┫┃╭━━┳━━┳╮╭┳━━╮╰╮╭╋━━╮┃┃╱╰╋━━┳━━┫╰━┫╰╯╯╭━━┳━━┳━━┳━━┳━╮
┃╰╯╰╯┃┃━┫┃┃╭━┫╭╮┃╰╯┃┃━┫╱┃┃┃╭╮┃┃┃╭━┫╭╮┃╭╮┃╭╮┃╭╮┃┃┃━┫┃━┫╭╮┃┃━┫╭╯
╰╮╭╮╭┫┃━┫╰┫╰━┫╰╯┃┃┃┃┃━┫╱┃╰┫╰╯┃┃╰┻━┃╰╯┃╰╯┃┃┃┃┃┃╰┫┃━┫┃━┫╰╯┃┃━┫┃
╱╰╯╰╯╰━━┻━┻━━┻━━┻┻┻┻━━╯╱╰━┻━━╯╰━━━┻━━┫╭━┻╯╰┻╯╰━┻━━┻━━┫╭━┻━━┻╯
╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱┃┃╱╱╱╱╱╱╱╱╱╱╱╱╱╱┃┃
╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╰╯╱╱╱╱╱╱╱╱╱╱╱╱╱╱╰╯`)

}

// SelectAuth is a function for user can select authorization method
func (d *Manager) SelectAuth() error {
	prompt := promptui.Select{
		Label: "Have an account?",
		Items: []string{"Registration", "Login", "Exit"},
	}

	_, res, err := prompt.Run()
	if err != nil {
		log.Println(err)
		return err
	}

	if res == "Exit" {
		os.Exit(0)
	}

	if v, ok := d.functions[res]; ok {
		err = v("user")
		if err != nil {
			return err
		}
		log.Println(myStyler(res))
	}

	return nil

}

// SelectFunc is a function for user can select one of function app
func (d *Manager) SelectFunc() error {

	prompt := promptui.Select{
		Label: "Выберте функцию",
		Items: []string{"Add", "Update", "Read", "Delete", "Delete an account", "Exit"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	if result == "Exit" {
		os.Exit(0)
	}

	prompt.Label = "Выберите тип данных"
	prompt.Items = []string{"password", "card", "binary data"}

	_, dataType, _ := prompt.Run()

	if v, ok := d.functions[result]; ok {
		return v(dataType)
	}

	return errors.New("unknown function")

}

// myPrompt is a function for printing control
func (d *Manager) myPrompt(label string) string {
	prompt := promptui.Prompt{
		Label: myStyler(myStyler(label)),
	}

	res, _ := prompt.Run()
	return res
}

// Add is a facade for adding new data into server
func (d *Manager) Add(dataType string) error {
	switch dataType {
	case "password":
		return d.addPassword(d.user.Login)
	case "card":
		return d.addCard(d.user.Login)
	case "binary data":
		return d.addBinData(d.user.Login)
	case "user":
		return d.addUser()
	default:
		return errors.New("unknown type")
	}
}

func (d *Manager) addUser() (err error) {

	var pass storage.User
	var code int

	pass.Login, err = d.e.Encrypt(d.myPrompt("Введите вашу почту"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Email, err = d.e.Encrypt(d.myPrompt("Введите ваш логин"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Password, err = d.e.Encrypt(d.myPrompt("Введите ваш пароль"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	d.user = &pass

	code, _, d.cookie, err = d.c.Send(&pass, "user", d.cookie, "/user/register")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 400 {
			fmt.Println(myStyler(myStyler("Неправильный формат данных")))
			return d.addUser()
		}
		return
	}

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) addPassword(login string) (err error) {

	var pass storage.Password
	var code int

	pass.LoginOwner = login

	pass.Service, err = d.e.Encrypt(d.myPrompt("Введите название сервиса"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Login, err = d.e.Encrypt(d.myPrompt("Введите ваш логин в сервису"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Password, err = d.e.Encrypt(d.myPrompt("Введите ваш пароль в сервису"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, _, d.cookie, err = d.c.Send(&pass, "password", d.cookie, "/user/add")
	if code != 200 {
		fmt.Println(myStyler("Что-то пошло не так"))
		return err
	}

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) addCard(login string) (err error) {

	var pass storage.Card
	var code int

	pass.LoginOwner = login

	pass.Bank, err = d.e.Encrypt(d.myPrompt("Введите название банка"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Number, err = d.e.Encrypt(d.myPrompt("Введите ваш номер карты"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.DataEnd, err = d.e.Encrypt(d.myPrompt("Введите ваш дату окончания карты"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.SecretCode, err = d.e.Encrypt(d.myPrompt("Введите ваш секретный код"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Owner, err = d.e.Encrypt(d.myPrompt("Введите владельца карты"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, _, d.cookie, err = d.c.Send(&pass, "card", d.cookie, "/user/add")
	if code != 200 {
		fmt.Println(myStyler("Что-то пошло не так: "), err)
		return err
	}

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) addBinData(login string) (err error) {

	var pass storage.BinaryData
	var code int

	pass.LoginOwner = login

	pass.Title, _ = d.e.Encrypt(d.myPrompt("Введите название файла"))
	path := d.myPrompt("Введите путь к файлу")

	pass.Data, err = os.ReadFile(path)
	if err != nil {
		fmt.Println(myStyler("Что-то пошло не так"))
		return err
	}

	data, err := d.e.Encrypt(string(pass.Data))
	if err != nil {
		fmt.Println(myStyler("Что-то пошло не так"))
		return err
	}

	pass.Data = []byte(data)

	code, _, d.cookie, err = d.c.Send(&pass, "bin", d.cookie, "/user/add")
	if code != 200 {
		fmt.Println(myStyler("Что-то пошло не так: "))
		return err
	}

	fmt.Println(myStyler("Готово"))
	return nil
}

// Read is a facade for reading data from server
func (d *Manager) Read(dataType string) error {
	switch dataType {
	case "password":
		return d.readPassword()
	case "card":
		return d.readCard()
	case "binary data":
		return d.readBinData()
	case "user":
		return d.readUser()
	default:
		return errors.New("unknown type")

	}
}

func (d *Manager) readUser() (err error) {

	var pass storage.User
	var code int

	pass.Login, err = d.e.Encrypt(d.myPrompt("Введите ваш логин"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Password, err = d.e.Encrypt(d.myPrompt("Введите ваш пароль"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	d.user = &pass

	code, _, d.cookie, err = d.c.Send(&pass, "user", d.cookie, "/user/login")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}
		return
	}

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) readPassword() (err error) {

	var code int
	var pass storage.Password
	var tmp any

	pass.LoginOwner = d.user.Login

	pass.Service, err = d.e.Encrypt(d.myPrompt("Введите название сервиса"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, tmp, d.cookie, err = d.c.Send(&pass, "password", d.cookie, "/user/read")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого сервиса"))
		return
	}

	pass.Service, _ = d.e.Decrypt(tmp.(storage.Password).Service)
	pass.Login, _ = d.e.Decrypt(tmp.(storage.Password).Login)
	pass.Password, _ = d.e.Decrypt(tmp.(storage.Password).Password)

	fmt.Printf("Название сервиса: %s\nЛогин: %s\nПароль: %s\n", pass.Service, pass.Login, pass.Password)

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) readCard() (err error) {

	var code int
	var pass storage.Card
	var tmp any

	pass.LoginOwner = d.user.Login

	pass.Bank, err = d.e.Encrypt(d.myPrompt("Введите название банка"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, tmp, d.cookie, err = d.c.Send(&pass, "card", d.cookie, "/user/read")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такой карты"))
		return
	}

	pass.Bank, _ = d.e.Decrypt(tmp.(storage.Card).Bank)
	pass.Number, _ = d.e.Decrypt(tmp.(storage.Card).Number)
	pass.DataEnd, _ = d.e.Decrypt(tmp.(storage.Card).DataEnd)
	pass.SecretCode, _ = d.e.Decrypt(tmp.(storage.Card).SecretCode)

	fmt.Printf("Название банка: %s\nНомер карты: %s\nДата окончания: %s\nСекретный код: %s\n",
		pass.Bank, pass.Number, pass.DataEnd, pass.SecretCode)

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) readBinData() (err error) {

	var code int
	var pass storage.BinaryData
	var tmp any

	pass.LoginOwner = d.user.Login

	pass.Title, err = d.e.Encrypt(d.myPrompt("Введите название файла"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, tmp, d.cookie, err = d.c.Send(&pass, "bin", d.cookie, "/user/read")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого файла"))
		return
	}

	pass.Title, _ = d.e.Decrypt(tmp.(storage.BinaryData).Title)
	data, _ := d.e.Decrypt(string(tmp.(storage.BinaryData).Data))

	pass.Data = []byte(data)

	fmt.Printf("Название название файла: %s\nСодержимое: %s\n",
		pass.Title, pass.Data)

	fmt.Println(myStyler("Готово"))
	return nil
}

// Update is a facade for updating data from server
func (d *Manager) Update(dataType string) error {
	switch dataType {
	case "password":
		return d.updatePassword()
	case "card":
		return d.updateCard()
	case "binary data":
		return d.updateBinData()
	default:
		return errors.New("unknown type")

	}
}

func (d *Manager) updatePassword() (err error) {

	var code int
	var pass storage.Password
	var tmp any

	pass.LoginOwner = d.user.Login

	pass.Service, err = d.e.Encrypt(d.myPrompt("Введите название сервиса"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Login, err = d.e.Encrypt(d.myPrompt("Введите новый логин (если он не изменилося введите старый)"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Password, err = d.e.Encrypt(d.myPrompt("Введите новый пароль (если он не изменилося введите старый)"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, tmp, d.cookie, err = d.c.Send(&pass, "password", d.cookie, "/user/update")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого файла"))
		return
	}

	pass.Service, _ = d.e.Decrypt(tmp.(storage.Password).Service)
	pass.Login, _ = d.e.Decrypt(tmp.(storage.Password).Login)
	pass.Password, _ = d.e.Decrypt(tmp.(storage.Password).Password)

	fmt.Printf("Название сервиса: %s\nЛогин: %s\nПароль: %s\n", pass.Service, pass.Login, pass.Password)

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) updateCard() (err error) {

	var code int
	var pass storage.Card
	var tmp any

	pass.LoginOwner = d.user.Login

	pass.Bank, err = d.e.Encrypt(d.myPrompt("Введите название банка"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.Number, err = d.e.Encrypt(d.myPrompt("Введите новый номер карты (если он не изменилося введите старый)"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.DataEnd, err = d.e.Encrypt(d.myPrompt("Введите новую дату окончания (если она не изменилося введите старый)"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.SecretCode, err = d.e.Encrypt(d.myPrompt("Введите новый секретный код (если он не изменилося введите старый)"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}
	pass.LoginOwner, err = d.e.Encrypt(d.myPrompt("Введите новое ФИО владельца (если оно не изменилося введите старый)"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, tmp, d.cookie, err = d.c.Send(&pass, "card", d.cookie, "/user/update")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 403 {
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такогой карты"))
		return
	}

	fmt.Println(myStyler("Изения внесены!"))

	pass.Bank, _ = d.e.Decrypt(tmp.(storage.Card).Bank)
	pass.Number, _ = d.e.Decrypt(tmp.(storage.Card).Number)
	pass.DataEnd, _ = d.e.Decrypt(tmp.(storage.Card).DataEnd)
	pass.SecretCode, _ = d.e.Decrypt(tmp.(storage.Card).SecretCode)

	fmt.Printf("Название банка: %s\nНомер карты: %s\nДата окончания: %s\nСекретный код: %s\n",
		pass.Bank, pass.Number, pass.DataEnd, pass.SecretCode)

	fmt.Println(myStyler("Готово"))
	return nil
}

func (d *Manager) updateBinData() (err error) {

	var code int
	var pass storage.BinaryData
	var tmp any

	pass.LoginOwner = d.user.Login

	pass.Title, err = d.e.Encrypt(d.myPrompt("Введите название файла"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	path := d.myPrompt("Введите путь к новому файлу")

	pass.Data, err = os.ReadFile(path)
	if err != nil {
		fmt.Println(myStyler("Что-то пошло не так"))
		return err

	}

	data, err := d.e.Encrypt(string(pass.Data))
	if err != nil {
		fmt.Println(myStyler("Что-то пошло не так"))
		return err
	}

	pass.Data = []byte(data)

	code, tmp, d.cookie, err = d.c.Send(&pass, "bin", d.cookie, "/user/update")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого файла"))
		return
	}

	fmt.Println(myStyler("Данные обновлены"))

	pass.Title, _ = d.e.Decrypt(tmp.(storage.BinaryData).Title)
	data, _ = d.e.Decrypt(string(tmp.(storage.BinaryData).Data))

	pass.Data = []byte(data)

	fmt.Printf("Название название файла: %s\nСодержимое: %s\n",
		pass.Title, pass.Data)

	fmt.Println(myStyler("Готово"))
	return nil
}

// Delete is a facade for deleting data from server
func (d *Manager) Delete(dataType string) error {
	switch dataType {
	case "password":
		return d.deletePassword()
	case "card":
		return d.deleteCard()
	case "binary data":
		return d.deleteBinData()
	case "delete an account":
		return d.deleteUser()
	default:
		return errors.New("unknown type")

	}
}

func (d *Manager) deletePassword() (err error) {

	var pass storage.Password
	var code int

	pass.LoginOwner = d.user.Login

	pass.Service, err = d.e.Encrypt(d.myPrompt("Введите название сервиса"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, _, d.cookie, err = d.c.Send(&pass, "password", d.cookie, "/user/delete")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return d.SelectFunc()
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого файла"), err)
		return d.SelectFunc()
	}

	fmt.Println(myStyler("Данные удалены"))

	return nil
}

func (d *Manager) deleteCard() (err error) {

	var pass storage.Card
	var code int

	pass.LoginOwner = d.user.Login

	pass.Bank, err = d.e.Encrypt(d.myPrompt("Введите название сервиса"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, _, d.cookie, err = d.c.Send(&pass, "card", d.cookie, "/user/delete")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")))
			return d.SelectFunc()
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого файла"), err)
		return d.SelectFunc()
	}

	fmt.Println(myStyler("Данные удалены"))

	return nil
}

func (d *Manager) deleteBinData() (err error) {

	var pass storage.BinaryData
	var code int

	pass.LoginOwner = d.user.Login

	pass.Title, err = d.e.Encrypt(d.myPrompt("Введите название файла"))
	if err != nil {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")), err)
			return
		}
	}

	code, _, d.cookie, err = d.c.Send(&pass, "bin", d.cookie, "/user/delete")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")))
			return d.SelectFunc()
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого файла"), err)
		return d.SelectFunc()
	}
	fmt.Println(myStyler("Данные удалены"))

	return nil
}

func (d *Manager) deleteUser() (err error) {

	var pass storage.User
	var code int

	if y := d.myPrompt("Осторожно! Вы действительно хотите удалить аккаунт? (y/n)"); y != "y" {
		return d.SelectFunc()
	}

	pass.Login = d.user.Login

	code, _, d.cookie, err = d.c.Send(&pass, "user", d.cookie, "/user/delete")
	if code != 200 {
		if err != nil {
			fmt.Println(myStyler(myStyler("Что-то пошло не так: ")))
			return d.SelectFunc()
		}
		if code == 403 {
			fmt.Println(myStyler(myStyler("Нет такого пользователя")))
			return d.SelectAuth()
		}

		fmt.Println(myStyler("Нет такого файла"), err)
		return d.SelectFunc()
	}

	fmt.Println(myStyler("Данные удалены"))

	return nil
}
