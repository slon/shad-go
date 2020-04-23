# Программа курса

1. Введение. Программа курса. Отчётность по курсу, критерии
   оценки. Философия дизайна. if, switch, for. Hello, world. Command
   line arguments. Word count. Animated gif. Fetching URL. Fetching
   URL concurrently. Web server. Tour of go. Local IDE
   setup. Submitting solutions to automated
   grading. gofmt. goimports. linting. Submitting PR's with bug fixes.

2. Базовые конструкции языка. names, declarations, variables,
   assignments. type declarations. packages and files. scope. Zero
   value. Выделение памяти. Стек vs куча. Basic data
   types. Constants. Composite data types. Arrays. Slices. Maps. Structs.
   JSON. text/template. string и []byte. Работа с unicode. Unicode
   replacement character.
   Функции. Функции с переменным числом аргументов. Анонимные функции. Ошибки.

3. Методы. Value receiver vs pointer receiver. Embedding. Method
   value. Encapsulation. Интерфейсы. Интерфейсы как
   контракты. io.Writer, io.Reader и их
   реализации. sort.Interface. error. http.Handler. Интерфейсы как
   перечисления. Type assertion. Type switch. The bigger the
   interface, the weaker the abstraction. Обработка ошибок. panic,
   defer, recover. errors.{Unwrap,Is,As}. fmt.Errorf. %w.

4. Горутины и каналы. clock server. echo server. Размер
   канала. Блокирующее и неблокирующее чтение. select
   statement. Channel axioms. `time.After`. `time.NewTicker`. Pipeline
   pattern. Cancellation. Parallel loop. sync.WaitGroup. Обработка
   ошибок в параллельном коде. errgroup.Group. Concurrent web
   crawler. Concurrent directory traversal.

5. Продвинутое тестирование. Subtests. *testing.B. (*T).Logf. (*T).Skipf. (*T).FailNow.
   testing.Short(), testing flags. Генерация моков. testify/{require,assert}. testify/suite. Test fixture.
   Интеграционные тесты. Goroutine leak detector. TestingMain. Coverage. Сравнение бенчмарков.

6. Concurrency with shared memory. sync.Mutex. sync.RWMutex. sync.Cond. atomic. sync.Once.
   Race detector. Async cache. Работа с базой данных. database/sql. sqlx.

7. Package context. Passing request-scoped data. http middleware. chi.Router. Request cancellation.
   Advanced concurrency patterns. Async cache. Graceful server shutdown. context.WithTimeout.
   Batching and cancellation.

8. database/sql, sqlx, работа с базами данных, redis.

9. Reflection. reflect.Type and reflect.Value. struct tags. net/rpc. encoding/gob.
   sync.Map. reflect.DeepEqual.

10. Пакет io, реализации Reader и Writer из стандартной библиотеки.
    Low-level programming. unsafe. Package binary. bytes.Buffer. cgo, syscall.

11. Архитектура GC. Write barrier. Stack growth. GC pause. GOGC. sync.Pool. Шедулер
    горутин. GOMACPROCS. Утечка тредов.

12. Go tooling. pprof. CPU and Memory profiling. Кросс-компиляция. GOOS, GOARCH. CGO_ENABLED=0.
    Build tags. go modules. godoc. x/analysis. Code generation.

13. Полезные библиотеки. CLI applications with cobra. Protobuf and
    GRPC. zap logging.

14. Запасная леция #1. Работа с crypto/* и x/crypto.

15. Запасная леция #2.
