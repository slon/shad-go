## dockertest

Это не настоящая задача, а заготовка на будущее.

### Что нужно сделать?

Установить docker и добиться успешного **локального** запуска тестов
```
go test -v ./dockertest/... -count=1
```

Только **после того, как тесты пройдут локально** можете запушить решение в систему.

### С чего начать?

<details>
<summary><b> Дополнительные шаги для Windows 10</b></summary>

<br/>

1. Установить WSL2 по [инструкции от Microsoft](https://docs.microsoft.com/en-us/windows/wsl/install-win10). <br/> 
**Важно:** нужна именно вторая версия - **WSL2**, проверьте, что она совместима с вашей системой. <br/>
Если шаг 5 не работает, включите опцию `Windows Hypervisor Platform` (Settings -> Apps -> Apps & features -> Optional features -> More Windows Features -> включить чекбокс Windows Hypervisor Platform).

2. Установите Docker Desktop по [инструкции](https://docs.docker.com/docker-for-windows/wsl/#download). Вероятно, выполнять шаги по активации поддержки WSL не потребуется, все подключится автоматически. 

3. Запустите Docker Desktop (никакие контейнеры запускать не надо, только если хотите убедиться в том, что все работает). <br/>
Запустите установленную в п.1 Linux OS через WSL2. <br/>
Далее используйте этот Linux для выполнения дальнейших шагов этого README.

_Замечание_: запущенный, но уже не использующийся Docker Desktop с бекэндом WSL занимает впустую много оперативной памяти, см [issue](https://github.com/microsoft/WSL/issues/4166) - можно ограничить максимальный доступный ему объем (см. [workaround](https://github.com/microsoft/WSL/issues/4166#issuecomment-526725261)), либо отключить автозапуск Docker Desktop и останавливать сервис, когда вы его не используете. 

</details>

#### Установить docker

https://docs.docker.com/engine/install/

После стандартной процедуры установки на Linux будет создана группа `docker`.
Чтобы использовать docker cli без sudo, нужно добавить себя в эту группу:
```
sudo groupadd docker
sudo usermod -aG docker $USER
```
После этого разлогиньтесь из os и залогиньтесь заново (или перезапустите систему).

Для проверки можно запустить
```
docker run hello-world
```

#### Установить docker-compose

https://docs.docker.com/compose/install/

#### Запустить контейнеры не через тесты

В директории `dockertest` выполнить
```
docker-compose up
```

### Что делать, если сразу не заработало?

Поискать решение проблемы в интернете.

Если решение найдено, и проблема выглядит общей, сделать merge request с улучшением README.

Если интернет не помог, спросить в чате.

### docker-compose cheat sheet

Запустить все контейнеры в daemon режиме пересобрав образы:
```
docker-compose up -d --build
```

Остановить все контейнеры:
```
docker-compose down
```

### Docker cheat sheet

Получить список образов
```
docker images
```

Список всех контейнеров:
```
docker ps -a
```

Остановить контейнер:
```
docker stop <NAME>
```

Удалить контейнер:
```
docker rm <NAME>
```

Удалить образ:
```
docker rmi <NAME>
```
