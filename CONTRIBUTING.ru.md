# Участие в NYXORA

Прежде всего, спасибо, что рассматриваете возможность участия в NYXORA! Мы приветствуем участие всех.

## Кодекс поведения

Этот проект следует [Кодексу поведения](CODE_OF_CONDUCT.md). Участвуя, вы должны соблюдать этот кодекс.

## Как я могу внести свой вклад?

### Сообщение об ошибке

Перед отправкой отчёта об ошибке:
- Проверьте [issues](https://github.com/nyxora-network/nyxora/issues), не сообщалось ли уже об этой проблеме
- Соберите информацию: версия ОС, версия Go, шаги для воспроизведения, вывод ошибки

**Отправьте отчёт об ошибке**, открыв [новый issue](https://github.com/nyxora-network/nyxora/issues/new?template=bug_report.md).

### Предложение функции

Откройте [запрос функции](https://github.com/nyxora-network/nyxora/issues/new?template=feature_request.md) с описанием:
- Проблема, которую вы решаете
- Как вы представляете решение
- Рассмотренные альтернативы

### Добавление нового транспорта

1. Создайте `internal/transport/<name>.go`, реализующий интерфейс `Transport`
2. Зарегистрируйте его в `internal/transport/registry.go`
3. Создайте `tunnels/<name>/` со скриптами установки и манифестом
4. Добавьте веса оценки в `internal/transport/scoring.go`
5. Напишите тесты и запустите `make test`

### Улучшение TUI

Интерактивный TUI находится в `internal/interactive/` и использует:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — фреймворк TUI
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — стилизация
- [Bubbles](https://github.com/charmbracelet/bubbles) — компоненты (textinput, spinner, progress)

Цвета темы используют значения Catppuccin TrueColor в файле `theme.go`.

## Процесс Pull Request

1. Fork репозиторий и создайте свою ветку из `main`
2. Запустите `make test` и `make vet` — оба должны пройти
3. Добавьте тесты для новой функциональности
4. Обновите документацию при необходимости
5. Убедитесь, что ваш код следует существующим соглашениям
6. Отправьте PR с понятным описанием

### Стиль коммитов

Используйте формат conventional commits:
- `feat:` — новая функция
- `fix:` — исправление ошибки
- `refactor:` — изменение кода без исправления/новой функции
- `docs:` — только документация
- `test:` — добавление/исправление тестов
- `style:` — форматирование, изменения стиля
- `chore:` — обслуживание, зависимости

## Настройка среды разработки

```bash
# Fork и clone
git clone https://github.com/YOUR_USERNAME/nyxora.git
cd nyxora

# Добавить upstream remote
git remote add upstream https://github.com/nyxora-network/nyxora.git

# Создать ветку для функции
git checkout -b feat/your-feature

# Внести изменения, затем:
make test
make vet
make build

# Коммит и push
git commit -m "feat: add your feature"
git push origin feat/your-feature
```

## Есть вопросы?

Откройте [обсуждение](https://github.com/nyxora-network/nyxora/discussions) или присоединяйтесь к [Telegram](https://t.me/NyxoraCore).
