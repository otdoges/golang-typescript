/**
 * Utility functions for TypeScript-Golang package
 */

import fs from 'fs-extra';
import path from 'path';
import { execSync } from 'child_process';

/**
 * File system utilities
 */
export class FileUtils {
  /**
   * Check if Go is installed
   */
  static isGoInstalled(): boolean {
    try {
      execSync('go version', { stdio: 'ignore' });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Check if gofmt is available
   */
  static isGofmtAvailable(): boolean {
    try {
      execSync('gofmt -help', { stdio: 'ignore' });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Find Go files in a directory
   */
  static async findGoFiles(directory: string): Promise<string[]> {
    const files: string[] = [];
    
    const scan = async (dir: string) => {
      const entries = await fs.readdir(dir, { withFileTypes: true });
      
      for (const entry of entries) {
        const fullPath = path.join(dir, entry.name);
        
        if (entry.isDirectory() && !entry.name.startsWith('.')) {
          await scan(fullPath);
        } else if (entry.isFile() && entry.name.endsWith('.go')) {
          files.push(fullPath);
        }
      }
    };
    
    await scan(directory);
    return files;
  }

  /**
   * Find TypeScript files in a directory
   */
  static async findTypeScriptFiles(directory: string): Promise<string[]> {
    const files: string[] = [];
    
    const scan = async (dir: string) => {
      const entries = await fs.readdir(dir, { withFileTypes: true });
      
      for (const entry of entries) {
        const fullPath = path.join(dir, entry.name);
        
        if (entry.isDirectory() && !entry.name.startsWith('.') && entry.name !== 'node_modules') {
          await scan(fullPath);
        } else if (entry.isFile() && (entry.name.endsWith('.ts') || entry.name.endsWith('.tsx'))) {
          files.push(fullPath);
        }
      }
    };
    
    await scan(directory);
    return files;
  }

  /**
   * Copy template files with variable substitution
   */
  static async copyTemplate(
    templatePath: string, 
    targetPath: string, 
    variables: Record<string, string> = {}
  ): Promise<void> {
    const stat = await fs.stat(templatePath);
    
    if (stat.isDirectory()) {
      await fs.ensureDir(targetPath);
      const entries = await fs.readdir(templatePath);
      
      for (const entry of entries) {
        const srcPath = path.join(templatePath, entry);
        const destPath = path.join(targetPath, entry);
        await this.copyTemplate(srcPath, destPath, variables);
      }
    } else {
      let content = await fs.readFile(templatePath, 'utf8');
      
      // Replace variables
      for (const [key, value] of Object.entries(variables)) {
        const regex = new RegExp(key, 'g');
        content = content.replace(regex, value);
      }
      
      await fs.writeFile(targetPath, content);
    }
  }

  /**
   * Create a backup of a file
   */
  static async backupFile(filePath: string): Promise<string> {
    const backupPath = `${filePath}.backup.${Date.now()}`;
    await fs.copy(filePath, backupPath);
    return backupPath;
  }

  /**
   * Restore a file from backup
   */
  static async restoreFromBackup(originalPath: string, backupPath: string): Promise<void> {
    await fs.copy(backupPath, originalPath);
    await fs.remove(backupPath);
  }
}

/**
 * Project utilities
 */
export class ProjectUtils {
  /**
   * Check if directory is a Go project
   */
  static async isGoProject(directory: string): Promise<boolean> {
    const goModPath = path.join(directory, 'go.mod');
    return await fs.pathExists(goModPath);
  }

  /**
   * Check if directory is a TypeScript project
   */
  static async isTypeScriptProject(directory: string): Promise<boolean> {
    const packageJsonPath = path.join(directory, 'package.json');
    const tsconfigPath = path.join(directory, 'tsconfig.json');
    
    return (await fs.pathExists(packageJsonPath)) || (await fs.pathExists(tsconfigPath));
  }

  /**
   * Get project name from go.mod
   */
  static async getGoModuleName(directory: string): Promise<string | null> {
    const goModPath = path.join(directory, 'go.mod');
    
    if (!(await fs.pathExists(goModPath))) {
      return null;
    }
    
    const content = await fs.readFile(goModPath, 'utf8');
    const match = content.match(/^module\s+(.+)$/m);
    
    return match ? match[1].trim() : null;
  }

  /**
   * Initialize a new Go module
   */
  static async initGoModule(directory: string, moduleName: string): Promise<void> {
    const goModPath = path.join(directory, 'go.mod');
    
    if (await fs.pathExists(goModPath)) {
      throw new Error('go.mod already exists');
    }
    
    const content = `module ${moduleName}\n\ngo 1.21\n`;
    await fs.writeFile(goModPath, content);
  }

  /**
   * Get project dependencies
   */
  static async getProjectDependencies(directory: string): Promise<{
    go: string[];
    npm: string[];
  }> {
    const result = { go: [] as string[], npm: [] as string[] };
    
    // Get Go dependencies
    const goModPath = path.join(directory, 'go.mod');
    if (await fs.pathExists(goModPath)) {
      const content = await fs.readFile(goModPath, 'utf8');
      const requireSection = content.match(/require\s*\(([\s\S]*?)\)/);
      
      if (requireSection) {
        const lines = requireSection[1].split('\n');
        for (const line of lines) {
          const match = line.trim().match(/^([^\s]+)/);
          if (match && !match[1].startsWith('//')) {
            result.go.push(match[1]);
          }
        }
      }
    }
    
    // Get NPM dependencies
    const packageJsonPath = path.join(directory, 'package.json');
    if (await fs.pathExists(packageJsonPath)) {
      const packageJson = JSON.parse(await fs.readFile(packageJsonPath, 'utf8'));
      
      if (packageJson.dependencies) {
        result.npm.push(...Object.keys(packageJson.dependencies));
      }
      
      if (packageJson.devDependencies) {
        result.npm.push(...Object.keys(packageJson.devDependencies));
      }
    }
    
    return result;
  }
}

/**
 * Code analysis utilities
 */
export class CodeAnalysis {
  /**
   * Extract interfaces from TypeScript code
   */
  static extractInterfaces(tsCode: string): Array<{
    name: string;
    fields: Array<{ name: string; type: string; optional: boolean }>;
  }> {
    const interfaces: Array<{
      name: string;
      fields: Array<{ name: string; type: string; optional: boolean }>;
    }> = [];
    
    const interfaceRegex = /interface\s+(\w+)\s*\{([^}]*)\}/g;
    let match;
    
    while ((match = interfaceRegex.exec(tsCode)) !== null) {
      const [, name, body] = match;
      const fields: Array<{ name: string; type: string; optional: boolean }> = [];
      
      const fieldLines = body.split('\n').map(line => line.trim()).filter(line => line);
      
      for (const line of fieldLines) {
        if (line.includes(':')) {
          const parts = line.replace(/[;,]$/, '').split(':');
          if (parts.length === 2) {
            const fieldName = parts[0].trim();
            const fieldType = parts[1].trim();
            const optional = fieldName.endsWith('?');
            
            fields.push({
              name: optional ? fieldName.slice(0, -1) : fieldName,
              type: fieldType,
              optional
            });
          }
        }
      }
      
      interfaces.push({ name, fields });
    }
    
    return interfaces;
  }

  /**
   * Extract classes from TypeScript code
   */
  static extractClasses(tsCode: string): Array<{
    name: string;
    extends?: string;
    methods: Array<{ name: string; parameters: string[]; returnType?: string }>;
    properties: Array<{ name: string; type: string }>;
  }> {
    const classes: Array<{
      name: string;
      extends?: string;
      methods: Array<{ name: string; parameters: string[]; returnType?: string }>;
      properties: Array<{ name: string; type: string }>;
    }> = [];
    
    const classRegex = /class\s+(\w+)(?:\s+extends\s+(\w+))?\s*\{([^}]*)\}/g;
    let match;
    
    while ((match = classRegex.exec(tsCode)) !== null) {
      const [, name, extendsClass, body] = match;
      const methods: Array<{ name: string; parameters: string[]; returnType?: string }> = [];
      const properties: Array<{ name: string; type: string }> = [];
      
      // This is a simplified parser - real implementation would be more robust
      const lines = body.split('\n').map(line => line.trim()).filter(line => line);
      
      for (const line of lines) {
        if (line.includes('(') && line.includes(')')) {
          // Looks like a method
          const methodMatch = line.match(/(\w+)\s*\(([^)]*)\)(?:\s*:\s*([^{]+))?/);
          if (methodMatch) {
            const [, methodName, params, returnType] = methodMatch;
            methods.push({
              name: methodName,
              parameters: params ? params.split(',').map(p => p.trim()) : [],
              returnType: returnType?.trim()
            });
          }
        } else if (line.includes(':') && !line.includes('(')) {
          // Looks like a property
          const propMatch = line.match(/(\w+)\s*:\s*([^;,]+)/);
          if (propMatch) {
            const [, propName, propType] = propMatch;
            properties.push({ name: propName, type: propType.trim() });
          }
        }
      }
      
      classes.push({
        name,
        extends: extendsClass,
        methods,
        properties
      });
    }
    
    return classes;
  }

  /**
   * Extract enums from TypeScript code
   */
  static extractEnums(tsCode: string): Array<{
    name: string;
    values: Array<{ name: string; value?: string | number }>;
    isStringEnum: boolean;
  }> {
    const enums: Array<{
      name: string;
      values: Array<{ name: string; value?: string | number }>;
      isStringEnum: boolean;
    }> = [];
    
    const enumRegex = /enum\s+(\w+)\s*\{([^}]*)\}/g;
    let match;
    
    while ((match = enumRegex.exec(tsCode)) !== null) {
      const [, name, body] = match;
      const values: Array<{ name: string; value?: string | number }> = [];
      
      const valueLines = body.split(',').map(line => line.trim()).filter(line => line);
      let isStringEnum = false;
      
      for (const line of valueLines) {
        if (line.includes('=')) {
          const [enumName, enumValue] = line.split('=').map(s => s.trim());
          const value = enumValue.startsWith('"') || enumValue.startsWith("'") 
            ? enumValue.slice(1, -1) 
            : parseInt(enumValue);
          
          if (typeof value === 'string') {
            isStringEnum = true;
          }
          
          values.push({ name: enumName, value });
        } else {
          values.push({ name: line });
        }
      }
      
      enums.push({ name, values, isStringEnum });
    }
    
    return enums;
  }
}

/**
 * Build utilities
 */
export class BuildUtils {
  /**
   * Run go build
   */
  static async goBuild(directory: string, outputPath?: string): Promise<{ success: boolean; output: string }> {
    try {
      const command = outputPath ? `go build -o ${outputPath}` : 'go build';
      const output = execSync(command, { 
        cwd: directory, 
        encoding: 'utf8',
        stdio: 'pipe'
      });
      
      return { success: true, output };
    } catch (error: any) {
      return { 
        success: false, 
        output: error.message || error.toString() 
      };
    }
  }

