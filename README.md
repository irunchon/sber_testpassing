# sber_testtask_1
Тестовое задание для Сбер devices 

## Описание  
Приложение позволяет пройти онлайн тест по заданному адресу. Приложение поддерживает возможность многопоточного прохождения теста - каждый поток (горутина) проходит независимо от остальных свой тест. Соответсвенно, каждый поток представляется веб-ресурсу как независимый клиент, чтобы получить свою копию теста.  
Количество параллельных потоков задаётся в качестве переменной окружения в файле .env, позволяя, таким образом,  конфигурировать количество параллельных потоков без перекомпиляции кода приложения.  

Каждый поток скрипта по умолчанию сообщает об успешном прохождении теста в виде JSON сообщения для логирования.  
### Версия Golang  
Приложение использует Golang версии 1.21.0  
### Используемые библиотеки  
Помимо стандартных библиотек используются следующие модули:
- [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html) - библиотека для парсинга HTML страниц  
- [github.com/stretchr/testify](https://github.com/stretchr/testify) - библиотека для тестирования   
### Настраиваемые параметры  
В файле .env задаются переменные окружения задаются параметры для конфигурации работы приложения.    
Перечень данных параметров и значения по умолчанию:
- Начальная страница онлайн-теста, к которой осуществляется подключение:  
`START_PAGE=http://147.78.65.149/start/`
- Финальная страница теста, сигнализирующая о его успешном завершении:  
`FINAL_PAGE=http://147.78.65.149/passed`
- Количество параллельных потоков (горутин):  
`QTY_OF_THREADS=5`
- Максимальное допустимое количество запросов к серверу в секунду (RPS):  
`MAX_SERVER_RPS=3`
- Уровень логирования:  
`LOG_LEVEL=INFO`
  Другие возможные варианты: Debug, Info, Warn, Error (используются в стандартной библиотеке log/slog go версии 1.21).  
## Запуск приложения  
Для удобства сборки и запуска приложения и тестов заданы соответствующие цели в Makefile.  

Сборка и запуск приложения:  
```bash
$ make
```  
или 
```bash
$ make run  
```  
Запуск тестов:
```bash
$ make test
```  
Информация о покрытии тестами:
```bash
$ make test-coverage
```  
