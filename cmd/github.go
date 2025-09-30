package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Set up GitHub with SSH keys",
	Long:  `Configure GitHub SSH keys and authentication for development`,
}

var githubSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up GitHub SSH keys and authentication",
	Long:  `Generate SSH keys, add to ssh-agent, and provide instructions for adding to GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		keyType, _ := cmd.Flags().GetString("key-type")
		skipAgent, _ := cmd.Flags().GetBool("skip-agent")

		if email == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter your GitHub email: ")
			email, _ = reader.ReadString('\n')
			email = strings.TrimSpace(email)
		}

		if email == "" {
			fmt.Println("❌ Email is required for SSH key generation")
			os.Exit(1)
		}

		fmt.Println("🔐 Setting up GitHub SSH authentication...")
		fmt.Printf("📧 Email: %s\n", email)
		fmt.Printf("🔑 Key type: %s\n", keyType)
		fmt.Println()

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		// Use private directory for SSH keys
		privateSshDir := filepath.Join(home, ".dotfiles", "private", ".ssh")
		keyPath := filepath.Join(privateSshDir, "id_"+keyType)
		pubKeyPath := keyPath + ".pub"

		// Create private .ssh directory if it doesn't exist
		if err := os.MkdirAll(privateSshDir, 0700); err != nil {
			fmt.Printf("❌ Error creating private .ssh directory: %v\n", err)
			os.Exit(1)
		}

		// Check if key already exists
		if _, err := os.Stat(keyPath); err == nil {
			fmt.Printf("🔑 SSH key already exists at %s\n", keyPath)

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Do you want to create a new key? (y/N): ")
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				fmt.Println("✅ Using existing SSH key")
				showPublicKey(pubKeyPath)
				return
			}
		}

		// Generate SSH key
		fmt.Println("🔨 Generating SSH key...")
		sshKeygenArgs := []string{
			"-t", keyType,
			"-C", email,
			"-f", keyPath,
			"-N", "", // No passphrase for simplicity
		}

		sshKeygenCmd := exec.Command("ssh-keygen", sshKeygenArgs...)
		if output, err := sshKeygenCmd.CombinedOutput(); err != nil {
			fmt.Printf("❌ Error generating SSH key: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			os.Exit(1)
		}

		fmt.Printf("✅ SSH key generated successfully!\n")
		fmt.Printf("🔑 Private key: %s\n", keyPath)
		fmt.Printf("🔑 Public key: %s\n", pubKeyPath)
		fmt.Println()

		// Set proper permissions
		os.Chmod(keyPath, 0600)
		os.Chmod(pubKeyPath, 0644)

		// Add to SSH agent
		if !skipAgent {
			fmt.Println("🔐 Adding key to SSH agent...")
			if err := addToSSHAgent(keyPath); err != nil {
				fmt.Printf("⚠️  Warning: Could not add key to SSH agent: %v\n", err)
			} else {
				fmt.Println("✅ Key added to SSH agent")
			}
			fmt.Println()
		}

		// Set up SSH stow package
		fmt.Println("🔗 Setting up SSH stow package...")
		if err := setupSSHStowPackage(home, privateSshDir); err != nil {
			fmt.Printf("⚠️  Warning: Could not set up SSH stow package: %v\n", err)
		} else {
			fmt.Println("✅ SSH stow package created")
		}

		// Show public key and next steps
		showPublicKey(pubKeyPath)
	},
}

var githubTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test GitHub SSH connection",
	Long:  `Test SSH connection to GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧪 Testing GitHub SSH connection...")

		sshCmd := exec.Command("ssh", "-T", "git@github.com")
		output, err := sshCmd.CombinedOutput()

		outputStr := string(output)
		if err != nil {
			if strings.Contains(outputStr, "successfully authenticated") {
				fmt.Println("✅ GitHub SSH connection successful!")
				fmt.Printf("Response: %s\n", outputStr)
			} else {
				fmt.Printf("❌ GitHub SSH connection failed: %v\n", err)
				fmt.Printf("Output: %s\n", outputStr)
				fmt.Println("\n💡 Make sure you've added your SSH key to GitHub:")
				fmt.Println("   https://github.com/settings/ssh/new")
			}
		} else {
			fmt.Println("✅ GitHub SSH connection successful!")
			fmt.Printf("Response: %s\n", outputStr)
		}
	},
}

