package main
import (
	"fmt"
	// "log"
	// "net/http"
	"os"
	// "os/exec"
	// "strings"

	"gopkg.in/yaml.v3"
)
type ServerConfig struct {
	WebServer string `yaml:"WebServer"`
	HTTPport  string `yaml:"HTTPport"`
	SSLcrt    string `yaml:"SSLcrt"`
	SSLkey    string `yaml:"SSLkey"`
}

// Load configuration file.
func LoadConfig(filePath string, config interface{}) (err error) {

        content, err := os.ReadFile(filePath)
        if err != nil {
                err = fmt.Errorf("os.ReadFile: %w", err)
                return err
        }

        content = []byte(os.ExpandEnv(string(content)))
        //      log.Printf("content=%s\n", content)

        if err := yaml.Unmarshal(content, config); err != nil {
                err = fmt.Errorf("yaml.Unmarshal(): %w", err)
                return err
        }

        //      log.Printf("\n")
        //      log.Printf("%+v\n", config)
        //      log.Printf("\n")

        return nil
}

