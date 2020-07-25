package main

import (
	"io"
	"reflect"
	"strings"
)

// InOutPlugins struct for holding references to plugins
type InOutPlugins struct {
	Inputs  []io.Reader
	Outputs []io.Writer
	All     []interface{}
}

// extractLimitOptions detects if plugin get called with limiter support
// Returns address and limit
func extractLimitOptions(options string) (string, string) {
	split := strings.Split(options, "|")

	if len(split) > 1 {
		return split[0], split[1]
	}

	return split[0], ""
}

// Automatically detects type of plugin and initialize it
//
// See this article if curious about reflect stuff below: http://blog.burntsushi.net/type-parametric-functions-golang
func (plugins *InOutPlugins) registerPlugin(constructor interface{}, options ...interface{}) {
	var path, limit string
	vc := reflect.ValueOf(constructor)

	// Pre-processing options to make it work with reflect
	vo := []reflect.Value{}
	for _, oi := range options {
		vo = append(vo, reflect.ValueOf(oi))
	}

	if len(vo) > 0 {
		// Removing limit options from path
		path, limit = extractLimitOptions(vo[0].String())

		// Writing value back without limiter "|" options
		vo[0] = reflect.ValueOf(path)
	}

	// Calling our constructor with list of given options
	plugin := vc.Call(vo)[0].Interface()
	pluginWrapper := plugin

	if limit != "" {
		pluginWrapper = NewLimiter(plugin, limit)
	} else {
		pluginWrapper = plugin
	}

	_, isR := plugin.(io.Reader)
	_, isW := plugin.(io.Writer)

	// Some of the output can be Readers as well because return responses
	if isR && !isW {
		plugins.Inputs = append(plugins.Inputs, pluginWrapper.(io.Reader))
	}

	if isW {
		plugins.Outputs = append(plugins.Outputs, pluginWrapper.(io.Writer))
	}

	plugins.All = append(plugins.All, plugin)
}

// NewPlugins specify and initialize all available plugins
func NewPlugins() *InOutPlugins {
	plugins := new(InOutPlugins)

	for _, options := range Settings.InputDummy {
		plugins.registerPlugin(NewDummyInput, options)
	}

	for range Settings.OutputDummy {
		plugins.registerPlugin(NewDummyOutput)
	}

	if Settings.OutputStdout {
		plugins.registerPlugin(NewDummyOutput)
	}

	if Settings.OutputNull {
		plugins.registerPlugin(NewNullOutput)
	}

	for _, options := range Settings.InputRAW {
		plugins.registerPlugin(NewRAWInput, options, Settings.RAWInputConfig)
	}

	for _, options := range Settings.InputTCP {
		plugins.registerPlugin(NewTCPInput, options, &Settings.InputTCPConfig)
	}

	for _, options := range Settings.OutputTCP {
		plugins.registerPlugin(NewTCPOutput, options, &Settings.OutputTCPConfig)
	}

	for _, options := range Settings.InputFile {
		plugins.registerPlugin(NewFileInput, options, Settings.InputFileLoop)
	}

	for _, path := range Settings.OutputFile {
		if strings.HasPrefix(path, "s3://") {
			plugins.registerPlugin(NewS3Output, path, &Settings.OutputFileConfig)
		} else {
			plugins.registerPlugin(NewFileOutput, path, &Settings.OutputFileConfig)
		}
	}

	for _, options := range Settings.InputHTTP {
		plugins.registerPlugin(NewHTTPInput, options)
	}

	// If we explicitly set Host header http output should not rewrite it
	// Fix: https://github.com/buger/gor/issues/174
	for _, header := range Settings.ModifierConfig.Headers {
		if header.Name == "Host" {
			Settings.OutputHTTPConfig.OriginalHost = true
			break
		}
	}

	for _, options := range Settings.OutputHTTP {
		plugins.registerPlugin(NewHTTPOutput, options, &Settings.OutputHTTPConfig)
	}

	for _, options := range Settings.OutputBinary {
		plugins.registerPlugin(NewBinaryOutput, options, &Settings.OutputBinaryConfig)
	}

	if Settings.OutputKafkaConfig.Host != "" && Settings.OutputKafkaConfig.Topic != "" {
		plugins.registerPlugin(NewKafkaOutput, "", &Settings.OutputKafkaConfig)
	}

	if Settings.InputKafkaConfig.Host != "" && Settings.InputKafkaConfig.Topic != "" {
		plugins.registerPlugin(NewKafkaInput, "", &Settings.InputKafkaConfig)
	}

	return plugins
}
