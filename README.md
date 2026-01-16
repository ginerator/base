# api-utils


**Easiest way**

GOPRIVATE=github.com/PlanToPack/api-utils go get github.com/PlanToPack/api-utils@latest

**Final implementation**

git config --global url."https://<username>:<token>@github.com/PlanToPack/".insteadOf "https://github.com/PlanToPack/"

**Steps to Configure Your Private Go Module**

1.  **Repository Setup (GitHub):**
    
    -   Ensure your Go code is in a dedicated GitHub repository.
        
    -   **Private Repository:** Make sure the repository is set to "Private" in your GitHub repository settings.
        
    -   **Go Module Initialization:** Inside your repository, initialize a Go module. If you haven't already, run:
        
        Bash
        
        ```
        go mod init github.com/your-github-username/your-private-repo-name
        
        ```
        
        -   Replace `your-github-username` and `your-private-repo-name` with your actual GitHub username and repository name.
    -   **Commit and Push:** Commit your `go.mod` file and push the changes to your GitHub repository.
        
2.  **Authentication Setup:**
    
    -   **Personal Access Token (PAT):**
        -   Go to your GitHub settings -> Developer settings -> Personal access tokens -> Tokens (classic).
        -   Click "Generate new token (classic)".
        -   Give your token a descriptive name (e.g., "Go Module Access").
        -   **Crucially, select the `repo` scope (or `read:packages` if you are using packages).** This grants the necessary permissions to access private repositories.
        -   Copy the generated token and store it securely. **You won't be able to see it again.**
3.  **Configuring Go to Use the PAT:**
    
    -   **`~/.netrc` or `~/_netrc` File:**
        -   Create a file named `.netrc` (or `_netrc` on Windows) in your home directory.
            
        -   Add the following lines to the file, replacing the placeholders with your GitHub username and PAT:
            
            ```
            machine github.com
              login your-github-username
              password your-personal-access-token
            
            ```
            
        -   **Important:** Set the file permissions to restrict access:
            
            Bash
            
            ```
            chmod 600 ~/.netrc # or chmod 600 ~/_netrc on windows using git bash.
            
            ```
            
            -   This ensures that only your user can read the file.
4.  **Importing the Private Module in Another Project:**
    
    -   **`go get` or `go mod tidy`:** In your other Go project, use `go get` or `go mod tidy` to add the private module as a dependency:
        
        Bash
        
        ```
        go get github.com/your-github-username/your-private-repo-name
        
        ```
        
        -   or if you are working on a project that already has dependencies, and you added the import statement to your go code, use:
        
        Bash
        
        ```
        go mod tidy
        
        ```
        
    -   **Import in Go Code:** In your Go code, import the module using the same import path:
        
        Go
        
        ```
        import "github.com/your-github-username/your-private-repo-name"
        ```