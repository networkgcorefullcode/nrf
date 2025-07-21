# NRF Development Environment

Este directorio contiene los archivos necesarios para crear un entorno de desarrollo completo para la aplicación NRF.

## Archivos incluidos

- `Dockerfile_dev` - Dockerfile optimizado para desarrollo
- `dev-container.ps1` - Script de gestión del contenedor para Windows PowerShell
- `dev-container.sh` - Script de gestión del contenedor para Linux/macOS Bash

## Características del entorno de desarrollo

### Herramientas incluidas

- **Git** - Control de versiones
- **Go 1.24.5** - Compilador Go
- **Herramientas de desarrollo Go**:
  - Air (hot reload)
  - golangci-lint (linting)
  - Delve (debugging)
- **Editores**: vim, nano
- **Herramientas de red**: curl, wget, netcat, tcpdump
- **Herramientas de debugging**: strace, lsof, htop
- **Utilidades**: tree, jq, make

### Scripts incluidos en el contenedor

- `update-code` - Actualiza el código desde git
- `build-app` - Compila la aplicación NRF
- `run-app` - Ejecuta la aplicación NRF
- `test-app` - Ejecuta los tests
- `lint-app` - Ejecuta el linter
- `clean-app` - Limpia artefactos de compilación
- `dev-info` - Muestra información del entorno

## Flujo de trabajo recomendado

### 1. Construir la imagen de desarrollo

**Windows PowerShell:**
```powershell
.\dev-container.ps1 build
```

**Linux/macOS:**
```bash
./dev-container.sh build
```

### 2. Iniciar el contenedor de desarrollo

**Windows PowerShell:**
```powershell
.\dev-container.ps1 run
```

**Linux/macOS:**
```bash
./dev-container.sh run
```

Esto iniciará un contenedor interactivo con:
- El código fuente montado en `/app`
- Puerto 29510 expuesto
- Configuración de Git (si está disponible en el host)

### 3. Desarrollo iterativo dentro del contenedor

Una vez dentro del contenedor:

```bash
# Actualizar código desde repositorio
update-code

# Compilar la aplicación
build-app

# Ejecutar la aplicación
run-app

# Ejecutar tests
test-app

# Ejecutar linter
lint-app
```

### 4. Flujo después de hacer push de nuevo código

Cuando hagas push de nuevo código desde tu máquina host:

```bash
# Dentro del contenedor, actualiza el código
update-code

# Recompila
build-app

# Ejecuta la nueva versión
run-app
```

## Comandos adicionales del script de gestión

### Entrar a un contenedor ya ejecutándose
**Windows PowerShell:**
```powershell
.\dev-container.ps1 exec
```

### Ver logs del contenedor
**Windows PowerShell:**
```powershell
.\dev-container.ps1 logs
```

### Detener el contenedor
**Windows PowerShell:**
```powershell
.\dev-container.ps1 stop
```

### Reiniciar el contenedor
**Windows PowerShell:**
```powershell
.\dev-container.ps1 restart
```

## Configuración de puertos

Por defecto, el puerto 29510 está expuesto para la aplicación NRF. Puedes modificar esto en los scripts de gestión si necesitas diferentes puertos.

## Volúmenes montados

- Código fuente: `./` → `/app`
- Configuración Git del host: `~/.gitconfig` → `/root/.gitconfig` (solo lectura)
- Claves SSH del host: `~/.ssh` → `/root/.ssh` (solo lectura)

## Troubleshooting

1. **Error de permisos con Git**: Asegúrate de que tu configuración de Git esté disponible
2. **Puerto ocupado**: Cambia el HOST_PORT en el script de gestión
3. **Problemas de compilación**: Ejecuta `clean-app` antes de `build-app`

## Personalización

Puedes personalizar el entorno modificando:
- `Dockerfile_dev` para añadir más herramientas
- Scripts en `/usr/local/bin/` dentro del contenedor
- Variables de entorno en el Dockerfile
