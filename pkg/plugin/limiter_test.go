//go:build !race

package plugin

// func TestOutputLimiter(t *testing.T) {
// 	wg := new(sync.WaitGroup)

// 	input := test.NewTestInput()
// 	output := NewLimiter(test.NewTestOutput(func(*plugin.Message) {
// 		wg.Done()
// 	}), "10")
// 	wg.Add(10)

// 	plugins := &InOutPlugins{
// 		Inputs:  []Reader{input},
// 		Outputs: []Writer{output},
// 	}
// 	plugins.All = append(plugins.All, input, output)

// 	emitter := NewEmitter()
// 	go emitter.Start(plugins, "")

// 	for i := 0; i < 100; i++ {
// 		input.EmitGET()
// 	}

// 	wg.Wait()
// 	emitter.Close()
// }

// func TestInputLimiter(t *testing.T) {
// 	wg := new(sync.WaitGroup)

// 	input := NewLimiter(test.NewTestInput(), "10")
// 	output := test.NewTestOutput(func(*Message) {
// 		wg.Done()
// 	})
// 	wg.Add(10)

// 	plugins := &InOutPlugins{
// 		Inputs:  []Reader{input},
// 		Outputs: []Writer{output},
// 	}
// 	plugins.All = append(plugins.All, input, output)

// 	emitter := NewEmitter()
// 	go emitter.Start(plugins, Settings.Middleware)

// 	for i := 0; i < 100; i++ {
// 		input.(*Limiter).plugin.(*TestInput).EmitGET()
// 	}

// 	wg.Wait()
// 	emitter.Close()
// }

// // Should limit all requests
// func TestPercentLimiter1(t *testing.T) {
// 	wg := new(sync.WaitGroup)

// 	input := test.NewTestInput()
// 	output := NewLimiter(NewTestOutput(func(*Message) {
// 		wg.Done()
// 	}), "0%")

// 	plugins := &InOutPlugins{
// 		Inputs:  []Reader{input},
// 		Outputs: []Writer{output},
// 	}
// 	plugins.All = append(plugins.All, input, output)

// 	emitter := NewEmitter()
// 	go emitter.Start(plugins, Settings.Middleware)

// 	for i := 0; i < 100; i++ {
// 		input.EmitGET()
// 	}

// 	wg.Wait()
// }

// // Should not limit at all
// func TestPercentLimiter2(t *testing.T) {
// 	wg := new(sync.WaitGroup)

// 	input := test.NewTestInput()
// 	output := NewLimiter(NewTestOutput(func(*Message) {
// 		wg.Done()
// 	}), "100%")
// 	wg.Add(100)

// 	plugins := &InOutPlugins{
// 		Inputs:  []Reader{input},
// 		Outputs: []Writer{output},
// 	}
// 	plugins.All = append(plugins.All, input, output)

// 	emitter := NewEmitter()
// 	go emitter.Start(plugins, Settings.Middleware)

// 	for i := 0; i < 100; i++ {
// 		input.EmitGET()
// 	}

// 	wg.Wait()
// }
