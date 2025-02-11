package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

// Константы
const (
	chromeDriverURL = "http://localhost:4444"
	loginURL        = "http://localhost:8080/loginPage"
	recipePageURL   = "http://localhost:8080/recipesPage"
)

// Тест страницы рецептов в Google Chrome
func TestRecipePageChrome(t *testing.T) {
	// Настройки Chrome WebDriver (без headless)
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeArgs := map[string]interface{}{"args": []string{}} // Убираем headless
	caps["goog:chromeOptions"] = chromeArgs

	// Подключение к WebDriver
	wd, err := selenium.NewRemote(caps, chromeDriverURL)
	if err != nil {
		t.Fatalf("Ошибка подключения к WebDriver: %v", err)
	}
	defer wd.Quit()

	// Открываем страницу логина
	if err := wd.Get(loginURL); err != nil {
		t.Fatalf("Ошибка загрузки страницы логина: %v", err)
	}

	// Ввод email
	emailInput, err := waitForElement(wd, selenium.ByID, "email", 5*time.Second)
	if err != nil {
		t.Fatalf("Поле email не найдено: %v", err)
	}
	if err := emailInput.SendKeys("adil_groud@mail.ru"); err != nil {
		t.Fatalf("Ошибка ввода email: %v", err)
	}

	// Ввод пароля
	passwordInput, err := waitForElement(wd, selenium.ByID, "password", 5*time.Second)
	if err != nil {
		t.Fatalf("Поле пароля не найдено: %v", err)
	}
	if err := passwordInput.SendKeys("1"); err != nil {
		t.Fatalf("Ошибка ввода пароля: %v", err)
	}

	// Клик по кнопке логина
	loginButton, err := waitForElement(wd, selenium.ByCSSSelector, "button.btn", 5*time.Second)
	if err != nil {
		t.Fatalf("Кнопка логина не найдена: %v", err)
	}
	if err := loginButton.Click(); err != nil {
		t.Fatalf("Ошибка клика по кнопке логина: %v", err)
	}

	// Ожидание алерта
	time.Sleep(2 * time.Second)
	alert, err := wd.AlertText()
	if err != nil {
		t.Fatalf("Алерт не найден: %v", err)
	}

	if alert != "You logged in successfully" {
		t.Fatalf("Неверный текст алерта: %s", alert)
	}

	// Закрытие алерта
	if err := wd.AcceptAlert(); err != nil {
		t.Fatalf("Ошибка закрытия алерта: %v", err)
	}

	time.Sleep(5 * time.Second)
	// Проверяем текущий URL
	currentURL, err := wd.CurrentURL()
	if err != nil {
		t.Fatalf("Ошибка получения текущего URL: %v", err)
	}

	// Если после логина мы на mainPage, то переходим на recipesPage
	if currentURL == "http://localhost:8080/mainPage" {
		recipeButton, err := waitForElement(wd, selenium.ByXPATH, "//button[contains(text(),'View Recipes')]", 5*time.Second)
		if err != nil {
			t.Fatalf("Кнопка 'View Recipes' не найдена: %v", err)
		}
		if err := recipeButton.Click(); err != nil {
			t.Fatalf("Ошибка клика по кнопке 'View Recipes': %v", err)
		}

		// Ждем редирект
		time.Sleep(3 * time.Second)

		currentURL, err = wd.CurrentURL()
		if err != nil {
			t.Fatalf("Ошибка получения текущего URL после клика: %v", err)
		}
		if currentURL != recipePageURL {
			t.Fatalf("Редирект не сработал, текущий URL: %s", currentURL)
		}
	}
	_, err = wd.FindElement(selenium.ByCSSSelector, `a[onclick="sortrecipes('name', 'asc')"]`)
	if err != nil {
		t.Errorf("Ошибка поиска элемента сортировки по имени: %v", err)
	}
	_, err = wd.FindElement(selenium.ByID, "recipeLevelFilter")
	if err != nil {
		t.Errorf("Ошибка поиска элемента фильтра по типу: %v", err)
	}
	// Проверяем загрузку рецептов
	_, err = waitForElement(wd, selenium.ByCSSSelector, ".recipe-card", 10*time.Second)
	if err != nil {
		t.Fatalf("Ошибка загрузки контейнера рецептов: %v", err)
	}
	fmt.Println("TESTING ETE TEST")
}

// Ожидание элемента с таймаутом
func waitForElement(wd selenium.WebDriver, by string, value string, timeout time.Duration) (selenium.WebElement, error) {
	endTime := time.Now().Add(timeout)
	var element selenium.WebElement
	var err error

	for time.Now().Before(endTime) {
		element, err = wd.FindElement(by, value)
		if err == nil {
			return element, nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil, err
}
