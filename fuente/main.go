/*
Herramienta de comandos con utilidades para
crear proyectos con el stack de Ch'aska.
*/
package main

import (
	"bytes"
	"errors"
	"slices"
	"strconv"

	"github.com/hernanatn/aplicacion.go"
	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"

	"chaskapp/data"
	"chaskapp/proyecto"
	"chaskapp/utiles"

	"fmt"

	"io"

	"os"
	"os/exec"

	"github.com/schollz/progressbar/v3"
)

func garantizarPython(a aplicacion.Aplicacion) {
	a.ImprimirSeparador()
	cmd := exec.Command("py", "-V")
	pyv_out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Para crear el proyecto se requiere tener un runtime de Python Instalado en el ambiente actual.\nVea python.org para más información\n%s\n", err)
		panic(1) // [HACER] Procesar Error
	}
	progreso := progressbar.Default(-1, cadena.Cadena("Garantizando Python").Italica().S())
	var buf bytes.Buffer
	io.Copy(io.MultiWriter(&buf, progreso), os.Stdout)
	a.BorrarLinea()
	a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Python Grantizado: %s", pyv_out)))

	progreso.Describe("Garantizado Pip")
	cmd = exec.Command("pip", "-V")
	pipv_out, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("py", "-m", "ensurepip")
		pip_ins, err2 := cmd.Output()
		fmt.Printf("Pip: %s\n", pip_ins)
		if err2 != nil {
			a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Se requiere el manejador de paquetes PIP, el cual no estaba presente en el ambiente actual y no se pudo instalar automáticamente por un error.\nPuede optar por instalarlo manualmente `py -m ensurepip`.\nVea python.org para más información\n%s\n", err2)))
			panic(2) // [HACER] Procesar Error
		}
	}
	a.BorrarLinea()
	a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Pip Grantizado: %s", pipv_out)))

	progreso.Describe("Actualizando Pip")
	cmd = exec.Command("py", "-m", "pip", "install", "--upgrade", "pip")
	pip_upt, err := cmd.Output()
	progreso.Clear()
	a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Pip Actualizado: %s\n", utiles.Limpiar(string(pip_upt)))))
	if err != nil {
		a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Se requiere el manejador de paquetes PIP, el cual se pudo instalar y pero no se pudo actualizar automáticamente.\nVea python.org para más información\n%s\n", err)))
		panic(3) // [HACER] Procesar Error
	}
	progreso.Describe(cadena.Cadena("Ambiente Python >=3.12 + pip garantizado correctamente.").Italica().Colorear(color.VerdeFuente).S())
	progreso.Finish()
	progreso.Clear()
}

func garantizarNPM(a aplicacion.Aplicacion) {
	a.ImprimirSeparador()
	cmd := exec.Command("node", "-v")
	npm_out, err := cmd.Output()
	if err != nil {
		a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Para crear el proyecto se requiere tener Node y npm instalado en el ambiente actual.\nVea https://nodejs.org/en/download/package-manager para más información\n%s\n", err)))
		panic(1) // [HACER] Procesar Error
	}
	progreso := progressbar.Default(-1, cadena.Cadena("Garantizando NPM").Italica().S())
	var buf bytes.Buffer
	io.Copy(io.MultiWriter(&buf, progreso), os.Stdout)
	a.BorrarLinea()
	a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("NPM Grantizado: %s", npm_out)))

	progreso.Describe("Garantizando lightningcss")
	cmd = exec.Command("npm", "install", "--save-dev", "lightningcss-cli")
	light_out, err := cmd.Output()
	progreso.Clear()
	a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("lightningcss instalado: %s", utiles.Limpiar(string(light_out)))))
	if err != nil {
		a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Se requiere la herramienta CLI de lightningcss para minificar y hacer compatible los estilos, pero no se pudo instalar via npm.\n%s\n", err)))
		panic(3) // [HACER] Procesar Error
	}
	progreso.Describe("Garantizando uglify-js")
	cmd = exec.Command("npm", "install", "--save-dev", "uglify-js", "-g")
	uglify_out, err := cmd.Output()
	progreso.Clear()
	a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("uglify-js instalado: %s\n", utiles.Limpiar(string(uglify_out)))))
	if err != nil {
		a.ImprimirCadena(aplicacion.Cadena(fmt.Sprintf("Se requiere la herramienta CLI de uglify para minificar los paquetes js, pero no se pudo instalar via npm.\n%s\n", err)))
		panic(3) // [HACER] Procesar Error
	}
	progreso.Describe(cadena.Cadena("Ambiente Node + npm garantizado correctamente.").Italica().Colorear(color.VerdeFuente).S())
	progreso.Finish()
	progreso.Clear()

}