  /**
   * Run go test
   */
  static async goTest(directory: string): Promise<{ success: boolean; output: string }> {
    try {
      const output = execSync('go test ./...', { 
        cwd: directory, 
        encoding: 'utf8',
        stdio: 'pipe'
      });
      
      return { success: true, output };
    } catch (error: any) {
      return { 
        success: false, 
        output: error.message || error.toString() 
      };
    }
  }

  /**
   * Run gofmt
   */
  static async goFormat(directory: string): Promise<{ success: boolean; output: string }> {
    try {
      const output = execSync('gofmt -w .', { 
        cwd: directory, 
        encoding: 'utf8',
        stdio: 'pipe'
      });
      
      return { success: true, output };
    } catch (error: any) {
      return { 
        success: false, 
        output: error.message || error.toString() 
      };
    }
  }

  /**
   * Run go mod tidy
   */
  static async goModTidy(directory: string): Promise<{ success: boolean; output: string }> {
    try {
      const output = execSync('go mod tidy', { 
        cwd: directory, 
        encoding: 'utf8',
        stdio: 'pipe'
      });
      
      return { success: true, output };
    } catch (error: any) {
      return { 
        success: false, 
        output: error.message || error.toString() 
      };
    }
  }
}

/**
 * Validation utilities
 */
export class ValidationUtils {
  /**
   * Validate project name
   */
  static validateProjectName(name: string): { valid: boolean; errors: string[] } {
    const errors: string[] = [];
    
    if (!name) {
      errors.push('Project name is required');
    }
    
    if (name.length < 2) {
      errors.push('Project name must be at least 2 characters');
    }
    
    if (!/^[a-zA-Z][a-zA-Z0-9-_]*$/.test(name)) {
      errors.push('Project name must start with a letter and contain only letters, numbers, hyphens, and underscores');
    }
    
    if (name.includes('..')) {
      errors.push('Project name cannot contain consecutive dots');
    }
    
    return { valid: errors.length === 0, errors };
  }

