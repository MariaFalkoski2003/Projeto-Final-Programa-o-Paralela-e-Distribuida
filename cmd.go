package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	RESET    = "\033[0m"
	VERDE    = "\033[32m"
	AZUL     = "\033[34m"
	VERMELHO = "\033[31m"
	AMARELO  = "\033[33m"
	CYAN     = "\033[36m"
	BRANCO   = "\033[37m"
)

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func limparTela() {
	fmt.Print("\033[2J\033[H")
}

type DadosCompartilhados struct {
	mu            sync.Mutex
	arrayOriginal []int
	arrayBubble   []int
	arrayCounting []int
	progressoBubble   int
	progressoCounting int
	passoBubble       int
	passoCounting     int
	comparandoBubble1 int
	comparandoBubble2 int
	trocandoBubble    bool
	etapaCounting      string
	indiceCounting     int
	valorAtualCounting int
	frequenciaCounting []int
	minCounting        int
	maxCounting        int
	bubbleFinalizado   bool
	countingFinalizado bool
	tempoBubble        time.Duration
	tempoCounting      time.Duration
}

// FUNÇÃO SIMPLES PARA DESENHAR BARRAS HORIZONTAIS
func desenharBarrasHorizontais(arr []int, cor string) {
	if len(arr) == 0 {
		fmt.Println("  Nenhum dado ainda")
		return
	}
	
	// Encontra o maior valor para escala
	maxValor := 0
	for _, v := range arr {
		if v > maxValor {
			maxValor = v
		}
	}
	if maxValor == 0 {
		maxValor = 1
	}
	
	// Largura máxima da barra (30 caracteres)
	larguraMax := 30
	
	// Desenha cada barra
	for i, valor := range arr {
		// Calcula tamanho da barra proporcional
		tamanho := (valor * larguraMax) / maxValor
		if tamanho < 1 && valor > 0 {
			tamanho = 1
		}
		
		// Mostra índice e barra
		fmt.Printf("  [%2d] ", i)
		for j := 0; j < tamanho; j++ {
			fmt.Print(cor + "█" + RESET)
		}
		// Mostra o valor
		fmt.Printf(" %d\n", valor)
	}
}

func desenharInfoCounting(dados *DadosCompartilhados) {
	switch dados.etapaCounting {
	case "contagem":
		fmt.Printf("Etapa: CONTAGEM\n")
		fmt.Printf("Elemento: %d\n", dados.valorAtualCounting)
		
		fmt.Printf("\nFrequencias:\n")
		for i := 0; i < len(dados.frequenciaCounting) && i < 10; i++ {
			valor := dados.minCounting + i
			if dados.frequenciaCounting[i] > 0 {
				// Desenha barra da frequência
				fmt.Printf("  %2d: ", valor)
				for j := 0; j < dados.frequenciaCounting[i]; j++ {
					fmt.Print(AZUL + "#" + RESET)
				}
				fmt.Printf(" (%d)\n", dados.frequenciaCounting[i])
			}
		}
		
	case "acumulacao":
		fmt.Printf("Etapa: ACUMULACAO\n")
		fmt.Printf("Posicao: %d\n", dados.indiceCounting)
		
		fmt.Printf("\nPosicoes:\n")
		for i := 0; i < len(dados.frequenciaCounting) && i < 10; i++ {
			valor := dados.minCounting + i
			if dados.frequenciaCounting[i] > 0 {
				fmt.Printf("  %2d: %d\n", valor, dados.frequenciaCounting[i])
			}
		}
		
	case "posicionamento":
		fmt.Printf("Etapa: POSICIONAMENTO\n")
		fmt.Printf("Colocando %d na posicao %d\n", dados.valorAtualCounting, dados.indiceCounting)
		fmt.Printf("Array: %v\n", dados.arrayCounting)
	}
}

