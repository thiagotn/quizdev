-- 002_seed_go_questions.sql
-- Seed: Golang questions for all levels

-- ─── BEGINNER ────────────────────────────────────────────────────────────────

-- Q1: Correct way to declare a variable
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'beginner', 'Qual trecho declara e inicializa corretamente uma variável inteira em Go?',
            'Em Go, `:=` é o operador de declaração curta. Ele infere o tipo automaticamente. `var x int = 10` também é válido, mas `:=` é idiomático dentro de funções.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('x := 10', true, 1),
    ('int x = 10', false, 2),
    ('var x = int(10)', false, 3),
    ('x = 10', false, 4)
) AS opt(code, correct, ord);

-- Q2: Correct for loop
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'beginner', 'Qual é a forma correta de um loop `for` que imprime números de 0 a 4?',
            'Go usa apenas `for` — não existe `while`. A sintaxe `for i := 0; i < 5; i++` é equivalente ao `for` clássico de C/Java.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('for i := 0; i < 5; i++ {
    fmt.Println(i)
}', true, 1),
    ('for (i := 0; i < 5; i++) {
    fmt.Println(i)
}', false, 2),
    ('while i < 5 {
    fmt.Println(i)
    i++
}', false, 3),
    ('loop i from 0 to 5 {
    fmt.Println(i)
}', false, 4)
) AS opt(code, correct, ord);

-- Q3: fmt.Println
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'beginner', 'Como imprimir "Hello, World!" no terminal em Go?',
            '`fmt.Println` imprime uma linha com quebra de linha automática. O pacote `fmt` precisa ser importado.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('fmt.Println("Hello, World!")', true, 1),
    ('print("Hello, World!")', false, 2),
    ('console.log("Hello, World!")', false, 3),
    ('System.out.println("Hello, World!")', false, 4)
) AS opt(code, correct, ord);

-- Q4: Function declaration
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'beginner', 'Qual trecho declara corretamente uma função que recebe dois inteiros e retorna sua soma?',
            'Em Go, o tipo do retorno vem após os parâmetros. A sintaxe é `func nome(params) retorno { ... }`.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('func sum(a int, b int) int {
    return a + b
}', true, 1),
    ('function sum(a, b int) int {
    return a + b
}', false, 2),
    ('func sum(a, b) int {
    return a + b
}', false, 3),
    ('int sum(int a, int b) {
    return a + b
}', false, 4)
) AS opt(code, correct, ord);

-- Q5: Slice vs Array
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'beginner', 'Qual trecho cria corretamente um slice de strings em Go?',
            '`[]string{}` cria um slice vazio. Arrays têm tamanho fixo: `[3]string{}`. Slices são dinâmicos e muito mais comuns em Go.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('names := []string{"Alice", "Bob"}', true, 1),
    ('names := [string]{"Alice", "Bob"}', false, 2),
    ('names := Array<string>{"Alice", "Bob"}', false, 3),
    ('names := string[]{"Alice", "Bob"}', false, 4)
) AS opt(code, correct, ord);

-- ─── INTERMEDIATE ─────────────────────────────────────────────────────────────

-- Q6: Error handling
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'intermediate', 'Qual é o padrão idiomático de tratamento de erro em Go?',
            'Go não usa exceções. Funções retornam um segundo valor do tipo `error`. A convenção é checar `if err != nil` imediatamente após a chamada.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('result, err := doSomething()
if err != nil {
    log.Fatal(err)
}', true, 1),
    ('try {
    result := doSomething()
} catch (err) {
    log.Fatal(err)
}', false, 2),
    ('result := doSomething()
if result.error != nil {
    log.Fatal(result.error)
}', false, 3),
    ('result, _ := doSomething()', false, 4)
) AS opt(code, correct, ord);

-- Q7: Goroutine
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'intermediate', 'Como iniciar uma goroutine em Go?',
            'A palavra-chave `go` antes de uma chamada de função inicia a execução concorrente dela como uma goroutine, que é leve e gerenciada pelo runtime de Go.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('go processData(data)', true, 1),
    ('goroutine processData(data)', false, 2),
    ('async processData(data)', false, 3),
    ('thread.start(processData, data)', false, 4)
) AS opt(code, correct, ord);

-- Q8: Struct + method
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'intermediate', 'Qual trecho define corretamente um método em uma struct em Go?',
            'Métodos em Go são funções com um receiver. O receiver aparece entre `func` e o nome do método: `func (r ReceiverType) MethodName()`.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('type Dog struct { Name string }

