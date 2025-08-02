#!/usr/bin/env node

import { Command } from 'commander';
import chalk from 'chalk';
import ora from 'ora';
import fs from 'fs-extra';
import path from 'path';
import { execSync } from 'child_process';

const program = new Command();

// Package info
const packageJson = JSON.parse(fs.readFileSync(path.join(__dirname, '../package.json'), 'utf8'));

// Helper functions
function getTemplatesDir(): string {
  return path.join(__dirname, '../templates');
}

function getLibrarySourceDir(): string {
  return path.join(__dirname, '../');
}

// Commands

/**
 * Initialize a new TypeScript-Go project
 */
async function initProject(projectName: string, options: any) {
  const spinner = ora('Creating new TypeScript-Go project...').start();
  
  try {
    const projectPath = path.resolve(projectName);
    
    // Check if directory exists
    if (fs.existsSync(projectPath)) {
      spinner.fail(`Directory ${projectName} already exists!`);
      return;
    }

    // Create project directory
    await fs.ensureDir(projectPath);
    
    // Copy Go library files
    const sourceDir = getLibrarySourceDir();
    const filesToCopy = [
      'types',
      'utils', 
      'async',
      'classes',
      'enums',
      'testing',
      'go.mod'
    ];

    for (const file of filesToCopy) {
      const srcPath = path.join(sourceDir, file);
      const destPath = path.join(projectPath, file);
      
      if (await fs.pathExists(srcPath)) {
        await fs.copy(srcPath, destPath);
      }
    }

    // Copy project template
    const templatePath = path.join(getTemplatesDir(), 'basic-project');
    if (await fs.pathExists(templatePath)) {
      await fs.copy(templatePath, projectPath, { overwrite: false });
    }

    // Update go.mod with project name
    const goModPath = path.join(projectPath, 'go.mod');
    if (await fs.pathExists(goModPath)) {
      let goModContent = await fs.readFile(goModPath, 'utf8');
      goModContent = goModContent.replace('typescript-golang', projectName);
      await fs.writeFile(goModPath, goModContent);
    }

    // Create main.go if it doesn't exist
    const mainGoPath = path.join(projectPath, 'main.go');
    if (!(await fs.pathExists(mainGoPath))) {
      const mainGoContent = `package main

import (
	"fmt"
	"${projectName}/types"
	"${projectName}/utils"
)

func main() {
	fmt.Println("🚀 Welcome to ${projectName}!")
	fmt.Println("TypeScript-like Go implementation ready to use!")
	
	// Example usage
	name := types.Some("TypeScript Developer")
	fmt.Printf("Hello, %s!\\n", name.GetOrDefault("Anonymous"))
	
	numbers := []int{1, 2, 3, 4, 5}
	doubled := utils.Map(numbers, func(x int) int { return x * 2 })
	fmt.Printf("Doubled numbers: %v\\n", doubled)
}
`;
      await fs.writeFile(mainGoPath, mainGoContent);
    }

    // Create README
    const readmePath = path.join(projectPath, 'README.md');
    if (!(await fs.pathExists(readmePath))) {
      const readmeContent = `# ${projectName}

A TypeScript-like Go project created with typescript-golang.

## Getting Started

\`\`\`bash
# Run the project
go run .

# Build the project  
go build -o ${projectName}

# Run tests
go test ./...
\`\`\`

## Features

This project includes the full TypeScript-like Go library:

- ✅ Optional Types
- ✅ Array Utilities
- ✅ Promise/Async Patterns
- ✅ Class-based OOP
- ✅ Enums
- ✅ Union Types
- ✅ String Utilities
- ✅ JSON Handling
- ✅ Collections (Map, Set)
- ✅ Event System
- ✅ Error Handling
- ✅ Testing Framework

## Documentation

See the [TypeScript-Golang Documentation](https://github.com/typescript-golang/typescript-golang) for complete API reference.

## Examples

### Optional Types
\`\`\`go
name := types.Some("John")
age := types.None[int]()

fmt.Println(name.GetOrDefault("Anonymous"))
fmt.Println(age.GetOrDefault(0))
\`\`\`

### Array Operations
\`\`\`go
numbers := []int{1, 2, 3, 4, 5}
doubled := utils.Map(numbers, func(x int) int { return x * 2 })
evens := utils.Filter(numbers, func(x int) bool { return x%2 == 0 })
\`\`\`

### Promises
\`\`\`go
promise := async.NewPromise(func() (string, error) {
    time.Sleep(1 * time.Second)
    return "Hello World", nil
})

result, err := promise.Await()
\`\`\`
`;
      await fs.writeFile(readmePath, readmeContent);
    }

    spinner.succeed(`Created ${chalk.green(projectName)} successfully!`);
    
    console.log(`
${chalk.blue('Next steps:')}
  ${chalk.gray('$')} cd ${projectName}
  ${chalk.gray('$')} go run .
  ${chalk.gray('$')} go mod tidy

${chalk.blue('Available commands:')}
  ${chalk.gray('$')} ts-go --help    ${chalk.dim('# Show all commands')}
  ${chalk.gray('$')} ts-go generate  ${chalk.dim('# Generate Go code from TypeScript')}
  ${chalk.gray('$')} ts-go examples  ${chalk.dim('# Show usage examples')}
`);

  } catch (error) {
    spinner.fail(`Failed to create project: ${error}`);
  }
}