func desenharComparacaoLadoALado(dados *DadosCompartilhados) {
	limparTela()
	
	fmt.Printf("%s\n", repeat("=", 70))
	fmt.Printf("EXECUCAO CONCORRENTE")
	fmt.Printf("%s\n", repeat("=", 70))
	
	// ===== BUBBLE SORT =====
	fmt.Printf("\n%s▶ BUBBLE SORT (Gráfico de Barras)%s\n", VERDE, RESET)
	
	if dados.bubbleFinalizado {
		fmt.Printf("  ✓ FINALIZADO - Tempo: %v\n", dados.tempoBubble)
	} else {
		fmt.Printf("  Passo: %d | Trocas: %d\n", dados.passoBubble, dados.progressoBubble)
	}
	
	if len(dados.arrayBubble) > 0 && dados.arrayBubble[0] != 0 {
		desenharBarrasHorizontais(dados.arrayBubble, VERDE)
	} else {
		// Mostra o array original enquanto não tem dados
		desenharBarrasHorizontais(dados.arrayOriginal, VERDE)
	}
	
	// ===== COUNTING SORT =====
	fmt.Printf("\n%s▶ COUNTING SORT (Passo a Passo)%s\n", AZUL, RESET)
	
	if dados.countingFinalizado {
    fmt.Printf("  ✓ FINALIZADO - Tempo: %v\n", dados.tempoCounting)
    desenharBarrasHorizontais(dados.arrayCounting, AZUL)
} else {
    desenharInfoCounting(dados)

    fmt.Printf("\nVisualização parcial:\n")
    if len(dados.arrayCounting) > 0 {
        desenharBarrasHorizontais(dados.arrayCounting, AZUL)
    } else {
        desenharBarrasHorizontais(dados.arrayOriginal, AZUL)
    }
}
	
	// ===== INFORMAÇÕES =====
	fmt.Printf("\n%s━ ARRAY ORIGINAL ━%s\n", CYAN, RESET)
	desenharBarrasHorizontais(dados.arrayOriginal, BRANCO)
	
	fmt.Printf("\n%s━ TEMPOS ━%s\n", AMARELO, RESET)
	if dados.bubbleFinalizado {
		fmt.Printf("  Bubble Sort:   %v\n", dados.tempoBubble)
	} else {
		fmt.Printf("  Bubble Sort:   (executando...)\n")
	}
	if dados.countingFinalizado {
		fmt.Printf("  Counting Sort: %v\n", dados.tempoCounting)
	} else {
		fmt.Printf("  Counting Sort: (executando...)\n")
	}
}

func bubbleSortConcorrente(dados *DadosCompartilhados, wg *sync.WaitGroup, delay time.Duration) {
	defer wg.Done()
	
	start := time.Now()
	
	n := len(dados.arrayOriginal)
	resultado := make([]int, n)
	copy(resultado, dados.arrayOriginal)
	
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			dados.mu.Lock()
			dados.comparandoBubble1 = j
			dados.comparandoBubble2 = j + 1
			dados.passoBubble = i + 1
			dados.arrayBubble = resultado
			dados.mu.Unlock()
			
			if resultado[j] > resultado[j+1] {
				resultado[j], resultado[j+1] = resultado[j+1], resultado[j]
				
				dados.mu.Lock()
				dados.trocandoBubble = true
				dados.arrayBubble = resultado
				dados.mu.Unlock()
			}
			
			dados.mu.Lock()
			dados.progressoBubble++
			dados.mu.Unlock()
			
			time.Sleep(delay)
		}
		
		dados.mu.Lock()
		dados.arrayBubble = resultado
		dados.mu.Unlock()
	}
	
	dados.mu.Lock()
	dados.arrayBubble = resultado
	dados.tempoBubble = time.Since(start)
	dados.bubbleFinalizado = true
	dados.mu.Unlock()
	
	fmt.Printf("\n✓ Bubble Sort FINALIZADO! Tempo: %v\n", dados.tempoBubble)
}

