package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// Define una estructura para las cámaras
type Camera struct {
	Number int // Número o identificador de la cámara
	Name   string
	URL    string
}

// Lista de cámaras para grabar
var cameras = []Camera{
	{1, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.11:8554/profile0"},
	{2, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.12:8554/profile0"},
	{3, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.13:8554/profile0"},
	{4, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.14:8554/profile0"},
	{5, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.15:8554/profile0"},
	{6, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.16:8554/profile0"},

	{7, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.17:8554/profile0"},
	{8, "CámaraVideos", "rtsp://admin:H2FuiDp4@10.33.1.18:8554/profile0"},
}

var rootCmd = &cobra.Command{
	Use:   "camcontrol",
	Short: "Controlador de grabación de cámaras RTSP",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia la grabación de todas las cámaras",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Iniciando grabaciones...")
		startRecordings()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func startRecordings() {
	var wg sync.WaitGroup
	for _, cam := range cameras {
		wg.Add(1)
		go func(camera Camera) {
			defer wg.Done()
			startRecordingForCamera(camera)
		}(cam)
	}
	wg.Wait() // Espera a que todas las grabaciones se inicien
}

func startRecordingForCamera(cam Camera) {
	currentTime := time.Now()
	// Reemplaza espacios en el nombre de la cámara con guiones bajos o cualquier otro caracter que prefieras
	dirName := fmt.Sprintf("Camara%d", cam.Number) // Asume que `Number` es un campo nuevo en tu estructura `Camera` que indica el número de la cámara.
	fileName := fmt.Sprintf("%s_%d-%02d-%02d_%02d-%02d-%02d.mp4",
		cam.Name,
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second())

	// Crea el directorio si no existe
	path := fmt.Sprintf("./videos/%s", dirName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm) // Usa MkdirAll para crear todos los directorios necesarios en la ruta
		if err != nil {
			fmt.Printf("❌ Error al crear el directorio para %s: %v\n", cam.Name, err)
			return
		}
	}

	fullPath := fmt.Sprintf("%s/%s", path, fileName)

	err := ffmpeg.Input(cam.URL).
		Output(fullPath, ffmpeg.KwArgs{"c:v": "copy"}). // Ajusta según sea necesario
		OverWriteOutput().
		Run()

	if err != nil {
		fmt.Printf("❌ Error al grabar %s: %v\n", cam.Name, err)
	} else {
		fmt.Printf("✅ Grabación completada para %s en %s\n", cam.Name, fullPath)
	}
}
