package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
)

// Environment configuration
type Environment struct {
	Name      string
	Namespace string
	RDSHost   string
}

// Global variables
var (
	environments = map[string]Environment{
		"STAGING": {
			Name:      "STAGING",
			Namespace: "stg-apps",
			RDSHost:   "reservas-stg-postgres.cbwcm8my4qns.sa-east-1.rds.amazonaws.com",
		},
		"PRODUCTION": {
			Name:      "PRODUCTION",
			Namespace: "prd-apps",
			RDSHost:   "reservas-prd-postgres.cbwcm8my4qns.sa-east-1.rds.amazonaws.com",
		},
	}

	selectedEnv    Environment
	localPort      string = "5432"
	podName        string
	portForwardCmd *exec.Cmd
	errorMessage   string
	showError      bool

	// Colors
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	blue   = color.New(color.FgBlue)
	cyan   = color.New(color.FgCyan)
	white  = color.New(color.FgWhite, color.Bold)
	gray   = color.New(color.FgHiBlack)
)

func main() {
	// Setup signal handling for cleanup
	setupSignalHandling()

	// Clear screen and show banner
	clearScreen()
	showBanner()

	// Main flow
	selectEnvironment()
	selectLocalPort()
	showConnectionSummary()
	establishConnection()
}

func setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanupResources(true)
	}()
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func showBanner() {
	green.Println("█████████████████████████████████████████████████████████████████████████")
	green.Println("█▌                                                                     ▐█")
	green.Println("█▌                                                                     ▐█")
	green.Println("█▌   _____ _____  _   _ ___________ _   _   ___   _   _ _____   ___    ▐█")
	green.Println("█▌  |  __ \\  _  || | | |  ___| ___ \\ \\ | | / _ \\ | \\ | /  __ \\ / _ \\   ▐█")
	green.Println("█▌  | |  \\/ | | || | | | |__ | |_/ /  \\| |/ /_\\ \\|  \\| | /  \\// /_\\ \\  ▐█")
	green.Println("█▌  | | __| | | || | | |  __||    /| . ` ||  _  || . ` | |    |  _  |  ▐█")
	green.Println("█▌  | |_\\ \\ \\_/ /\\ \\_/ / |___| |\\ \\| |\\  || | | || |\\  | \\__/\\| | | |  ▐█")
	green.Println("█▌   \\____/\\___/  \\___/\\____/\\_| \\_\\_| \\_/\\_| |_/\\_| \\_/\\____/\\_| |_/  ▐█")
	green.Println("█▌                                                                     ▐█")
	green.Println("█▌                                                                     ▐█")
	green.Println("█████████████████████████████████████████████████████████████████████████")
	color.NoColor = false
	cyan.Println("          PostgreSQL Database Connector v2.0")
	gray.Println("                   Governança Tech")
	fmt.Println()
}

func showGoodbyeBanner() {
	clearScreen()
	red.Println("█████████████████████████████████████████████████████████████████████████")
	red.Println("█▌                                                                     ▐█")
	red.Println("█▌                                                                     ▐█")
	red.Println("█▌   _____ _____  _   _ ___________ _   _   ___   _   _ _____   ___    ▐█")
	red.Println("█▌  |  __ \\  _  || | | |  ___| ___ \\ \\ | | / _ \\ | \\ | /  __ \\ / _ \\   ▐█")
	red.Println("█▌  | |  \\/ | | || | | | |__ | |_/ /  \\| |/ /_\\ \\|  \\| | /  \\// /_\\ \\  ▐█")
	red.Println("█▌  | | __| | | || | | |  __||    /| . ` ||  _  || . ` | |    |  _  |  ▐█")
	red.Println("█▌  | |_\\ \\ \\_/ /\\ \\_/ / |___| |\\ \\| |\\  || | | || |\\  | \\__/\\| | | |  ▐█")
	red.Println("█▌   \\____/\\___/  \\___/\\____/\\_| \\_\\_| \\_/\\_| |_/\\_| \\_/\\____/\\_| |_/  ▐█")
	red.Println("█▌                                                                     ▐█")
	red.Println("█▌                                                                     ▐█")
	red.Println("█████████████████████████████████████████████████████████████████████████")
	color.NoColor = false
	yellow.Println("              Conexão Encerrada")
	gray.Println("                 Recursos Limpos")
	fmt.Println()
}

