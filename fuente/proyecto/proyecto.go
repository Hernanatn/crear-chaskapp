package proyecto

import (
	//"os/exec"
	//"io"
	"bufio"
	"bytes"
	"chaskapp/data"
	"chaskapp/utiles"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/hernanatn/aplicacion.go/consola"
	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"

	"github.com/schollz/progressbar/v3"
)

type Consola = consola.Consola
type DominioProyecto int

const (
	BASICO   DominioProyecto = iota
	CORREO   DominioProyecto = iota
	ADMIN    DominioProyecto = iota
	COMERCIO DominioProyecto = iota
)

// Interfaz
type Proyecto struct {
	Nombre       string
	Ruta         string
	Dominio      DominioProyecto
	Dependencias []string
}

func NuevoProyecto(n string, r string, d DominioProyecto) *Proyecto {
	esteProy := new(Proyecto)
	esteProy.Nombre = n
	esteProy.Ruta = r
	esteProy.Dominio = d
	esteProy.Dependencias = make([]string, 0, 15)

	var depsTemp []string

	archivoDeps, err := data.DEPENDENCIAS.ReadFile("dependencias/basico.txt")
	if err != nil {
		log.Fatal("No se pudo leer el archivo de dependencias básicas. \n", err)
		panic(1) // [HACER] Procesar Error
	}

	switch esteProy.Dominio {
	case COMERCIO:
		ad, err := data.DEPENDENCIAS.ReadFile("dependencias/COMERCIO/comercio.txt")
		if err != nil {
			log.Fatal("No se pudo leer el archivo de dependencias. \n", err)
			panic(1) // [HACER] Procesar Error
		}
		deps := strings.Split(string(ad), "\n")
		depsTemp = append(depsTemp, deps...)
		fallthrough
	case ADMIN:
		ad, err := data.DEPENDENCIAS.ReadFile("dependencias/ADMIN/admin.txt")
		if err != nil {
			log.Fatal("No se pudo leer el archivo de dependencias. \n", err)
			panic(1) // [HACER] Procesar Error
		}
		deps := strings.Split(string(ad), "\n")
		depsTemp = append(depsTemp, deps...)
		fallthrough
	case CORREO:
		ad, err := data.DEPENDENCIAS.ReadFile("dependencias/CORREO/correo.txt")
		if err != nil {
			log.Fatal("No se pudo leer el archivo de dependencias. \n", err)
			panic(1) // [HACER] Procesar Error
		}
		deps := strings.Split(string(ad), "\n")
		depsTemp = append(depsTemp, deps...)
		fallthrough
	case BASICO:
		fallthrough
	default:
		depsBasicas := strings.Split(string(archivoDeps), "\n")
		depsTemp = append(depsTemp, depsBasicas...)
	}

	for _, d := range depsTemp {
		esteProy.AgregarDependencia(d)
	}
	// [HACER] Agregar dependencias base según dominio
	return esteProy
}

// Métodos
func (p *Proyecto) Inicializar(n string, r string, d DominioProyecto) *Proyecto {
	p.Nombre = n
	p.Ruta = r
	p.Dominio = d
	p.Dependencias = make([]string, 0, 15)

	var depsTemp []string

	archivoDeps, err := data.DEPENDENCIAS.ReadFile("dependencias/requeriments.txt")
	if err != nil {
		log.Fatal("No se pudo leer el archivo de dependencias básicas. \n", err)
		panic(1) // [HACER] Procesar Error
	}

	depsBasicas := strings.Split(string(archivoDeps), "\n")

	depsTemp = append(depsTemp, depsBasicas...)

	switch p.Dominio {
	case COMERCIO:
		ad, err := data.DEPENDENCIAS.ReadFile("dependencias/COMERCIO/comercio.txt")
		if err != nil {
			log.Fatal("No se pudo leer el archivo de dependencias. \n", err)
			panic(1) // [HACER] Procesar Error
		}
		deps := strings.Split(string(ad), "\n")
		depsTemp = append(depsTemp, deps...)
		fallthrough
	case ADMIN:
		ad, err := data.DEPENDENCIAS.ReadFile("dependencias/ADMIN/admin.txt")
		if err != nil {
			log.Fatal("No se pudo leer el archivo de dependencias. \n", err)
			panic(1) // [HACER] Procesar Error
		}
		deps := strings.Split(string(ad), "\n")
		depsTemp = append(depsTemp, deps...)
		fallthrough
	case CORREO:
		ad, err := data.DEPENDENCIAS.ReadFile("dependencias/CORREO/correo.txt")
		if err != nil {
			log.Fatal("No se pudo leer el archivo de dependencias. \n", err)
			panic(1) // [HACER] Procesar Error
		}
		deps := strings.Split(string(ad), "\n")
		depsTemp = append(depsTemp, deps...)
	}

	for _, d := range depsTemp {
		p.AgregarDependencia(d)
	}
	// [HACER] Agregar dependencias base según dominio
	return p
}

