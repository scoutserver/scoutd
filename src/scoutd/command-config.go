package scoutd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
)

const yamlTemplate = `
# generated by 'scout config'
{{ if .current.AccountKey }}account_key: {{ .current.AccountKey }}{{ end }}
{{ if ne .current.HostName .default.HostName }}hostname: {{ .current.HostName }}{{ end }}
{{ if ne .current.RunDir .default.RunDir }}run_dir: {{ .current.RunDir }}{{ end }}
{{ if ne .current.LogFile .default.LogFile }}log_file: {{ .current.LogFile }}{{ end }}
{{ if ne .current.RubyPath .default.RubyPath }}ruby_path: {{ .current.RubyPath }}{{ end }}
{{ if ne .current.AgentRubyBin .default.AgentRubyBin }}agent_ruby_bin: {{ .current.AgentRubyBin }}{{ end }}
{{ if ne .current.AgentEnv .default.AgentEnv }}environment: {{ .current.AgentEnv }}{{ end }}
{{ if ne .current.AgentRoles .default.AgentRoles }}roles: {{ .current.AgentRoles }}{{ end }}
{{ if ne .current.AgentDisplayName .default.AgentDisplayName }}display_name: {{ .current.AgentDisplayName }}{{ end }}
{{ if ne .current.AgentDataFile .default.AgentDataFile }}agent_data_file: {{ .current.AgentDataFile }}{{ end }}
{{ if ne .current.HttpProxyUrl .default.HttpProxyUrl }}http_proxy: {{ .current.HttpProxyUrl }}{{ end }}
{{ if ne .current.HttpsProxyUrl .default.HttpsProxyUrl }}https_proxy: {{ .current.HttpsProxyUrl }}{{ end }}
{{ if .statsd }}statsd:{{ end }}
{{ if .statsd }}  enabled: {{ .statsd.Statsd.Enabled }}{{ end }}
{{ if .statsd }}  addr: {{ .statsd.Statsd.Addr }}{{ end }}
{{ if ne .current.ReportingServerUrl .default.ReportingServerUrl }}reporting_server_url: {{ .current.ReportingServerUrl }}{{ end }}
`

func GenConfig(cfg ScoutConfig) {
	var buf bytes.Buffer
	var defaultCfg = LoadDefaults()
	t := template.Must(template.New("config").Parse(yamlTemplate))
	configMap := map[string]ScoutConfig{
		"current": cfg,
		"default": defaultCfg,
	}
	if !reflect.DeepEqual(defaultCfg.Statsd, cfg.Statsd) {
		configMap["statsd"] = cfg
	}
	err := t.Execute(&buf, configMap)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
	// Use a temporary buffer to remove any empty lines generated from the template
	var tmpbuf bytes.Buffer
	var readErr error
	for line := ""; readErr == nil; {
		line, readErr = buf.ReadString('\n')
		if line != "\n" {
			tmpbuf.WriteString(line)
		}
	}
	var s = tmpbuf.String() // The final configration string to write out
	if genCfgOptions.Outfile != "" {
		var out string
		if genCfgOptions.Outfile == "DEFAULT_VALUE" {
			out = cfg.ConfigFile
		} else {
			out = genCfgOptions.Outfile
		}
		if _, err := os.Stat(out); os.IsNotExist(err) {
			WriteConfig(out, s)
		} else {
			if genCfgOptions.AssumeYes {
				WriteConfig(out, s)
			} else {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("File '%s' exists. Overwrite? [y/N] ", out)
				resp, _ := reader.ReadString('\n')
				if strings.ToUpper(string(resp[0])) == "Y" {
					WriteConfig(out, s)
				}
			}
		}
	} else {
		fmt.Printf("\n\n%s\n\n", s)
	}
}

func WriteConfig(filePath string, cfg string) {
	err := ioutil.WriteFile(filePath, []byte(cfg), 0660)
	if err != nil {
		log.Fatalf("Error writing to %s: %s", filePath, err)
	}
}