func printHeader(title string) {
	blue.Println("═══════════════════════════════════════════════════════════════════")
	white.Printf("                          %s\n", title)
	blue.Println("═══════════════════════════════════════════════════════════════════")
	fmt.Println()
}

func printStep(step, description string) {
	cyan.Printf("[%s] %s\n", step, description)
}

func printSuccess(message string) {
	green.Printf("✓ %s\n", message)
}

func printError(message string) {
	red.Printf("✗ %s\n", message)
}

func printWarning(message string) {
	yellow.Printf("⚠ %s\n", message)
}

func printInfo(message string) {
	blue.Printf("ℹ %s\n", message)
}

func selectEnvironment() {
	clearScreen()
	showBanner()
	printHeader("SELEÇÃO DE AMBIENTE")

	white.Println("Selecione o ambiente para conectar ao RDS:")
	fmt.Println()

	envKeys := []string{"STAGING", "PRODUCTION"}

	for i, envKey := range envKeys {
		env := environments[envKey]
		switch envKey {
		case "STAGING":
			yellow.Printf("  [%d] %s ", i+1, env.Name)
			gray.Printf("(namespace: %s)\n", env.Namespace)
			gray.Printf("      RDS: %s\n", env.RDSHost)
		case "PRODUCTION":
			red.Printf("  [%d] %s ", i+1, env.Name)
			gray.Printf("(namespace: %s)\n", env.Namespace)
			gray.Printf("      RDS: %s\n", env.RDSHost)
		}
		fmt.Println()
	}

	gray.Println("  [0] Sair")
	fmt.Println()

	for {
		cyan.Printf("Digite sua escolha [0-%d]: ", len(envKeys))
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		if choice == "0" {
			printWarning("Operação cancelada pelo usuário.")
			os.Exit(0)
		}

		if choiceNum, err := strconv.Atoi(choice); err == nil && choiceNum >= 1 && choiceNum <= len(envKeys) {
			selectedEnv = environments[envKeys[choiceNum-1]]
			printSuccess(fmt.Sprintf("Ambiente selecionado: %s", selectedEnv.Name))
			printInfo(fmt.Sprintf("Namespace: %s", selectedEnv.Namespace))
			printInfo(fmt.Sprintf("RDS Host: %s", selectedEnv.RDSHost))
			fmt.Println()
			time.Sleep(1 * time.Second)
			break
		} else {
			printError(fmt.Sprintf("Opção inválida. Digite um número entre 0 e %d.", len(envKeys)))
		}
	}
}

