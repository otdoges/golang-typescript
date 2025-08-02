/**
 * Code generation utilities for TypeScript to Go conversion
 */

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

/**
 * Advanced TypeScript to Go code generator
 */
export class TypeScriptToGoGenerator {
  private options: ConversionOptions;
  private errors: string[] = [];
  private warnings: string[] = [];
  private generatedFiles: string[] = [];

  constructor(options: ConversionOptions = {}) {
    this.options = {
      preserveComments: true,
      generateTests: false,
      addTypeAnnotations: true,
      outputPackage: 'main',
      ...options
    };
  }

  /**
   * Convert TypeScript code to Go
   */
  async convertCode(tsCode: string): Promise<ConversionResult> {
    this.errors = [];
    this.warnings = [];
    this.generatedFiles = [];

    let goCode = this.processCode(tsCode);
    
    // Add package declaration
    if (!goCode.startsWith('package ')) {
      goCode = `package ${this.options.outputPackage}\n\n${goCode}`;
    }

    // Add necessary imports
    const imports = this.generateImports(goCode);
    if (imports) {
      goCode = goCode.replace(/^package \w+\n\n/, `package ${this.options.outputPackage}\n\n${imports}\n\n`);
    }

    return {
      goCode,
      errors: this.errors,
      warnings: this.warnings,
      generatedFiles: this.generatedFiles
    };
  }

  private processCode(code: string): string {
    let result = code;

    // Convert interfaces to structs
    result = this.convertInterfaces(result);
    
    // Convert classes to structs with methods
    result = this.convertClasses(result);
    
    // Convert enums
    result = this.convertEnums(result);
    
    // Convert functions
    result = this.convertFunctions(result);
    
    // Convert type definitions
    result = this.convertTypeDefinitions(result);
    
    // Convert variable declarations
    result = this.convertVariableDeclarations(result);
    
    // Convert array/object literals
    result = this.convertLiterals(result);
    
    // Convert method calls to Go style
    result = this.convertMethodCalls(result);

    return result;
  }

  private convertInterfaces(code: string): string {
    const interfaceRegex = /interface\s+(\w+)\s*\{([^}]*)\}/g;
    