func (p Proyecto) crearDirectorios(con Consola) Proyecto {
	var err error
	repo := path.Clean(p.Ruta)
	err = os.MkdirAll(repo, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio del repositorio para el proyecto. \n"+cadena.Cadena(repo), err)
		panic(1) // [HACER] Procesar Error
	}

	raiz := path.Join(repo, p.Nombre)
	err = os.MkdirAll(raiz, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio raíz del proyecto. \n"+cadena.Cadena(raiz), err)
		panic(1) // [HACER] Procesar Error
	}
	fuente := path.Join(raiz, "fuente")
	err = os.MkdirAll(fuente, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio fuente del proyecto. \n"+cadena.Cadena(fuente), err)
		panic(1) // [HACER] Procesar Error
	}
	proy := path.Join(fuente, p.Nombre)
	err = os.MkdirAll(proy, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio fuente del proyecto. \n"+cadena.Cadena(proy), err)
		panic(1) // [HACER] Procesar Error
	}
	bdd := path.Join(proy, "base_de_datos")
	err = os.MkdirAll(bdd, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio bdd del proyecto. \n"+cadena.Cadena(bdd), err)
		panic(1) // [HACER] Procesar Error
	}
	vistas := path.Join(proy, "vistas")
	err = os.MkdirAll(vistas, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio vistas del proyecto. \n"+cadena.Cadena(vistas), err)
		panic(1) // [HACER] Procesar Error
	}
	estatico := path.Join(proy, "estatico")
	err = os.MkdirAll(estatico, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio estatico del proyecto. \n"+cadena.Cadena(estatico), err)
		panic(1) // [HACER] Procesar Error
	}
	js := path.Join(estatico, "js")
	err = os.MkdirAll(js, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio estatico del proyecto. \n"+cadena.Cadena(js), err)
		panic(1) // [HACER] Procesar Error
	}
	minjs := path.Join(estatico, "min-js")
	err = os.MkdirAll(minjs, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio estatico del proyecto. \n"+cadena.Cadena(minjs), err)
		panic(1) // [HACER] Procesar Error
	}
	css := path.Join(estatico, "css")
	err = os.MkdirAll(css, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio estatico del proyecto. \n"+cadena.Cadena(css), err)
		panic(1) // [HACER] Procesar Error
	}
	mincss := path.Join(estatico, "min-css")
	err = os.MkdirAll(mincss, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio estatico del proyecto. \n"+cadena.Cadena(mincss), err)
		panic(1) // [HACER] Procesar Error
	}
	html := path.Join(estatico, "html")
	err = os.MkdirAll(html, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio estatico del proyecto. \n"+cadena.Cadena(html), err)
		panic(1) // [HACER] Procesar Error
	}
	media := path.Join(estatico, "media")
	err = os.MkdirAll(media, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio estatico del proyecto. \n"+cadena.Cadena(media), err)
		panic(1) // [HACER] Procesar Error
	}
	plantillas := path.Join(proy, "plantillas")
	err = os.MkdirAll(plantillas, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio plantillas del proyecto. \n"+cadena.Cadena(plantillas), err)
		panic(1) // [HACER] Procesar Error
	}
	utiles := path.Join(proy, "utiles")
	err = os.MkdirAll(utiles, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio utiles del proyecto. \n"+cadena.Cadena(utiles), err)
		panic(1) // [HACER] Procesar Error
	}
	api := path.Join(proy, "api")
	err = os.MkdirAll(api, fs.ModePerm)
	if err != nil {
		con.ImprimirFatal("No se pudo crear el directorio api del proyecto. \n"+cadena.Cadena(api), err)
		panic(1) // [HACER] Procesar Error
	}

	return p
}