func selectLocalPort() {
	clearScreen()
	showBanner()
	printHeader("CONFIGURAÇÃO DE PORTA LOCAL")

	white.Println("Configure a porta local para o port-forward:")
	fmt.Println()
	green.Print("Recomendado: 5432 ")
	gray.Println("(porta padrão do PostgreSQL)")
	fmt.Println()

	for {
		cyan.Print("Digite a porta local [5432] ou 0 para sair: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		portInput := strings.TrimSpace(input)

		if portInput == "0" {
			printWarning("Operação cancelada pelo usuário.")
			os.Exit(0)
		}

		if portInput == "" {
			localPort = "5432"
			printSuccess(fmt.Sprintf("Usando porta padrão: %s", localPort))
			break
		}

		if port, err := strconv.Atoi(portInput); err == nil && port >= 1024 && port <= 65535 {
			localPort = portInput
			printSuccess(fmt.Sprintf("Porta configurada: %s", localPort))
			break
		} else {
			printError("Porta inválida. Use um número entre 1024 e 65535, ou 0 para sair.")
		}
	}

	// Check if port is free
	if isPortInUse(localPort) {
		printWarning(fmt.Sprintf("Porta %s já está em uso!", localPort))
		cyan.Print("Deseja continuar mesmo assim? [y/N/0 para sair]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		confirm := strings.TrimSpace(strings.ToLower(input))

		if confirm == "0" {
			printWarning("Operação cancelada pelo usuário.")
			os.Exit(0)
		} else if confirm != "y" && confirm != "yes" {
			printError("Operação cancelada.")
			os.Exit(1)
		}
	}
	fmt.Println()
	time.Sleep(1 * time.Second)
}

func showConnectionSummary() {
	clearScreen()
	showBanner()
	printHeader("RESUMO DA CONEXÃO")

	podName = fmt.Sprintf("bastion-rds-%s-%d", os.Getenv("USER"), time.Now().Unix())
	if os.Getenv("USER") == "" {
		podName = fmt.Sprintf("bastion-rds-user-%d", time.Now().Unix())
	}

	white.Println("Configuração da conexão:")
	fmt.Println()
	cyan.Printf("  Ambiente:       %s\n", selectedEnv.Name)
	cyan.Printf("  Namespace:      %s\n", selectedEnv.Namespace)
	cyan.Printf("  RDS Host:       %s\n", selectedEnv.RDSHost)
	cyan.Printf("  Porta Local:    %s\n", localPort)
	cyan.Printf("  Pod Name:       %s\n", podName)
	fmt.Println()

	yellow.Print("Confirma a criação da conexão? [Y/n/0 para sair]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	confirm := strings.TrimSpace(strings.ToLower(input))

	if confirm == "0" {
		printWarning("Operação cancelada pelo usuário.")
		os.Exit(0)
	} else if confirm == "n" || confirm == "no" {
		printError("Operação cancelada pelo usuário.")
		os.Exit(0)
	}
	fmt.Println()
}

func establishConnection() {
	printHeader("ESTABELECENDO CONEXÃO")

	checkPrerequisites()
	createBastionPod()
	waitPodReady()
	testConnectivity()
	startPortForward()
}

func checkPrerequisites() {
	printStep("CHECK", "Verificando pré-requisitos...")

	// Check kubectl
	if !commandExists("kubectl") {
		errorExit("kubectl não encontrado. Instale o kubectl e tente novamente.")
	}

	// Check cluster connectivity
	cmd := exec.Command("kubectl", "cluster-info")
	if err := cmd.Run(); err != nil {
		errorExit("Não foi possível conectar ao cluster Kubernetes. Verifique suas credenciais e conectividade.")
	}

	// Check namespace
	cmd = exec.Command("kubectl", "get", "namespace", selectedEnv.Namespace)
	if err := cmd.Run(); err != nil {
		errorExit(fmt.Sprintf("Namespace '%s' não encontrado no cluster. Verifique se o namespace existe.", selectedEnv.Namespace))
	}

	printSuccess("Pré-requisitos verificados")
}

func createBastionPod() {
	printStep("CREATE", "Criando bastion pod...")

	manifest := fmt.Sprintf(`
apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
  labels:
    app: bastion-rds
    created-by: %s
    environment: %s
spec:
  containers:
  - name: socat-proxy
    image: alpine/socat
    command: 
    - sh
    - -c
    - |
      echo "Iniciando proxy para %s:5432"
      echo "Escutando na porta 5432"
      socat TCP-LISTEN:5432,fork,reuseaddr TCP:%s:5432
    ports:
    - containerPort: 5432
    resources:
      requests:
        cpu: 10m
        memory: 16Mi
      limits:
        cpu: 50m
        memory: 32Mi
  restartPolicy: Never
  activeDeadlineSeconds: 3600
`, podName, selectedEnv.Namespace, os.Getenv("USER"), selectedEnv.Name, selectedEnv.RDSHost, selectedEnv.RDSHost)

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(manifest)
	if err := cmd.Run(); err != nil {
		errorExit(fmt.Sprintf("Falha ao criar bastion pod: %v", err))
	}

	printSuccess(fmt.Sprintf("Pod criado: %s", podName))
}

func waitPodReady() {
	printStep("WAIT", "Aguardando pod ficar pronto...")

	cmd := exec.Command("kubectl", "wait", "--for=condition=Ready",
		fmt.Sprintf("pod/%s", podName), "-n", selectedEnv.Namespace, "--timeout=60s")
	if err := cmd.Run(); err != nil {
		errorExit("Pod não ficou pronto dentro do tempo limite (60s). Possível problema de recursos ou permissões.")
	}

	printSuccess("Pod pronto!")
}

func testConnectivity() {
	printStep("TEST", "Testando conectividade...")

	// Wait for socat to initialize
	time.Sleep(3 * time.Second)

	// Test RDS connectivity
	cmd := exec.Command("kubectl", "exec", podName, "-n", selectedEnv.Namespace,
		"--", "nc", "-zv", selectedEnv.RDSHost, "5432")
	if err := cmd.Run(); err != nil {
		errorExit(fmt.Sprintf("Não foi possível conectar ao RDS %s:5432. Verifique security groups e conectividade de rede.", selectedEnv.RDSHost))
	}

	printSuccess("Conectividade ao RDS OK")
	printSuccess("Proxy configurado e pronto!")
}

func startPortForward() {
	printStep("FORWARD", "Iniciando port-forward...")

	// Start port-forward
	portForwardCmd = exec.Command("kubectl", "port-forward",
		fmt.Sprintf("pod/%s", podName), fmt.Sprintf("%s:5432", localPort), "-n", selectedEnv.Namespace)

	if err := portForwardCmd.Start(); err != nil {
		errorExit(fmt.Sprintf("Port-forward falhou ao iniciar. Porta %s pode estar em uso ou sem permissões.", localPort))
	}

	time.Sleep(2 * time.Second)

	showConnectionInfo()

	// Wait for port-forward to finish
	portForwardCmd.Wait()
}

func showConnectionInfo() {
	printHeader("CONEXÃO ESTABELECIDA")

	green.Println("Conexão estabelecida com sucesso!")
	fmt.Println()
	white.Println("Como conectar:")
	fmt.Println()

	cyan.Println("🔧 DBeaver/pgAdmin:")
	white.Printf("   Host:     localhost\n")
	white.Printf("   Port:     %s\n", localPort)
	white.Printf("   Database: [seu_database]\n")
	white.Printf("   Username: [seu_usuario]\n")
	fmt.Println()

	cyan.Println("💻 psql (linha de comando):")
	white.Printf("   psql -h localhost -p %s -U postgres\n", localPort)
	fmt.Println()

	cyan.Println("🌐 Outras ferramentas:")
	white.Printf("   Qualquer cliente PostgreSQL pode conectar em:\n")
	white.Printf("   localhost:%s\n", localPort)
	fmt.Println()

	red.Println("🛑 Para parar: Pressione Ctrl+C")
	yellow.Println("⚠  Tempo limite: 1 hora")
	fmt.Println()
}

func cleanup() {
	cleanupResources(false)
}

func cleanupResources(shouldExit bool) {
	fmt.Println()

	// Show error if any
	if showError && errorMessage != "" {
		red.Println("╔════════════════════════════════════════════════════════════════════╗")
		red.Println("║                             ERRO                                   ║")
		red.Println("╚════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		printError(errorMessage)
		fmt.Println()
		yellow.Println("Detalhes técnicos:")

		// Show pod logs if exists
		if podName != "" {
			cmd := exec.Command("kubectl", "get", "pod", podName, "-n", selectedEnv.Namespace)
			if cmd.Run() == nil {
				gray.Println("Logs do pod:")
				logCmd := exec.Command("kubectl", "logs", podName, "-n", selectedEnv.Namespace, "--tail=10")
				logCmd.Stdout = os.Stdout
				logCmd.Run()
				fmt.Println()

				gray.Println("Status do pod:")
				statusCmd := exec.Command("kubectl", "describe", "pod", podName, "-n", selectedEnv.Namespace)
				statusCmd.Stdout = os.Stdout
				statusCmd.Run()
			}
		}
		fmt.Println()

		if shouldExit {
			cyan.Print("Pressione Enter para continuar com a limpeza...")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}
	}

	printStep("CLEANUP", "Limpando recursos...")

	// Kill port-forward if running
	if portForwardCmd != nil && portForwardCmd.Process != nil {
		portForwardCmd.Process.Kill()
		printSuccess("Port-forward parado")
	}

	// Delete pod
	if podName != "" {
		cmd := exec.Command("kubectl", "get", "pod", podName, "-n", selectedEnv.Namespace)
		if cmd.Run() == nil {
			printStep("DELETE", fmt.Sprintf("Removendo pod %s...", podName))
			delCmd := exec.Command("kubectl", "delete", "pod", podName, "-n", selectedEnv.Namespace, "--grace-period=10")
			delCmd.Run()
			printSuccess("Pod removido")
		}
	}

	if shouldExit {
		showGoodbyeBanner()
		os.Exit(0)
	}
}

func errorExit(message string) {
	errorMessage = message
	showError = true

	// Mostrar o erro antes de sair
	fmt.Println()
	red.Println("╔════════════════════════════════════════════════════════════════════╗")
	red.Println("║                             ERRO                                   ║")
	red.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	printError(errorMessage)
	fmt.Println()

	// Executar cleanup para mostrar detalhes técnicos
	cleanupResources(false)

	// Se cleanup não sair, sair aqui
	os.Exit(1)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func isPortInUse(port string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("netstat", "-an")
	} else {
		cmd = exec.Command("lsof", "-i:"+port)
	}
	// lsof retorna 0 quando encontra a porta em uso, erro quando não encontra
	return cmd.Run() == nil
}
