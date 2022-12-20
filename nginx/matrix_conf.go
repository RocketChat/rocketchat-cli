package nginx

import (
	"errors"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"log"
	"os"
	"rc-cli/filesystem"
	"strconv"
)

type MatrixConfLocationWellKnownReturnClause struct {
	Code    int
	Content string
}

type MatrixConfLocationWellKnown struct {
	AccessLog    string
	AddHeader    string
	DefaultType  string
	ReturnClause MatrixConfLocationWellKnownReturnClause
}

type MatrixConfLocationMatrix struct {
	ProxySetHeader   string
	ProxyPass        string
	ProxyReadTimeout string
}

type MatrixConfLocations struct {
	Matrix          MatrixConfLocationMatrix
	WellknownServer MatrixConfLocationWellKnown
	WellknownClient MatrixConfLocationWellKnown
}

type MatrixConfDef struct {
	Listen     int
	ServerName string
	Locations  MatrixConfLocations
}

var MatrixConf MatrixConfDef

// Internal file handler
var nginxConfigFile *gonginx.Config

func ReadMatrixConfFile(basePath string) {
	// Initialize file
	dirPath, filePath := getMatrixConfFilePath(basePath)

	// Ensure directory
	err := filesystem.EnsureDir(dirPath)
	if err != nil {
		log.Fatal("[ReadMatrixConfFile] Could not create directories: ", err)
	}

	// Create the default file if needed
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		_ = os.WriteFile(filePath, []byte(matrixConfDefaults), 0644)
	}

	// Read the config file
	p, err := parser.NewParser(filePath)
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	nginxConfigFile = p.Parse()

	serverBlock := nginxConfigFile.FindDirectives("server")[0].GetBlock()

	// Set the struct
	MatrixConf = MatrixConfDef{
		Listen:     nginxConfigGetInt(serverBlock, "Listen"),
		ServerName: nginxConfigGetString(serverBlock, "server_name"),
	}

	// Locations
	for _, locationBlock := range serverBlock.FindDirectives("location") {
		switch locationBlock.GetParameters()[0] {
		case "/.well-known/Matrix/server":
			MatrixConf.Locations.WellknownServer = MatrixConfLocationWellKnown{
				AccessLog:    nginxConfigGetString(locationBlock.GetBlock(), "access_log"),
				AddHeader:    nginxConfigGetString(locationBlock.GetBlock(), "add_header"),
				DefaultType:  nginxConfigGetString(locationBlock.GetBlock(), "default_type"),
				ReturnClause: nginxConfigGetReturnClause(locationBlock.GetBlock(), "return"),
			}
		case "/.well-known/Matrix/client":
			MatrixConf.Locations.WellknownClient = MatrixConfLocationWellKnown{
				AccessLog:    nginxConfigGetString(locationBlock.GetBlock(), "access_log"),
				AddHeader:    nginxConfigGetString(locationBlock.GetBlock(), "add_header"),
				DefaultType:  nginxConfigGetString(locationBlock.GetBlock(), "default_type"),
				ReturnClause: nginxConfigGetReturnClause(locationBlock.GetBlock(), "return"),
			}
		case "/":
			MatrixConf.Locations.Matrix = MatrixConfLocationMatrix{
				ProxySetHeader:   nginxConfigGetString(locationBlock.GetBlock(), "proxy_set_header"),
				ProxyPass:        nginxConfigGetString(locationBlock.GetBlock(), "proxy_pass"),
				ProxyReadTimeout: nginxConfigGetString(locationBlock.GetBlock(), "proxy_read_timeout"),
			}
		}
	}
}

func WriteMatrixConfFile(basePath string) {
	serverBlock := nginxConfigFile.FindDirectives("server")[0].GetBlock()

	// Base
	nginxConfigGetDirective(serverBlock, "Listen").GetParameters()[0] = strconv.Itoa(MatrixConf.Listen)
	nginxConfigGetDirective(serverBlock, "server_name").GetParameters()[0] = MatrixConf.ServerName
	// Location wellknown server
	locationServerBlock := serverBlock.FindDirectives("location")[0].GetBlock()
	nginxConfigGetDirective(locationServerBlock, "access_log").GetParameters()[0] = MatrixConf.Locations.WellknownServer.AccessLog
	nginxConfigGetDirective(locationServerBlock, "add_header").GetParameters()[0] = MatrixConf.Locations.WellknownServer.AddHeader
	nginxConfigGetDirective(locationServerBlock, "default_type").GetParameters()[0] = MatrixConf.Locations.WellknownServer.DefaultType
	nginxConfigGetDirective(locationServerBlock, "return").GetParameters()[0] = strconv.Itoa(MatrixConf.Locations.WellknownServer.ReturnClause.Code)
	nginxConfigGetDirective(locationServerBlock, "return").GetParameters()[1] = MatrixConf.Locations.WellknownServer.ReturnClause.Content
	// Location wellknown client
	locationClientBlock := serverBlock.FindDirectives("location")[1].GetBlock()
	nginxConfigGetDirective(locationClientBlock, "access_log").GetParameters()[0] = MatrixConf.Locations.WellknownClient.AccessLog
	nginxConfigGetDirective(locationClientBlock, "add_header").GetParameters()[0] = MatrixConf.Locations.WellknownClient.AddHeader
	nginxConfigGetDirective(locationClientBlock, "default_type").GetParameters()[0] = MatrixConf.Locations.WellknownClient.DefaultType
	nginxConfigGetDirective(locationClientBlock, "return").GetParameters()[0] = strconv.Itoa(MatrixConf.Locations.WellknownClient.ReturnClause.Code)
	nginxConfigGetDirective(locationClientBlock, "return").GetParameters()[1] = MatrixConf.Locations.WellknownClient.ReturnClause.Content
	// Location root
	locationRootBlock := serverBlock.FindDirectives("location")[2].GetBlock()
	nginxConfigGetDirective(locationRootBlock, "proxy_set_header").GetParameters()[0] = MatrixConf.Locations.Matrix.ProxySetHeader
	nginxConfigGetDirective(locationRootBlock, "proxy_pass").GetParameters()[0] = MatrixConf.Locations.Matrix.ProxyPass
	nginxConfigGetDirective(locationRootBlock, "proxy_read_timeout").GetParameters()[0] = MatrixConf.Locations.Matrix.ProxyReadTimeout

	// Write the config
	_, filePath := getMatrixConfFilePath(basePath)
	err := os.WriteFile(filePath, []byte(gonginx.DumpBlock(nginxConfigFile, gonginx.IndentedStyle)), 0644)
	if err != nil {
		log.Fatal("[WriteMatrixConfFile] Could not write file: ", err)
		return
	}
}