func countingSortConcorrente(dados *DadosCompartilhados, wg *sync.WaitGroup, delay time.Duration) {
	defer wg.Done()
	
	start := time.Now()
	
	if len(dados.arrayOriginal) == 0 {
		return
	}
	
	min := dados.arrayOriginal[0]
	max := dados.arrayOriginal[0]
	for _, v := range dados.arrayOriginal {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		time.Sleep(delay / 2)
	}
	
	dados.mu.Lock()
	dados.minCounting = min
	dados.maxCounting = max
	dados.mu.Unlock()
	
	intervalo := max - min + 1
	count := make([]int, intervalo)
	
	// PASSO 1: CONTAGEM
	dados.mu.Lock()
	dados.etapaCounting = "contagem"
	dados.frequenciaCounting = count
	dados.mu.Unlock()
	
	for _, v := range dados.arrayOriginal {
		count[v-min]++
		
		dados.mu.Lock()
		dados.progressoCounting++
		dados.valorAtualCounting = v
		dados.frequenciaCounting = count
		dados.mu.Unlock()
		
		time.Sleep(delay)
	}
	
	// PASSO 2: ACUMULACAO
	dados.mu.Lock()
	dados.etapaCounting = "acumulacao"
	dados.mu.Unlock()
	
	for i := 1; i < intervalo; i++ {
		count[i] += count[i-1]
		
		dados.mu.Lock()
		dados.progressoCounting++
		dados.indiceCounting = i
		dados.frequenciaCounting = count
		dados.mu.Unlock()
		
		time.Sleep(delay)
	}
	
	// PASSO 3: POSICIONAMENTO
	dados.mu.Lock()
	dados.etapaCounting = "posicionamento"
	dados.mu.Unlock()
	
	resultado := make([]int, len(dados.arrayOriginal))
	for i := len(dados.arrayOriginal) - 1; i >= 0; i-- {
		valor := dados.arrayOriginal[i]
		posicao := count[valor-min] - 1
		resultado[posicao] = valor
		count[valor-min]--
		
		dados.mu.Lock()
		dados.progressoCounting++
		dados.valorAtualCounting = valor
		dados.indiceCounting = posicao
		dados.arrayCounting = resultado
		dados.mu.Unlock()
		
		time.Sleep(delay)
	}
	
	dados.mu.Lock()
	dados.arrayCounting = resultado
	dados.tempoCounting = time.Since(start)
	dados.countingFinalizado = true
	dados.mu.Unlock()
	
	fmt.Printf("\n✓ Counting Sort FINALIZADO! Tempo: %v\n", dados.tempoCounting)
}

func visualizacaoLadoALado(dados *DadosCompartilhados, wg *sync.WaitGroup, delay time.Duration) {
	defer wg.Done()
	
	for {
		dados.mu.Lock()
		bubbleFinalizado := dados.bubbleFinalizado
		countingFinalizado := dados.countingFinalizado
		dados.mu.Unlock()
		
		if bubbleFinalizado && countingFinalizado {
			break
		}
		
		dados.mu.Lock()
		desenharComparacaoLadoALado(dados)
		dados.mu.Unlock()
		
		time.Sleep(delay)
	}
}

func demonstrarRaceCondition() {
	fmt.Printf("\n%s━ DEMONSTRANDO RACE CONDITION ━%s\n", AMARELO, RESET)
	
	// SEM MUTEX - ERRADO
	var contadorSemMutex int
	var wg sync.WaitGroup
	
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			contadorSemMutex++ // SEM MUTEX
		}()
	}
	wg.Wait()
	
	// COM MUTEX - CORRETO
	var contadorComMutex int
	var mutex sync.Mutex
	
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()
			contadorComMutex++
			mutex.Unlock()
		}()
	}
	wg.Wait()
	
	fmt.Printf("  Sem Mutex: %d (esperado 1000) %s✗ ERRADO%s\n", contadorSemMutex, VERMELHO, RESET)
	fmt.Printf("  Com Mutex: %d (esperado 1000) %s✓ CORRETO%s\n", contadorComMutex, VERDE, RESET)
}