func (p Proyecto) escribrirArchivos(con Consola) Proyecto {
	var err error
	p, err = p.escribirEstaticos()
	if err != nil {
		con.ImprimirFatal("No se pudo escribir los archivos estáticos base del proyecto. \n", err)
		panic(1) // [HACER] Procesar Error
	}
	return p
}
func (p Proyecto) escribirEstaticos() (Proyecto, error) {
	var err error
	rutaEstaticos := path.Join(path.Join(path.Join(path.Join(path.Clean(p.Ruta), p.Nombre), "fuente"), p.Nombre), "estatico")

	js := make([]byte, 0)
	leerJS := func(ruta string) func(string, fs.DirEntry, error) error {
		return func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("err: %v\n", err)
				// return err
			} else {
				if filepath.Ext(p) == ".js" {
					archivo, err := data.DEPENDENCIAS.ReadFile(path.Join(p))
					if err != nil {
						fmt.Printf("No se pudo abrir el archivo %s: %v\n", p, err)
						return err
					}
					js = append(js, "\n"...)
					js = append(js, archivo...)
				}
			}
			return nil
		}
	}

	var ruta string
	switch p.Dominio {
	case COMERCIO:
		ruta = "dependencias/COMERCIO/js"
		err = fs.WalkDir(data.DEPENDENCIAS, ruta, leerJS(ruta))
		if err != nil {
			return p, err
		}
		fallthrough
	case ADMIN:
		ruta = "dependencias/ADMIN/js"
		err = fs.WalkDir(data.DEPENDENCIAS, ruta, leerJS(ruta))
		if err != nil {
			return p, err
		}
		fallthrough
	case CORREO:
		ruta = "dependencias/CORREO/js"
		err = fs.WalkDir(data.DEPENDENCIAS, ruta, leerJS(ruta))
		if err != nil {
			return p, err
		}
	}
	os.WriteFile(path.Join(rutaEstaticos, "js", "paquete.js"), js, fs.ModePerm)

	err = fs.WalkDir(
		data.ESTATICO,
		"estatico/css",
		func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("err: %v\n", err)
				// return err
			} else {
				if filepath.Ext(p) == ".css" {
					fmt.Println(p)
					archivo, err := data.ESTATICO.ReadFile(path.Join(p))
					if err != nil {
						fmt.Printf("No se pudo abrir el archivo %s: %v\n", p, err)
						return err
					}
					err = os.WriteFile(path.Join(rutaEstaticos, "css", path.Base(p)), archivo, fs.ModePerm)
					if err != nil {
						fmt.Printf("no se pudo escribir el archivo %s: %v\n", p, err)
						return err
					}
				}
			}
			return nil
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return p, err
	}

	return p, nil
}

func (p Proyecto) instalarDependencias(con Consola) Proyecto {
	progreso := progressbar.Default(100, cadena.Cadena("Instalando dependencias: ").Italica().Colorear(color.GrisFuente).S())
	var buf bytes.Buffer
	escritor := bufio.NewWriter(io.MultiWriter(&buf, progreso))
	io.Copy(escritor, con.FSalida())
	for _, d := range p.Dependencias {
		if utiles.Limpiar(d) == "" || utiles.Limpiar(d) == " " {
			continue
		}
		progreso.Describe(cadena.Cadena(fmt.Sprintf("Instalando %s: ", d)).Italica().Colorear(color.GrisFuente).S())
		cmd := exec.Command("pip", "install", d)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("No se pudo instalar la dependencia: ", d, ".\n", err)
			//panic(err)
		}
		progreso.Clear()
		con.ImprimirCadena(cadena.Cadena(fmt.Sprintf("\t %s instalada.\n", d)))
		progreso.Add(100 / len(p.Dependencias))
		con.BorrarLinea()
	}

	progreso.Describe("Dependencias instaladas correctamente.")
	progreso.Finish()
	progreso.Clear()
	return p

}
func (p Proyecto) Crear(con Consola) Proyecto {

	return p.
		crearDirectorios(con).
		escribrirArchivos(con).
		instalarDependencias(con)

}

func (p *Proyecto) AgregarDependencias(deps []string) {
	p.Dependencias = append(p.Dependencias, deps...)
}

// Toma una cadena `d` con formato `[dependencia]==[version]`
func (p *Proyecto) AgregarDependencia(d string) {
	dependencia := utiles.Limpiar(d)
	p.Dependencias = append(p.Dependencias, dependencia)
}

func (p *Proyecto) Minificar() *Proyecto {
	ruta := path.Join(path.Join(path.Join(path.Join(path.Clean(p.Ruta), p.Nombre), "fuente"), p.Nombre), "estatico")
	err := fs.WalkDir(
		os.DirFS(ruta),
		"css",
		func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("err: %v\n", err)
				// return err
			} else {
				if filepath.Ext(p) == ".css" {
					cmd := exec.Command("npx", "lightningcss", "--minify", "--bundle", "--targets", ">= 0.25% and last 25 versions", filepath.Join(ruta, "css", path.Base(p)), "-o", filepath.Join(ruta, "min-css", path.Base(p)))
					min, err := cmd.Output()
					if err != nil {
						fmt.Printf("No se pudo minificar el css.\n%s\n", err)
						panic(1) // [HACER] Procesar Error
					}
					fmt.Println(min)
				}
			}
			return nil
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(420)
	}
	err = fs.WalkDir(
		os.DirFS(ruta),
		"js",
		func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("err: %v\n", err)
				// return err
			} else {
				if filepath.Ext(p) == ".js" {
					cmd := exec.Command("npx", "uglifyjs", filepath.Join(ruta, "js", path.Base(p)), "-o", filepath.Join(ruta, "min-js", path.Base(p)), "--compress", "--webkit")
					min, err := cmd.Output()
					if err != nil {
						fmt.Printf("No se pudo minificar el js.\n%s\n", err)
						panic(1) // [HACER] Procesar Error
					}
					fmt.Println(min)
				}
			}
			return nil
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	return p
}
