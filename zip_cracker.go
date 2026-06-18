package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alexmullins/zip"
	"github.com/schollz/progressbar/v3"
)

const banner = `
██████╗ ███████╗    ███████╗██╗██████╗ 
██╔══██╗██╔════╝    ╚══███╔╝██║██╔══██╗
██████╔╝█████╗        ███╔╝ ██║██████╔╝
██╔══██╗██╔══╝       ███╔╝  ██║██╔═══╝ 
██████╔╝██║         ███████╗██║██║     
╚═════╝ ╚═╝         ╚══════╝╚═╝╚═╝     
ZIP Password Cracker v1.0
Author: GhostGTR666 - Gagaltotal666
Github: github.com/gagaltotal/zipfile-cracker-tot
`
const separator = "─────────────────────────────────────────────────"

func main() {
	fmt.Print(banner)
	fmt.Println(separator)
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// VALIDASI ARGUMEN
	// ═══════════════════════════════════════════════════════════
	if len(os.Args) != 3 {
		fmt.Println("┌─────────────────────────────────────────┐")
		fmt.Println("│  Usage: zippcrack <zip_file> <wordlist> │")
		fmt.Println("│  Example: zippcrack secret.zip word.txt │")
		fmt.Println("└─────────────────────────────────────────┘")
		os.Exit(1)
	}

	zipPath := os.Args[1]
	wordlistPath := os.Args[2]

	// ═══════════════════════════════════════════════════════════
	// VALIDASI FILE
	// ═══════════════════════════════════════════════════════════
	if err := validateFile(zipPath, "Zip file"); err != nil {
		fmt.Printf("[!] %v\n", err)
		os.Exit(1)
	}

	if err := validateFile(wordlistPath, "Wordlist"); err != nil {
		fmt.Printf("[!] %v\n", err)
		os.Exit(1)
	}

	// ═══════════════════════════════════════════════════════════
	// OPEN ZIP FILE
	// ═══════════════════════════════════════════════════════════
	fmt.Printf("[*] Target    : %s\n", zipPath)
	fmt.Printf("[*] Wordlist  : %s\n", wordlistPath)

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Printf("[!] Error opening zip file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := zipReader.Close(); err != nil {
			fmt.Printf("[!] Warning: error closing zip file: %v\n", err)
		}
	}()

	// ═══════════════════════════════════════════════════════════
	// CARI FILE TERENKRIPSI
	// ═══════════════════════════════════════════════════════════
	testFile := findEncryptedFile(zipReader.File)
	if testFile == nil {
		fmt.Println("\n[!] No encrypted files found in the archive.")
		fmt.Println("    The zip file might not be password protected.")
		os.Exit(1)
	}
	fmt.Printf("[*] Testing on: %s\n", testFile.Name)

	// ═══════════════════════════════════════════════════════════
	// HITUNG JUMLAH PASSWORD DI WORDLIST
	// ═══════════════════════════════════════════════════════════
	nWords, err := countLines(wordlistPath)
	if err != nil {
		fmt.Printf("[!] Error counting wordlist lines: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[*] Passwords : %d\n\n", nWords)

	// ═══════════════════════════════════════════════════════════
	// OPEN WORDLIST
	// ═══════════════════════════════════════════════════════════
	wordlist, err := os.Open(wordlistPath)
	if err != nil {
		fmt.Printf("[!] Error opening wordlist: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := wordlist.Close(); err != nil {
			fmt.Printf("[!] Warning: error closing wordlist: %v\n", err)
		}
	}()

	// ═══════════════════════════════════════════════════════════
	// PROGRESS BAR
	// ═══════════════════════════════════════════════════════════
	bar := progressbar.NewOptions(nWords,
		progressbar.OptionSetDescription("Cracking"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	// ═══════════════════════════════════════════════════════════
	// MAIN CRACKING LOOP
	// ═══════════════════════════════════════════════════════════
	scanner := bufio.NewScanner(wordlist)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		password := strings.TrimRight(scanner.Text(), "\r\n")

		if tryPassword(testFile, password) {
			bar.Finish()
			fmt.Println(separator)
			fmt.Printf("\n  [+] ═══════════════════════════════════ [+]\n")
			fmt.Printf("  [+]   PASSWORD FOUND: %-18s[+]\n", password)
			fmt.Printf("  [+] ═══════════════════════════════════ [+]\n\n")
			os.Exit(0)
		}

		bar.Add(1)
	}

	// ═══════════════════════════════════════════════════════════
	// HANDLE SCANNER ERROR
	// ═══════════════════════════════════════════════════════════
	if err := scanner.Err(); err != nil {
		bar.Finish()
		fmt.Printf("\n[!] Error reading wordlist: %v\n", err)
		os.Exit(1)
	}

	// ═══════════════════════════════════════════════════════════
	// PASSWORD TIDAK DITEMUKAN
	// ═══════════════════════════════════════════════════════════
	bar.Finish()
	fmt.Println(separator)
	fmt.Println("\n  [-] PASSWORD NOT FOUND")
	fmt.Println("    Try using a different wordlist.")
}

// ═══════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ═══════════════════════════════════════════════════════════════
func validateFile(path, description string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s '%s' not found", description, path)
	}
	if err != nil {
		return fmt.Errorf("error accessing %s '%s': %v", description, path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s '%s' is a directory, not a file", description, path)
	}
	if info.Size() == 0 {
		return fmt.Errorf("%s '%s' is empty", description, path)
	}
	return nil
}

func findEncryptedFile(files []*zip.File) *zip.File {
	for _, f := range files {
		if f.FileInfo().IsDir() {
			continue
		}
		if f.IsEncrypted() {
			return f
		}
	}
	return nil
}

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r\n")
		if line != "" {
			count++
		}
	}

	return count, scanner.Err()
}

func tryPassword(f *zip.File, password string) bool {
	if password == "" {
		return false
	}

	f.SetPassword(password)
	rc, err := f.Open()
	if err != nil {
		return false
	}

	if err := rc.Close(); err != nil {

	}
	return true
}
