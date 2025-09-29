// Modern Dotfiles Setup Wizard
// Clean, modular architecture with proper state management

class SetupWizard {
    constructor() {
        this.currentStep = 0;
        this.totalSteps = 7;
        this.config = {
            personal: {
                name: '',
                email: '',
                editor: 'nvim'
            },
            installation: {
                homebrew: true,
                brewfile: true,
                dotfiles: true,
                macos_defaults: true,
                npm_packages: true
            },
            system: {
                appearance: {
                    dark_mode: true,
                    enable_24_hour_time: false
                },
                dock: {
                    autohide: false,
                    position: 'bottom',
                    tile_size: 50,
                    magnification: false
                }
            },
            development: {
                git: {
                    default_branch: 'main',
                    pull_rebase: true,
                    push_default: 'simple'
                },
                languages: {},
                frameworks: {},
                tools: {},
                shell: {
                    theme: 'powerlevel10k',
                    terminal_theme: 'dark',
                    plugins: {}
                },
                aliases: {}
            },
            tools: {
                git: false,
                stow: false,
                basic: false,
                modern_cli: false,
                json: false
            },
            shell: {
                zsh: false
            },
            packages: {
                extra_brews: [],
                extra_casks: [],
                extra_taps: [],
                npm_globals: []
            },
            metadata: {
                title: '',
                description: '',
                version: '1.0.0',
                created_by: 'web-wizard'
            }
        };

        this.languages = [
            {
                key: 'javascript',
                name: 'JavaScript',
                icon: 'fab fa-js-square',
                color: 'text-yellow-500',
                description: 'Essential for web development',
                popularity: 95
            },
            {
                key: 'typescript',
                name: 'TypeScript',
                icon: 'fas fa-code',
                color: 'text-blue-600',
                description: 'Type-safe JavaScript',
                popularity: 85
            },
            {
                key: 'python',
                name: 'Python',
                icon: 'fab fa-python',
                color: 'text-blue-500',
                description: 'Data science & automation',
                popularity: 90
            },
            {
                key: 'go',
                name: 'Go',
                icon: 'fas fa-bolt',
                color: 'text-cyan-500',
                description: 'Fast backend development',
                popularity: 75
            },
            {
                key: 'rust',
                name: 'Rust',
                icon: 'fas fa-cog',
                color: 'text-orange-500',
                description: 'Systems programming',
                popularity: 70
            },
            {
                key: 'java',
                name: 'Java',
                icon: 'fab fa-java',
                color: 'text-red-500',
                description: 'Enterprise applications',
                popularity: 80
            },
            {
                key: 'swift',
                name: 'Swift',
                icon: 'fab fa-swift',
                color: 'text-orange-600',
                description: 'iOS & macOS development',
                popularity: 65
            },
            {
                key: 'kotlin',
                name: 'Kotlin',
                icon: 'fas fa-mobile-alt',
                color: 'text-purple-500',
                description: 'Android development',
                popularity: 60
            },
            {
                key: 'csharp',
                name: 'C#',
                icon: 'fas fa-hashtag',
                color: 'text-purple-600',
                description: '.NET development',
                popularity: 70
            }
        ];

        this.tools = [
            {
                key: 'git',
                name: 'Git',
                icon: 'fab fa-git-alt',
                color: 'text-orange-500',
                description: 'Version control system',
                recommended: true
            },
            {
                key: 'stow',
                name: 'GNU Stow',
                icon: 'fas fa-link',
                color: 'text-blue-500',
                description: 'Dotfiles management',
                recommended: true
            },
            {
                key: 'docker',
                name: 'Docker',
                icon: 'fab fa-docker',
                color: 'text-blue-600',
                description: 'Containerization platform'
            },
            {
                key: 'vscode',
                name: 'VS Code',
                icon: 'fas fa-code',
                color: 'text-blue-500',
                description: 'Popular code editor'
            },
            {
                key: 'postman',
                name: 'Postman',
                icon: 'fas fa-paper-plane',
                color: 'text-orange-500',
                description: 'API development environment'
            },
            {
                key: 'basic',
                name: 'Basic Tools',
                icon: 'fas fa-toolbox',
                color: 'text-gray-600',
                description: 'curl, wget, tree, unzip'
            }
        ];

        this.presets = {
            frontend: ['javascript', 'typescript'],
            backend: ['python', 'go', 'java'],
            fullstack: ['javascript', 'typescript', 'python'],
            mobile: ['swift', 'kotlin', 'java']
        };

        this.popularPackages = {
            development: [
                { brew: 'git', name: 'Git', description: 'Distributed version control system', category: 'Essential' },
                { brew: 'gh', name: 'GitHub CLI', description: 'GitHub command line tool', category: 'Git' },
                { brew: 'lazygit', name: 'LazyGit', description: 'Simple terminal UI for git commands', category: 'Git' },
                { brew: 'node', name: 'Node.js', description: 'JavaScript runtime environment', category: 'Runtime' },
                { brew: 'python@3.12', name: 'Python 3.12', description: 'Python programming language', category: 'Runtime' },
                { brew: 'go', name: 'Go', description: 'Go programming language', category: 'Runtime' },
                { brew: 'rust', name: 'Rust', description: 'Rust programming language', category: 'Runtime' },
                { brew: 'docker', name: 'Docker', description: 'Container platform', category: 'DevOps' },
                { brew: 'docker-compose', name: 'Docker Compose', description: 'Multi-container Docker applications', category: 'DevOps' },
                { brew: 'terraform', name: 'Terraform', description: 'Infrastructure as code tool', category: 'DevOps' },
                { brew: 'kubectl', name: 'kubectl', description: 'Kubernetes command-line tool', category: 'DevOps' },
                { brew: 'helm', name: 'Helm', description: 'Kubernetes package manager', category: 'DevOps' },
                { brew: 'ansible', name: 'Ansible', description: 'Configuration management tool', category: 'DevOps' },
                { brew: 'vim', name: 'Vim', description: 'Vi IMproved text editor', category: 'Editor' },
                { brew: 'neovim', name: 'Neovim', description: 'Hyperextensible Vim-based text editor', category: 'Editor' },
                { brew: 'tmux', name: 'tmux', description: 'Terminal multiplexer', category: 'Terminal' },
                { brew: 'zsh', name: 'Zsh', description: 'Z shell with many improvements', category: 'Shell' },
                { brew: 'fish', name: 'Fish', description: 'User-friendly command line shell', category: 'Shell' },
                { brew: 'starship', name: 'Starship', description: 'Cross-shell prompt', category: 'Shell' },
                { brew: 'mysql', name: 'MySQL', description: 'Open source relational database', category: 'Database' },
                { brew: 'postgresql', name: 'PostgreSQL', description: 'Object-relational database system', category: 'Database' },
                { brew: 'redis', name: 'Redis', description: 'In-memory data structure store', category: 'Database' },
                { brew: 'mongodb/brew/mongodb-community', name: 'MongoDB', description: 'Document database', category: 'Database' }
            ],
            utilities: [
                { brew: 'wget', name: 'wget', description: 'Internet file retriever', category: 'Network' },
                { brew: 'curl', name: 'cURL', description: 'Command line tool for transferring data', category: 'Network' },
                { brew: 'httpie', name: 'HTTPie', description: 'Modern command line HTTP client', category: 'Network' },
                { brew: 'speedtest-cli', name: 'Speedtest CLI', description: 'Command line internet speed test', category: 'Network' },
                { brew: 'jq', name: 'jq', description: 'Lightweight JSON processor', category: 'JSON' },
                { brew: 'yq', name: 'yq', description: 'Process YAML like jq processes JSON', category: 'YAML' },
                { brew: 'tree', name: 'tree', description: 'Display directories as trees', category: 'File' },
                { brew: 'htop', name: 'htop', description: 'Improved top (interactive process viewer)', category: 'System' },
                { brew: 'btop', name: 'btop', description: 'Modern system monitor', category: 'System' },
                { brew: 'bat', name: 'bat', description: 'Clone of cat with syntax highlighting', category: 'File' },
                { brew: 'eza', name: 'eza', description: 'Modern replacement for ls', category: 'File' },
                { brew: 'fd', name: 'fd', description: 'Simple, fast alternative to find', category: 'File' },
                { brew: 'ripgrep', name: 'ripgrep', description: 'Search tool like grep and ag', category: 'Search' },
                { brew: 'fzf', name: 'fzf', description: 'Command-line fuzzy finder', category: 'Search' },
                { brew: 'ag', name: 'The Silver Searcher', description: 'Code-searching tool similar to ack', category: 'Search' },
                { brew: 'dust', name: 'dust', description: 'More intuitive version of du', category: 'System' },
                { brew: 'duf', name: 'duf', description: 'Better df alternative', category: 'System' },
                { brew: 'ncdu', name: 'ncdu', description: 'NCurses disk usage analyzer', category: 'System' },
                { brew: 'watch', name: 'watch', description: 'Execute a program periodically', category: 'System' },
                { brew: 'rsync', name: 'rsync', description: 'File transfer and synchronization', category: 'File' },
                { brew: 'unzip', name: 'unzip', description: 'Extract compressed files', category: 'Archive' },
                { brew: 'p7zip', name: '7-Zip', description: 'File archiver with high compression', category: 'Archive' },
                { brew: 'mas', name: 'Mac App Store CLI', description: 'Install Mac App Store apps', category: 'macOS' },
                { brew: 'mackup', name: 'Mackup', description: 'Sync application settings', category: 'macOS' }
            ],
            media: [
                { brew: 'ffmpeg', name: 'FFmpeg', description: 'Audio and video processing tool', category: 'Media' },
                { brew: 'imagemagick', name: 'ImageMagick', description: 'Image editing and conversion', category: 'Images' },
                { brew: 'youtube-dl', name: 'youtube-dl', description: 'Download videos from YouTube', category: 'Download' },
                { brew: 'yt-dlp', name: 'yt-dlp', description: 'Enhanced YouTube downloader', category: 'Download' },
                { brew: 'sox', name: 'SoX', description: 'Sound processing library', category: 'Audio' },
                { brew: 'gifsicle', name: 'Gifsicle', description: 'GIF manipulation utility', category: 'Images' }
            ]
        };

        this.popularCasks = {
            development: [
                { cask: 'visual-studio-code', name: 'VS Code', description: 'Source code editor', category: 'Editor' },
                { cask: 'sublime-text', name: 'Sublime Text', description: 'Sophisticated text editor', category: 'Editor' },
                { cask: 'jetbrains-toolbox', name: 'JetBrains Toolbox', description: 'Manage JetBrains IDEs', category: 'Editor' },
                { cask: 'cursor', name: 'Cursor', description: 'AI-powered code editor', category: 'Editor' },
                { cask: 'zed', name: 'Zed', description: 'High-performance code editor', category: 'Editor' },
                { cask: 'github-desktop', name: 'GitHub Desktop', description: 'Git GUI client', category: 'Development' },
                { cask: 'sourcetree', name: 'Sourcetree', description: 'Git GUI by Atlassian', category: 'Development' },
                { cask: 'postman', name: 'Postman', description: 'API development environment', category: 'Development' },
                { cask: 'insomnia', name: 'Insomnia', description: 'REST client', category: 'Development' },
                { cask: 'tableplus', name: 'TablePlus', description: 'Database management tool', category: 'Development' },
                { cask: 'sequel-pro', name: 'Sequel Pro', description: 'MySQL database management', category: 'Development' },
                { cask: 'docker', name: 'Docker Desktop', description: 'Containerization platform', category: 'Development' },
                { cask: 'iterm2', name: 'iTerm2', description: 'Terminal emulator for macOS', category: 'Terminal' },
                { cask: 'warp', name: 'Warp', description: 'Modern terminal with AI features', category: 'Terminal' },
                { cask: 'hyper', name: 'Hyper', description: 'Electron-based terminal', category: 'Terminal' }
            ],
            productivity: [
                { cask: 'google-chrome', name: 'Google Chrome', description: 'Web browser by Google', category: 'Browser' },
                { cask: 'firefox', name: 'Firefox', description: 'Web browser by Mozilla', category: 'Browser' },
                { cask: 'arc', name: 'Arc', description: 'Modern web browser', category: 'Browser' },
                { cask: 'brave-browser', name: 'Brave', description: 'Privacy-focused browser', category: 'Browser' },
                { cask: 'slack', name: 'Slack', description: 'Team collaboration hub', category: 'Communication' },
                { cask: 'discord', name: 'Discord', description: 'Voice and text chat', category: 'Communication' },
                { cask: 'zoom', name: 'Zoom', description: 'Video conferencing', category: 'Communication' },
                { cask: 'microsoft-teams', name: 'Microsoft Teams', description: 'Team collaboration platform', category: 'Communication' },
                { cask: 'notion', name: 'Notion', description: 'All-in-one workspace', category: 'Productivity' },
                { cask: 'obsidian', name: 'Obsidian', description: 'Knowledge base on local Markdown files', category: 'Productivity' },
                { cask: 'raycast', name: 'Raycast', description: 'Productivity launcher', category: 'Productivity' },
                { cask: 'alfred', name: 'Alfred', description: 'Application launcher', category: 'Productivity' },
                { cask: 'rectangle', name: 'Rectangle', description: 'Window management utility', category: 'Utility' },
                { cask: 'bartender-4', name: 'Bartender 4', description: 'Menu bar organization', category: 'Utility' },
                { cask: 'cleanmymac', name: 'CleanMyMac X', description: 'System cleaner and optimizer', category: 'Utility' },
                { cask: '1password', name: '1Password', description: 'Password manager', category: 'Security' },
                { cask: 'bitwarden', name: 'Bitwarden', description: 'Open source password manager', category: 'Security' }
            ],
            design: [
                { cask: 'figma', name: 'Figma', description: 'Collaborative design tool', category: 'Design' },
                { cask: 'sketch', name: 'Sketch', description: 'Vector graphics editor', category: 'Design' },
                { cask: 'canva', name: 'Canva', description: 'Graphic design platform', category: 'Design' },
                { cask: 'adobe-creative-cloud', name: 'Adobe Creative Cloud', description: 'Adobe creative suite', category: 'Design' },
                { cask: 'affinity-designer', name: 'Affinity Designer', description: 'Professional graphic design software', category: 'Design' },
                { cask: 'affinity-photo', name: 'Affinity Photo', description: 'Professional photo editing software', category: 'Design' }
            ]
        };

        this.npmPackages = [
            { name: '@vue/cli', description: 'Vue.js development tools', category: 'Framework' },
            { name: 'create-react-app', description: 'Create React applications', category: 'Framework' },
            { name: '@angular/cli', description: 'Angular development tools', category: 'Framework' },
            { name: 'create-next-app', description: 'Create Next.js apps', category: 'Framework' },
            { name: 'vite', description: 'Fast build tool for modern web projects', category: 'Build' },
            { name: 'webpack', description: 'Module bundler for JavaScript', category: 'Build' },
            { name: 'typescript', description: 'Typed superset of JavaScript', category: 'Language' },
            { name: 'prettier', description: 'Opinionated code formatter', category: 'Tool' },
            { name: 'eslint', description: 'JavaScript linting utility', category: 'Tool' },
            { name: 'nodemon', description: 'Monitor for changes in Node.js apps', category: 'Tool' },
            { name: 'pm2', description: 'Production process manager for Node.js', category: 'Tool' },
            { name: 'serve', description: 'Static file serving', category: 'Tool' },
            { name: 'http-server', description: 'Simple HTTP server', category: 'Tool' },
            { name: 'yarn', description: 'Fast, reliable package manager', category: 'Package Manager' },
            { name: 'pnpm', description: 'Fast, disk space efficient package manager', category: 'Package Manager' },
            { name: 'express-generator', description: 'Express application generator', category: 'Framework' },
            { name: 'nest', description: 'NestJS command line interface', category: 'Framework' },
            { name: 'gatsby-cli', description: 'Gatsby command line interface', category: 'Framework' },
            { name: 'vercel', description: 'Vercel deployment platform CLI', category: 'Deployment' },
            { name: 'netlify-cli', description: 'Netlify deployment platform CLI', category: 'Deployment' },
            { name: 'firebase-tools', description: 'Firebase command line tools', category: 'Deployment' },
            { name: 'aws-cli', description: 'Amazon Web Services CLI', category: 'Cloud' },
            { name: 'serverless', description: 'Build serverless applications', category: 'Cloud' }
        ];

        this.selectedPackages = {
            brews: new Set(),
            casks: new Set(),
            npm: new Set()
        };
        this.currentPackageTab = 'development';
        this.currentSearchTerm = '';

        this.init();
    }

