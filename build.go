package main

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	windows	= "windows"
	linux	= "linux"
)

var (
	goarch			string
	goos			string
	gocc			string
	cgo			bool
	pkgArch			string
	version			string	= "v1"
	linuxPackageVersion	string	= "v1"
	linuxPackageIteration	string	= ""
	race			bool
	phjsToRelease		string
	workingDir		string
	includeBuildId		bool		= true
	buildId			string		= "0"
	binaries		[]string	= []string{"grafana-server", "grafana-cli"}
	isDev			bool		= false
	enterprise		bool		= false
)

func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	ensureGoPath()
	var buildIdRaw string
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "GOARCH")
	flag.StringVar(&goos, "goos", runtime.GOOS, "GOOS")
	flag.StringVar(&gocc, "cc", "", "CC")
	flag.BoolVar(&cgo, "cgo-enabled", cgo, "Enable cgo")
	flag.StringVar(&pkgArch, "pkg-arch", "", "PKG ARCH")
	flag.StringVar(&phjsToRelease, "phjs", "", "PhantomJS binary")
	flag.BoolVar(&race, "race", race, "Use race detector")
	flag.BoolVar(&includeBuildId, "includeBuildId", includeBuildId, "IncludeBuildId in package name")
	flag.BoolVar(&enterprise, "enterprise", enterprise, "Build enterprise version of Grafana")
	flag.StringVar(&buildIdRaw, "buildId", "0", "Build ID from CI system")
	flag.BoolVar(&isDev, "dev", isDev, "optimal for development, skips certain steps")
	flag.Parse()
	buildId = shortenBuildId(buildIdRaw)
	readVersionFromPackageJson()
	if pkgArch == "" {
		pkgArch = goarch
	}
	log.Printf("Version: %s, Linux Version: %s, Package Iteration: %s\n", version, linuxPackageVersion, linuxPackageIteration)
	if flag.NArg() == 0 {
		log.Println("Usage: go run build.go build")
		return
	}
	workingDir, _ = os.Getwd()
	for _, cmd := range flag.Args() {
		switch cmd {
		case "setup":
			setup()
		case "build-srv":
			clean()
			build("grafana-server", "./pkg/cmd/grafana-server", []string{})
		case "build-cli":
			clean()
			build("grafana-cli", "./pkg/cmd/grafana-cli", []string{})
		case "build-server":
			clean()
			build("grafana-server", "./pkg/cmd/grafana-server", []string{})
		case "build":
			for _, binary := range binaries {
				build(binary, "./pkg/cmd/"+binary, []string{})
			}
		case "build-frontend":
			grunt(gruntBuildArg("build")...)
		case "test":
			test("./pkg/...")
			grunt("test")
		case "package":
			grunt(gruntBuildArg("build")...)
			grunt(gruntBuildArg("package")...)
			if goos == linux {
				createLinuxPackages()
			}
		case "package-only":
			grunt(gruntBuildArg("package")...)
			if goos == linux {
				createLinuxPackages()
			}
		case "pkg-archive":
			grunt(gruntBuildArg("package")...)
		case "pkg-rpm":
			grunt(gruntBuildArg("release")...)
			createRpmPackages()
		case "pkg-deb":
			grunt(gruntBuildArg("release")...)
			createDebPackages()
		case "sha-dist":
			shaFilesInDist()
		case "latest":
			makeLatestDistCopies()
		case "clean":
			clean()
		default:
			log.Fatalf("Unknown command %q", cmd)
		}
	}
}
func makeLatestDistCopies() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	files, err := ioutil.ReadDir("dist")
	if err != nil {
		log.Fatalf("failed to create latest copies. Cannot read from /dist")
	}
	latestMapping := map[string]string{"_amd64.deb": "dist/grafana_latest_amd64.deb", ".x86_64.rpm": "dist/grafana-latest-1.x86_64.rpm", ".linux-amd64.tar.gz": "dist/grafana-latest.linux-x64.tar.gz", ".linux-armv7.tar.gz": "dist/grafana-latest.linux-armv7.tar.gz", ".linux-arm64.tar.gz": "dist/grafana-latest.linux-arm64.tar.gz"}
	for _, file := range files {
		for extension, fullName := range latestMapping {
			if strings.HasSuffix(file.Name(), extension) {
				runError("cp", path.Join("dist", file.Name()), fullName)
			}
		}
	}
}
func readVersionFromPackageJson() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reader, err := os.Open("package.json")
	if err != nil {
		log.Fatal("Failed to open package.json")
		return
	}
	defer reader.Close()
	jsonObj := map[string]interface{}{}
	jsonParser := json.NewDecoder(reader)
	if err := jsonParser.Decode(&jsonObj); err != nil {
		log.Fatal("Failed to decode package.json")
	}
	version = jsonObj["version"].(string)
	linuxPackageVersion = version
	linuxPackageIteration = ""
	parts := strings.Split(version, "-")
	if len(parts) > 1 {
		linuxPackageVersion = parts[0]
		linuxPackageIteration = parts[1]
	}
	if includeBuildId {
		if buildId != "0" {
			linuxPackageIteration = fmt.Sprintf("%s%s", buildId, linuxPackageIteration)
		} else {
			linuxPackageIteration = fmt.Sprintf("%d%s", time.Now().Unix(), linuxPackageIteration)
		}
	}
}