  /**
   * Validate Go module name
   */
  static validateGoModuleName(name: string): { valid: boolean; errors: string[] } {
    const errors: string[] = [];
    
    if (!name) {
      errors.push('Module name is required');
    }
    
    // Go module names should be valid import paths
    if (!/^[a-zA-Z0-9.\-_/]+$/.test(name)) {
      errors.push('Module name contains invalid characters');
    }
    
    if (name.startsWith('/') || name.endsWith('/')) {
      errors.push('Module name cannot start or end with /');
    }
    
    return { valid: errors.length === 0, errors };
  }

  /**
   * Check if directory is empty or safe to overwrite
   */
  static async validateTargetDirectory(directory: string): Promise<{ 
    valid: boolean; 
    errors: string[]; 
    warnings: string[] 
  }> {
    const errors: string[] = [];
    const warnings: string[] = [];
    
    try {
      const exists = await fs.pathExists(directory);
      
      if (exists) {
        const files = await fs.readdir(directory);
        
        if (files.length > 0) {
          warnings.push(`Directory ${directory} is not empty`);
          
          const importantFiles = ['go.mod', 'package.json', 'main.go'];
          const hasImportantFiles = files.some(file => importantFiles.includes(file));
          
          if (hasImportantFiles) {
            errors.push(`Directory ${directory} contains important files that might be overwritten`);
          }
        }
      }
    } catch (error) {
      errors.push(`Cannot access directory ${directory}: ${error}`);
    }
    
    return { valid: errors.length === 0, errors, warnings };
  }
}

/**
 * Logger utility
 */
export class Logger {
  private static isVerbose = false;

  static setVerbose(verbose: boolean): void {
    this.isVerbose = verbose;
  }

  static info(message: string): void {
    console.log(`ℹ️  ${message}`);
  }

  static success(message: string): void {
    console.log(`✅ ${message}`);
  }

  static warning(message: string): void {
    console.log(`⚠️  ${message}`);
  }

  static error(message: string): void {
    console.error(`❌ ${message}`);
  }

  static debug(message: string): void {
    if (this.isVerbose) {
      console.log(`🐛 ${message}`);
    }
  }

  static step(message: string): void {
    console.log(`📝 ${message}`);
  }

  static progress(message: string): void {
    if (this.isVerbose) {
      console.log(`⚙️  ${message}`);
    }
  }
}