    init() {
        this.setupEventListeners();
        this.renderLanguages();
        this.renderTools();
        this.renderPackages();
        this.updateUI();
        this.updateProgress();
        this.animateOnLoad();
    }

    setupEventListeners() {
        // Navigation
        window.nextStep = () => this.nextStep();
        window.previousStep = () => this.previousStep();
        window.saveConfiguration = () => this.saveConfiguration();
        window.selectPreset = (preset) => this.selectPreset(preset);

        // Form inputs
        this.setupFormInputs();
        this.setupPackageInputs();
        this.setupInstallationInputs();
        // Skip git/ssh/development inputs since we simplified the UI

        // Checkbox change handlers
        this.setupCheckboxHandlers();
    }

    setupFormInputs() {
        const inputs = ['userName', 'userEmail', 'userEditor', 'configTitle', 'configDescription'];
        inputs.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.addEventListener('input', () => this.updatePersonalInfo());
            }
        });

        // System preferences
        const systemInputs = ['darkMode', 'time24Hour', 'dockAutohide', 'dockPosition'];
        systemInputs.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.addEventListener('change', () => this.updateSystemConfig());
            }
        });
    }

    setupCheckboxHandlers() {
        // Modern CLI and Zsh
        const modernCli = document.getElementById('modernCli');
        const zshShell = document.getElementById('zshShell');

        if (modernCli) {
            modernCli.addEventListener('change', () => {
                this.config.shell.modern_cli = modernCli.checked;
                this.updateCardSelection('modernCliCard', modernCli.checked);
                this.updatePreview();
            });
        }

        if (zshShell) {
            zshShell.addEventListener('change', () => {
                this.config.shell.zsh = zshShell.checked;
                this.updateCardSelection('zshCard', zshShell.checked);
                this.updatePreview();
            });
        }
    }

    setupGitSshInputs() {
        // Git configuration
        const gitInputs = ['gitSigningKey', 'gitDefaultBranch', 'gitAutoSignCommits', 'gitAliases'];
        gitInputs.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.addEventListener('change', () => this.updateGitConfig());
            }
        });

        // SSH configuration
        const sshInputs = ['sshKeyType', 'addToSshAgent', 'setupGitHub', 'setupGitLab', 'setupBitbucket'];
        sshInputs.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.addEventListener('change', () => this.updateSshConfig());
            }
        });
    }

    setupDevelopmentInputs() {
        // Terminal and development environment
        const devInputs = [
            'terminalTheme', 'terminalFont', 'enableLigatures',
            'commonAliases', 'dockerAliases', 'kubernetesAliases',
            'homebrewBundle', 'npmGlobals', 'pythonSetup',
            'dotfilesRepo', 'backupExisting'
        ];
        devInputs.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.addEventListener('change', () => this.updateDevelopmentConfig());
            }
        });

        // Editor change should update IDE extensions
        const editorSelect = document.getElementById('userEditor');
        if (editorSelect) {
            editorSelect.addEventListener('change', () => {
                this.updatePersonalInfo();
                this.updateIdeExtensions();
            });
        }
    }

    updateCardSelection(cardId, selected) {
        const card = document.getElementById(cardId);
        if (card) {
            if (selected) {
                card.classList.add('selected');
            } else {
                card.classList.remove('selected');
            }
        }
    }

    renderLanguages() {
        const grid = document.getElementById('languageGrid');
        if (!grid) return;

        grid.innerHTML = this.languages.map(lang => `
            <div class="selection-card language-card" data-key="${lang.key}" onclick="toggleLanguage('${lang.key}')">
                <div class="language-icon">•</div>
                <h3 style="font-weight: 600; margin: 0 0 5px 0; font-size: 14px;">${lang.name}</h3>
                <p style="font-size: 12px; color: #666; margin: 0;">${lang.description}</p>
            </div>
        `).join('');

        // Setup click handlers
        window.toggleLanguage = (key) => this.toggleLanguage(key);
        window.togglePackage = (id, type) => this.togglePackage(id, type);
    }

    renderTools() {
        const grid = document.getElementById('coreToolsGrid');
        if (!grid) return;

        grid.innerHTML = this.tools.map(tool => `
            <div class="tool-item" data-key="${tool.key}" onclick="toggleTool('${tool.key}')">
                <input type="checkbox" class="tool-checkbox" onchange="event.stopPropagation()">
                <div style="flex: 1;">
                    <h4 style="font-weight: 600; margin: 0 0 3px 0; font-size: 14px;">${tool.name}</h4>
                    <p style="font-size: 12px; color: #666; margin: 0;">${tool.description}</p>
                </div>
                ${tool.recommended ? '<span style="background: #fff3cd; color: #856404; font-size: 11px; padding: 2px 6px; border-radius: 3px;">Recommended</span>' : ''}
            </div>
        `).join('');

        // Setup click handlers
        window.toggleTool = (key) => this.toggleTool(key);
        window.generateGpgKey = () => this.generateGpgKey();
        window.generateSshKey = () => this.generateSshKey();
        window.addExistingKey = () => this.addExistingKey();

        // Update IDE extensions when editor changes
        this.updateIdeExtensions();
    }

    toggleLanguage(key) {
        this.config.development.languages[key] = !this.config.development.languages[key];
        this.updateLanguageCard(key);
        this.updatePreview();
    }

    toggleTool(key) {
        this.config.tools[key] = !this.config.tools[key];
        this.updateToolCard(key);
        this.updatePreview();
    }

    updateLanguageCard(key) {
        const card = document.querySelector(`[data-key="${key}"].language-card`);
        if (!card) return;

        const isSelected = this.config.development.languages[key];
        const icon = card.querySelector('.language-icon');
        const indicator = card.querySelector('.selection-indicator');
        const progressBar = card.querySelector('.popularity-bar .h-1');

        if (isSelected) {
            card.classList.add('selected');
        } else {
            card.classList.remove('selected');
        }
    }

    updateToolCard(key) {
        const card = document.querySelector(`[data-key="${key}"].tool-item`);
        if (!card) return;

        const checkbox = card.querySelector('.tool-checkbox');
        const isSelected = this.config.tools[key];

        checkbox.checked = isSelected;

        if (isSelected) {
            card.classList.add('selected');
        } else {
            card.classList.remove('selected');
        }
    }

    selectPreset(preset) {
        // Reset all languages
        Object.keys(this.config.development.languages).forEach(key => {
            this.config.development.languages[key] = false;
        });

        // Apply preset
        if (this.presets[preset]) {
            this.presets[preset].forEach(lang => {
                this.config.development.languages[lang] = true;
            });
        }

        // Update UI
        this.languages.forEach(lang => {
            this.updateLanguageCard(lang.key);
        });
        this.updatePreview();
        this.showNotification(`🎯 ${preset.charAt(0).toUpperCase() + preset.slice(1)} preset applied!`, 'success');
    }

    updatePersonalInfo() {
        this.config.personal.name = document.getElementById('userName')?.value || '';
        this.config.personal.email = document.getElementById('userEmail')?.value || '';
        this.config.personal.editor = document.getElementById('userEditor')?.value || 'nvim';
        this.config.metadata.title = document.getElementById('configTitle')?.value || '';
        this.config.metadata.description = document.getElementById('configDescription')?.value || '';
        this.updatePreview();
    }

    updateSystemConfig() {
        this.config.system.appearance.dark_mode = document.getElementById('darkMode')?.checked || false;
        this.config.system.appearance.enable_24_hour_time = document.getElementById('time24Hour')?.checked || false;
        this.config.system.dock.autohide = document.getElementById('dockAutohide')?.checked || false;
        this.config.system.dock.position = document.getElementById('dockPosition')?.value || 'bottom';
        this.updatePreview();
    }

    nextStep() {
        if (this.currentStep < this.totalSteps - 1) {
            // Validate current step
            if (!this.validateCurrentStep()) {
                return;
            }

            this.currentStep++;
            this.updateUI();
            this.updateProgress();
            this.animateStepTransition();
        }
    }

    previousStep() {
        if (this.currentStep > 0) {
            this.currentStep--;
            this.updateUI();
            this.updateProgress();
            this.animateStepTransition();
        }
    }

    validateCurrentStep() {
        if (this.currentStep === 0) {
            const name = document.getElementById('userName')?.value.trim();
            const email = document.getElementById('userEmail')?.value.trim();

            if (!name || !email) {
                this.showNotification('❌ Please fill in your name and email', 'error');
                return false;
            }

            // Basic email validation
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(email)) {
                this.showNotification('❌ Please enter a valid email address', 'error');
                return false;
            }
        }

        // Skip Git/SSH step validation since we simplified the UI
        if (this.currentStep === 1) {
            // No validation needed for language step
        }

        // Language step
        if (this.currentStep === 1) {
            const hasLanguages = Object.values(this.config.development.languages).some(Boolean);
            if (!hasLanguages) {
                this.showNotification('💡 Consider selecting at least one programming language', 'warning');
                // Don't block, just warn
            }
        }

        // Packages step
        if (this.currentStep === 3) {
            this.updatePackageConfig();
        }

        // System step
        if (this.currentStep === 4) {
            this.updateSystemConfig();
        }

        // Extra settings step
        if (this.currentStep === 5) {
            this.updatePackageConfig();
            this.updateInstallationConfig();
        }

        return true;
    }

    updateUI() {
        // Hide all step contents
        document.querySelectorAll('.step-content').forEach(content => {
            content.style.display = 'none';
        });

        // Show current step
        const currentContent = document.getElementById(`step-content-${this.currentStep}`);
        if (currentContent) {
            currentContent.style.display = 'block';
        }

        // Update step indicators
        document.querySelectorAll('.step').forEach((step, index) => {
            step.classList.remove('active', 'completed');
            if (index < this.currentStep) {
                step.classList.add('completed');
            } else if (index === this.currentStep) {
                step.classList.add('active');
            }
        });

        // Update navigation buttons
        const prevBtn = document.getElementById('prevBtn');
        const nextBtn = document.getElementById('nextBtn');
        const saveBtn = document.getElementById('saveBtn');

        if (prevBtn) {
            prevBtn.style.display = this.currentStep > 0 ? 'inline-flex' : 'none';
        }

        if (nextBtn) {
            nextBtn.style.display = this.currentStep < this.totalSteps - 1 ? 'inline-flex' : 'none';
        }

        if (saveBtn) {
            saveBtn.style.display = this.currentStep === this.totalSteps - 1 ? 'inline-flex' : 'none';
        }

        this.updatePreview();
        this.updateConfigSummary();
    }

    updateProgress() {
        const progress = ((this.currentStep + 1) / this.totalSteps) * 100;
        const progressBar = document.getElementById('progressBar');
        const progressText = document.getElementById('progressText');

        if (progressBar) {
            progressBar.style.width = `${progress}%`;
        }

        if (progressText) {
            progressText.textContent = `Step ${this.currentStep + 1} of ${this.totalSteps}`;
        }
    }

    updatePreview() {
        const preview = document.getElementById('sidebarPreview');
        if (!preview) return;

        const enabledLanguages = this.getEnabledLanguages();
        const enabledTools = this.getEnabledTools();
        const systemFeatures = this.getEnabledSystemFeatures();

        if (enabledLanguages.length === 0 && enabledTools.length === 0 && systemFeatures.length === 0) {
            preview.innerHTML = `
                <div class="text-center text-gray-500 py-8">
                    <i class="fas fa-rocket text-4xl mb-4 opacity-50"></i>
                    <p>Start configuring to see preview</p>
                </div>
            `;
            return;
        }

        preview.innerHTML = `
            <div class="space-y-4">
                ${this.renderPreviewSection('Languages', enabledLanguages, 'fa-code', 'text-blue-500')}
                ${this.renderPreviewSection('Tools', enabledTools, 'fa-tools', 'text-green-500')}
                ${this.renderPreviewSection('System', systemFeatures, 'fa-cog', 'text-purple-500')}

                <div class="border-t pt-4 mt-4">
                    <div class="flex justify-between text-sm mb-2">
                        <span class="text-gray-600">Estimated Size:</span>
                        <strong class="text-orange-600">${this.getEstimatedSize()}</strong>
                    </div>
                    <div class="flex justify-between text-sm">
                        <span class="text-gray-600">Setup Time:</span>
                        <strong class="text-red-600">${this.getEstimatedTime()}</strong>
                    </div>
                </div>
            </div>
        `;

        // Update live preview on system step
        const livePreview = document.getElementById('livePreview');
        if (livePreview && this.currentStep === 5) {
            livePreview.innerHTML = `
                <div class="flex justify-between">
                    <span class="text-gray-600">Selected Languages:</span>
                    <span class="font-semibold text-blue-600">${enabledLanguages.length}</span>
                </div>
                <div class="flex justify-between">
                    <span class="text-gray-600">Essential Tools:</span>
                    <span class="font-semibold text-green-600">${enabledTools.length}</span>
                </div>
                <div class="flex justify-between">
                    <span class="text-gray-600">System Features:</span>
                    <span class="font-semibold text-purple-600">${systemFeatures.length}</span>
                </div>
                <hr class="my-4">
                <div class="flex justify-between">
                    <span class="text-gray-600">Estimated Install Size:</span>
                    <span class="font-semibold text-orange-600">${this.getEstimatedSize()}</span>
                </div>
                <div class="flex justify-between">
                    <span class="text-gray-600">Setup Time:</span>
                    <span class="font-semibold text-red-600">${this.getEstimatedTime()}</span>
                </div>
            `;
        }
    }

    renderPreviewSection(title, items, icon, color) {
        if (items.length === 0) return '';

        return `
            <div class="mb-4">
                <h4 class="font-semibold text-gray-700 mb-2 flex items-center">
                    <i class="fas ${icon} ${color} mr-2"></i>${title}
                </h4>
                <div class="flex flex-wrap gap-1">
                    ${items.map(item => `
                        <span class="inline-block bg-blue-500 text-white text-xs px-2 py-1 rounded-full">
                            ${item}
                        </span>
                    `).join('')}
                </div>
            </div>
        `;
    }

    updateConfigSummary() {
        const summary = document.getElementById('configSummary');
        const estimates = document.getElementById('installEstimates');

        if (summary && this.currentStep === 6) {
            const enabledLanguages = this.getEnabledLanguages();
            const enabledTools = this.getEnabledTools();
            const systemFeatures = this.getEnabledSystemFeatures();
            const extraPackages = this.getEnabledPackages();

            summary.innerHTML = `
                <div style="display: flex; flex-direction: column; gap: 10px;">
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                        <span>Programming Languages</span>
                        <span style="background: #007acc; color: white; padding: 3px 8px; border-radius: 12px; font-size: 12px;">
                            ${enabledLanguages.length} selected
                        </span>
                    </div>
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                        <span>Development Tools</span>
                        <span style="background: #28a745; color: white; padding: 3px 8px; border-radius: 12px; font-size: 12px;">
                            ${enabledTools.length} selected
                        </span>
                    </div>
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                        <span>Extra Packages</span>
                        <span style="background: #f0ad4e; color: white; padding: 3px 8px; border-radius: 12px; font-size: 12px;">
                            ${extraPackages.length} packages
                        </span>
                    </div>
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                        <span>System Settings</span>
                        <span style="background: #6f42c1; color: white; padding: 3px 8px; border-radius: 12px; font-size: 12px;">
                            ${systemFeatures.length} configured
                        </span>
                    </div>
                </div>
            `;
        }

        if (estimates && this.currentStep === 6) {
            estimates.innerHTML = `
                Estimated install time: <strong>${this.getEstimatedTime()}</strong><br>
                Estimated disk usage: <strong>${this.getEstimatedSize()}</strong>
            `;
        }
    }

    getEnabledLanguages() {
        return Object.keys(this.config.development.languages)
            .filter(key => this.config.development.languages[key])
            .map(key => this.languages.find(l => l.key === key)?.name || key);
    }

    getEnabledTools() {
        const tools = Object.keys(this.config.tools)
            .filter(key => this.config.tools[key])
            .map(key => this.tools.find(t => t.key === key)?.name || key);

        if (this.config.tools.modern_cli) tools.push('Modern CLI Bundle');
        if (this.config.shell.zsh) tools.push('Zsh + Oh My Zsh');

        return tools;
    }

    getEnabledSystemFeatures() {
        const features = [];

        // Git features
        if (this.config.development.git.pull_rebase) features.push('Git Pull Rebase');
        if (this.config.development.git.default_branch !== 'main') features.push(`Default Branch: ${this.config.development.git.default_branch}`);

        // Development features
        if (this.config.development.aliases.git) features.push('Git Aliases');
        if (this.config.development.aliases.docker) features.push('Docker Shortcuts');
        if (this.config.development.aliases.system) features.push('System Shortcuts');
        if (this.config.development.shell.theme !== 'powerlevel10k') features.push('Custom Shell Theme');

        // System features
        if (this.config.system.appearance.dark_mode) features.push('Dark Mode');
        if (this.config.system.appearance.enable_24_hour_time) features.push('24-Hour Time');
        if (this.config.system.dock.autohide) features.push('Auto-hide Dock');
        if (this.config.system.dock.position !== 'bottom') features.push(`Dock ${this.config.system.dock.position}`);

        return features;
    }

    getEnabledPackages() {
        const packages = [];
        if (this.config.packages.extra_brews) packages.push(...this.config.packages.extra_brews);
        if (this.config.packages.extra_casks) packages.push(...this.config.packages.extra_casks);
        if (this.config.packages.npm_globals) packages.push(...this.config.packages.npm_globals);
        return packages;
    }

    getEstimatedSize() {
        const langCount = this.getEnabledLanguages().length;
        const toolCount = this.getEnabledTools().length;
        const packageCount = this.getEnabledPackages().length;
        const size = (langCount * 200) + (toolCount * 50) + (packageCount * 30) + 300;
        return size > 1000 ? `${(size/1000).toFixed(1)}GB` : `${size}MB`;
    }

    getEstimatedTime() {
        const items = this.getEnabledLanguages().length + this.getEnabledTools().length;
        const minutes = Math.max(5, items * 2);
        return `${minutes} minutes`;
    }

    async saveConfiguration() {
        const loadingOverlay = document.getElementById('loadingOverlay');
        if (loadingOverlay) {
            loadingOverlay.style.display = 'flex';
        }

        try {
            // Update all configs before saving
            this.updatePersonalInfo();
            this.updateGitConfig();
            this.updateSshConfig();
            this.updateDevelopmentConfig();
            this.updateSystemConfig();

            const response = await fetch('/api/save', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(this.config)
            });

            if (response.ok) {
                this.showNotification('🎉 Configuration saved successfully!', 'success');
                setTimeout(() => {
                    this.showNotification('💻 Run "dotfiles install" in your terminal to apply changes', 'info');
                }, 2000);
                setTimeout(() => {
                    this.showNotification('🔧 Your development environment will be fully configured!', 'info');
                }, 4000);
            } else {
                const error = await response.text();
                throw new Error(error || 'Failed to save configuration');
            }
        } catch (error) {
            console.error('Save error:', error);
            this.showNotification('❌ Error saving configuration: ' + error.message, 'error');
        } finally {
            if (loadingOverlay) {
                setTimeout(() => {
                    loadingOverlay.style.display = 'none';
                }, 1000);
            }
        }
    }

    animateOnLoad() {
        // Simplified - no animations
    }

    animateStepTransition() {
        // Simplified - no animations
    }

    updateGitConfig() {
        this.config.development.git.default_branch = document.getElementById('gitDefaultBranch')?.value || 'main';
        this.config.development.git.pull_rebase = document.getElementById('gitPullRebase')?.checked || false;
        this.config.development.git.push_default = document.getElementById('gitPushDefault')?.value || 'simple';
        this.updatePreview();
    }

    updateSshConfig() {
        // SSH config is not part of the current Go struct, so we'll skip this for now
        // or integrate it into a different part of the config if needed
        this.updatePreview();
    }

    updateDevelopmentConfig() {
        this.config.development.shell.theme = document.getElementById('terminalTheme')?.value || 'powerlevel10k';
        this.config.development.shell.terminal_theme = document.getElementById('terminalFont')?.value || 'dark';

        this.config.development.aliases.git = document.getElementById('gitAliases')?.checked || false;
        this.config.development.aliases.docker = document.getElementById('dockerAliases')?.checked || false;
        this.config.development.aliases.system = document.getElementById('systemAliases')?.checked || false;

        this.config.installation.homebrew = document.getElementById('homebrewBundle')?.checked || true;
        this.config.installation.npm_packages = document.getElementById('npmGlobals')?.checked || true;

        this.updatePreview();
    }

    updateIdeExtensions() {
        const editor = this.config.personal.editor;
        const container = document.getElementById('ideExtensions');
        if (!container || !this.ideExtensions[editor]) return;

        container.innerHTML = this.ideExtensions[editor].map(ext => `
            <label class="flex items-center justify-between">
                <div>
                    <div class="font-medium text-gray-800">${ext.name}</div>
                    <div class="text-sm text-gray-600">${ext.description}</div>
                </div>
                <input type="checkbox" class="h-5 w-5 text-blue-600 rounded"
                       onchange="updateIdeExtension('${ext.name}', this.checked)">
            </label>
        `).join('');

        // Setup global function for extension toggling
        window.updateIdeExtension = (name, checked) => {
            this.config.development.ide_extensions[name] = checked;
            this.updatePreview();
        };
    }

    generateGpgKey() {
        this.showNotification('🔑 GPG key generation guide will be provided in terminal', 'info');
        // This would trigger backend GPG key generation
    }

    generateSshKey() {
        this.showNotification('🔐 SSH key generation will be handled during installation', 'info');
        // This would trigger backend SSH key generation
    }

    addExistingKey() {
        this.showNotification('📁 Existing SSH key import will be available in terminal', 'info');
        // This would open a file dialog or provide instructions
    }

    showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        const colors = {
            success: 'bg-green-500 text-white',
            error: 'bg-red-500 text-white',
            info: 'bg-blue-500 text-white',
            warning: 'bg-yellow-500 text-black'
        };

        notification.className = `fixed top-4 right-4 p-4 rounded-lg shadow-lg z-50 transition-all transform translate-x-full ${colors[type] || colors.info}`;
        notification.textContent = message;
        document.body.appendChild(notification);

        // Slide in
        setTimeout(() => notification.classList.remove('translate-x-full'), 100);

        // Slide out and remove
        setTimeout(() => {
            notification.classList.add('translate-x-full');
            setTimeout(() => notification.remove(), 300);
        }, 5000);
    }

    setupPackageInputs() {
        // Setup package tabs
        document.querySelectorAll('.package-tab').forEach(tab => {
            tab.addEventListener('click', (e) => {
                this.switchPackageTab(e.target.dataset.tab);
            });
        });

        // Setup custom package inputs
        const customBrews = document.getElementById('customBrews');
        const customCasks = document.getElementById('customCasks');
        const customNpmGlobals = document.getElementById('customNpmGlobals');

        if (customBrews) {
            customBrews.addEventListener('input', () => this.updateCustomPackages());
        }
        if (customCasks) {
            customCasks.addEventListener('input', () => this.updateCustomPackages());
        }
        if (customNpmGlobals) {
            customNpmGlobals.addEventListener('input', () => this.updateCustomPackages());
        }
    }

    setupInstallationInputs() {
        const installInputs = [
            'installHomebrew', 'installBrewfile', 'installDotfiles',
            'installMacosDefaults', 'installNpmPackages'
        ];
        installInputs.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.addEventListener('change', () => this.updateInstallationConfig());
            }
        });
    }

    renderPackages() {
        this.switchPackageTab('development');
    }

    switchPackageTab(tab) {
        this.currentPackageTab = tab;

        // Update tab appearance
        document.querySelectorAll('.package-tab').forEach(t => {
            t.classList.remove('active');
            t.style.borderBottom = 'none';
        });
        const activeTab = document.querySelector(`[data-tab="${tab}"]`);
        if (activeTab) {
            activeTab.classList.add('active');
            activeTab.style.borderBottom = '2px solid #007acc';
        }

        // Render content
        const content = document.getElementById('packageTabContent');
        if (!content) return;

        switch (tab) {
            case 'development':
                content.innerHTML = this.renderPackageGrid(this.filterPackageList(this.popularPackages.development), 'brew');
                break;
            case 'utilities':
                content.innerHTML = this.renderPackageGrid(this.filterPackageList(this.popularPackages.utilities), 'brew');
                break;
            case 'apps':
                content.innerHTML = this.renderPackageGrid(this.filterPackageList(this.popularCasks.development.concat(this.popularCasks.productivity, this.popularCasks.design)), 'cask');
                break;
            case 'npm':
                content.innerHTML = this.renderPackageGrid(this.filterPackageList(this.npmPackages), 'npm');
                break;
            case 'media':
                content.innerHTML = this.renderPackageGrid(this.filterPackageList(this.popularPackages.media), 'brew');
                break;
        }
    }

    renderPackageGrid(packages, type) {
        return `
            <div class="package-grid">
                ${packages.map(pkg => {
                    const id = type === 'npm' ? pkg.name : (pkg.brew || pkg.cask);
                    const isSelected = this.selectedPackages[type === 'cask' ? 'casks' : (type === 'npm' ? 'npm' : 'brews')].has(id);
                    return `
                        <div class="package-card ${isSelected ? 'selected' : ''}" onclick="togglePackage('${id}', '${type}')">
                            <div style="display: flex; align-items: center; margin-bottom: 8px;">
                                <input type="checkbox" ${isSelected ? 'checked' : ''} style="margin-right: 10px;" onchange="event.stopPropagation()">
                                <strong style="flex: 1;">${pkg.name}</strong>
                                <span class="package-category">${pkg.category}</span>
                            </div>
                            <p style="margin: 0; font-size: 13px; color: #666; line-height: 1.4;">${pkg.description}</p>
                            <p style="margin: 5px 0 0 0; font-size: 11px; color: #999;">brew install ${type === 'cask' ? '--cask ' : ''}${id}</p>
                        </div>
                    `;
                }).join('')}
            </div>
        `;
    }

    updatePackageConfig() {
        // Combine selected packages with custom ones
        const customBrews = document.getElementById('customBrews')?.value || '';
        const customCasks = document.getElementById('customCasks')?.value || '';
        const customNpmGlobals = document.getElementById('customNpmGlobals')?.value || '';

        const customBrewList = customBrews.split(',').map(s => s.trim()).filter(s => s.length > 0);
        const customCaskList = customCasks.split(',').map(s => s.trim()).filter(s => s.length > 0);
        const npmList = customNpmGlobals.split(',').map(s => s.trim()).filter(s => s.length > 0);

        this.config.packages.extra_brews = [...this.selectedPackages.brews, ...customBrewList];
        this.config.packages.extra_casks = [...this.selectedPackages.casks, ...customCaskList];
        this.config.packages.npm_globals = [...this.selectedPackages.npm, ...npmList];
        this.updatePreview();
    }

    updateCustomPackages() {
        this.updatePackageConfig();
    }

    updateInstallationConfig() {
        this.config.installation.homebrew = document.getElementById('installHomebrew')?.checked || false;
        this.config.installation.brewfile = document.getElementById('installBrewfile')?.checked || false;
        this.config.installation.dotfiles = document.getElementById('installDotfiles')?.checked || false;
        this.config.installation.macos_defaults = document.getElementById('installMacosDefaults')?.checked || false;
        this.config.installation.npm_packages = document.getElementById('installNpmPackages')?.checked || false;
        this.updatePreview();
        this.updateInstallationConfig();
    }

    togglePackage(id, type) {
        const packageSet = type === 'cask' ? this.selectedPackages.casks :
                          type === 'npm' ? this.selectedPackages.npm :
                          this.selectedPackages.brews;

        if (packageSet.has(id)) {
            packageSet.delete(id);
        } else {
            packageSet.add(id);
        }

        // Re-render current tab to update selection
        this.switchPackageTab(this.currentPackageTab);
        this.updatePackageConfig();
    }

    filterPackages(searchTerm) {
        this.currentSearchTerm = searchTerm.toLowerCase();
        this.switchPackageTab(this.currentPackageTab);
    }

    filterPackageList(packages) {
        if (!this.currentSearchTerm) return packages;

        return packages.filter(pkg => {
            const searchText = `${pkg.name} ${pkg.description} ${pkg.category}`.toLowerCase();
            return searchText.includes(this.currentSearchTerm);
        });
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.wizard = new SetupWizard();

    // Make togglePackage globally accessible
    window.togglePackage = (id, type) => {
        window.wizard.togglePackage(id, type);
    };
});