/*
func NuevoProyecto(a aplicacion.Aplicacion, argumentos []string) *proyecto.Proyecto {
	a.ImprimirSeparador()
	var err error

	var n string
	if !slices.Contains(argumentos, "-p") {

		a.Escribir(cadena.Negrita("Ingrese el nombre del proyecto: "))
		n, err = a.Leer()
		if err != nil {
			log.Fatal(err)
			panic(1) // [HACER] Procesar Error
		}
	} else {
		n = argumentos[slices.Index(argumentos, "-p")+1]
	}
	nombre := utiles.Limpiar(n)

	var r string
	if !slices.Contains(argumentos, "-r") {
		fmt.Print(cadena.Negrita("Defina ruta raíz del proyecto (omita si desea usar el directoria actual): "))
		r, err = lector.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			panic(1) // [HACER] Procesar Error
		}
	} else {
		r = argumentos[slices.Index(argumentos, "-r")+1]
	}
	ruta := path.Clean(utiles.Limpiar(r))

	var d string = "0"
	var dependenciasAdicionales string
	if !slices.Contains(argumentos, "--no-deps") {
		if !slices.Contains(argumentos, "-d") {
			fmt.Println(cadena.Negrita("Elija dependencias según dominio:"))
			fmt.Println(cadena.Negrita("\t0."), "\tBásico")
			fmt.Println(cadena.Negrita("\t1."), "\tContacto/Correos")
			fmt.Println(cadena.Negrita("\t2."), "\tAdministrador")
			fmt.Println(cadena.Negrita("\t3."), "\tComercio Electónico")
			fmt.Printf("Su opción (0 a 3): ")
			d, err = lector.ReadString('\n')
			if err != nil {
				panic(1) // [HACER] Procesar Error
			}
		} else {
			d = argumentos[slices.Index(argumentos, "-d")+1]
		}

		if !slices.Contains(argumentos, "-deps") {
			fmt.Println(aplicacion.Negrita("Indique dependencias adicionales"), aplicacion.Italica("\n(indique ruta a archivo requeriments.txt o lista separada por espacios de dependencias en formato `[dependencia](==[version])`.):"))
			dependenciasAdicionales, err = lector.ReadString('\n')
			if err != nil {
				log.Fatal(err)
				panic(1) // [HACER] Procesar Error
			}
		} else {
			dependenciasAdicionales = argumentos[slices.Index(argumentos, "-deps")+1]
		}
	}
	dom, err := strconv.Atoi(utiles.Limpiar(d))
	if err != nil {
		log.Fatal(err)
		panic(err) // [HACER] Procesar Error
	}
	dominio := proyecto.DominioProyecto(dom)

	aplicacion.ImprimirSeparador()
	fmt.Println("Creando")
	fmt.Printf("Proyecto: %s\n", nombre)
	fmt.Printf("Ruta: %s\n", ruta)
	fmt.Printf("Dominio: %d\n", dominio)
	fmt.Printf("Dependencias Adicionales: %s", dependenciasAdicionales)
	aplicacion.ImprimirSeparador()
	fmt.Printf("¿Desea continuar? (S/N):")
	confirmacion, err := lector.ReadString('\n')
	if err != nil {
		panic(1) // [HACER] Procesar Error
	}

	var esteProyecto *proyecto.Proyecto
	switch utiles.Limpiar(confirmacion) {
	case "S", "s":
		esteProyecto = proyecto.NuevoProyecto(nombre, ruta, dominio)
		fmt.Printf("%s creado!\n", esteProyecto.Nombre)
	default:
		panic(69) //[HACER] editar info proyecto / volver a empezar.

	}

	if strings.Contains(dependenciasAdicionales, ".txt") {
		req := path.Clean(utiles.Limpiar(dependenciasAdicionales))
		requerimientos, err := os.Open(req)
		if err != nil {
			log.Fatal(err)
			panic(1) // [HACER] Procesar Error
		}
		defer requerimientos.Close()

		escaner := bufio.NewScanner(requerimientos)

		for escaner.Scan() {
			esteProyecto.AgregarDependencia(escaner.Text())
		}

		if err := escaner.Err(); err != nil {
			log.Fatal(err)
			panic(1) // [HACER] Procesar Error
		}

	} else {
		depsRecibidas := strings.Split(utiles.Limpiar(dependenciasAdicionales), " ")
		for _, d := range depsRecibidas {
			esteProyecto.AgregarDependencia(d)
		}
	}

	esteProyecto.
		Crear().
		InstalarDependencias()

	return esteProyecto
	//esteProyecto.agregarDependencias()

	//switch entrada {
	//case "":
	//}
}*/