func desenharComparacaoFinal(original, bubble, counting []int, tempoBubble, tempoCounting time.Duration) {
	fmt.Printf("\n%s\n", repeat("=", 70))
	fmt.Printf("RESULTADOS FINAIS")
	fmt.Printf("\n%s\n", repeat("=", 70))
	
	// Mostra os arrays
	fmt.Printf("\n%s▶ BUBBLE SORT%s\n", VERDE, RESET)
	desenharBarrasHorizontais(bubble, VERDE)
	
	fmt.Printf("\n%s▶ COUNTING SORT%s\n", AZUL, RESET)
	desenharBarrasHorizontais(counting, AZUL)
	
	fmt.Printf("\n%s▶ ORIGINAL%s\n", CYAN, RESET)
	desenharBarrasHorizontais(original, BRANCO)
	
	fmt.Printf("\n%s━ TEMPOS ━%s\n", AMARELO, RESET)
	fmt.Printf("  Bubble Sort:   %v\n", tempoBubble)
	fmt.Printf("  Counting Sort: %v\n", tempoCounting)
	
	// Comparação
	if fmt.Sprint(bubble) == fmt.Sprint(counting) {
		fmt.Printf("\n%s✓ RESULTADO: Ordenações IGUAIS!%s\n", VERDE, RESET)
	} else {
		fmt.Printf("\n%s✗ RESULTADO: Ordenações DIFERENTES!%s\n", VERMELHO, RESET)
	}
}

func main() {
	limparTela()
	
	fmt.Printf("%s\n", repeat("=", 70))
	fmt.Printf("SORTS CONCORRENTES - GOROUTINES E MUTEX")
	fmt.Printf("%s\n\n", repeat("=", 70))
	
	var tamanho int
	var delay int
	
	fmt.Print("Digite o tamanho do array (5-10): ")
	fmt.Scan(&tamanho)
	
	fmt.Print("Digite o delay em ms (300-500): ")
	fmt.Scan(&delay)
	delayDuration := time.Duration(delay) * time.Millisecond
	
	rand.Seed(time.Now().UnixNano())
	arrOriginal := make([]int, tamanho)
	for i := 0; i < tamanho; i++ {
		arrOriginal[i] = rand.Intn(20) + 1
	}
	
	fmt.Printf("\nArray original: %v\n", arrOriginal)
	fmt.Print("\nPressione ENTER para iniciar...")
	fmt.Scanln()
	fmt.Scanln()
	
	dados := &DadosCompartilhados{
		arrayOriginal: arrOriginal,
		arrayBubble:   make([]int, tamanho),
		arrayCounting: append([]int{}, arrOriginal...),
		comparandoBubble1: -1,
		comparandoBubble2: -1,
		etapaCounting:     "contagem",
		frequenciaCounting: make([]int, 0),
		bubbleFinalizado:   false,
		countingFinalizado: false,
	}
	
	var wg sync.WaitGroup
	
	limparTela()
	fmt.Printf("INICIANDO THREADS CONCORRENTES...\n\n")
	
	start := time.Now()
	
	wg.Add(3)
	go bubbleSortConcorrente(dados, &wg, delayDuration)
	go countingSortConcorrente(dados, &wg, delayDuration)
	go visualizacaoLadoALado(dados, &wg, delayDuration/2)
	
	wg.Wait()
	elapsedTotal := time.Since(start)
	
	limparTela()
	
	desenharComparacaoFinal(arrOriginal, dados.arrayBubble, dados.arrayCounting, dados.tempoBubble, dados.tempoCounting)
	
	fmt.Printf("\nTEMPO TOTAL CONCORRENTE: %v\n", elapsedTotal)
	
	demonstrarRaceCondition()
	
	fmt.Printf("\n%sCONCLUIDO!%s\n", CYAN, RESET)
}