/**
 * Generate Go code from TypeScript
 */
async function generateGoCode(inputFile: string, options: any) {
  const spinner = ora('Generating Go code from TypeScript...').start();
  
  try {
    if (!fs.existsSync(inputFile)) {
      spinner.fail(`Input file ${inputFile} not found!`);
      return;
    }

    const outputFile = options.output || inputFile.replace('.ts', '.go');
    const tsContent = await fs.readFile(inputFile, 'utf8');
    
    // Basic TypeScript to Go conversion patterns
    let goContent = tsContent
      // Convert interface to Go struct
      .replace(/interface\s+(\w+)\s*{([^}]*)}/g, (match, name, body) => {
        const fields = body.trim().split('\n').map(line => {
          const trimmed = line.trim();
          if (trimmed.endsWith(';')) {
            const parts = trimmed.slice(0, -1).split(':');
            if (parts.length === 2) {
              const fieldName = parts[0].trim();
              const fieldType = parts[1].trim();
              const goType = convertTypeScriptTypeToGo(fieldType);
              return `\t${capitalize(fieldName)} ${goType} \`json:"${fieldName}"\``;
            }
          }
          return '';
        }).filter(Boolean).join('\n');
        
        return `type ${name} struct {\n${fields}\n}`;
      })
      // Convert enum to Go constants
      .replace(/enum\s+(\w+)\s*{([^}]*)}/g, (match, name, body) => {
        const values = body.trim().split(',').map(v => v.trim()).filter(Boolean);
        let result = `type ${name} int\n\nconst (\n`;
        values.forEach((value, index) => {
          const parts = value.split('=');
          const enumName = parts[0].trim();
          if (index === 0) {
            result += `\t${enumName} ${name} = iota\n`;
          } else {
            result += `\t${enumName}\n`;
          }
        });
        result += ')';
        return result;
      })
      // Convert function declarations
      .replace(/function\s+(\w+)\s*\(([^)]*)\)\s*:\s*(\w+)/g, (match, name, params, returnType) => {
        const goParams = convertParameters(params);
        const goReturnType = convertTypeScriptTypeToGo(returnType);
        return `func ${name}(${goParams}) ${goReturnType}`;
      })
      // Convert class to struct with methods
      .replace(/class\s+(\w+)\s*{([^}]*)}/g, (match, name, body) => {
        // This is a simplified conversion - real implementation would be more complex
        return `type ${name} struct {\n\t// TODO: Convert class properties\n}\n\n// TODO: Convert class methods`;
      });

    // Add package declaration if not present
    if (!goContent.startsWith('package ')) {
      goContent = 'package main\n\n' + goContent;
    }

    await fs.writeFile(outputFile, goContent);
    
    spinner.succeed(`Generated Go code: ${chalk.green(outputFile)}`);
    
    console.log(`
${chalk.blue('Generated:')} ${outputFile}
${chalk.yellow('Note:')} This is a basic conversion. Manual review and adjustments may be needed.

${chalk.blue('Next steps:')}
  ${chalk.gray('$')} go fmt ${outputFile}
  ${chalk.gray('$')} go build ${outputFile}
`);

  } catch (error) {
    spinner.fail(`Failed to generate Go code: ${error}`);
  }
}

/**
 * Show usage examples
 */