type linuxPackageOptions struct {
	packageType		string
	homeDir			string
	binPath			string
	serverBinPath		string
	cliBinPath		string
	configDir		string
	ldapFilePath		string
	etcDefaultPath		string
	etcDefaultFilePath	string
	initdScriptFilePath	string
	systemdServiceFilePath	string
	postinstSrc		string
	initdScriptSrc		string
	defaultFileSrc		string
	systemdFileSrc		string
	depends			[]string
}

func createDebPackages() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	previousPkgArch := pkgArch
	if pkgArch == "armv7" {
		pkgArch = "armhf"
	}
	createPackage(linuxPackageOptions{packageType: "deb", homeDir: "/usr/share/grafana", binPath: "/usr/sbin", configDir: "/etc/grafana", etcDefaultPath: "/etc/default", etcDefaultFilePath: "/etc/default/grafana-server", initdScriptFilePath: "/etc/init.d/grafana-server", systemdServiceFilePath: "/usr/lib/systemd/system/grafana-server.service", postinstSrc: "packaging/deb/control/postinst", initdScriptSrc: "packaging/deb/init.d/grafana-server", defaultFileSrc: "packaging/deb/default/grafana-server", systemdFileSrc: "packaging/deb/systemd/grafana-server.service", depends: []string{"adduser", "libfontconfig"}})
	pkgArch = previousPkgArch
}
func createRpmPackages() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	previousPkgArch := pkgArch
	switch {
	case pkgArch == "armv7":
		pkgArch = "armhfp"
	case pkgArch == "arm64":
		pkgArch = "aarch64"
	}
	createPackage(linuxPackageOptions{packageType: "rpm", homeDir: "/usr/share/grafana", binPath: "/usr/sbin", configDir: "/etc/grafana", etcDefaultPath: "/etc/sysconfig", etcDefaultFilePath: "/etc/sysconfig/grafana-server", initdScriptFilePath: "/etc/init.d/grafana-server", systemdServiceFilePath: "/usr/lib/systemd/system/grafana-server.service", postinstSrc: "packaging/rpm/control/postinst", initdScriptSrc: "packaging/rpm/init.d/grafana-server", defaultFileSrc: "packaging/rpm/sysconfig/grafana-server", systemdFileSrc: "packaging/rpm/systemd/grafana-server.service", depends: []string{"/sbin/service", "fontconfig", "freetype", "urw-fonts"}})
	pkgArch = previousPkgArch
}
func createLinuxPackages() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	createDebPackages()
	createRpmPackages()
}
func createPackage(options linuxPackageOptions) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	packageRoot, _ := ioutil.TempDir("", "grafana-linux-pack")
	runPrint("mkdir", "-p", filepath.Join(packageRoot, options.homeDir))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, options.configDir))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, "/etc/init.d"))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, options.etcDefaultPath))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, "/usr/lib/systemd/system"))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, "/usr/sbin"))
	for _, binary := range binaries {
		runPrint("cp", "-p", filepath.Join(workingDir, "tmp/bin/"+binary), filepath.Join(packageRoot, "/usr/sbin/"+binary))
	}
	runPrint("cp", "-p", options.initdScriptSrc, filepath.Join(packageRoot, options.initdScriptFilePath))
	runPrint("cp", "-p", options.defaultFileSrc, filepath.Join(packageRoot, options.etcDefaultFilePath))
	runPrint("cp", "-p", options.systemdFileSrc, filepath.Join(packageRoot, options.systemdServiceFilePath))
	runPrint("cp", "-a", filepath.Join(workingDir, "tmp")+"/.", filepath.Join(packageRoot, options.homeDir))
	runPrint("rm", "-rf", filepath.Join(packageRoot, options.homeDir, "bin"))
	args := []string{"-s", "dir", "--description", "Grafana", "-C", packageRoot, "--url", "https://grafana.com", "--maintainer", "contact@grafana.com", "--config-files", options.initdScriptFilePath, "--config-files", options.etcDefaultFilePath, "--config-files", options.systemdServiceFilePath, "--after-install", options.postinstSrc, "--version", linuxPackageVersion, "-p", "./dist"}
	name := "grafana"
	if enterprise {
		name += "-enterprise"
		args = append(args, "--replaces", "grafana")
	}
	args = append(args, "--name", name)
	description := "Grafana"
	if enterprise {
		description += " Enterprise"
	}
	args = append(args, "--vendor", description)
	if !enterprise {
		args = append(args, "--license", "\"Apache 2.0\"")
	}
	if options.packageType == "rpm" {
		args = append(args, "--rpm-posttrans", "packaging/rpm/control/posttrans")
	}
	if options.packageType == "deb" {
		args = append(args, "--deb-no-default-config-files")
	}
	if pkgArch != "" {
		args = append(args, "-a", pkgArch)
	}
	if linuxPackageIteration != "" {
		args = append(args, "--iteration", linuxPackageIteration)
	}
	for _, dep := range options.depends {
		args = append(args, "--depends", dep)
	}
	args = append(args, ".")
	fmt.Println("Creating package: ", options.packageType)
	runPrint("fpm", append([]string{"-t", options.packageType}, args...)...)
}
func ensureGoPath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if os.Getenv("GOPATH") == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		gopath := filepath.Clean(filepath.Join(cwd, "../../../../"))
		log.Println("GOPATH is", gopath)
		os.Setenv("GOPATH", gopath)
	}
}
func grunt(params ...string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if runtime.GOOS == windows {
		runPrint(`.\node_modules\.bin\grunt`, params...)
	} else {
		runPrint("./node_modules/.bin/grunt", params...)
	}
}
func gruntBuildArg(task string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := []string{task}
	if includeBuildId {
		args = append(args, fmt.Sprintf("--pkgVer=%v-%v", linuxPackageVersion, linuxPackageIteration))
	} else {
		args = append(args, fmt.Sprintf("--pkgVer=%v", version))
	}
	if pkgArch != "" {
		args = append(args, fmt.Sprintf("--arch=%v", pkgArch))
	}
	if phjsToRelease != "" {
		args = append(args, fmt.Sprintf("--phjsToRelease=%v", phjsToRelease))
	}
	if enterprise {
		args = append(args, "--enterprise")
	}
	args = append(args, fmt.Sprintf("--platform=%v", goos))
	return args
}
func setup() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	runPrint("go", "get", "-v", "github.com/golang/dep")
	runPrint("go", "install", "-v", "./pkg/cmd/grafana-server")
}
func test(pkg string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	setBuildEnv()
	runPrint("go", "test", "-short", "-timeout", "60s", pkg)
}
func build(binaryName, pkg string, tags []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	binary := fmt.Sprintf("./bin/%s-%s/%s", goos, goarch, binaryName)
	if isDev {
		binary = fmt.Sprintf("./bin/%s", binaryName)
	}
	if goos == windows {
		binary += ".exe"
	}
	if !isDev {
		rmr(binary, binary+".md5")
	}
	args := []string{"build", "-ldflags", ldflags()}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	if race {
		args = append(args, "-race")
	}
	args = append(args, "-o", binary)
	args = append(args, pkg)
	if !isDev {
		setBuildEnv()
		runPrint("go", "version")
		fmt.Printf("Targeting %s/%s\n", goos, goarch)
	}
	runPrint("go", args...)
	if !isDev {
		err := md5File(binary)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func ldflags() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var b bytes.Buffer
	b.WriteString("-w")
	b.WriteString(fmt.Sprintf(" -X main.version=%s", version))
	b.WriteString(fmt.Sprintf(" -X main.commit=%s", getGitSha()))
	b.WriteString(fmt.Sprintf(" -X main.buildstamp=%d", buildStamp()))
	b.WriteString(fmt.Sprintf(" -X main.buildBranch=%s", getGitBranch()))
	return b.String()
}
func rmr(paths ...string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, path := range paths {
		log.Println("rm -r", path)
		os.RemoveAll(path)
	}
}
func clean() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if isDev {
		return
	}
	rmr("dist")
	rmr("tmp")
	rmr(filepath.Join(os.Getenv("GOPATH"), fmt.Sprintf("pkg/%s_%s/github.com/grafana", goos, goarch)))
}
func setBuildEnv() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os.Setenv("GOOS", goos)
	if goos == windows {
		os.Setenv("CGO_CFLAGS", "-D_WIN32_WINNT=0x0601")
	}
	if goarch != "amd64" || goos != linux {
		cgo = true
	}
	if strings.HasPrefix(goarch, "armv") {
		os.Setenv("GOARCH", "arm")
		os.Setenv("GOARM", goarch[4:])
	} else {
		os.Setenv("GOARCH", goarch)
	}
	if goarch == "386" {
		os.Setenv("GO386", "387")
	}
	if cgo {
		os.Setenv("CGO_ENABLED", "1")
	}
	if gocc != "" {
		os.Setenv("CC", gocc)
	}
}
func getGitBranch() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	v, err := runError("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "master"
	}
	return string(v)
}
func getGitSha() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	v, err := runError("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "unknown-dev"
	}
	return string(v)
}
func buildStamp() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bs, err := runError("git", "show", "-s", "--format=%ct")
	if err != nil {
		return time.Now().Unix()
	}
	s, _ := strconv.ParseInt(string(bs), 10, 64)
	return s
}
func runError(cmd string, args ...string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ecmd := exec.Command(cmd, args...)
	bs, err := ecmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(bs), nil
}
func runPrint(cmd string, args ...string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log.Println(cmd, strings.Join(args, " "))
	ecmd := exec.Command(cmd, args...)
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	err := ecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func md5File(file string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	h := md5.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return err
	}
	out, err := os.Create(file + ".md5")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%x\n", h.Sum(nil))
	if err != nil {
		return err
	}
	return out.Close()
}
func shaFilesInDist() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	filepath.Walk("./dist", func(path string, f os.FileInfo, err error) error {
		if path == "./dist" {
			return nil
		}
		if !strings.Contains(path, ".sha256") {
			err := shaFile(path)
			if err != nil {
				log.Printf("Failed to create sha file. error: %v\n", err)
			}
		}
		return nil
	})
}
func shaFile(file string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	h := sha256.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return err
	}
	out, err := os.Create(file + ".sha256")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%x\n", h.Sum(nil))
	if err != nil {
		return err
	}
	return out.Close()
}
func shortenBuildId(buildId string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	buildId = strings.Replace(buildId, "-", "", -1)
	if len(buildId) < 9 {
		return buildId
	}
	return buildId[0:8]
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