func addToSSHAgent(keyPath string) error {
	// Start ssh-agent if not running
	if os.Getenv("SSH_AUTH_SOCK") == "" {
		fmt.Println("🔄 Starting ssh-agent...")
		evalCmd := exec.Command("bash", "-c", "eval $(ssh-agent -s)")
		if err := evalCmd.Run(); err != nil {
			return fmt.Errorf("failed to start ssh-agent: %v", err)
		}
	}

	// Add key to agent
	addCmd := exec.Command("ssh-add", keyPath)
	return addCmd.Run()
}

func showPublicKey(pubKeyPath string) {
	fmt.Println("📋 Your public SSH key:")
	fmt.Println("=" + strings.Repeat("=", 50))

	pubKeyContent, err := os.ReadFile(pubKeyPath)
	if err != nil {
		fmt.Printf("❌ Error reading public key: %v\n", err)
		return
	}

	fmt.Print(string(pubKeyContent))
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	fmt.Println("📌 Next steps:")
	fmt.Println("1. Copy the above public key")
	fmt.Println("2. Go to: https://github.com/settings/ssh/new")
	fmt.Println("3. Add a title (e.g., 'My Development Machine')")
	fmt.Println("4. Paste the public key")
	fmt.Println("5. Click 'Add SSH key'")
	fmt.Println("6. Run: dotfiles stow ssh (to symlink SSH keys to home directory)")
	fmt.Println("7. Test with: dotfiles github test")
	fmt.Println()

	// Try to copy to clipboard
	if err := copyToClipboard(string(pubKeyContent)); err == nil {
		fmt.Println("📋 Public key copied to clipboard!")
	} else {
		fmt.Printf("⚠️  Could not copy to clipboard: %v\n", err)
	}
}

func setupSSHStowPackage(home, privateSshDir string) error {
	stowDir := filepath.Join(home, ".dotfiles", "stow")
	sshStowDir := filepath.Join(stowDir, "ssh")

	// Create ssh stow package directory
	if err := os.MkdirAll(sshStowDir, 0755); err != nil {
		return fmt.Errorf("failed to create ssh stow directory: %v", err)
	}

	// Create relative symlink from stow package to private .ssh directory
	stowSshLink := filepath.Join(sshStowDir, ".ssh")
	relativePrivatePath := filepath.Join("..", "..", "private", ".ssh")

	// Remove existing symlink if it exists
	if _, err := os.Lstat(stowSshLink); err == nil {
		if err := os.Remove(stowSshLink); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}

	// Create the symlink
	if err := os.Symlink(relativePrivatePath, stowSshLink); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}

	fmt.Printf("   Created symlink: %s -> %s\n", stowSshLink, relativePrivatePath)
	fmt.Printf("   💡 Run 'dotfiles stow ssh' to symlink SSH keys to home directory\n")

	return nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	// Try different clipboard commands based on OS
	if _, err := exec.LookPath("pbcopy"); err == nil {
		// macOS
		cmd = exec.Command("pbcopy")
	} else if _, err := exec.LookPath("xclip"); err == nil {
		// Linux with xclip
		cmd = exec.Command("xclip", "-selection", "clipboard")
	} else if _, err := exec.LookPath("xsel"); err == nil {
		// Linux with xsel
		cmd = exec.Command("xsel", "--clipboard", "--input")
	} else {
		return fmt.Errorf("no clipboard utility found")
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func init() {
	githubSetupCmd.Flags().StringP("email", "e", "", "Email for SSH key")
	githubSetupCmd.Flags().String("key-type", "ed25519", "SSH key type (ed25519, rsa)")
	githubSetupCmd.Flags().Bool("skip-agent", false, "Skip adding key to SSH agent")

	githubCmd.AddCommand(githubSetupCmd)
	githubCmd.AddCommand(githubTestCmd)
	rootCmd.AddCommand(githubCmd)
}