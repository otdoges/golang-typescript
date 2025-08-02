/**
 * TypeScript definitions for typescript-golang package
 */

export interface TypeScriptGolangOptions {
  projectName?: string;
  template?: string;
  outputDir?: string;
  verbose?: boolean;
}

export interface ConversionOptions {
  preserveComments?: boolean;
  generateTests?: boolean;
  addTypeAnnotations?: boolean;
  outputPackage?: string;
}

export interface ConversionResult {
  goCode: string;
  errors: string[];
  warnings: string[];
  generatedFiles: string[];
}

export interface ProjectTemplate {
  name: string;
  description: string;
  author: string;
  version: string;
  features: string[];
  dependencies?: string[];
}

export interface BuildResult {
  success: boolean;
  output: string;
}

export interface ValidationResult {
  valid: boolean;
  errors: string[];
  warnings?: string[];
}

// CLI Functions
export declare function createProject(options: TypeScriptGolangOptions): Promise<void>;
export declare function generateGoCode(inputFile: string, outputFile?: string): Promise<void>;
export declare function getTemplates(): Promise<string[]>;

// Main Classes
export declare class GoUtils {
  static formatCode(code: string): Promise<string>;
  static validateSyntax(code: string): Promise<{ valid: boolean; errors?: string }>;
  static convertType(tsType: string): string;
}

export declare class TemplateUtils {
  static getTemplateInfo(templateName: string): Promise<ProjectTemplate | null>;
  static listTemplatesWithInfo(): Promise<Array<{ name: string; info: ProjectTemplate | null }>>;
}

export declare class CodeGenerator {
  static generateProject(config: {
    name: string;
    features: string[];
    template?: string;
  }): Promise<string[]>;
  static generateTypeDefinitions(goCode: string): string;
}

export declare class TypeScriptToGoGenerator {
  constructor(options?: ConversionOptions);
  convertCode(tsCode: string): Promise<ConversionResult>;
}

export declare class FileUtils {
  static isGoInstalled(): boolean;
  static isGofmtAvailable(): boolean;
  static findGoFiles(directory: string): Promise<string[]>;
  static findTypeScriptFiles(directory: string): Promise<string[]>;
  static copyTemplate(
    templatePath: string,
    targetPath: string,
    variables?: Record<string, string>
  ): Promise<void>;
  static backupFile(filePath: string): Promise<string>;
  static restoreFromBackup(originalPath: string, backupPath: string): Promise<void>;
}

export declare class ProjectUtils {
  static isGoProject(directory: string): Promise<boolean>;
  static isTypeScriptProject(directory: string): Promise<boolean>;
  static getGoModuleName(directory: string): Promise<string | null>;
  static initGoModule(directory: string, moduleName: string): Promise<void>;
  static getProjectDependencies(directory: string): Promise<{
    go: string[];
    npm: string[];
  }>;
}

export declare class CodeAnalysis {
  static extractInterfaces(tsCode: string): Array<{
    name: string;
    fields: Array<{ name: string; type: string; optional: boolean }>;
  }>;
  static extractClasses(tsCode: string): Array<{
    name: string;
    extends?: string;
    methods: Array<{ name: string; parameters: string[]; returnType?: string }>;
    properties: Array<{ name: string; type: string }>;
  }>;
  static extractEnums(tsCode: string): Array<{
    name: string;
    values: Array<{ name: string; value?: string | number }>;
    isStringEnum: boolean;
  }>;
}

export declare class BuildUtils {
  static goBuild(directory: string, outputPath?: string): Promise<BuildResult>;
  static goTest(directory: string): Promise<BuildResult>;
  static goFormat(directory: string): Promise<BuildResult>;
  static goModTidy(directory: string): Promise<BuildResult>;
}

export declare class ValidationUtils {
  static validateProjectName(name: string): ValidationResult;
  static validateGoModuleName(name: string): ValidationResult;
  static validateTargetDirectory(directory: string): Promise<ValidationResult>;
}

export declare class Logger {
  static setVerbose(verbose: boolean): void;
  static info(message: string): void;
  static success(message: string): void;
  static warning(message: string): void;
  static error(message: string): void;
  static debug(message: string): void;
  static step(message: string): void;
  static progress(message: string): void;
}

// Package constants
export declare const version: string;
export declare const name: string;

// Default export
declare const typeScriptGolang: {
  createProject: typeof createProject;
  generateGoCode: typeof generateGoCode;
  getTemplates: typeof getTemplates;
  GoUtils: typeof GoUtils;
  TemplateUtils: typeof TemplateUtils;
  CodeGenerator: typeof CodeGenerator;
  TypeScriptToGoGenerator: typeof TypeScriptToGoGenerator;
  FileUtils: typeof FileUtils;
  ProjectUtils: typeof ProjectUtils;
  CodeAnalysis: typeof CodeAnalysis;
  BuildUtils: typeof BuildUtils;
  ValidationUtils: typeof ValidationUtils;
  Logger: typeof Logger;
  version: string;
  name: string;
};

export default typeScriptGolang;