function showExamples() {
  console.log(`
${chalk.blue.bold('TypeScript-Golang Usage Examples')}
${'='.repeat(40)}

${chalk.green('1. Create a new project:')}
  ${chalk.gray('$')} ts-go init my-project
  ${chalk.gray('$')} cd my-project
  ${chalk.gray('$')} go run .

${chalk.green('2. Generate Go from TypeScript:')}
  ${chalk.gray('$')} ts-go generate types.ts --output types.go

${chalk.green('3. Optional Types:')}
  ${chalk.dim('TypeScript:')} const name: string | undefined = getName();
  ${chalk.dim('Go:')}        name := types.Some("John")

${chalk.green('4. Array Operations:')}
  ${chalk.dim('TypeScript:')} const doubled = numbers.map(x => x * 2);
  ${chalk.dim('Go:')}        doubled := utils.Map(numbers, func(x int) int { return x * 2 })

${chalk.green('5. Promises:')}
  ${chalk.dim('TypeScript:')} const result = await fetchData();
  ${chalk.dim('Go:')}        result, err := promise.Await()

${chalk.green('6. Classes:')}
  ${chalk.dim('TypeScript:')} class Person { constructor(public name: string) {} }
  ${chalk.dim('Go:')}        person := NewPerson("Alice")

${chalk.green('7. Events:')}
  ${chalk.dim('TypeScript:')} emitter.on('data', (data) => console.log(data));
  ${chalk.dim('Go:')}        emitter.On("data", func(data string) { fmt.Println(data) })

${chalk.blue('For complete documentation:')}
  https://github.com/typescript-golang/typescript-golang
`);
}

/**
 * List available templates
 */
async function listTemplates() {
  const templatesDir = getTemplatesDir();
  
  try {
    const templates = await fs.readdir(templatesDir);
    
    console.log(`
${chalk.blue.bold('Available Project Templates')}
${'='.repeat(35)}
`);
    
    for (const template of templates) {
      const templatePath = path.join(templatesDir, template);
      const stat = await fs.stat(templatePath);
      
      if (stat.isDirectory()) {
        // Try to read template description
        const descPath = path.join(templatePath, 'template.json');
        let description = 'Basic project template';
        
        if (await fs.pathExists(descPath)) {
          const templateInfo = JSON.parse(await fs.readFile(descPath, 'utf8'));
          description = templateInfo.description || description;
        }
        
        console.log(`  ${chalk.green(template.padEnd(20))} ${chalk.dim(description)}`);
      }
    }
    
    console.log(`
${chalk.blue('Usage:')}
  ${chalk.gray('$')} ts-go init my-project --template <template-name>
`);
    
  } catch (error) {
    console.log(chalk.red('No templates found.'));
  }
}

// Utility functions
function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

function convertTypeScriptTypeToGo(tsType: string): string {
  const typeMap: Record<string, string> = {
    'string': 'string',
    'number': 'float64',
    'boolean': 'bool',
    'void': '',
    'any': 'interface{}',
    'object': 'interface{}',
    'string[]': '[]string',
    'number[]': '[]float64',
    'boolean[]': '[]bool'
  };
  
  return typeMap[tsType] || tsType;
}

function convertParameters(params: string): string {
  if (!params.trim()) return '';
  
  return params.split(',').map(param => {
    const parts = param.trim().split(':');
    if (parts.length === 2) {
      const name = parts[0].trim();
      const type = parts[1].trim();
      const goType = convertTypeScriptTypeToGo(type);
      return `${name} ${goType}`;
    }
    return param.trim();
  }).join(', ');
}

// CLI Setup
program
  .name('ts-go')
  .description('TypeScript-like Go development toolkit')
  .version(packageJson.version);

program
  .command('init <project-name>')
  .description('Initialize a new TypeScript-Go project')
  .option('-t, --template <name>', 'Project template to use', 'basic-project')
  .action(initProject);

program
  .command('generate <input-file>')
  .description('Generate Go code from TypeScript')
  .option('-o, --output <file>', 'Output file path')
  .action(generateGoCode);

program
  .command('examples')
  .description('Show usage examples')
  .action(showExamples);

program
  .command('templates')
  .description('List available project templates')
  .action(listTemplates);

program
  .command('demo')
  .description('Run the TypeScript-Golang demo')
  .action(() => {
    try {
      console.log(chalk.blue('Running TypeScript-Golang demo...'));
      execSync('go run .', { stdio: 'inherit', cwd: __dirname + '/..' });
    } catch (error) {
      console.error(chalk.red('Failed to run demo:'), error);
    }
  });

// Parse CLI arguments
program.parse(process.argv);

// Show help if no command provided
if (!process.argv.slice(2).length) {
  program.outputHelp();
}