    return code.replace(interfaceRegex, (_match: string, name: string, body: string) => {
      const fields = this.parseInterfaceFields(body);
      let structFields = '';
      
      fields.forEach(field => {
        const goType = this.convertType(field.type);
        const jsonTag = field.optional ? `,omitempty` : '';
        structFields += `\t${this.capitalize(field.name)} ${goType} \`json:"${field.name}${jsonTag}"\`\n`;
      });

      return `type ${name} struct {\n${structFields}}`;
    });
  }

  private convertClasses(code: string): string {
    const classRegex = /class\s+(\w+)(?:\s+extends\s+(\w+))?\s*\{([^}]*)\}/g;
    
    return code.replace(classRegex, (_match: string, className: string, baseClass: string, body: string) => {
      const members = this.parseClassMembers(body);
      let result = '';

      // Generate struct
      let structFields = '';
      if (baseClass) {
        structFields += `\t*${baseClass}\n`;
      }
      
      members.properties.forEach(prop => {
        const goType = this.convertType(prop.type);
        structFields += `\t${this.capitalize(prop.name)} ${goType} \`json:"${prop.name}"\`\n`;
      });

      result += `type ${className} struct {\n${structFields}}\n\n`;

      // Generate constructor
      if (members.constructor) {
        const params = members.constructor.parameters.map(p => 
          `${p.name} ${this.convertType(p.type)}`
        ).join(', ');
        
        result += `func New${className}(${params}) *${className} {\n`;
        result += `\treturn &${className}{\n`;
        
        if (baseClass) {
          result += `\t\t${baseClass}: New${baseClass}(),\n`;
        }
        
        members.constructor.parameters.forEach(p => {
          result += `\t\t${this.capitalize(p.name)}: ${p.name},\n`;
        });
        
        result += `\t}\n}\n\n`;
      }

      // Generate methods
      members.methods.forEach(method => {
        const params = method.parameters.map(p => 
          `${p.name} ${this.convertType(p.type)}`
        ).join(', ');
        
        const receiver = `(${className.toLowerCase()[0]} *${className})`;
        const returnType = method.returnType ? this.convertType(method.returnType) : '';
        
        result += `func ${receiver} ${this.capitalize(method.name)}(${params})`;
        if (returnType) {
          result += ` ${returnType}`;
        }
        result += ` {\n\t// TODO: Implement method\n`;
        if (returnType && returnType !== '') {
          result += `\tvar zero ${returnType}\n\treturn zero\n`;
        }
        result += `}\n\n`;
      });

      return result;
    });
  }

  private convertEnums(code: string): string {
    const enumRegex = /enum\s+(\w+)\s*\{([^}]*)\}/g;
    
    return code.replace(enumRegex, (_match: string, name: string, body: string) => {
      const values = this.parseEnumValues(body);
      let result = '';

      // Determine if string or numeric enum
      const isStringEnum = values.some(v => v.value && typeof v.value === 'string');

      if (isStringEnum) {
        result += `type ${name} string\n\nconst (\n`;
        values.forEach(value => {
          const enumValue = value.value ? `"${value.value}"` : `"${value.name.toLowerCase()}"`;
          result += `\t${value.name} ${name} = ${enumValue}\n`;
        });
      } else {
        result += `type ${name} int\n\nconst (\n`;
        values.forEach((value, index) => {
          if (index === 0) {
            result += `\t${value.name} ${name} = iota\n`;
          } else {
            result += `\t${value.name}\n`;
          }
        });
      }

      result += ')\n\n';

      // Add String() method for enums
      result += `func (${name.toLowerCase()[0]} ${name}) String() string {\n`;
      result += `\tswitch ${name.toLowerCase()[0]} {\n`;
      values.forEach(value => {
        result += `\tcase ${value.name}:\n\t\treturn "${value.name}"\n`;
      });
      result += `\tdefault:\n\t\treturn "Unknown"\n\t}\n}\n\n`;

      return result;
    });
  }

  private convertFunctions(code: string): string {
    const functionRegex = /function\s+(\w+)\s*\(([^)]*)\)\s*:\s*([^{]+)\s*\{([^}]*)\}/g;
    
    return code.replace(functionRegex, (_match: string, name: string, params: string, returnType: string, body: string) => {
      const goParams = this.convertParameters(params);
      const goReturnType = returnType.trim() !== 'void' ? this.convertType(returnType.trim()) : '';
      const goBody = this.convertFunctionBody(body);

      let result = `func ${this.capitalize(name)}(${goParams})`;
      if (goReturnType) {
        result += ` ${goReturnType}`;
      }
      result += ` {\n${goBody}\n}`;

      return result;
    });
  }

  private convertTypeDefinitions(code: string): string {
    const typeRegex = /type\s+(\w+)\s*=\s*([^;]+);?/g;
    
    return code.replace(typeRegex, (_match: string, name: string, definition: string) => {
      // Handle union types
      if (definition.includes('|')) {
        this.warnings.push(`Union type ${name} converted to interface{} - manual review needed`);
        return `type ${name} interface{}`;
      }
      
      const goType = this.convertType(definition.trim());
      return `type ${name} = ${goType}`;
    });
  }

  private convertVariableDeclarations(code: string): string {
    // Convert const/let/var declarations
    const varRegex = /(const|let|var)\s+(\w+)(?:\s*:\s*([^=]+))?\s*=\s*([^;]+);?/g;
    
    return code.replace(varRegex, (_match: string, _keyword: string, name: string, type: string, value: string) => {
      const goValue = this.convertValue(value.trim());
      
      if (type) {
        const goType = this.convertType(type.trim());
        return `var ${name} ${goType} = ${goValue}`;
      } else {
        return `${name} := ${goValue}`;
      }
    });
  }

  private convertLiterals(code: string): string {
    // Convert object literals
    code = code.replace(/\{([^}]*)\}/g, (_match: string, content: string) => {
      if (content.includes(':')) {
        // Looks like an object literal
        const pairs = content.split(',').map((pair: string) => {
          const [key, value] = pair.split(':').map((s: string) => s.trim());
          return `${key}: ${this.convertValue(value)}`;
        });
        return `{${pairs.join(', ')}}`;
      }
      return _match;
    });

    // Convert array literals
    code = code.replace(/\[([^\]]*)\]/g, (_match: string, content: string) => {
      if (content.trim()) {
        const items = content.split(',').map((item: string) => this.convertValue(item.trim()));
        return `[]interface{}{${items.join(', ')}}`;
      }
      return '[]interface{}{}';
    });

    return code;
  }

  private convertMethodCalls(code: string): string {
    // Convert .map() calls
    code = code.replace(/(\w+)\.map\(([^)]+)\)/g, (_match: string, array: string, fn: string) => {
      return `utils.Map(${array}, ${fn})`;
    });

    // Convert .filter() calls
    code = code.replace(/(\w+)\.filter\(([^)]+)\)/g, (_match: string, array: string, fn: string) => {
      return `utils.Filter(${array}, ${fn})`;
    });

    // Convert .reduce() calls
    code = code.replace(/(\w+)\.reduce\(([^)]+)\)/g, (_match: string, array: string, args: string) => {
      return `utils.Reduce(${array}, ${args})`;
    });

    // Convert Promise constructors
    code = code.replace(/new Promise<([^>]+)>\(([^)]+)\)/g, (_match: string, type: string, executor: string) => {
      const goType = this.convertType(type);
      return `async.NewPromise[${goType}](${executor})`;
    });

    return code;
  }

  private parseInterfaceFields(body: string): Array<{name: string, type: string, optional: boolean}> {
    const fields: Array<{name: string, type: string, optional: boolean}> = [];
    const lines = body.split('\n').map(line => line.trim()).filter(line => line);

    lines.forEach(line => {
      if (line.includes(':')) {
        const parts = line.replace(/[;,]$/, '').split(':');
        if (parts.length === 2) {
          const name = parts[0].trim();
          const type = parts[1].trim();
          const optional = name.endsWith('?');
          
          fields.push({
            name: optional ? name.slice(0, -1) : name,
            type,
            optional
          });
        }
      }
    });

    return fields;
  }

  private parseClassMembers(body: string): {
    properties: Array<{name: string, type: string}>,
    methods: Array<{name: string, parameters: Array<{name: string, type: string}>, returnType?: string}>,
    constructor?: {parameters: Array<{name: string, type: string}>}
  } {
    const result = {
      properties: [] as Array<{name: string, type: string}>,
      methods: [] as Array<{name: string, parameters: Array<{name: string, type: string}>, returnType?: string}>,
      constructor: undefined as {parameters: Array<{name: string, type: string}>} | undefined
    };

    // This is a simplified parser - a real implementation would use a proper TypeScript parser
    const lines = body.split('\n').map(line => line.trim()).filter(line => line);

    lines.forEach(line => {
      if (line.startsWith('constructor(')) {
        const params = this.extractParameters(line);
        result.constructor = { parameters: params };
      } else if (line.includes('(') && line.includes(')')) {
        // Looks like a method
        const methodMatch = line.match(/(\w+)\s*\(([^)]*)\)(?:\s*:\s*([^{]+))?/);
        if (methodMatch) {
          const [, name, params, returnType] = methodMatch;
          result.methods.push({
            name,
            parameters: this.extractParameters(params),
            returnType: returnType?.trim()
          });
        }
      } else if (line.includes(':') && !line.includes('(')) {
        // Looks like a property
        const propMatch = line.match(/(\w+)\s*:\s*([^;,]+)/);
        if (propMatch) {
          const [, name, type] = propMatch;
          result.properties.push({ name, type: type.trim() });
        }
      }
    });

    return result;
  }

  private parseEnumValues(body: string): Array<{name: string, value?: string | number}> {
    const values: Array<{name: string, value?: string | number}> = [];
    const lines = body.split(',').map(line => line.trim()).filter(line => line);

    lines.forEach(line => {
      if (line.includes('=')) {
        const [name, value] = line.split('=').map(s => s.trim());
        values.push({ 
          name, 
          value: value.startsWith('"') ? value.slice(1, -1) : parseInt(value) 
        });
      } else {
        values.push({ name: line });
      }
    });

    return values;
  }

  private extractParameters(paramStr: string): Array<{name: string, type: string}> {
    const params: Array<{name: string, type: string}> = [];
    
    if (!paramStr.trim()) return params;

    const paramList = paramStr.split(',').map(p => p.trim());
    
    paramList.forEach(param => {
      const parts = param.split(':');
      if (parts.length === 2) {
        params.push({
          name: parts[0].trim(),
          type: parts[1].trim()
        });
      }
    });

    return params;
  }

  private convertParameters(params: string): string {
    if (!params.trim()) return '';
    
    return params.split(',').map(param => {
      const parts = param.trim().split(':');
      if (parts.length === 2) {
        const name = parts[0].trim();
        const type = parts[1].trim();
        const goType = this.convertType(type);
        return `${name} ${goType}`;
      }
      return param.trim();
    }).join(', ');
  }

  private convertFunctionBody(body: string): string {
    // This is a very basic conversion - real implementation would be more sophisticated
    return `\t// TODO: Convert function body\n\t${body.trim()}`;
  }

  private convertType(tsType: string): string {
    const typeMap: Record<string, string> = {
      'string': 'string',
      'number': 'float64',
      'boolean': 'bool',
      'void': '',
      'any': 'interface{}',
      'object': 'interface{}',
      'null': 'nil',
      'undefined': 'nil'
    };

    // Handle array types
    if (tsType.endsWith('[]')) {
      const elementType = tsType.slice(0, -2);
      return `[]${this.convertType(elementType)}`;
    }

    // Handle generic types
    if (tsType.includes('<') && tsType.includes('>')) {
      const baseType = tsType.substring(0, tsType.indexOf('<'));
      const genericType = tsType.substring(tsType.indexOf('<') + 1, tsType.lastIndexOf('>'));
      
      if (baseType === 'Array') {
        return `[]${this.convertType(genericType)}`;
      } else if (baseType === 'Promise') {
        return `*async.Promise[${this.convertType(genericType)}]`;
      } else if (baseType === 'Map') {
        const types = genericType.split(',').map(t => t.trim());
        if (types.length === 2) {
          return `*types.Map[${this.convertType(types[0])}, ${this.convertType(types[1])}]`;
        }
      } else if (baseType === 'Set') {
        return `*types.Set[${this.convertType(genericType)}]`;
      } else if (baseType === 'Optional') {
        return `types.Optional[${this.convertType(genericType)}]`;
      }
    }

    return typeMap[tsType] || tsType;
  }

  private convertValue(value: string): string {
    // Convert boolean literals
    if (value === 'true' || value === 'false') {
      return value;
    }
    
    // Convert null/undefined
    if (value === 'null' || value === 'undefined') {
      return 'nil';
    }
    
    // Convert string literals
    if (value.startsWith('"') || value.startsWith("'")) {
      return value.replace(/'/g, '"');
    }
    
    // Convert number literals
    if (/^\d+(\.\d+)?$/.test(value)) {
      return value;
    }

    return value;
  }

  private generateImports(goCode: string): string {
    const imports = new Set<string>();

    // Check for usage of our types
    if (goCode.includes('types.')) {
      imports.add('"PROJECT_NAME/types"');
    }
    if (goCode.includes('utils.')) {
      imports.add('"PROJECT_NAME/utils"');
    }
    if (goCode.includes('async.')) {
      imports.add('"PROJECT_NAME/async"');
    }
    if (goCode.includes('classes.')) {
      imports.add('"PROJECT_NAME/classes"');
    }
    if (goCode.includes('enums.')) {
      imports.add('"PROJECT_NAME/enums"');
    }

    // Standard library imports
    if (goCode.includes('fmt.')) {
      imports.add('"fmt"');
    }
    if (goCode.includes('time.')) {
      imports.add('"time"');
    }

    if (imports.size === 0) return '';

    const importList = Array.from(imports).join('\n\t');
    return `import (\n\t${importList}\n)`;
  }

  private capitalize(str: string): string {
    return str.charAt(0).toUpperCase() + str.slice(1);
  }
}