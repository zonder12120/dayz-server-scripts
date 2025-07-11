# Скрипты для сервера DayZ

## Содержание
- [Краткое описание](#краткое-описание)
- [Настройка](#настройка)
- [Запуск](#запуск)
- [Бэкап](#бэкапы)


## Краткое описание

Скрипты, упрощающие конфигурацию DayZ сервера. Сами скрипты изначально сгенерированы нейросетью, но слегка отредактированы мной.

ВНИМАНИЕ! Для корректной работы скриптов сервер должен был быть успешно запущен хотя бы 1 раз (и полностью проинициализирован)! Для скриптов конфигураций модов (внезапно) нужны сами моды, а так же сервак тоже должен быть успешно запущен с этими модами хотя бы 1 раз!

1. loot-settings - скрипт, который на данный момент только и делает, что контролирует количество ВСЕГО лута одним множителем.

2. expansion-ai-patrol - скрипт для сервера, использующий мод @DayZ-Expansion-AI, который через соответствующие параметры конфигурации изменяет как количество бойцов в патруле, так и в общем количество патрулей и их путевых точек.

## Настройка

Для запуска скриптов необходимо установить язык программирования Go на устройство, на котором будеешь запускать сервак. Люди, которые запускают сервер на Linux справятся и без ссылок, но для домашнего использования на Windows оставляю ссылку здесь: https://go.dev/doc/

Каждый отдельный скрипт - это директория (папка) внутри этого репозитория, каждый скрипт - это файл main.go

Внутри каждой директории скрипта настройки находятся в файле config.yml, есть комментарии к конфигурации, которые помогут настроить всё на свой вкус. 

Файлы конфигурации .yml можно открыть любым блокнотом, менять значения нужно ровно в том же формате, как они указаны в исходной конфигурации.

Если ты умеешь работать с кодом, то можешь открыть файлы main.go, чтобы как-то модифицировать скрипты. Некоторые комментарии к коду есть, они помогут сориентироваться.

## Запуск

Для запуска скрипта нам нужно будет работать с командной строкой.

Люди, работающие на Linux и так разберутся, а для всех, кто работает на Windows ниже напишу шаги:

1. Скачиваем скрипты из репозитория (справа будет зелёная кнопка Code, нажав на которую будет всплывающая менюшка, где нужно тыкнуть на "Download ZIP"), либо клонируй репозиторий, если умеешь
2. Если не клонировали репозиторий, то разархивируем архив в любое удобное место
3. Запускаем PowerShell
4. Копируем из проводника полный путь к скачанным скриптам
5. Пишем в PowerShell следующее
```text
cd сюда_вставляем_скопированный_путь
```
Например, вот так:
```text
cd C:\Users\daniil\GolandProjects\dayz-server-scripts
```
6. Далее, исходя из нужного нам скрипта, запускаем его через команду в PowerShell
```text
go run директория_конкретного_скрипта/main.go
```
Например, вот так:
```text
go run loot-settings/main.go
```

Если сделал всё правильно, то должен получить соответствующую надпись в консоле в конце, где будет текст "🎉 Готово!"

## Бэкапы

По итогу работы скрипта у нас по указанному пути в конфигурации будут сохраняться файлы backup в директорию backups (у файлов будет соответствующий суффикс в названии). Далее, если захотим восстановить исходные параметры, то просто копируем с заменой оттуда интересующий бэкап, убирая лишний суффикс, чтобы название файла соответствовало исходному.