var chaskapp aplicacion.Aplicacion

func ini(a aplicacion.Aplicacion, args ...string) error {
	a.ImprimirLinea(cadena.Cadena(data.HECHO_POR_CHASKA).Negrita())
	if len(args) <= 0 || len(args) > 0 && !slices.Contains(args, "--sin-ini") {
		garantizarPython(a)
		garantizarNPM(a)
	}

	return nil
}

func fin(a aplicacion.Aplicacion, args ...string) error {
	a.ImprimirLinea(aplicacion.Cadena("¡Adiós!"))
	return nil
}

func init() {
	chaskapp = aplicacion.NuevaAplicacion(
		"Ch'askapp",
		"chaskapp",
		"chaskapp / v 0.1",
		make([]string, 0),
		aplicacion.NuevaConsola(os.Stdin, os.Stdout),
	).
		RegistrarInicio(ini).
		RegistrarLimpieza(fin).
		RegistrarFinal(fin)
}

func main() {
	var menuDominio = aplicacion.NuevoMenu(
		chaskapp,
		'@',
	)

	menuDominio.
		RegistrarOpcion(
			&aplicacion.OpcionMenu{
				Nombre: "1. Básico",
				Accion: aplicacion.Accion(
					func(con aplicacion.Consola, opciones aplicacion.Opciones, parametros aplicacion.Parametros, argumentos ...any) (res any, cod aplicacion.CodigoError, err error) {
						return proyecto.BASICO, aplicacion.EXITO, nil
					}),
			}).
		RegistrarOpcion(
			&aplicacion.OpcionMenu{
				Nombre: "2. Correo",
				Accion: aplicacion.Accion(
					func(con aplicacion.Consola, opciones aplicacion.Opciones, parametros aplicacion.Parametros, argumentos ...any) (res any, cod aplicacion.CodigoError, err error) {
						return proyecto.CORREO, aplicacion.EXITO, nil
					}),
			}).
		RegistrarOpcion(
			&aplicacion.OpcionMenu{
				Nombre: "3. Admin",
				Accion: aplicacion.Accion(
					func(con aplicacion.Consola, opciones aplicacion.Opciones, parametros aplicacion.Parametros, argumentos ...any) (res any, cod aplicacion.CodigoError, err error) {
						return proyecto.ADMIN, aplicacion.EXITO, nil
					}),
			}).
		RegistrarOpcion(
			&aplicacion.OpcionMenu{
				Nombre: "4. Comercio",
				Accion: aplicacion.Accion(
					func(con aplicacion.Consola, opciones aplicacion.Opciones, parametros aplicacion.Parametros, argumentos ...any) (res any, cod aplicacion.CodigoError, err error) {
						return proyecto.COMERCIO, aplicacion.EXITO, nil
					}),
			})

	var crear aplicacion.Comando = aplicacion.NuevoComando(
		"crear",
		"crear -p <PROYECTO> -r <RAIZ PROYECTO> -d <DOMINIO> [OPCIONES]",
		[]string{"create"},
		"Crea un nuevo proyecto web con el stack de Ch'aska y organiza las dependencias básicas.",
		aplicacion.Accion(
			func(con aplicacion.Consola, opciones aplicacion.Opciones, parametros aplicacion.Parametros, argumentos ...any) (res any, cod aplicacion.CodigoError, err error) {

				for ll, va := range parametros {
					con.ImprimirLinea(cadena.Cadena(fmt.Sprintf("%s,%s", ll, va)))
				}
				for _, s := range opciones {
					fmt.Println(s)
					con.ImprimirLinea(cadena.Cadena(s))
				}
				var e error
				var ok bool
				var errores []error

				var p any
				var nom aplicacion.Cadena
				var nombre string
				p, existe := parametros["-p"]
				if !existe {
					p, e = con.Leer("Indique nombre del proyecto")
					if e != nil {
						con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer el nombre del proyecto. Se utilizó valor por defecto.", e)))
						nom = aplicacion.Cadena("proyecto-chaska")
						errores = append(errores, e)
					}
				} else {
					nom, _ := p.(string)
					p = aplicacion.Cadena(nom)
				}
				nom, ok = p.(aplicacion.Cadena)
				if !ok {
					con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer el nombre del proyecto. Se utilizó valor por defecto.", e)))
					nom = aplicacion.Cadena("proyecto-chaska")
					errores = append(errores, e)
				}
				nombre = nom.S()

				var r any
				var ru aplicacion.Cadena
				var ruta string
				r, existe = parametros["-r"]
				if !existe {
					r, e = con.Leer("Indique ruta raiz del proyecto")
					if e != nil {
						con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer la ruta del proyecto. Se utilizó valor por defecto.", e)))
						r = aplicacion.Cadena(".")
						errores = append(errores, e)
					}
				} else {
					ru, _ := r.(string)
					r = aplicacion.Cadena(ru)
				}
				ru, ok = r.(aplicacion.Cadena)
				if !ok {
					con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer la ruta del proyecto. Se utilizó valor por defecto.", e)))
					r = aplicacion.Cadena(".")
					errores = append(errores, e)
				}

				ruta = ru.S()

				var d any
				var dominio proyecto.DominioProyecto
				d, existe = parametros["-d"]
				if !existe {
					e = con.ImprimirLinea("Indique el dominio del proyecto:")
					if e != nil {
						con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer el dominio del proyecto.", e)))
						//[HACER] Agregar valor por defecto
						errores = append(errores, e)
					}
					o, e := menuDominio.Correr()
					if e != nil {
						con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer el dominio del proyecto.", e)))
						//[HACER] Agregar valor por defecto
						errores = append(errores, e)
					} else {
						d, _, e = o.Accion(con, []string{}, aplicacion.Parametros{})
						if e != nil {
							con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer el dominio del proyecto.", e)))
							//[HACER] Agregar valor por defecto
							errores = append(errores, e)
						}
					}
				} else {
					dom, _ := d.(string)
					do, e := strconv.Atoi(dom)
					if e != nil {
						con.EscribirLinea(cadena.Cadena(cadena.Error("No se pudo leer el dominio del proyecto.", e)))
						//[HACER] Agregar valor por defecto
						errores = append(errores, e)
					}
					d = proyecto.DominioProyecto(do)
				}
				dominio, _ = d.(proyecto.DominioProyecto)

				con.ImprimirLinea(cadena.Cadena(fmt.Sprintf("%s%s%d", nombre, ruta, dominio)))

				if len(errores) > 0 {
					return nil, aplicacion.ERROR, errors.Join(errores...)
				}
				var esteProyecto *proyecto.Proyecto
				esteProyecto = proyecto.NuevoProyecto(nombre, ruta, dominio)
				esteProyecto.Crear(chaskapp)
				return esteProyecto, aplicacion.EXITO, nil
			}),
		[]string{"--sin-deps", "--por-defecto"},
	)
	var minificar aplicacion.Comando = aplicacion.NuevoComando(
		"minificar",
		"minificar -r <RAIZ PROYECTO> -o <OBJETIVOS> [OPCIONES]",
		[]string{"minify"},
		"Minifica los estáticos del proyecto indicado (o del directorio actual, en su defecto).",
		aplicacion.Accion(
			func(con aplicacion.Consola, opciones aplicacion.Opciones, parametros aplicacion.Parametros, argumentos ...any) (res any, cod aplicacion.CodigoError, err error) {
				con.ImprimirLinea(aplicacion.Cadena("minificar"))
				return nil, aplicacion.EXITO, nil
			}),
		[]string{"--sin-objetivos"},
	)
	var correr aplicacion.Comando = aplicacion.NuevoComando(
		"correr",
		"correr -p <PROYOECTO> -r <RAIZ PROYECTO> [OPCIONES] --> [OPCIONES flask]",
		[]string{"run"},
		"Corre el servidor con flask.",
		aplicacion.Accion(
			func(con aplicacion.Consola, opciones aplicacion.Opciones, parametros aplicacion.Parametros, argumentos ...any) (res any, cod aplicacion.CodigoError, err error) {
				con.ImprimirLinea(aplicacion.Cadena("correr"))
				return nil, aplicacion.EXITO, nil
			}),
		make([]string, 0),
	)
	/*
		crear := (aplicacion.Comando{
			Nombre:      "crear",
			Uso:         "crear [OPCIONES]",
			Descripcion: "Crea un nuevo proyecto web con el stack de Ch'aska y organiza las dependencias básicas.",
			Accion: aplicacion.Accion(
				func(salida *aplicacion.Salida, parametros aplicacion.Parametros, opciones ...string) (res any, cod aplicacion.CodigoError, err error) {
					salida.WriteString("crear\n")
					salida.Flush()
					return aplicacion.EXITO, nil
				}),
			Opciones: []string{"-p", "-r", "-d", "--sin-deps", "--por-defecto"},
		}).RegistrarComando((&aplicacion.Comando{
			Nombre:      "prueba",
			Uso:         "prueba [OPCIONES]",
			Descripcion: "subcomando de prueba",
			Accion: aplicacion.Accion(
				func(salida *aplicacion.Salida, parametros aplicacion.Parametros, opciones ...string) (res any, cod aplicacion.CodigoError, err error) {
					salida.WriteString("prueba\n")
					salida.Flush()
					return aplicacion.EXITO, nil
				}),
			Opciones: []string{"-p", "-r", "-d", "--sin-deps", "--por-defecto"},
		}).RegistrarComando(&aplicacion.Comando{
			Nombre:      "jaja",
			Uso:         "jaja [OPCIONES]",
			Descripcion: "subcomando de jaja",
			Accion: aplicacion.Accion(
				func(salida *aplicacion.Salida, parametros aplicacion.Parametros, opciones ...string) (res any, cod aplicacion.CodigoError, err error) {
					salida.WriteString("jaja\n")
					salida.Flush()
					return aplicacion.EXITO, nil
				}),
			Opciones: []string{"-p", "-r", "-d", "--sin-deps", "--por-defecto"},
		}))
		minificar := &aplicacion.Comando{
			Nombre:      "minificar",
			Uso:         "minificar [OPCIONES]",
			Descripcion: "Minifica los estáticos del proyecto indicado (o del directorio actual, en su defecto).",
			Accion: aplicacion.Accion(
				func(salida *aplicacion.Salida, parametros aplicacion.Parametros, opciones ...string) (res any, cod aplicacion.CodigoError, err error) {
					salida.WriteString("minificar\n")
					salida.Flush()
					return aplicacion.EXITO, nil
				}),
			Opciones: []string{"-p", "-r"},
		}
		correr := &aplicacion.Comando{
			Nombre:      "correr",
			Uso:         "correr [OPCIONES] --> [OPCIONES flask]",
			Descripcion: "Corre el servidor con flask.",
			Accion: aplicacion.Accion(
				func(salida *aplicacion.Salida, parametros aplicacion.Parametros, opciones ...string) (res any, cod aplicacion.CodigoError, err error) {
					salida.WriteString("correr\n")
					salida.Flush()
					return aplicacion.EXITO, nil
				}),
			Opciones: []string{"-p", "-r"},
		}
	*/
	res, err := chaskapp.
		RegistrarComando(crear).
		RegistrarComando(minificar).
		RegistrarComando(correr).
		Correr(os.Args[1:]...)

	if err != nil {
		fmt.Print(cadena.Fatal(" ", err))
	}
	if res != nil {
		// [HACER] ver si sirve que suba el output hasta main...
	}

	/*
			lector := bufio.NewReader(os.Stdin)
			var esteProyecto *proyecto.Proyecto

			argumentos := os.Args[1:]
			if len(argumentos) > 0 {
				switch argumentos[0] {
				case "proyecto":
					goto Proyecto
				case "minificar":
					goto Minificar
				case "ayuda", "-a", "-h":
					goto Ayuda
				default:
					goto Ambiente
				}
			}

		Ambiente:
			garantizarPython()
			garantizarNPM()

		Proyecto:
			esteProyecto = NuevoProyecto(*lector, argumentos)

		Minificar:
			if esteProyecto != nil {
				esteProyecto.Minificar()
			} else {
				if len(argumentos) < 2 {
					panic(2)
				}
			}

			goto Cerrar

		Ayuda:
			aplicacion.ImprimirTitulo("Ayuda")
			return

		Cerrar:
			aplicacion.ImprimirSeparador()
			aplicacion.Imprimir(aplicacion.Cadena{"Juan"})
			aplicacion.ImprimirSubtitulo("¡Todo listo!")
			return
	*/
}
