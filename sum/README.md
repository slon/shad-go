# sum

В этой задаче вам нужно научиться сдавать решения в тестовую систему.

0. **(Один раз)** Зарегистрируйтесь в тестовой системе

   Если вы еще не зарегистрировались в тестовой системе, сделайте это сейчас.
   Система создаст для вас личный репозиторий.
   Перейдите в него по ссылке My Repo на https://go.manytask.org/

1. **(Один раз)** [Настройте](https://gitlab.manytask.org/-/profile/keys) ssh ключ. Если вы не знаете как это сделать,
   воспользуйтесь инструкцией на странице по ссылке.

2. **(Один раз)** Склонируйте ваш личный репозиторий

   ```shell
   # Нажмите на синюю кнопку clone и скопируйте адрес и "Clone with SSH"
   # Выполните в консоли команду, заменив последний аргумент на ваш адрес
   git clone git@gitlab.manytask.org:go/students-2024-spring/USERNAME.git .
   ```

3. Откройте файл `sum.go` и реализуйте функцию сложения двух чисел.

4. Проверьте, что ваше решение проходит тесты локально.

   ```shell
   # Из корня репозитория.
   go test ./sum/...
   ```

5. Проверьте, что код проходит линтер. Линтер нужно установить [по инструкции](https://github.com/golangci/golangci-lint#binary).

   ```shell
   # Из корня репозитория.
   golangci-lint run ./sum/...
   ```

6. Добавьте ваши изменения в гит и сделайте коммит.

   ```shell
   git add .
   git commit -m "Solved sum"
   ```

7. Сделайте пуш.

   ```shell
   git push
   ```
   
   **NOTE:** Система тестирует только те задачи, которые изменялись в последнем коммите. Если вы
   сделаете несколько коммитов подряд, и затем один пуш, то протестирован будет только последний коммит.
   Если вы хотите перезапустить тестирование в коммите, вы можете нажать на кнопку Retry на странице
   с логом тестирования, или сделать новый коммит с незначительными изменениями и запушить его.

8. Посмотрите как проходит тестирование, пройдя по ссылке My Submits со страницы https://go.manytask.org/

9. Убедитесь, что ваша оценка появилась в [таблице](https://docs.google.com/spreadsheets/d/1j4s6QLTjm-bUJplz0R2hOlhWipRBE9MOZYJlEw1iFbk).

### Примечание

Мы периодически вносим разные изменения в тесты и readme.

Чтобы ваш репозиторий был синхронизирован с публичным, предлагаем каждый раз, когда вы садитесь за задачи, пуллить публичный репозиторий.

0) Проверьте, привязан ли у вас `upstream` репозиторий как `remote`:
   ```shell
   git remote -v
   # origin  git@gitlab.manytask.org:go/students-2024-spring/USERNAME.git (fetch)
   # origin  git@gitlab.manytask.org:go/students-2024-spring/USERNAME.git (push)
   # upstream        git@gitlab.manytask.org:go/public-2024-spring.git (fetch)
   # upstream        git@gitlab.manytask.org:go/public-2024-spring.git (push)
   ```
1) Если upstream не привязан, добавьте его:

   ```shell
   git remote add upstream git@gitlab.manytask.org:go/public-2024-spring.git
   ```
2) Получите изменения из upstream:

   ```shell
   git fetch upstream
   ```
3) Переключитесь на main и выполните rebase:
   ```shell
   git checkout main
   git rebase upstream/main
   ```
4) Отправьте изменения на ваш форк (origin):
   ```shell
   git push origin main --force-with-lease
   ```
Это обновит ваш форк на GitLab последними изменениями из оригинального репозитория, сохраняя предыдущие изменения.