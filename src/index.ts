/**
 * TypeScript-Golang: A comprehensive Go implementation that provides TypeScript-like language features
 * @packageDocumentation
 */

export * from './cli';
export * from './generator';
export * from './utils';

// Package information
export const version = require('../package.json').version;
export const name = 'typescript-golang';

/**
 * Main package exports for programmatic usage
 */
export interface TypeScriptGolangOptions {
  projectName?: string;
  template?: string;
  outputDir?: string;
  verbose?: boolean;
}

/**
 * Create a new TypeScript-Go project programmatically
 */
export async function createProject(options: TypeScriptGolangOptions): Promise<void> {
  const { initProject } = await import('./cli');
  await initProject(options.projectName || 'my-project', options);
}

/**
 * Generate Go code from TypeScript
 */
export async function generateGoCode(inputFile: string, outputFile?: string): Promise<void> {
  const { generateGoCode } = await import('./cli');
  await generateGoCode(inputFile, { output: outputFile });
}

/**
 * Get available project templates
 */
export async function getTemplates(): Promise<string[]> {
  const fs = await import('fs-extra');
  const path = await import('path');
  
  const templatesDir = path.join(__dirname, '../templates');
  
  try {
    const templates = await fs.readdir(templatesDir);
    return templates.filter(async (template) => {
      const templatePath = path.join(templatesDir, template);
      const stat = await fs.stat(templatePath);
      return stat.isDirectory();
    });
  } catch {
    return [];
  }
}

/**
 * Utility functions for working with Go code
 */
export class GoUtils {
  /**
   * Format Go code (requires gofmt to be installed)
   */
  static async formatCode(code: string): Promise<string> {
    const { execSync } = await import('child_process');
    const fs = await import('fs-extra');
    const path = await import('path');
    const os = await import('os');
    
    const tempFile = path.join(os.tmpdir(), `temp-${Date.now()}.go`);
    
    try {
      await fs.writeFile(tempFile, code);
      const formatted = execSync(`gofmt ${tempFile}`, { encoding: 'utf8' });
      return formatted;
    } catch (error) {
      // If gofmt fails, return original code
      return code;
    } finally {
      try {
        await fs.unlink(tempFile);
      } catch {
        // Ignore cleanup errors
      }
    }
  }

  /**
   * Validate Go syntax (requires go to be installed)
   */
  static async validateSyntax(code: string): Promise<{ valid: boolean; errors?: string }> {
    const { execSync } = await import('child_process');
    const fs = await import('fs-extra');
    const path = await import('path');
    const os = await import('os');
    
    const tempFile = path.join(os.tmpdir(), `validate-${Date.now()}.go`);
    
    try {
      await fs.writeFile(tempFile, code);
      execSync(`go build ${tempFile}`, { encoding: 'utf8' });
      return { valid: true };
    } catch (error: any) {
      return { 
        valid: false, 
        errors: error.message || 'Syntax validation failed'
      };
    } finally {
      try {
        await fs.unlink(tempFile);
      } catch {
        // Ignore cleanup errors
      }
    }
  }

  /**
   * Convert TypeScript type to Go type
   */
  static convertType(tsType: string): string {
    const typeMap: Record<string, string> = {
      'string': 'string',
      'number': 'float64',
      'boolean': 'bool',
      'void': '',
      'any': 'interface{}',
      'object': 'interface{}',
      'null': 'nil',
      'undefined': 'nil',
      'string[]': '[]string',
      'number[]': '[]float64',
      'boolean[]': '[]bool',
      'Array<string>': '[]string',
      'Array<number>': '[]float64',
      'Array<boolean>': '[]bool',
      'Promise<string>': '*async.Promise[string]',
      'Promise<number>': '*async.Promise[float64]',
      'Promise<boolean>': '*async.Promise[bool]',
      'Map<string, string>': '*types.Map[string, string]',
      'Set<string>': '*types.Set[string]',
      'Optional<string>': 'types.Optional[string]',
      'Optional<number>': 'types.Optional[float64]',
      'Optional<boolean>': 'types.Optional[bool]',
    };
    
    return typeMap[tsType] || tsType;
  }
}

/**
 * Template utilities
 */
export class TemplateUtils {
  /**
   * Get template information
   */
  static async getTemplateInfo(templateName: string): Promise<any> {
    const fs = await import('fs-extra');
    const path = await import('path');
    
    const templatePath = path.join(__dirname, '../templates', templateName, 'template.json');
    
    try {
      const content = await fs.readFile(templatePath, 'utf8');
      return JSON.parse(content);
    } catch {
      return null;
    }
  }

  /**
   * List all available templates with their info
   */
  static async listTemplatesWithInfo(): Promise<Array<{ name: string; info: any }>> {
    const templates = await getTemplates();
    const templatesWithInfo = [];
    
    for (const template of templates) {
      const info = await this.getTemplateInfo(template);
      templatesWithInfo.push({ name: template, info });
    }
    
    return templatesWithInfo;
  }
}

/**
 * Code generation utilities
 */
export class CodeGenerator {
  /**
   * Generate a complete Go project structure
   */
  static async generateProject(config: {
    name: string;
    features: string[];
    template?: string;
  }): Promise<string[]> {
    const generatedFiles: string[] = [];
    
    // This would contain more sophisticated code generation logic
    // For now, we'll return a basic structure
    
    generatedFiles.push('main.go');
    generatedFiles.push('go.mod');
    
    if (config.features.includes('web-api')) {
      generatedFiles.push('handlers.go');
      generatedFiles.push('routes.go');
    }
    
    if (config.features.includes('cli')) {
      generatedFiles.push('cmd/root.go');
      generatedFiles.push('cmd/commands.go');
    }
    
    return generatedFiles;
  }

  /**
   * Generate TypeScript type definitions for Go structs
   */
  static generateTypeDefinitions(goCode: string): string {
    // Basic implementation - would be more sophisticated in practice
    const lines = goCode.split('\n');
    let tsDefinitions = '';
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i].trim();
      
      if (line.startsWith('type ') && line.includes('struct')) {
        const typeName = line.split(' ')[1];
        tsDefinitions += `interface ${typeName} {\n`;
        
        // Look for struct fields
        i++;
        while (i < lines.length && !lines[i].includes('}')) {
          const fieldLine = lines[i].trim();
          if (fieldLine && !fieldLine.startsWith('//')) {
            // Parse field (simplified)
            const parts = fieldLine.split(' ');
            if (parts.length >= 2) {
              const fieldName = parts[0].toLowerCase();
              const fieldType = this.convertGoTypeToTypeScript(parts[1]);
              tsDefinitions += `  ${fieldName}: ${fieldType};\n`;
            }
          }
          i++;
        }
        
        tsDefinitions += '}\n\n';
      }
    }
    
    return tsDefinitions;
  }

  private static convertGoTypeToTypeScript(goType: string): string {
    const typeMap: Record<string, string> = {
      'string': 'string',
      'int': 'number',
      'int64': 'number', 
      'float64': 'number',
      'bool': 'boolean',
      '[]string': 'string[]',
      '[]int': 'number[]',
      '[]float64': 'number[]',
      '[]bool': 'boolean[]',
      'interface{}': 'any',
    };
    
    return typeMap[goType] || 'any';
  }
}

// Default export
export default {
  createProject,
  generateGoCode,
  getTemplates,
  GoUtils,
  TemplateUtils,
  CodeGenerator,
  version,
  name
};