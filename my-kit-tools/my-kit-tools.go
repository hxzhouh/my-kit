package main

import (
	"bytes"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"html/template"
	"log"
)

// register
func init() {
	generator.RegisterPlugin(new(netRpcPlugin))
}

type netRpcPlugin struct {
	*generator.Generator
}

//首先Name()方法返回插件的名字
func (t *netRpcPlugin) Name() string {
	return "netrpc"
}

// Init is called once after data structures are built but before
// code generation begins. Init()初始化的时候用参数g进行初始化，因此插件是从参数g对象继承了全部的公有方法。
func (t *netRpcPlugin) Init(g *generator.Generator) {
	t.Generator = g
}

// Generate produces the code generated by the plugin for this file,
// except for the imports, by calling the generator's methods P, In, and Out.
// Generate()方法调用自定义的genServiceCode()方法生成每个服务的代码。
func (t *netRpcPlugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		t.genServiceCode(svc)
	}
}

// GenerateImports produces the import declarations for this file.
// It is called after Generate. 其中GenerateImports()方法调用自定义的genImportCode()方法生成导入代码.
func (t *netRpcPlugin) GenerateImports(file *generator.FileDescriptor) {
	if len(file.Service) > 0 {
		t.genImportCode(file)
	}
}

func (t *netRpcPlugin) genImportCode(file *generator.FileDescriptor) {
	//t.P(`import "net/rpc"`)
}

func (t *netRpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
	spec := t.buildServiceSpec(svc)
	var buf bytes.Buffer
	p := template.Must(template.New("").Parse(tmplService))
	err := p.Execute(&buf, spec)
	if err != nil {
		log.Fatal(err)
	}
	t.P(buf.String())
}

//其中输入参数是*descriptor.ServiceDescriptorProto类型，完整描述了一个服务的所有信息。然后通过svc.GetName()就可以获取Protobuf文件中定义的服务的名字。
// Protobuf文件中的名字转为Go语言的名字后，需要通过generator.CamelCase()函数进行一次转换。
// 类似地，在for循环中我们通过m.GetName()获取方法的名字，然后再转为Go语言中对应的名字。
// 比较复杂的是对输入和输出参数名字的解析：首先需要通过m.GetInputType()获取输入参数的类型，然后通过p.ObjectNamed()获取类型对应的类对象信息，最后获取类对象的名字。
func (p *netRpcPlugin) buildServiceSpec(
	svc *descriptor.ServiceDescriptorProto) *ServiceSpec {
	spec := &ServiceSpec{
		ServiceName: generator.CamelCase(svc.GetName()),
	}
	for _, m := range svc.Method {
		spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
			MethodName:     generator.CamelCase(m.GetName()),
			InputTypeName:  p.TypeName(p.ObjectNamed(m.GetInputType())),
			OutputTypeName: p.TypeName(p.ObjectNamed(m.GetOutputType())),
		})
	}
	return spec
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

type ServiceSpec struct {
	ServiceName string
	MethodList  []ServiceMethodSpec
}

const tmplService = `
{{$root := .}}
type {{.ServiceName}}Interface interface {
    {{- range $_, $m := .MethodList}}
    {{$m.MethodName}}(*{{$m.InputTypeName}}, *{{$m.OutputTypeName}}) error
    {{- end}}
}
func Register{{.ServiceName}}(
    srv *rpc.Server, x {{.ServiceName}}Interface,
) error {
    if err := srv.RegisterName("{{.ServiceName}}", x); err != nil {
        return err
    }
    return nil
}
type {{.ServiceName}}Client struct {
    *rpc.Client
}
var _ {{.ServiceName}}Interface = (*{{.ServiceName}}Client)(nil)
func Dial{{.ServiceName}}(network, address string) (
    *{{.ServiceName}}Client, error,
) {
    c, err := rpc.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &{{.ServiceName}}Client{Client: c}, nil
}
{{range $_, $m := .MethodList}}
func (p *{{$root.ServiceName}}Client) {{$m.MethodName}}(
    in *{{$m.InputTypeName}}, out *{{$m.OutputTypeName}},
) error {
    return p.Client.Call("{{$root.ServiceName}}.{{$m.MethodName}}", in, out)
}
{{end}}
`