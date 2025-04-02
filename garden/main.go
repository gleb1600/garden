package main

import (
	"bufio"
	"context"
	"fmt"
	plantsfunc "garden/plants"
	dbfunc "garden/storage"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	connString := "postgres://garden:secret@localhost:5431/gardendb?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		fmt.Printf("Ошибка подключения: %v\n", err)
	}
	defer pool.Close()
	db := dbfunc.NewPostgresDB(ctx, pool)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	commands := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		printHelp()
		for {
			fmt.Print("> ")
			scanner.Scan()
			commands <- scanner.Text()
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Виртуальный сад открыт!")

	for {
		select {
		case <-ticker.C:
			plants, err := db.GetAllPlants()
			if err != nil {
				fmt.Printf("Ошибка взятия из БД: %v\n", err)
				continue
			}

			for i := range plants {
				plantsfunc.UpdatePlantState(&plants[i])
				if err := db.UpdatePlant(&plants[i]); err != nil {
					fmt.Printf("Ошибка обновления: %v\n", err)
				}
			}

		case cmd := <-commands:
			handleCommand(db, cmd)

		case <-done:
			fmt.Println("\nСад закрывается...")
			return
		}
	}
}

func printHelp() {
	fmt.Println(`
Доступные команды:
  list                 - Показать все растения
  add [species] [name] - Посадить новое растение
  remove [id]          - Удалить растение
  water [id]           - Полить растение
  fertilize [id]       - Удобрить растение
  help                 - Показать помощь
  exit                 - Выйти из программы`)
}

func handleCommand(db *dbfunc.PostgresDB, cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "help":
		printHelp()

	case "list":
		plants, err := db.GetAllPlants()
		if err != nil {
			fmt.Printf("Ошибка взятия из БД: %v\n", err)
			return
		}
		plantsfunc.PrintGardenState(plants)

	case "water":
		notValid, errorText, id := idIsValid(db, parts)
		if notValid {
			fmt.Println(errorText)
			return
		}
		if err := db.WaterPlant(id); err != nil {
			fmt.Printf("Ошибка полива: %v\n", err)
		} else {
			fmt.Printf("Растение %d полито!\n", id)
		}

	case "fertilize":
		notValid, errorText, id := idIsValid(db, parts)
		if notValid {
			fmt.Println(errorText)
			return
		}
		if err := db.FertilizePlant(id); err != nil {
			fmt.Printf("Ошибка удобрения: %v\n", err)
		} else {
			fmt.Printf("Растение %d удобрено!\n", id)
		}

	case "remove":
		notValid, errorText, id := idIsValid(db, parts)
		if notValid {
			fmt.Println(errorText)
			return
		}
		err := db.DeletePlant(id)
		if err != nil {
			fmt.Printf("Ошибка удаления растения: %v\n", err)
			return
		}
		fmt.Printf("растение с ID %d удалено", id)

	case "add":
		if len(parts) < 3 {
			fmt.Println("Укажите [вид] название растения")
			return
		}
		species := parts[1]
		name := parts[2]
		if len(name) > 255 || len(species) > 255 {
			fmt.Println("Название и вид не должны превышать 255 символов")
			return
		}
		exists, err := db.PlantExistsByName(name)
		if err != nil {
			fmt.Printf("Ошибка проверки на имя: %v\n", err)
			return
		}
		if exists {
			fmt.Println("Растение с таким именем уже существует")
			return
		}
		err = db.CreatePlantWithParams(name, species)
		if err != nil {
			fmt.Printf("Ошибка создания растения: %v\n", err)
			return
		}
		fmt.Printf("Растение [%s] %s добавлено!\n", species, name)
	case "exit":
		os.Exit(0)

	default:
		fmt.Println("Неизвестная команда. Введите 'help' для помощи")
	}
}

func idIsValid(db *dbfunc.PostgresDB, parts []string) (bool, string, int) {
	if len(parts) < 2 {
		fmt.Println("Укажите ID растения")
		return true, "Укажите ID растения", 0
	}
	id, err := strconv.Atoi(parts[1])
	if id_exist, errdb := db.PlantExists(id); !id_exist || errdb != nil || err != nil {
		fmt.Println("Неверный ID растения")
		return true, "Неверный ID растения", 0
	} else {
		return false, "", id
	}
}