func (d Dog) Bark() string {
    return d.Name + " says: Woof!"
}', true, 1),
    ('type Dog struct { Name string }

func Dog.Bark() string {
    return Dog.Name + " says: Woof!"
}', false, 2),
    ('type Dog struct { Name string }

method Bark(d Dog) string {
    return d.Name + " says: Woof!"
}', false, 3),
    ('class Dog {
    Name string
    func Bark() string {
        return this.Name + " says: Woof!"
    }
}', false, 4)
) AS opt(code, correct, ord);

-- Q9: Interface
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'intermediate', 'Como uma struct implementa uma interface em Go?',
            'Go usa implementação implícita de interfaces. Se uma struct tem todos os métodos definidos pela interface, ela a implementa automaticamente — sem `implements` ou declaração explícita.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('type Speaker interface { Speak() string }

type Human struct{}

func (h Human) Speak() string {
    return "Hello!"
}', true, 1),
    ('type Speaker interface { Speak() string }

type Human struct{} implements Speaker {
    func Speak() string { return "Hello!" }
}', false, 2),
    ('type Human struct{}

func (h Human) implements Speaker() string {
    return "Hello!"
}', false, 3),
    ('type Human struct : Speaker {
    Speak() string { return "Hello!" }
}', false, 4)
) AS opt(code, correct, ord);

-- Q10: Channel
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'intermediate', 'Qual trecho cria um channel e envia um valor corretamente?',
            'Channels são criados com `make(chan tipo)`. O operador `<-` é usado tanto para enviar (`ch <- valor`) quanto para receber (`valor := <-ch`).')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('ch := make(chan int)
go func() { ch <- 42 }()
val := <-ch', true, 1),
    ('ch := chan int
ch.send(42)
val := ch.receive()', false, 2),
    ('ch := new(chan int)
ch <- 42
val := ch', false, 3),
    ('ch := channel(int)
send(ch, 42)
val := recv(ch)', false, 4)
) AS opt(code, correct, ord);

-- ─── ADVANCED ─────────────────────────────────────────────────────────────────

-- Q11: Context cancellation
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'advanced', 'Qual trecho cria corretamente um contexto com cancelamento?',
            '`context.WithCancel` retorna um contexto filho e uma função `cancel`. É essencial chamar `defer cancel()` para liberar recursos quando o contexto não for mais necessário.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('ctx, cancel := context.WithCancel(context.Background())
defer cancel()
doWork(ctx)', true, 1),
    ('ctx := context.New()
ctx.Cancel()
doWork(ctx)', false, 2),
    ('ctx := context.Background()
ctx.WithCancel()
doWork(ctx)', false, 3),
    ('ctx, cancel := context.Timeout(5 * time.Second)
cancel()
doWork(ctx)', false, 4)
) AS opt(code, correct, ord);

-- Q12: Defer order
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'advanced', 'Qual é a saída do código abaixo?\n\nfunc main() {\n    defer fmt.Println("A")\n    defer fmt.Println("B")\n    defer fmt.Println("C")\n}',
            'Defers em Go são executados em ordem LIFO (Last In, First Out) — o último `defer` registrado é o primeiro a executar. Isso é útil para cleanup em ordem inversa de aquisição de recursos.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('C
B
A', true, 1),
    ('A
B
C', false, 2),
    ('B
A
C', false, 3),
    ('Erro de compilação', false, 4)
) AS opt(code, correct, ord);

-- Q13: Generics
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'advanced', 'Qual trecho usa corretamente Generics (Go 1.18+) para uma função que retorna o maior de dois valores?',
            'Go 1.18 introduziu generics com parâmetros de tipo em `[]`. A constraint `cmp.Ordered` (ou `constraints.Ordered`) garante que o tipo suporta operadores de comparação.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}', true, 1),
    ('func Max<T>(a, b T) T {
    if a > b { return a }
    return b
}', false, 2),
    ('func Max(T any)(a, b T) T {
    if a > b { return a }
    return b
}', false, 3),
    ('generic func Max(a, b any) any {
    if a > b { return a }
    return b
}', false, 4)
) AS opt(code, correct, ord);

-- Q14: Mutex
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'advanced', 'Qual trecho protege corretamente uma variável compartilhada com Mutex?',
            '`sync.Mutex` é usada com `Lock()` antes de acessar o recurso e `Unlock()` depois. `defer mu.Unlock()` garante o desbloqueio mesmo em caso de panic.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('var mu sync.Mutex
var count int

func increment() {
    mu.Lock()
    defer mu.Unlock()
    count++
}', true, 1),
    ('var count int

func increment() {
    synchronized {
        count++
    }
}', false, 2),
    ('var mu sync.Mutex
var count int

func increment() {
    count++ // mutex é automático
}', false, 3),
    ('var mu sync.RWMutex
var count int

func increment() {
    mu.RLock()
    defer mu.RUnlock()
    count++
}', false, 4)
) AS opt(code, correct, ord);

-- Q15: Select with channels
WITH q AS (
    INSERT INTO questions (language, level, title, explanation)
    VALUES ('go', 'advanced', 'O que o `select` faz no contexto de channels em Go?',
            '`select` bloqueia até que um dos cases esteja pronto. Se múltiplos estiverem prontos simultaneamente, um é escolhido aleatoriamente. É o mecanismo central para multiplexação de channels.')
    RETURNING id
)
INSERT INTO question_options (question_id, code_snippet, is_correct, display_order)
SELECT q.id, opt.code, opt.correct, opt.ord FROM q,
(VALUES
    ('select {
case msg := <-ch1:
    fmt.Println("ch1:", msg)
case msg := <-ch2:
    fmt.Println("ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
}', true, 1),
    ('switch ch1 | ch2 {
case ch1:
    fmt.Println("ch1")
case ch2:
    fmt.Println("ch2")
}', false, 2),
    ('listen(ch1, ch2) {
case ch1 -> msg:
    fmt.Println(msg)
}', false, 3),
    ('select ch1, ch2 {
    on ch1: fmt.Println("ch1")
    on ch2: fmt.Println("ch2")
}', false, 4)
) AS opt(code, correct, ord);
